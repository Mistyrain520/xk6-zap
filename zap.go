package zaplogger

import (
	"sync"

	"go.k6.io/k6/js/modules"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// init is called by the Go runtime at application startup.
func init() {
	modules.Register("k6/x/zaplogger", new(RootModule))
}

type RootModule struct{}
type ZapLogger struct {
	vu modules.VU
	mu sync.Mutex
}

var (
	_ modules.Module   = &RootModule{}
	_ modules.Instance = &ZapLogger{}
)

func (*RootModule) NewModuleInstance(vu modules.VU) modules.Instance {
	return &ZapLogger{vu: vu}
}

func (zaplogger *ZapLogger) Exports() modules.Exports {
	return modules.Exports{Default: zaplogger}
}
func (z *ZapLogger) InitLogger(path string) *zap.SugaredLogger {
	// file, _ := os.Create(path)
	writeSyncer := zapcore.AddSync(&lumberjack.Logger{
		Filename:   path,
		MaxSize:    500, // megabytes
		MaxBackups: 3,
		MaxAge:     28, // days
	})
	encoder := getEncoder()
	// core := zapcore.NewCore(encoder, zapcore.NewMultiWriteSyncer(writeSyncer, zapcore.AddSync(os.Stdout)), zapcore.DebugLevel)
	//不在控制台输出，去掉zapcore.AddSync(os.Stdout)
	core := zapcore.NewCore(encoder, zapcore.NewMultiWriteSyncer(writeSyncer), zapcore.DebugLevel)

	logger := zap.New(core)
	sugarLogger := logger.Sugar()
	return sugarLogger
}

//	func (z *ZapLogger) Sync(sugar *zap.SugaredLogger) {
//		sugar.Sync()
//	}

//	func (z *ZapLogger) Debugw(sugar *zap.SugaredLogger, msg string, keysAndValues ...interface{}) {
//		sugar.Debugw(msg, keysAndValues...)
//	}
//
//	func (z *ZapLogger) Warnw(sugar *zap.SugaredLogger, msg string, keysAndValues ...interface{}) {
//		sugar.Warnw(msg, keysAndValues...)
//	}
//
//	func (z *ZapLogger) Errorw(sugar *zap.SugaredLogger, msg string, keysAndValues ...interface{}) {
//		sugar.Errorw(msg, keysAndValues...)
//	}
//
//	func (z *ZapLogger) Fatalw(sugar *zap.SugaredLogger, msg string, keysAndValues ...interface{}) {
//		sugar.Fatalw(msg, keysAndValues...)
//	}
func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder //指定时间格式
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoder := zapcore.NewJSONEncoder(encoderConfig)
	return encoder
}
