package logger

import (
	"fmt"
	"os"
	"path"
	"time"
)

// FileLogger 主结构体
type FileLogger struct {
	Level         LogLevel
	filePath      string
	fileName      string
	errorFileName string
	maxFileSize   int64
	fileObj       *os.File
	errFileObj    *os.File
	logChan       chan *logMsg
	async         bool
}

// 单条日志chan结构体
type logMsg struct {
	level     LogLevel
	msg       string
	funcName  string
	fileName  string
	timestamp string
	line      int
}

// 构造函数
func NewFileLogger(levelStr, path string, async bool) *FileLogger {
	logLevel, err := parseLogLevel(levelStr)
	if err != nil {
		panic(err)
	}

	fl := &FileLogger{
		Level:       logLevel,
		filePath:    path,
		fileName:    time.Now().Format("2006-01-02") + ".log",
		maxFileSize: 1024 * 1024 * 10,
		logChan:     make(chan *logMsg, 1000),
		async:       async,
	}
	err = fl.initFile()
	if err != nil {
		panic(err)
	}
	return fl
}

// 初始化本地文件
func (f *FileLogger) initFile() error {
	err := os.MkdirAll(f.filePath, 0644)
	if err != nil {
		fmt.Println(err)
		return err
	}

	fullName := path.Join(f.filePath, f.fileName)

	fileObj, err := os.OpenFile(fullName, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println(err)
		return err
	}

	errFileObj, err := os.OpenFile(fullName+".err", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println(err)
		return err
	}

	//2个日志变量赋值
	f.fileObj = fileObj
	f.errFileObj = errFileObj

	//设置5个消费者，监听logChan
	for i := 0; i < 5; i++ {
		go f.asyncWriter()
	}

	return nil
}

// 关闭文件
func (f *FileLogger) Close() {
	f.fileObj.Close()
	f.errFileObj.Close()
}

// 检查文件大小是否超出，超出则需要切割
func (f *FileLogger) checkSize(file *os.File) bool {
	fileInfo, err := file.Stat()
	if err != nil {
		return false
	}
	return fileInfo.Size() >= f.maxFileSize
}

// 切割文件
func (f *FileLogger) splitFile(file *os.File) (*os.File, error) {
	//需要切割日志文件
	nowStr := time.Now().Format("20060102150405")
	fileInfo, err := file.Stat()
	if err != nil {
		return nil, err
	}
	oldLogName := path.Join(f.filePath, fileInfo.Name())
	newLogName := fmt.Sprintf("%s.bak%s", oldLogName, nowStr)

	//1.关闭当前日志文件
	file.Close()

	//2.备份旧文件 rename xx.log.back201908311709
	os.Rename(oldLogName, newLogName)

	//3.打开一个新的日志文件
	fileObj, err := os.OpenFile(oldLogName, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	//4.将新的文件对象赋值给f.fileObj
	return fileObj, nil
}

// 同步写入
func (f *FileLogger) log(lv LogLevel, format string, a ...interface{}) {
	if f.enable(lv) {
		funcName, fileName, lineNo := getInfo(3)

		if f.checkSize(f.fileObj) {
			newFile, err := f.splitFile(f.fileObj)
			if err != nil {
				fmt.Println(err)
				return
			}
			f.fileObj = newFile
		}
		fmt.Fprintf(f.fileObj,
			"[%s] [%s] [%s:%s:%d] %s\n",
			getTime(),
			getLogString(lv),
			fileName,
			funcName,
			lineNo,
			fmt.Sprintf(format, a...))

		//单独记录错误日志
		if lv >= ERROR {
			if f.checkSize(f.errFileObj) {
				newFile, err := f.splitFile(f.errFileObj)
				if err != nil {
					return
				}
				f.errFileObj = newFile
			}
			fmt.Fprintf(f.errFileObj,
				"[%s] [%s] [%s:%s:%d] %s\n",
				getTime(),
				getLogString(lv),
				fileName,
				funcName,
				lineNo,
				fmt.Sprintf(format, a...))
		}
	}
}

// 写入logChan
func (f *FileLogger) logAsync(lv LogLevel, format string, a ...interface{}) {
	if f.enable(lv) {
		funcName, fileName, lineNo := getInfo(3)

		//先把日志发送到logChan中,用select可以防止通道满了之后阻塞
		select {
		case f.logChan <- &logMsg{
			level:     lv,
			msg:       fmt.Sprintf(format, a...),
			funcName:  funcName,
			fileName:  fileName,
			timestamp: getTime(),
			line:      lineNo,
		}:
		default:
		}
	}
}

//启动异步消费者，监听logChan
func (f *FileLogger) asyncWriter() {
	for {
		select {
		case logTmp := <-f.logChan:
			if f.checkSize(f.fileObj) {
				newFile, err := f.splitFile(f.fileObj)
				if err != nil {
					fmt.Println(err)
					return
				}
				f.fileObj = newFile
			}
			fmt.Fprintf(f.fileObj,
				"[%s] [%s] [%s:%s:%d] %s\n",
				logTmp.timestamp,
				getLogString(logTmp.level),
				logTmp.fileName,
				logTmp.funcName,
				logTmp.line,
				logTmp.msg)

			//单独记录错误日志
			if logTmp.level >= ERROR {
				if f.checkSize(f.errFileObj) {
					newFile, err := f.splitFile(f.errFileObj)
					if err != nil {
						return
					}
					f.errFileObj = newFile
				}
				fmt.Fprintf(f.errFileObj,
					"[%s] [%s] [%s:%s:%d] %s\n",
					logTmp.timestamp,
					getLogString(logTmp.level),
					logTmp.fileName,
					logTmp.funcName,
					logTmp.line,
					logTmp.msg)
			}
		default:
			//取不到日志先休息500ms，sleep会让出cpu，不会占用
			time.Sleep(time.Millisecond * 500)
		}
	}
}

// 判断是否需要记录日志
func (f *FileLogger) enable(logLevel LogLevel) bool {
	return logLevel >= f.Level
}

func (f *FileLogger) Debug(format string, a ...interface{}) {
	if f.async {
		f.logAsync(DEBUG, format, a...)
	} else {
		f.log(DEBUG, format, a...)
	}
}

func (f *FileLogger) Info(format string, a ...interface{}) {
	if f.async {
		f.logAsync(INFO, format, a...)
	} else {
		f.log(INFO, format, a...)
	}
}

func (f *FileLogger) Warning(format string, a ...interface{}) {
	if f.async {
		f.logAsync(WARNING, format, a...)
	} else {
		f.log(WARNING, format, a...)
	}
}

func (f *FileLogger) Error(format string, a ...interface{}) {
	if f.async {
		f.logAsync(ERROR, format, a...)
	} else {
		f.log(ERROR, format, a...)
	}
}
