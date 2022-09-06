package logger

import (
	nested "github.com/antonfisher/nested-logrus-formatter"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	log "github.com/sirupsen/logrus"
	"time"
)

//func New(level logrus.Level, reportCaller bool,filePath string) (*logrus.Logger,error) {
//	logger := logrus.New()
//	logger.SetLevel(level)
//	logger.SetReportCaller(reportCaller)
//	logger.SetFormatter(&nested.Formatter{
//		FieldsOrder:     []string{"component", "category"},
//		TimestampFormat: time.RFC3339,
//	})
//
//	file,err := os.OpenFile(filePath,os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
//	if err == nil {
//		logger.SetOutput(file)
//	} else {
//		logger.Fatal("Failed to log to file, using default stderr")
//	}
//	logger.SetOutput(os.Stdout)
//	//logger.SetFormatter(&logrus.TextFormatter{
//	//	TimestampFormat: "2006-01-02 15:03:04",
//	//	PadLevelText:true,
//	//	CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
//	//		//处理文件名
//	//		fileName := path.Base(frame.File)
//	//		return frame.Function, fileName
//	//	},
//	//
//	//})
//
//	//logger.SetFormatter(&logrus.JSONFormatter{
//	//	//PrettyPrint:true,
//	//	TimestampFormat: "2006-01-02 15:03:04",
//	//
//	//	CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
//	//		//处理文件名
//	//		fileName := path.Base(frame.File)
//	//		return frame.Function, fileName
//	//	},
//	//})
//
//	return logger,nil
//}

type Logger struct {
	*log.Logger
	rotateLogs *rotatelogs.RotateLogs
}

func New(config *Config) (*Logger, error) {
	logger := log.New()
	formatter := &nested.Formatter{
		FieldsOrder:     []string{"component", "category"},
		TimestampFormat: time.RFC3339,
		//NoColors:        true,
	}

	logger.SetFormatter(formatter)
	logger.SetLevel(config.Level)
	logger.SetReportCaller(config.ReportCaller)
	writer, err := rotatelogs.New(
		config.FilePath+":%Y-%m-%d %H:%M:%S.log",
		rotatelogs.WithLinkName(config.FilePath),
		rotatelogs.WithMaxAge(config.MaxAge),
		rotatelogs.WithRotationTime(config.RotationTime))
	if err != nil {
		return nil, err
	}
	lfHook := lfshook.NewHook(lfshook.WriterMap{
		log.InfoLevel:  writer,
		log.DebugLevel: writer,
		log.ErrorLevel: writer,
	}, formatter)
	//logger.SetOutput(os.Stdout)
	logger.AddHook(lfHook)

	return &Logger{
		Logger:     logger,
		rotateLogs: writer,
	}, nil
}

func (l *Logger) Close() {
	_ = l.rotateLogs.Close()
	_ = l.rotateLogs.Close()
}

type Config struct {
	Level        log.Level
	ReportCaller bool
	FilePath     string
	MaxAge       time.Duration
	RotationTime time.Duration
}
