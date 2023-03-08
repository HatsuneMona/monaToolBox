package bootstrap

import (
	"fmt"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"monaToolBox/global"
	"monaToolBox/utils"
	"os"
	"time"
)

var (
	level   zapcore.Level // zap 日志等级
	options []zap.Option  // zap 配置项
)

// InitLog 初始化zapLog
func InitLog() *zap.Logger {

	// 创建Log目录
	createLogRootDir()

	// 设置日志等级
	setLogLevel()

	// 展示log记录详细信息，例如调用Log位置的行号信息
	if global.Config.Log.ShowLine {
		options = append(options, zap.AddCaller())
	}

	// 初始化 zap
	return zap.New(getZapCore(), options...)

}

// createLogRootDir 根据用户配置的log目录，创建目录
func createLogRootDir() {
	logDir := global.Config.Log.RootDir
	if ok, _ := utils.PathExists(logDir); !ok {
		err := os.MkdirAll(logDir, 0755)
		if err != nil {
			panic(fmt.Errorf("log目录初始化失败，err:%s. \n", err))
		}
	}
}

// setLogLevel 设置日志记录等级
func setLogLevel() {
	switch global.Config.Log.Level {
	case "debug":
		level = zap.DebugLevel
		options = append(options, zap.AddStacktrace(level))
	case "info":
		level = zap.InfoLevel
	case "warn":
		level = zap.WarnLevel
	case "error":
		level = zap.ErrorLevel
		options = append(options, zap.AddStacktrace(level))
	case "dpanic":
		level = zap.DPanicLevel
	case "panic":
		level = zap.PanicLevel
	case "fatal":
		level = zap.FatalLevel
	default:
		level = zap.InfoLevel
	}
}

// getZapCore 扩展 Zap
func getZapCore() zapcore.Core {
	var encoder zapcore.Encoder

	// 调整编码器默认配置
	encoderConfig := zap.NewProductionEncoderConfig() // 这个函数返回了一个默认的Config，根据需求对这个配置里的字段进行修改。
	encoderConfig.EncodeTime = func(time time.Time, encoder zapcore.PrimitiveArrayEncoder) {
		encoder.AppendString(time.Format("[" + "2006-01-02 15:04:05.000" + "]"))
	}

	// encoderConfig.EncodeLevel = func(l zapcore.Level, encoder zapcore.PrimitiveArrayEncoder) {
	// 	encoder.AppendString(global.App.Config.App.Env + "." + l.String())
	// }

	// 设置编码器
	if global.Config.Log.Format == "json" {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}

	return zapcore.NewCore(encoder, getLogWriter(), level)
}

// 使用 lumberjack 作为日志写入器
func getLogWriter() zapcore.WriteSyncer {
	file := &lumberjack.Logger{
		Filename:   global.Config.Log.RootDir + "/" + global.Config.Log.Filename,
		MaxSize:    global.Config.Log.MaxSize,
		MaxBackups: global.Config.Log.MaxBackups,
		MaxAge:     global.Config.Log.MaxAge,
		Compress:   global.Config.Log.Compress,
	}

	return zapcore.AddSync(file)
}
