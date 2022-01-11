// Package logger 负责日志收敛，给日志添加前缀和独立文件输出，以便在本框架被集成之后的日志可读
package logger

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// 日志级别定义
const (
	TraceLevel    = 6 // 可详细跟踪
	DebugLevel    = 5 // 可调式
	InfoLevel     = 4 // 日常使用关键信息
	WarningLevel  = 3 // 警告级别信息
	ErrorLevel    = 2 // 错误级别信息
	CriticalLevel = 1 // 严重错误
)

// defaultPrefix 默认日志前缀
const defaultPrefix = "[goom-mocker]"

// openDebugEnv 开启debug日志
const openDebugEnv = "GOOM_DEBUG"

var (
	// LogLevel 日志级别
	// level总共分5个级别：debug < info< warning< error< critical
	LogLevel = InfoLevel
	// ConsoleLevel 控制台打印级别
	ConsoleLevel = WarningLevel
	// ShowError2Console 把错误同步打印到控制台
	ShowError2Console = false
	// Logger 独立日志文件
	Logger io.Writer = os.Stdout
	// EnableLogColor 开启并发日志染色
	EnableLogColor = false
	// colorGetter 日志染色器, 用于并发测试区分协程ID
	colorGetter func() string
	// logFile 日志路径
	logFile *os.File
)

// levelMap 日志级别-级别名称映射
var levelMap = map[int]string{
	TraceLevel:    "trace",
	DebugLevel:    "debug",
	InfoLevel:     "info",
	WarningLevel:  "warn",
	ErrorLevel:    "error",
	CriticalLevel: "critical",
}

// init 初始化
func init() {
	if d := os.Getenv(openDebugEnv); d != "" {
		OpenDebug()
	}

	loggerPath, err := loggerPath()
	if err != nil {
		fmt.Println("loggerPath error:", err)
		return
	}

	logFile, err = os.OpenFile(filepath.Join(loggerPath, "goom-mocker.log"),
		os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766)
	if err != nil {
		fmt.Println("init log file error:", err)
		return
	}

	Logger = logFile
}

// OpenDebug 开启debug模式
func OpenDebug() {
	ConsoleLevel = DebugLevel
}

// CloseDebug 关闭debug模式
func CloseDebug() {
	ConsoleLevel = WarningLevel
}

// IsDebugOpen 是否开启debug模式
func IsDebugOpen() bool {
	return ConsoleLevel >= DebugLevel
}

// OpenTrace 打开日志跟踪
func OpenTrace() {
	OpenDebug()
	SetLog2Console(true)
	LogLevel = TraceLevel
}

// CloseTrace 关闭日志跟踪
func CloseTrace() {
	CloseDebug()
	LogLevel = InfoLevel
	SetLog2Console(false)
}

// SetLog2Console 设置是否打印到控制台
func SetLog2Console(b bool) {
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
		_, _ = Logger.Write(withPrefix("info", v))
	}
}

// LogInfof 打印info日志
func LogInfof(format string, a ...interface{}) {
	if LogLevel >= InfoLevel {
		_, _ = Logger.Write(withPrefixStr("info", format, a...))
	}
}

// LogWarning 打印warning日志
func LogWarning(v ...interface{}) {
	if LogLevel >= WarningLevel {
		line := withPrefix("warning", v)

		_, _ = Logger.Write(line)

		write2Console(line)
	}
}

// LogWarningf 打印warning日志
func LogWarningf(format string, a ...interface{}) {
	if LogLevel >= WarningLevel {
		line := withPrefixStr("warning", format, a...)
		_, _ = Logger.Write(line)

		write2Console(line)
	}
}

// LogImportant 打印重要的日志
func LogImportant(v ...interface{}) {
	line := withPrefix("info", v)
	_, _ = Logger.Write(line)

	write2Console(line)
}

// LogImportantf 打印重要的日志
func LogImportantf(format string, a ...interface{}) {
	line := withPrefixStr("info", format, a...)
	_, _ = Logger.Write(line)

	write2Console(line)
}

// LogError 打印error日志
func LogError(v ...interface{}) {
	if LogLevel >= ErrorLevel {
		line := withPrefix("error", v)
		_, _ = Logger.Write(line)

		write2Console(line)
	}
}

// LogErrorf 打印error日志
func LogErrorf(format string, a ...interface{}) {
	if LogLevel >= ErrorLevel {
		line := withPrefixStr("error", format, a...)
		_, _ = Logger.Write(line)

		write2Console(line)
	}
}

// Log2Console 打印日志到控制台
func Log2Console(level int, s string) {
	if level <= ConsoleLevel {
		os.Stdout.Write([]byte(s))
	}
}

// Log2Consolef 打印日志到控制台
func Log2Consolef(level int, format string, a ...interface{}) {
	if level <= ConsoleLevel {
		line := withPrefixStr(levelMap[level], format, a...)
		os.Stdout.Write(line)
	}
}

// write2Console 输出到控制台，如果ShowError2Console==true时
func write2Console(line []byte) {
	if ShowError2Console {
		os.Stdout.Write(line)
	}
}

// withPrefix 给日志带上[goom]前缀, 方便与业务日志区分
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

// withPrefixStr 给日志带上[goom]前缀
func withPrefixStr(level, format string, a ...interface{}) []byte {
	time := time.Now().Format("2006-01-02 15:04:05")

	if EnableLogColor {
		return []byte(time + " " + defaultPrefix + "[" + level + "]: " +
			colorGetter() + " " + fmt.Sprintf(format, a...) + "\n")
	}

	return []byte(time + " " + defaultPrefix + "[" + level + "]: " +
		fmt.Sprintf(format, a...) + "\n")
}

// loggerPath 获取日志存储路径
func loggerPath() (string, error) {
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
