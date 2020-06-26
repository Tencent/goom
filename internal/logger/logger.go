package logger

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// 默认日志前缀
const defaultPrefix = "[goom-mocker]"

// ShowError2Console 把错误同步打印到控制台
var ShowError2Console = false

// 日志级别
// level总共分5个级别：debug < info< warning< error< critical
var LogLevel = 6

// Logger 独立日志文件
var Logger io.Writer = os.Stdout

// EnableLogColor 开启并发日志染色
var EnableLogColor = false
var colorGetter func() string

// logFile 日志路径
var logFile *os.File

// 日志级别定义
var (
	TraceLevel    = 6
	DebugLevel    = 5
	InfoLevel     = 4
	WarningLevel  = 3
	ErrorLevel    = 2
	CriticalLevel = 1
)

func init() {
	loggerPath, err := getLoggerPath()
	if err != nil {
		fmt.Println("getLoggerPath error:", err)
		return
	}
	logFile, err = os.OpenFile(filepath.Join(loggerPath, "goom-mocker.log"), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766)
	if err != nil {
		fmt.Println("init log file error:", err)
		return
	}
	Logger = logFile
}

// Log2Console 打印到控制台
func Log2Console(b bool) {
	if b {
		Logger = os.Stdout
	} else {
		Logger = logFile
	}
}

// SetLogColor 设置日志染色
func SetLogColor(enable bool, getter func() string) {
	if enable && getter != nil {
		colorGetter = getter
		EnableLogColor = true
	} else {
		EnableLogColor = false
	}
}

// LogTraceEnable 是否开启trace
func LogTraceEnable() bool {
	return LogLevel >= TraceLevel
}

// LogTrace 打印trace日志
func LogTrace(v ...interface{}) {
	if LogLevel >= TraceLevel {
		_, _ = Logger.Write(withPrefix("trace", v))
	}
}

// LogTracef 打印trace日志
func LogTracef(format string, a ...interface{}) {
	if LogLevel >= TraceLevel {
		_, _ = Logger.Write(withPrefixStr("trace", format, a...))
	}
}

// LogDebugEnable 是否开启debug
func LogDebugEnable() bool {
	return LogLevel >= DebugLevel
}

// LogDebug 打印debug日志
func LogDebug(v ...interface{}) {
	if LogLevel >= DebugLevel {
		_, _ = Logger.Write(withPrefix("debug", v))
	}
}

// LogDebugf 打印debug日志
func LogDebugf(format string, a ...interface{}) {
	if LogLevel >= DebugLevel {
		_, _ = Logger.Write(withPrefixStr("debug", format, a...))
	}
}

// LogInfo 打印info日志
func LogInfo(v ...interface{}) {
	if LogLevel >= InfoLevel {
		Logger.Write(withPrefix("info", v))
	}
}

// LogInfof 打印info日志
func LogInfof(format string, a ...interface{}) {
	if LogLevel >= InfoLevel {
		Logger.Write(withPrefixStr("info", format, a...))
	}
}

// LogWarning 打印warning日志
func LogWarning(v ...interface{}) {
	if LogLevel >= WarningLevel {
		line := withPrefix("warning", v)
		Logger.Write(line)
		if ShowError2Console {
			os.Stdout.Write(line)
		}
	}
}

// LogWarningf 打印warning日志
func LogWarningf(format string, a ...interface{}) {
	if LogLevel >= WarningLevel {
		line := withPrefixStr("warning", format, a...)
		Logger.Write(line)
		if ShowError2Console {
			os.Stdout.Write(line)
		}
	}
}

// LogImportant 打印重要的日志
func LogImportant(v ...interface{}) {
	line := withPrefix("info", v)
	Logger.Write(line)
	if ShowError2Console {
		os.Stdout.Write(line)
	}
}

// LogImportantf 打印重要的日志
func LogImportantf(format string, a ...interface{}) {
	line := withPrefixStr("info", format, a...)
	Logger.Write(line)
	if ShowError2Console {
		os.Stdout.Write(line)
	}
}

// LogError 打印error日志
func LogError(v ...interface{}) {
	if LogLevel >= ErrorLevel {
		line := withPrefix("error", v)
		Logger.Write(line)
		if ShowError2Console {
			os.Stdout.Write(line)
		}
	}
}

// LogErrorf 打印error日志
func LogErrorf(format string, a ...interface{}) {
	if LogLevel >= ErrorLevel {
		line := withPrefixStr("error", format, a...)
		Logger.Write(line)
		if ShowError2Console {
			os.Stdout.Write(line)
		}
	}
}

func withPrefix(level string, v []interface{}) []byte {
	arr := make([]string, 0, len(v)+1)
	arr = append(arr, time.Now().Format("2006-01-02 15:04:05"))
	arr = append(arr, defaultPrefix, "["+level+"]:")
	if EnableLogColor {
		arr = append(arr, colorGetter())
	}
	for _, a := range v {
		arr = append(arr, fmt.Sprintf("%s", a))
	}
	arr = append(arr, "\n")
	return []byte(strings.Join(arr, " "))
}

func withPrefixStr(level, format string, a ...interface{}) []byte {
	time := time.Now().Format("2006-01-02 15:04:05")
	if EnableLogColor {
		return []byte(time + " " + defaultPrefix + "[" + level + "]: " + colorGetter() + " " + fmt.Sprintf(format, a...) + "\n")
	} else {
		return []byte(time + " " + defaultPrefix + "[" + level + "]: " + fmt.Sprintf(format, a...) + "\n")
	}
}

// 获取日志存储路径
func getLoggerPath() (string, error) {
	var logFileLocation = "."
	// 获取当前目录
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err == nil && "/" != dir {
		logFileLocation = filepath.Join(dir, "logs")
	}

	// 判断文件夹是否存在
	_, err = os.Stat(logFileLocation)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.Mkdir(logFileLocation, os.ModePerm)
			if err != nil {
				fmt.Println("init log file error:", err)
				return ".", err
			}
		} else {
			return ".", err
		}
	}

	fmt.Println("goom-mocker logFileLocation:", logFileLocation)
	return logFileLocation, err
}
