package logger

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const CONFIG_FILE = "logger/config.json"

var configReader ConfigReader

type ConfigReader struct {
	Logging struct {
		PrintLogs bool   `json:"print_logs"`
		LogFile   string `json:"log_file"`
	} `json:"logging"`
}

func readConfig() {
	jsonFile, err := ioutil.ReadFile(CONFIG_FILE)
	if err != nil {
		fmt.Println(err)
		panic("Error while reading the config file")
	}
	json.Unmarshal(jsonFile, &configReader)
}

var sugarLogger *zap.SugaredLogger

func InitLogger() *zap.SugaredLogger {
	readConfig()
	writeSyncer := getLogWriter()
	encoder := getEncoder()
	core := zapcore.NewCore(encoder, writeSyncer, zapcore.DebugLevel)

	logger := zap.New(core, zap.AddCaller())
	sugarLogger = logger.Sugar()
	return sugarLogger
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func getLogWriter() zapcore.WriteSyncer {

	var logFileName string
	logFileName = ""
	if configReader.Logging.PrintLogs {
		logFileName = "./" + configReader.Logging.LogFile
	}

	lumberJackLogger := &lumberjack.Logger{
		Filename:   logFileName,
		MaxSize:    100,
		MaxBackups: 5,
		MaxAge:     30,
		Compress:   false,
	}
	return zapcore.AddSync(lumberJackLogger)
}
