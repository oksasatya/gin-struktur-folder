package utils

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/natefinch/lumberjack"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm/logger"
)

// FileHook is a custom hook for logging to a file with a different format
type FileHook struct {
	Writer    io.Writer
	Formatter logrus.Formatter
	LevelsVal []logrus.Level
}

func NewFileHook(levels []logrus.Level, writer io.Writer, formatter logrus.Formatter) *FileHook {
	return &FileHook{
		Writer:    writer,
		Formatter: formatter,
		LevelsVal: levels,
	}
}

func (hook *FileHook) Levels() []logrus.Level {
	return hook.LevelsVal
}

func (hook *FileHook) Fire(entry *logrus.Entry) error {
	if os.Getenv("PORT") == "8080" {
		entry.Data["Environment"] = "Development"
	} else {
		entry.Data["Environment"] = "Production"
	}

	line, err := hook.Formatter.Format(entry)
	if err != nil {
		logrus.Errorf("Error formatting log entry for file: %v", err)
		return err
	}

	// Write the formatted entry to the writer
	_, err = hook.Writer.Write(line)
	if err != nil {
		logrus.Errorf("Error writing log entry to file: %v", err)
		return err
	}
	return nil
}

// SetupLogger initializes the logger with both terminal and file logging
func SetupLogger() {
	logDir := "logs"
	if err := os.MkdirAll(logDir, 0755); err != nil {
		logrus.Fatalf("Failed to create log directory: %v", err)
	}

	logFilePath := filepath.Join(logDir, "app.log")

	fileLogger := &lumberjack.Logger{
		Filename:   logFilePath,
		MaxSize:    10, // Megabytes
		MaxBackups: 3,
		MaxAge:     28, // Days
		Compress:   true,
	}

	// Set up terminal logger (this will be the default output)
	logrus.SetOutput(os.Stdout)
	logrus.SetFormatter(&logrus.TextFormatter{
		TimestampFormat:        "2006-01-02 15:04:05",
		FullTimestamp:          true,
		ForceColors:            true,  // Enable colors for terminal output
		DisableColors:          false, // Keep colors in terminal
		QuoteEmptyFields:       true,
		DisableQuote:           true,
		DisableLevelTruncation: true,
		PadLevelText:           false,
	})

	// Set log level
	logrus.SetLevel(logrus.InfoLevel)

	// Add custom hook for file logging with a different formatter (no colors)
	logrus.AddHook(NewFileHook(logrus.AllLevels, fileLogger, &logrus.TextFormatter{
		TimestampFormat:  "2006-01-02 15:04:05",
		FullTimestamp:    true,
		ForceColors:      false, // Disable colors for file output
		DisableColors:    true,
		QuoteEmptyFields: true,
	}))
}

// LogrusLogger is a custom logger for GIN that uses logrus
func LogrusLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// save start time
		start := time.Now()

		// Process request
		c.Next()

		latency := time.Since(start)
		req := c.Request
		res := c.Writer

		logrus.WithFields(logrus.Fields{
			"method":         req.Method,
			"uri":            req.RequestURI,
			"status":         res.Status(),
			"latency":        latency,
			"ip":             c.ClientIP(),
			"user_agent":     req.UserAgent(),
			"host":           req.Host,
			"referer":        req.Referer(),
			"protocol":       req.Proto,
			"content_length": req.ContentLength,
			"query_params":   c.Request.URL.Query().Encode(),
			"response_size":  res.Size(),
			"cookies":        req.Cookies(),
			"handler_name":   runtime.FuncForPC(reflect.ValueOf(c.Handler()).Pointer()).Name(),
		}).Info("HTTP request")
	}
}

// LogrusGormLogger is a custom logger for GORM that uses logrus
type LogrusGormLogger struct {
	LogLevel logger.LogLevel
}

// NewLogrusGormLogger is a constructor for LogrusGormLogger
func NewLogrusGormLogger(logLevel logger.LogLevel) *LogrusGormLogger {
	return &LogrusGormLogger{
		LogLevel: logLevel,
	}
}

// LogMode is implement of logger.Interface
func (l *LogrusGormLogger) LogMode(level logger.LogLevel) logger.Interface {
	return &LogrusGormLogger{
		LogLevel: level,
	}
}

// Info is implement of logger.Interface
func (l *LogrusGormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Info {
		logrus.WithContext(ctx).Infof(msg, data...)
	}
}

// Warn is implement of logger.Interface
func (l *LogrusGormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Warn {
		logrus.WithContext(ctx).Warnf(msg, data...)
	}
}

// Error is implement of logger.Interface
func (l *LogrusGormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Error {
		logrus.WithContext(ctx).Errorf(msg, data...)
	}
}

// Trace is implement of logger.Interface
func (l *LogrusGormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.LogLevel <= 0 {
		return
	}
	elapsed := time.Since(begin)
	sql, rows := fc()

	// Skip logging for ErrRecordNotFound
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		entry := logrus.WithContext(ctx).WithFields(logrus.Fields{
			"elapsed": elapsed,
			"rows":    rows,
			"sql":     sql,
		})

		entry.WithField("error", err).Error("database operation failed")
	} else if elapsed > 200*time.Millisecond { // log if execution time is longer than 200ms
		logrus.WithContext(ctx).Warnf("slow database operation: %s", sql)
	} else {
		logrus.WithContext(ctx).Infof("database operation: %s", sql)
	}
}
