package log

import (
	"io"
	"os"
	"path"

	logging "github.com/op/go-logging"
)

/**
 *日志级别
 */
const (
	critical int = iota
	err_or
	warning
	notice
	info
	debug
)

var (
	// 日志管理器
	mlog *logger
	/**
	 * 日志输出格式
	 */
	logFormat = []string{
		`%{shortfunc} ▶ %{level} %{message}`,
		`%{time:15:04:05.00} %{shortfile} ▶ %{level} %{id:03x} %{message}`,
		`%{color}%{time:15:04:05.00} %{shortfunc} %{shortfile} ▶ %{level} %{id:03x}%{color:reset} %{message}`,
	}

	/**
	 * 日志级别与 string类型映射
	 */
	logLevelMap = map[string]int{
		"CRITICAL": critical,
		"ERROR":    err_or,
		"WARNING":  warning,
		"NOTICE":   notice,
		"INFO":     info,
		"DEBUG":    debug,
	}
)

type logger struct {
	log      *logging.Logger
	level    string
	filePath string
	ModeName string
	format   int
}

func init() {
	Init("DEBUG")
}

/**
 * 初始化日志
 * @param logLevel The arguments could be INFO, DEGUE, ERROR
 */
func Init(logLevel string) {
	mlog = newLog(logLevel)
	return
}

func newLog(level string) *logger {
	log := new(logger)
	log.level = level
	log.filePath = "./log"
	log.ModeName = "GameServer"
	log.format = 1
	log.log = logging.MustGetLogger(log.ModeName)
	log.AddLogBackend()
	return log
}

func (l *logger) AddLogBackend() {
	l.log.ExtraCalldepth = 2
	// backend1 := l.getFileBackend()
	backend2 := l.getStdOutBackend()
	logging.SetBackend(backend2)
	return
}

func SetLevel(level string) {
	logging.SetLevel(logging.Level(logLevelMap[level]), "")
	return
}

func (l *logger) getFileBackend() logging.LeveledBackend {
	//判断是否存在该文件夹
	if err := os.MkdirAll(l.filePath, 0777); err != nil {
		panic(err)
	}
	// 打开一个文件
	file, err := os.OpenFile(path.Join(l.filePath, l.ModeName+"_info.log"), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	backend := l.getLogBackend(file, logLevelMap[l.level])
	logging.SetBackend(backend)
	return backend
}

func (l *logger) getStdOutBackend() logging.LeveledBackend {
	bked := l.getLogBackend(os.Stderr, logLevelMap[l.level])
	return bked
}

/**
 * 获取终端
 */
func (l *logger) getLogBackend(out io.Writer, level int) logging.LeveledBackend {
	backend := logging.NewLogBackend(out, "", 1)
	format := logging.MustStringFormatter(logFormat[l.format])
	backendFormatter := logging.NewBackendFormatter(backend, format)
	backendLeveled := logging.AddModuleLevel(backendFormatter)
	backendLeveled.SetLevel(logging.Level(level), "")
	return backendLeveled
}

func (l *logger) logI(infmt string, args ...interface{}) {
	l.log.Infof(infmt, args...)
	return
}
func (l *logger) logE(infmt string, args ...interface{}) {
	l.log.Errorf(infmt, args...)
	return
}
func (l *logger) logW(infmt string, args ...interface{}) {
	l.log.Warningf(infmt, args...)
	return
}
func (l *logger) logD(infmt string, args ...interface{}) {
	l.log.Debugf(infmt, args...)
	return
}

func LogI(fmtstr string, args ...interface{}) {
	mlog.logI(fmtstr, args...)
	return
}

func LogW(fmtstr string, args ...interface{}) {
	mlog.logW(fmtstr, args...)
	return
}

func LogE(fmtstr string, args ...interface{}) {
	mlog.logE(fmtstr, args...)
	return
}

func LogD(fmtstr string, args ...interface{}) {
	mlog.logD(fmtstr, args...)
	return
}
