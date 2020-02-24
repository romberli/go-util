package mylog

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"sort"
	"strings"

	"github.com/pingcap/errors"
	zaplog "github.com/pingcap/log"
	log "github.com/sirupsen/logrus"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"

	"github.com/romber2001/go-util/common"
)

const (
	DefaultLogFileName   = "run.log"
	defaultLogTimeFormat = "2006-01-02T15:04:05.000"
	// DefaultLogMaxSize is the default size of log files.
	DefaultLogMaxSize    = 100 // MB
	DefaultLogMaxBackups = 5
	DefaultLogMaxDays    = 7
	// DefaultLogFormat is the default format of the log.
	DefaultLogFormat = "text"
	defaultLogLevel  = log.InfoLevel
)

var MyLogger *zap.Logger

// EmptyFileLogConfig is an empty FileLogConfig.
var EmptyFileLogConfig = FileLogConfig{}

// FileLogConfig serializes file log related config in toml/json.
type FileLogConfig struct {
	zaplog.FileLogConfig
}

// NewFileLogConfig creates a FileLogConfig.
func NewFileLogConfig(fileName string, maxSize int, maxDays int, maxBackups int) (fileLogConfig *FileLogConfig, err error) {
	var baseDir string
	var logDir string
	var exists bool

	fileName = strings.TrimSpace(fileName)

	if fileName == "" {
		if baseDir, err = os.Getwd(); err != nil {
			return nil, err
		}

		logDir = path.Join(baseDir, "log")

		if exists, err = common.PathExistsLocal(logDir); err != nil {
			return nil, err
		}

		if !exists {
			if _, err := os.Create(logDir); err != nil {
				return nil, err
			}
		}

		fileName = path.Join(logDir, DefaultLogFileName)
	} else {
		logDir = path.Dir(fileName)

		if exists, err = common.PathExistsLocal(logDir); err != nil {
			return nil, err
		}
	}

	fileLogConfig = &FileLogConfig{
		FileLogConfig: zaplog.FileLogConfig{
			Filename:   fileName,
			MaxSize:    maxSize,
			MaxDays:    maxDays,
			MaxBackups: maxBackups,
		},
	}

	return fileLogConfig, nil
}

// LogConfig serializes log related config in toml/json.
type LogConfig struct {
	zaplog.Config
}

// NewLogConfig creates a LogConfig.
func NewLogConfig(level, format string, fileCfg FileLogConfig, disableTimestamp bool, opts ...func(*zaplog.Config)) *LogConfig {
	c := &LogConfig{
		Config: zaplog.Config{
			Level:            level,
			Format:           format,
			DisableTimestamp: disableTimestamp,
			File:             fileCfg.FileLogConfig,
		},
	}
	for _, opt := range opts {
		opt(&c.Config)
	}
	return c
}

func stringToLogLevel(level string) log.Level {
	switch strings.ToLower(level) {
	case "fatal":
		return log.FatalLevel
	case "error":
		return log.ErrorLevel
	case "warn", "warning":
		return log.WarnLevel
	case "debug":
		return log.DebugLevel
	case "info":
		return log.InfoLevel
	}

	return defaultLogLevel
}

// textFormatter is for compatibility with ngaut/log
type textFormatter struct {
	DisableTimestamp bool
	EnableEntryOrder bool
}

// Format implements logrus.Formatter
func (f *textFormatter) Format(entry *log.Entry) ([]byte, error) {
	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}

	if !f.DisableTimestamp {
		fmt.Fprintf(b, "%s ", entry.Time.Format(defaultLogTimeFormat))
	}
	if file, ok := entry.Data["file"]; ok {
		fmt.Fprintf(b, "%s:%v:", file, entry.Data["line"])
	}
	fmt.Fprintf(b, " [%s] %s", entry.Level.String(), entry.Message)

	if f.EnableEntryOrder {
		keys := make([]string, 0, len(entry.Data))
		for k := range entry.Data {
			if k != "file" && k != "line" {
				keys = append(keys, k)
			}
		}
		sort.Strings(keys)
		for _, k := range keys {
			fmt.Fprintf(b, " %v=%v", k, entry.Data[k])
		}
	} else {
		for k, v := range entry.Data {
			if k != "file" && k != "line" {
				fmt.Fprintf(b, " %v=%v", k, v)
			}
		}
	}

	b.WriteByte('\n')

	return b.Bytes(), nil
}

func stringToLogFormatter(format string, disableTimestamp bool) log.Formatter {
	switch strings.ToLower(format) {
	case "text":
		return &textFormatter{
			DisableTimestamp: disableTimestamp,
		}
	default:
		return &textFormatter{}
	}
}

// initFileLog initializes file based logging options.
func initFileLog(cfg *zaplog.FileLogConfig, logger *log.Logger) error {
	if st, err := os.Stat(cfg.Filename); err == nil {
		if st.IsDir() {
			return errors.New("can't use directory as log file name")
		}
	}
	if cfg.MaxSize <= 0 {
		cfg.MaxSize = DefaultLogMaxSize
	}
	if cfg.MaxBackups <= 0 {
		cfg.MaxBackups = DefaultLogMaxBackups
	}
	if cfg.MaxDays <= 0 {
		cfg.MaxDays = DefaultLogMaxDays
	}

	// use lumberjack to log rotate
	output := &lumberjack.Logger{
		Filename:   cfg.Filename,
		MaxSize:    cfg.MaxSize,
		MaxAge:     cfg.MaxDays,
		MaxBackups: cfg.MaxBackups,
		LocalTime:  true,
	}

	if logger == nil {
		log.SetOutput(output)
	} else {
		logger.Out = output
	}

	return nil
}

// InitLogger initializes logger.
func InitLogger(cfg *LogConfig) error {
	log.SetLevel(stringToLogLevel(cfg.Level))

	if cfg.Format == "" {
		cfg.Format = DefaultLogFormat
	}
	formatter := stringToLogFormatter(cfg.Format, cfg.DisableTimestamp)
	log.SetFormatter(formatter)

	if len(cfg.File.Filename) != 0 {
		if err := initFileLog(&cfg.File, nil); err != nil {
			return errors.Trace(err)
		}
	}

	return nil
}

// newZapLogger returns a zap logger
func newZapLogger() (*zap.Logger, *zaplog.ZapProperties) {
	var props *zaplog.ZapProperties

	conf := &zaplog.Config{Level: "info", File: zaplog.FileLogConfig{}}
	MyLogger, props, _ = zaplog.InitLogger(conf)
	return MyLogger, props
}

// InitZapLogger initializes a zap logger with cfg.
func InitZapLogger(cfg *LogConfig) error {
	var (
		props *zaplog.ZapProperties
		err   error
	)

	if MyLogger, props, err = zaplog.InitLogger(
		&cfg.Config,
		zap.AddStacktrace(zapcore.ErrorLevel),
		zap.AddCaller(),
		zap.Development(),
	); err != nil {
		return errors.Trace(err)
	}

	ReplaceGlobals(MyLogger, props)

	return nil
}

// L returns the global Logger, which can be reconfigured with ReplaceGlobals.
// It's safe for concurrent use.
func L() *zap.Logger {
	return _globalL
}

// S returns the global SugaredLogger, which can be reconfigured with
// ReplaceGlobals. It's safe for concurrent use.
func S() *zap.SugaredLogger {
	return _globalS
}

func ReplaceGlobals(logger *zap.Logger, props *zaplog.ZapProperties) {
	_globalL = logger
	_globalS = logger.Sugar()
	_globalP = props
}

var (
	_globalL, _globalP = newZapLogger()
	_globalS           = _globalL.Sugar()
)
