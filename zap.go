package zaplogger

import (
	"go.k6.io/k6/v2/js/modules"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
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

type DynamicObject map[string]interface{}

func (d DynamicObject) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	for k, v := range d {
		switch v := v.(type) {
		case int:
			enc.AddInt(k, v)
		case float64:
			enc.AddFloat64(k, v)
		case string:
			enc.AddString(k, v)
		default:
			enc.AddReflected(k, v)
		}
	}
	return nil
}
func (*RootModule) NewModuleInstance(vu modules.VU) modules.Instance {
	return &ZapLogger{vu: vu}
}

func (zaplogger *ZapLogger) Exports() modules.Exports {
	return modules.Exports{Default: zaplogger}
}
func (z *ZapLogger) InitLogger(path string, overwrite bool, args ...int) *zap.SugaredLogger {
	defaultArgs := []int{500, 3, 28}
	for i := len(args); i < len(defaultArgs); i++ {
		args = append(args, defaultArgs[i])
	}

	config := writerConfig{
		MaxSize:    args[0],
		MaxBackups: args[1],
		MaxAge:     args[2],
	}
	return getSharedLogger(path, overwrite, config)
}

func (z *ZapLogger) CloseLogger(path string) error {
	return closeSharedLogger(path)
}

func (z *ZapLogger) CloseAll() int {
	return closeAllSharedLoggers()
}
func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoder := zapcore.NewJSONEncoder(encoderConfig)
	return encoder
}
func (z *ZapLogger) CreateDynamicObject(args ...interface{}) DynamicObject {
	obj := make(DynamicObject)
	for i := 0; i < len(args); i += 2 {
		key, _ := args[i].(string)
		obj[key] = args[i+1]
	}
	return obj
}
func (z *ZapLogger) ZapObject(key string, args ...interface{}) zapcore.Field {
	obj := make(DynamicObject)
	for i := 0; i < len(args); i += 2 {
		key, _ := args[i].(string)
		obj[key] = args[i+1]
	}
	return zap.Object(key, obj)
}
