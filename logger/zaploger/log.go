package zaploger

import (
	"context"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/helays/utils/v2/tools"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gorm.io/gorm/logger"
)

// LogConfig defines the configuration for each log level.
type LogConfig struct {
	FileName   string `json:"file_name" yaml:"file_name" int:"file_name"`       // 日志文件名称
	FilePath   string `json:"file_path" yaml:"file_path" int:"file_path"`       // 日志路径
	MaxSize    int    `json:"max_size" yaml:"max_size" int:"max_size"`          // 文件最大容量 单位MB
	MaxBackups int    `json:"max_backups" yaml:"max_backups" int:"max_backups"` // 文件最多数量
	MaxAge     int    `json:"max_age" yaml:"max_age" int:"max_age"`             // 日志保留最大时长 单位天，0为不限制
	Compress   bool   `json:"compress" yaml:"compress" int:"compress"`          // 是否压缩
	ToStdout   bool   `json:"to_stdout" yaml:"to_stdout" ini:"to_stdout"`       // 是否输出到标准输出
}

// Config holds the application configuration.
type Config struct {
	ConsoleSeparator string               `json:"console_separator" yaml:"console_separator" ini:"console_separator"` // 控制台分隔符
	LogFormat        string               `json:"log_format" yaml:"log_format" ini:"log_format"`                      // 日志格式 默认普通控制台格式，支持json格式
	LogLevel         string               `json:"log_level" yaml:"log_level" ini:"log_level"`
	LogLevelConfigs  map[string]LogConfig `json:"log_level_configs" yaml:"log_level_configs" ini:"log_level_configs"` // key is log level, value is its configuration
}

// Logger implements the gorm.Logger interface.
type Logger struct {
	logger *zap.Logger
	level  zapcore.LevelEnabler
}

// customTimeEncoder formats time as "2006-01-02 15:04:05".
func customTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format(time.DateTime))
}

// New creates a new instance of Logger with given configuration.
func New(cfg *Config) (*Logger, error) {
	defalutLevel := zapcore.DebugLevel
	if cfg.LogLevel != "" {
		level, err := zapcore.ParseLevel(cfg.LogLevel)
		if err != nil {
			return nil, err
		}
		defalutLevel = level
	}
	var (
		cores         []zapcore.Core
		encoderConfig = zapcore.EncoderConfig{
			TimeKey:        "time",                         // 指定时间戳字段的键名。例如，设置为 "time" 会使得每条日志包含一个名为 time 的字段来表示日志的时间戳。
			LevelKey:       "level",                        // 指定日志级别字段的键名。例如，设置为 "level" 会使得每条日志包含一个名为 level 的字段来表示日志的级别（如 info, error 等）。
			NameKey:        "logger",                       // 指定日志记录器名称字段的键名。如果使用命名日志记录器，此字段将显示日志记录器的名称。
			MessageKey:     "msg",                          // 指定消息字段的键名。这是实际的日志消息文本。
			StacktraceKey:  "stacktrace",                   // 指定堆栈跟踪字段的键名。当发生错误时，可以包含完整的堆栈跟踪信息。
			LineEnding:     zapcore.DefaultLineEnding,      // 指定每条日志记录结束时使用的行尾字符，默认是 \n。对于某些特殊需求（如 Windows 环境），可能需要设置为 \r\n。
			EncodeLevel:    zapcore.LowercaseLevelEncoder,  // 定义如何将日志级别编码为字符串。Zap 提供了几种内置的编码方式，如 LowercaseLevelEncoder（小写字母）、CapitalLevelEncoder（大写字母）、CapitalColorLevelEncoder（带颜色的大写字母）等。
			EncodeTime:     customTimeEncoder,              // 使用自定义的时间编码器
			EncodeDuration: zapcore.SecondsDurationEncoder, // 定义如何格式化持续时间。默认情况下，Zap 使用秒作为单位进行编码，但也可以选择毫秒、微秒或纳秒。
			EncodeCaller:   zapcore.ShortCallerEncoder,     // 定义如何格式化调用者信息。可以选择完整路径（FullCallerEncoder）、短路径（ShortCallerEncoder）或其他自定义格式。
		}
		encoder zapcore.Encoder
	)
	if strings.ToUpper(cfg.LogFormat) == "json" {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		encoderConfig.ConsoleSeparator = tools.Ternary(cfg.ConsoleSeparator == "", " ", cfg.ConsoleSeparator)
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}
	for levelStrs, config := range cfg.LogLevelConfigs {
		for _, levelStr := range strings.Split(levelStrs, ",") {
			level, err := zapcore.ParseLevel(levelStr)
			if err != nil {
				return nil, err
			}
			writers := make([]zapcore.WriteSyncer, 0)
			if config.FilePath != "" {
				fame := filepath.Join(tools.Fileabs(config.FilePath), levelStr+"_"+tools.Ternary(config.FileName == "", "log", config.FileName))
				lumberJackLogger := &lumberjack.Logger{
					Filename:   fame,
					MaxSize:    config.MaxSize,
					MaxAge:     config.MaxAge,
					MaxBackups: config.MaxBackups,
					LocalTime:  true,
					Compress:   config.Compress,
				}
				writers = append(writers, zapcore.AddSync(lumberJackLogger))
			}
			if config.ToStdout {
				writers = append(writers, zapcore.AddSync(os.Stdout))
			}
			// 定义一个自定义的日志级别过滤器，确保每个核心只处理特定范围的日志级别
			levelEnabler := func(minLevel, maxLevel zapcore.Level) zap.LevelEnablerFunc {
				return func(lvl zapcore.Level) bool {
					return lvl >= minLevel && lvl < maxLevel
				}
			}
			core := zapcore.NewCore(encoder, zapcore.NewMultiWriteSyncer(writers...), levelEnabler(level, level+1))
			cores = append(cores, core)
		}
	}
	combinedCore := zapcore.NewTee(cores...)
	return &Logger{
		logger: zap.New(combinedCore, zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel)),
		level:  defalutLevel, // Default to DebugLevel
	}, nil
}

func (l *Logger) LogMode(level logger.LogLevel) logger.Interface {
	enabler := convertLogLevel(level)
	return &Logger{
		logger: l.logger,
		level:  enabler,
	}
}

// ConvertLogLevel converts GORM log levels to Zap LevelEnabler.
func convertLogLevel(level logger.LogLevel) zapcore.LevelEnabler {
	switch level {
	case logger.Silent:
		return zapcore.PanicLevel
	case logger.Error:
		return zapcore.ErrorLevel
	case logger.Warn:
		return zapcore.WarnLevel
	default:
		return zapcore.DebugLevel // Use DebugLevel for Info and other levels
	}
}

func (l *Logger) Info(ctx context.Context, msg string, data ...any) {
	if l.level.Enabled(zapcore.InfoLevel) {
		l.logger.Info(msg, input2Field(data...)...)
	}
}

func (l *Logger) Warn(ctx context.Context, msg string, data ...any) {
	if l.level.Enabled(zapcore.WarnLevel) {
		l.logger.Warn(msg, input2Field(data...)...)
	}
}

func (l *Logger) Error(ctx context.Context, msg string, data ...any) {
	if l.level.Enabled(zapcore.ErrorLevel) {
		l.logger.Error(msg, input2Field(data...)...)
	}
}

func (l *Logger) Debug(ctx context.Context, msg string, data ...any) {
	if l.level.Enabled(zapcore.DebugLevel) {
		l.logger.Debug(msg, input2Field(data...)...)
	}
}

func (l *Logger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.level.Enabled(zapcore.DebugLevel) {
		sql, rows := fc()
		l.Debug(
			ctx,
			"跟踪",
			zap.String("sql", sql),
			zap.Int64("rows", rows),
			zap.Duration("took", time.Since(begin)),
			zap.Error(err),
		)
	}
}

func input2Field(data ...any) (fields []zap.Field) {
	for i, d := range data {
		switch t := d.(type) {
		case []zapcore.Field:
			fields = append(fields, t...)
		case zapcore.Field:
			fields = append(fields, t)
		default:
			fields = append(fields, zap.Any(strconv.Itoa(i), d))
		}
	}
	return
}

// Auto2Field 自动将数据转换为zap.Field
func Auto2Field(data ...any) []zap.Field {
	n := len(data)
	if n == 0 {
		return nil
	}
	count := n / 2
	fields := make([]zap.Field, 0, count)
	for i := 0; i < n-1; i += 2 {
		fields = append(fields, zap.Any(tools.Any2string(data[i]), data[i+1]))
	}

	return fields
}
