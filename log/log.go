package log

/* 参考 ULR https://blog.csdn.net/weixin_45565886/article/details/132520476 */
import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

// 基于log库自定义实现logger
var (
	infoLogger  *log.Logger
	debugLogger *log.Logger
	errorLogger *log.Logger

	logFilePath   string
	logFileHandle *os.File

	logLevel   int
	currentDay int          //每天生成一个日志文件
	fileLock   sync.RWMutex //读写锁，保证同一时间只有一个协程重命名文件
)

const (
	DEBUGLEVEL = iota
	INFOLEVEL
	ERRORLEVEL
)

func init() {
	logLevel = DEBUGLEVEL
}

func SetLevel(level int) {
	logLevel = level
}

func SetFile(file string) {
	var err error
	logFileHandle, err = os.OpenFile(file, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0664)
	if err != nil {
		panic(err)
	} else {
		//初始化自定义的Logger（基于golang中的log库）
		//log.LstdFlags表示时间格式等
		//log.Llongfile表示文件名及调用代码的位置,log.Llongfile=》改为通过getCallTrace获取前缀
		currentDay = time.Now().YearDay()
		infoLogger = log.New(os.Stdout, "[INFO ]", log.LstdFlags)
		debugLogger = log.New(os.Stderr, "[DEBUG] ", log.LstdFlags)
		errorLogger = log.New(os.Stderr, "[ERROR] ", log.LstdFlags)
		//infoLogger = log.New(logFileHandle, "[INFO] ", log.LstdFlags)
		//debugLogger = log.New(logFileHandle, "[DEBUG] ", log.LstdFlags)
		//errorLogger = log.New(logFileHandle, "[ERROR] ", log.LstdFlags)
		logFilePath = file
	}
}

func isDayChanged() {
	fileLock.Lock()
	defer fileLock.Unlock()
	{ // 不单独存档到文件
		return
	}

	day := time.Now().YearDay()
	if day == currentDay {
		return
	}

	//关闭之前的文件，重命名，并生成新的文件
	logFileHandle.Close()
	postFix := time.Now().Add(-24 * time.Hour).Format("20060102")
	err := os.Rename(logFilePath, logFilePath+"."+postFix)
	if err != nil {
		//TODO 重命名日志文件失败，根据自身情况做处理
	}
	logFileHandle, err = os.OpenFile(logFilePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0664)
	if err != nil {
		//TODO 打开新的日志文件失败，根据自己业务需求做处理
	}
	infoLogger = log.New(logFileHandle, "[INFO ]", log.LstdFlags)
	debugLogger = log.New(logFileHandle, "[DEBUG]", log.LstdFlags)
	currentDay = day
}

func Printf(format string, v ...interface{}) {
	if logLevel <= DEBUGLEVEL {
		isDayChanged()
		var msg string
		/* [INFO] 2024/05/20 21:30:21 utils/log/log_test.go:12 utils.log test!%!!(MISSING)(EXTRA []interface {}=[])
		 * 解决 args 为NULL 引起的格式问题，下同
		 */
		if len(v) == 0 {
			msg = fmt.Sprintf("%s%s", getPrefix(), format)
		} else {
			msg = fmt.Sprintf(getPrefix()+format, v)
		}

		debugLogger.Print(msg)
	}
}

// Debug golang中的any相当于interface{}空接口
func Debug(format string, args ...interface{}) {
	if logLevel <= DEBUGLEVEL {
		isDayChanged()
		var msg string
		if len(args) == 0 {
			msg = fmt.Sprintf("%s%s", getPrefix(), format)
		} else {
			msg = fmt.Sprintf(getPrefix()+format, args)
		}

		debugLogger.Print(msg)
	}
}

func Info(format string, args ...interface{}) {
	if logLevel <= INFOLEVEL {
		isDayChanged()
		var msg string
		if len(args) == 0 {
			msg = fmt.Sprintf("%s%s", getPrefix(), format)
		} else {
			msg = fmt.Sprintf(getPrefix()+format, args)
		}

		infoLogger.Print(msg)
	}
}

func Error(format string, args ...interface{}) {
	if logLevel <= ERRORLEVEL {
		isDayChanged()
		var msg string
		if len(args) == 0 {
			msg = fmt.Sprintf("%s%s", getPrefix(), format)
		} else {
			msg = fmt.Sprintf(getPrefix()+format, args)
		}

		errorLogger.Print(msg)
	}
}

// 获取函数调用栈关系：拿到调用Info或者Debug所在的文件名及代码行数（runtime包）
func getCallTrace() (string, int) {
	_, file, lineNo, ok := runtime.Caller(3)
	if ok {
		return file, lineNo
	} else {
		return "", 0
	}
}

// 获取调用Info、Debug代码所在行数，文件名只获取最后三级
func getPrefix() string {
	file, lineNo := getCallTrace()
	path := strings.Split(file, "/")
	if len(path) > 3 {
		file = strings.Join(path[len(path)-3:], "/")
	}
	return file + ":" + strconv.Itoa(lineNo) + " "
}
