package zaplogger

import (
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
func (z *ZapLogger) InitLogger(path string, args ...int) *zap.SugaredLogger {
	// MaxSize:    500, // megabytes
	// MaxBackups: 3,
	// MaxAge:     28, // days
	defaultArgs := []int{500, 3, 28}
	for i := len(args); i < len(defaultArgs); i++ {
		args = append(args, defaultArgs[i])
	}
	writeSyncer := zapcore.AddSync(&lumberjack.Logger{
		Filename:   path,
		MaxSize:    args[0],
		MaxBackups: args[1],
		MaxAge:     args[2],
	})
	encoder := getEncoder()
	// core := zapcore.NewCore(encoder, safeW, zapcore.NewMultiWriteSyncer(writeSyncer, zapcore.AddSync(os.Stdout)), zapcore.DebugLevel)
	//不在控制台输出，去掉zapcore.AddSync(os.Stdout)
	core := zapcore.NewCore(encoder, zapcore.NewMultiWriteSyncer(writeSyncer), zapcore.DebugLevel)

	logger := zap.New(core)
	sugarLogger := logger.Sugar()
	return sugarLogger
}
func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoder := zapcore.NewJSONEncoder(encoderConfig)
	return encoder
}
