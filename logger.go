package utils

/* 参考 ULR https://blog.csdn.net/weixin_45565886/article/details/132520476 */
import (
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

	logOut     *os.File
	logLevel   int
	currentDay int //每天生成一个日志文件
	logFile    string
	fileLock   sync.RWMutex //读写锁，保证同一时间只有一个协程重命名文件
)

const (
	DebugLevel = iota //0
	InfoLevel         //1
)

func SetLevel(level int) {
	logLevel = level
}

func init() {
	fileLock = sync.RWMutex{}
}

func SetFile(file string) {
	var err error
	logOut, err = os.OpenFile(file, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0664)
	if err != nil {
		panic(err)
	} else {
		//初始化自定义的Logger（基于golang中的log库）
		//log.LstdFlags表示时间格式等
		//log.Llongfile表示文件名及调用代码的位置,log.Llongfile=》改为通过getCallTrace获取前缀
		currentDay = time.Now().YearDay()
		infoLogger = log.New(logOut, "[INFO] ", log.LstdFlags)
		debugLogger = log.New(logOut, "[DEBUG] ", log.LstdFlags)
		logFile = file
	}
}

func checkIfDayChange() {
	fileLock.Lock()
	defer fileLock.Unlock()
	day := time.Now().YearDay()
	if day == currentDay {
		return
	} else {
		//关闭之前的文件，重命名，并生成新的文件
		logOut.Close()
		postFix := time.Now().Add(-24 * time.Hour).Format("20060102")
		err := os.Rename(logFile, logFile+"."+postFix)
		if err != nil {
			//TODO 重命名日志文件失败，根据自身情况做处理
		}
		logOut, err = os.OpenFile(logFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0664)
		if err != nil {
			//TODO 打开新的日志文件失败，根据自己业务需求做处理
		}
		infoLogger = log.New(logOut, "[INFO] ", log.LstdFlags)
		debugLogger = log.New(logOut, "[DEBUG] ", log.LstdFlags)
		currentDay = day
	}
}

// Debug golang中的any相当于interface{}空接口
func Debug(format string, v ...any) {
	if logLevel <= DebugLevel {
		checkIfDayChange()
		debugLogger.Printf(getPrefix()+format, v)
	}
}

func Info(format string, v ...any) {
	if logLevel <= InfoLevel {
		checkIfDayChange()
		infoLogger.Printf(getPrefix()+format, v)
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
