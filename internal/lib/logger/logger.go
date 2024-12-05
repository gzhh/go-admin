package logger

import (
	"context"
	"go-admin/internal/lib/config"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// TODO: optimize if necessary write gin log
var WriteSyncer zapcore.WriteSyncer

var Logger *zap.Logger
var SugaredLogger *zap.SugaredLogger

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewJSONEncoder(encoderConfig)
}

func getWriterSyncer(filePath string) zapcore.WriteSyncer {
	// // 直接使用单一文件写入
	// file, _ := os.Create(filePath)
	// return zapcore.AddSync(file)

	// 引入第三方库 Lumberjack 加入日志切割功能
	lumberWriteSyncer := &lumberjack.Logger{
		Filename:   filePath,
		MaxSize:    1,    // megabytes
		MaxBackups: 100,  // number
		MaxAge:     30,   // days
		Compress:   true, // 确定是否应该使用gzip压缩已旋转的日志文件。默认值是不执行压缩。
	}
	return zapcore.AddSync(lumberWriteSyncer)

}

func Setup() {
	// handle log config
	var log config.LogConfig
	if config.Settings.AdminServer != nil {
		log = config.Settings.AdminServer.Log
	} else {
		panic("Settings AdminServer not found")
	}

	LogSavePath = log.LogSavePath
	LogSaveName = log.LogSaveName
	LogFileExt = log.LogFileExt
	filePath := getLogFileFullPath()

	// init logger
	WriteSyncer = getWriterSyncer(filePath)
	encoder := getEncoder()
	core := zapcore.NewCore(encoder, WriteSyncer, zap.DebugLevel)
	Logger = zap.New(core, zap.AddCaller())
	SugaredLogger = Logger.Sugar()

}

func setPrefix(ctx context.Context) *zap.SugaredLogger {
	var traceID string
	if ctx != nil {
		if v := ctx.Value(TraceID); v != nil {
			traceID = v.(string)
		}
	}

	return SugaredLogger.With(zap.String("trace_id", traceID))
}

// Logger methods
func Debug(ctx context.Context, v ...interface{}) {
	setPrefix(ctx).Debug(v...)
}

func Info(ctx context.Context, v ...interface{}) {
	setPrefix(ctx).Info(v...)
}

func Warn(ctx context.Context, v ...interface{}) {
	setPrefix(ctx).Warn(v...)
}

func Error(ctx context.Context, v ...interface{}) {
	setPrefix(ctx).Error(v...)
}

func Panic(ctx context.Context, v ...interface{}) {
	setPrefix(ctx).Panic(v...)
}

func Fatal(ctx context.Context, v ...interface{}) {
	setPrefix(ctx).Fatalln(v...)
}

func Debugf(ctx context.Context, format string, v ...interface{}) {
	setPrefix(ctx).Debugf(format, v...)
}

func Infof(ctx context.Context, format string, v ...interface{}) {
	setPrefix(ctx).Infof(format, v...)
}

func Warnf(ctx context.Context, format string, v ...interface{}) {
	setPrefix(ctx).Warnf(format, v...)
}

func Errorf(ctx context.Context, format string, v ...interface{}) {
	setPrefix(ctx).Errorf(format, v...)
}

func Panicf(ctx context.Context, format string, v ...interface{}) {
	setPrefix(ctx).Panicf(format, v...)
}

func Fatalf(ctx context.Context, format string, v ...interface{}) {
	setPrefix(ctx).Fatalf(format, v...)
}
