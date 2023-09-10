package logger

import (
	"encoding/csv"
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/fatih/color"
	"gopkg.in/yaml.v2"
)

type LogLevel int

const (
	LogLevelDefault LogLevel = iota
	LogLevelInfo
	LogLevelWarning
	LogLevelError
	LogLevelDanger
)

type LogSettings struct {
	LogLevel         string `yaml:"log_level"`
	LogFileBool      bool   `yaml:"log_file_bool"`
	LogFileName      string `yaml:"log_file_name"`
	LogFileExtension string `yaml:"log_file_extension"`
	LogFileMaxBytes  int    `yaml:"log_file_max_bytes"`
	LogFormat        string `yaml:"log_format"`
	LogDateFormat    string `yaml:"log_date_format"`
}

type LogConfig struct {
	LogSettings LogSettings `yaml:"log_settings"`
	LogColors   struct {
		InfoColor    string `yaml:"info_color"`
		WarningColor string `yaml:"warning_color"`
		ErrorColor   string `yaml:"error_color"`
		DangerColor  string `yaml:"danger_color"`
		DefaultColor string `yaml:"default_color"`
	} `yaml:"log_colors"`
}

type Logger struct {
	logLevel    LogLevel
	logFile     *os.File
	csvWriter   *csv.Writer
	logColors   map[LogLevel]*color.Color
	logFileName string
}

func InitialiseLogger(configFile string) (*Logger, error) {
	if configFile == "" {
		// default config file
		configFile = "config.yaml"
	}
	config, err := readConfig(configFile)
	if err != nil {
		return nil, err
	}

	logger := &Logger{
		logLevel:  getLogLevel(config.LogSettings.LogLevel),
		logColors: map[LogLevel]*color.Color{},
	}

	logger.initLogColors(config)

	initialisingMessage := fmt.Sprintf("[%s]\n%s%s%s%s\n%s%s\n",
		color.HiYellowString("Logger initialized"),
		"\x1b[0m",
		color.BlueString("Configuration: "),
		color.WhiteString(configFile),
		"\x1b[0m",
		color.HiBlueString("Log file: "),
		color.WhiteString(config.LogSettings.LogFileName+config.LogSettings.LogFileExtension),
	)

	if config.LogSettings.LogFileBool {
		logFileName := config.LogSettings.LogFileName + config.LogSettings.LogFileExtension
		logFile, err := os.OpenFile(logFileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			return nil, fmt.Errorf("%v", err)
		}
		logger.logFile = logFile
		logger.csvWriter = csv.NewWriter(logFile)
		logger.logFileName = logFileName
		fmt.Print(initialisingMessage)

	} else {
		initialisingMessage_nlf := initialisingMessage + fmt.Sprintf("%s%s%s\n",
			color.HiBlueString("Log file: "), color.WhiteString("False"), "\x1b[0m")
		fmt.Print(initialisingMessage_nlf)
	}

	return logger, nil
}

func (l *Logger) initLogColors(config LogConfig) {
	l.logColors[LogLevelInfo] = getColorFromName(config.LogColors.InfoColor)
	l.logColors[LogLevelWarning] = getColorFromName(config.LogColors.WarningColor)
	l.logColors[LogLevelError] = getColorFromName(config.LogColors.ErrorColor)
	l.logColors[LogLevelDanger] = getColorFromName(config.LogColors.DangerColor)
	l.logColors[LogLevelDefault] = getColorFromName(config.LogColors.DefaultColor)
}

var colorMap = map[string]color.Attribute{
	"blue":    color.BgBlue,
	"yellow":  color.BgYellow,
	"red":     color.BgRed,
	"magenta": color.BgHiMagenta,
	"cyan":    color.BgHiCyan,
}

func getColorFromName(colorName string) *color.Color {
	if colorName == "" {
		return color.New()
	}

	attr, exists := colorMap[colorName]
	if !exists {
		return color.New()
	}

	return color.New(attr)
}

func (l *Logger) LogMessage (level LogLevel, message string) {
	if level >= l.logLevel {
		timestamp := time.Now().Format(time.RFC3339)
		sourceFile, line := getSourceLocation()
		levelStr := level.String()
		logEntry := []string{timestamp, levelStr, message, sourceFile, line}
		l.writeToLogFile(logEntry)
		l.printToConsole(logEntry, level)
	}
}

func (l *Logger) writeToLogFile(logEntry []string) {
	if l.logFile != nil {
		l.csvWriter.Write(logEntry)
		l.csvWriter.Flush()
	}
}

func (l *Logger) printToConsole(logEntry []string, level LogLevel) {
	colorizer := l.logColors[level]
	timestamp := logEntry[0]
	levelStr := logEntry[1]
	message := logEntry[2]
	file := logEntry[3]
	fmt.Printf("[%s] %s %s -> %s\n", timestamp, colorizer.Sprint(levelStr), message, file)
}

func (l *Logger) Close() {
	if l.logFile != nil {
		l.logFile.Close()
	}
}

func getLogLevel(levelName string) LogLevel {
	switch strings.ToLower(levelName) {
	case "info":
		return LogLevelInfo
	case "warning":
		return LogLevelWarning
	case "error":
		return LogLevelError
	case "danger":
		return LogLevelDanger
	default:
		return LogLevelDefault
	}
}

func getSourceLocation() (string, string) {
	_, sourceFile, line, ok := runtime.Caller(2)
	if ok {
		return sourceFile, fmt.Sprintf("%d", line)
	}
	return "unknown", "0"
}

func (level LogLevel) String() string {
	switch level {
	case LogLevelInfo:
		return "INFO"
	case LogLevelWarning:
		return "WARNING"
	case LogLevelError:
		return "ERROR"
	case LogLevelDanger:
		return "DANGER"
	default:
		return "LOG"
	}
}

func readConfig(configFile string) (LogConfig, error) {
	var config LogConfig
	var path string = "config/"
	file, err := os.Open(path + configFile)
	if err != nil {
		return config, err
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(&config)
	return config, err
}
