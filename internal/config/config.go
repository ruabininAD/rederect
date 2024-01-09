package config

import (
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"strconv"
	"strings"
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
	//mustInitConfigFile()
	mustInitEnvFile()
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

func mustInitEnvFile() {
	// Загрузка переменных окружения из файла .env
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Ошибка при загрузке файла .env: %s", err)
	}
	Cfg.Port = os.Getenv("APP_LISTEN_PORT")

	Cfg.DB.Ip = os.Getenv("DATABASE_IP")
	Cfg.DB.Port = os.Getenv("DATABASE_PORT")
	Cfg.DB.NameDB = os.Getenv("DATABASE_NAME")
	Cfg.DB.Password = os.Getenv("DATABASE_PASSWORD")

	Cfg.Logger.Level = os.Getenv("LOGGER_LEVEL")
	Cfg.Logger.OutputPaths = strings.Split(os.Getenv("LOGGER_OUTPUTPATHS"), ",")
	Cfg.Logger.Pretty, _ = strconv.ParseBool(os.Getenv("LOGGER_PRETTYLOG")) //fixme
	Cfg.Logger.Output = os.Getenv("LOGGER_OUTPUT")
	Cfg.Logger.Format = os.Getenv("LOGGER_FORMAT")
	Cfg.Logger.StackTrace, _ = strconv.ParseBool(os.Getenv("LOGGER_STACKTRACE"))

}
