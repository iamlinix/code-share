package zaps

import (
	"cnpc.com.cn/cnpc/dserver/common"

	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var errorLogger *zap.SugaredLogger

func TimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {

	enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
}

func NewEncoderConfig() zapcore.EncoderConfig {

	return zapcore.EncoderConfig{
		TimeKey:     "ts",
		LevelKey:    "level",
		MessageKey:  "msg",
		LineEnding:  zapcore.DefaultLineEnding,
		EncodeLevel: zapcore.CapitalLevelEncoder,
		EncodeTime:  TimeEncoder,
		//EncodeTime:	zapcore.ISO8601TimeEncoder,
		//EncodeDuration: zapcore.StringDurationEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
	}
}

func InitZaps() {

	var wWriter io.Writer
	var eWriter io.Writer
	var fWriter io.Writer

	path := viper.GetString(common.ConfKeyLogPath)
	iPath := viper.GetString(common.ConfKeyLogPathInfo)
	wPath := viper.GetString(common.ConfKeyLogPathWarn)
	ePath := viper.GetString(common.ConfKeyLogPathError)
	fPath := viper.GetString(common.ConfKeyLogPathFatal)

	fmt.Println("path: ", path)
	fmt.Println("ipath: ", iPath)
	fmt.Println("wpath: ", wPath)
	fmt.Println("epath: ", ePath)
	fmt.Println("fpath: ", fPath)
	if iPath == "" || wPath == "" || ePath == "" || fPath == "" {
		log.Fatal("get log path failed")
		os.Exit(-1)
	}

	//init encoder
	encoder := zapcore.NewJSONEncoder(NewEncoderConfig())

	iLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl == zapcore.InfoLevel
	})

	wLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl == zapcore.WarnLevel
	})

	eLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl == zapcore.ErrorLevel
	})

	fLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl == zapcore.FatalLevel
	})

	iWriter := RWriter(iPath)

	if common.Exists(wPath) {
		wWriter, _ = os.OpenFile(wPath,
			os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	} else {
		wWriter, _ = os.Create(wPath)
	}

	if common.Exists(ePath) {
		eWriter, _ = os.OpenFile(ePath,
			os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	} else {
		eWriter, _ = os.Create(ePath)
	}

	if common.Exists(fPath) {
		fWriter, _ = os.OpenFile(fPath,
			os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	} else {
		fWriter, _ = os.Create(fPath)
	}

	core := zapcore.NewTee(
		zapcore.NewCore(encoder, zapcore.AddSync(iWriter), iLevel),
		zapcore.NewCore(encoder, zapcore.AddSync(wWriter), wLevel),
		zapcore.NewCore(encoder, zapcore.AddSync(eWriter), eLevel),
		zapcore.NewCore(encoder, zapcore.AddSync(fWriter), fLevel),
	)

	log := zap.New(core)
	errorLogger = log.Sugar()
}

func RWriter(filename string) io.Writer {

	hook, err := rotatelogs.New(
		strings.Replace(filename, ".log", "", -1)+"-%Y%m%d.log",
		rotatelogs.WithLinkName(filename),
		rotatelogs.WithMaxAge(7*24*time.Hour),
		rotatelogs.WithRotationTime(24*time.Hour),
	)

	if err != nil {
		panic(err)
	}

	return hook
}

func Debug(args ...interface{}) {
	errorLogger.Debug(args...)
}

func Debugf(template string, args ...interface{}) {
	errorLogger.Debugf(template, args...)
}

func Info(args ...interface{}) {
	errorLogger.Info(args...)
}

func Infof(template string, args ...interface{}) {
	errorLogger.Infof(template, args...)
}

func Warn(args ...interface{}) {
	errorLogger.Warn(args...)
}

func Warnf(template string, args ...interface{}) {
	errorLogger.Warnf(template, args...)
}

func Error(args ...interface{}) {
	errorLogger.Error(args...)
}

func Errorf(template string, args ...interface{}) {
	errorLogger.Errorf(template, args...)
}

func Fatal(args ...interface{}) {
	errorLogger.Fatal(args...)
}

func Fatalf(template string, args ...interface{}) {
	errorLogger.Fatalf(template, args...)
}
