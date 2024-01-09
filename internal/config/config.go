package config

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/yaml.v3"
	"os"
)

var Log *zap.Logger
var Cfg Config

type Config struct {
	DB struct {
		Port     string `yaml:"port"`
		Password string `yaml:"password"`
		NameDB   string `yaml:"name"`
		Ip       string `yaml:"ip"`
	} `yaml:"dataBase"`

	Logger struct {
		Pretty      bool     `yaml:"prettyLog"`
		Format      string   `yaml:"format"`
		Level       string   `yaml:"level"`
		Output      string   `yaml:"output"`
		StackTrace  bool     `yaml:"stacktrace"`
		OutputPaths []string `yaml:"outputPaths"`
	} `yaml:"logger"`

	Port string `yaml:"port"`
}

func Init() {
	mustInitConfigFile()

	initLog()
}

func mustInitConfigFile() {
	// Чтение файла конфигурации
	data, err := os.ReadFile("config.yaml")
	if err != nil {
		Log.Panic("Ошибка чтения файла конфигурации:" + err.Error())
	}

	// Разбор YAML и запись в структуру
	err = yaml.Unmarshal(data, &Cfg)
	if err != nil {
		Log.Panic("Ошибка разбора YAML: " + err.Error())
	}
}

func initLog() {

	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = "timestamp"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder

	logLevel := zap.InfoLevel

	switch Cfg.Logger.Level {
	case "info":
		logLevel = zap.InfoLevel
	case "debug":
		logLevel = zap.DebugLevel
	case "error":
		logLevel = zap.ErrorLevel
	case "warn":
		logLevel = zap.WarnLevel
	case "panic":
		logLevel = zap.PanicLevel
	case "fatal":
		logLevel = zap.FatalLevel
	}

	encoding := Cfg.Logger.Format
	if Cfg.Logger.Pretty {
		encoding = "console"
	}

	logConfig := zap.Config{
		Level:             zap.NewAtomicLevelAt(logLevel),
		Development:       false,
		DisableCaller:     false,
		DisableStacktrace: Cfg.Logger.StackTrace,
		Sampling:          nil,
		Encoding:          encoding,
		EncoderConfig:     encoderCfg,
		OutputPaths:       Cfg.Logger.OutputPaths,
		ErrorOutputPaths:  Cfg.Logger.OutputPaths,
		InitialFields:     map[string]interface{}{
			//"pid": os.Getpid(),
		},
	}
	Log = zap.Must(logConfig.Build())

	Log.Debug("log init")
}
