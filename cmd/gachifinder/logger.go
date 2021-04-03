package main

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/natefinch/lumberjack"
	logger "github.com/sirupsen/logrus"

	"github.com/seversky/gachifinder"
)

func setLogger(config *gachifinder.Config) error {
	var timeFormat string
	if len(config.Global.Log.GoTimeFormat) < 1 {
		timeFormat = time.RFC3339
	} else {
		timeFormat = config.Global.Log.GoTimeFormat
	}

	if config.Global.Log.Format == gachifinder.TEXT {
		logger.SetFormatter(&logger.TextFormatter{
			ForceColors: config.Global.Log.ForceColors,
			FullTimestamp: true,
			TimestampFormat: timeFormat})
	} else {
		logger.SetFormatter(&logger.JSONFormatter{TimestampFormat: timeFormat})
	}

	var logMultiWriter io.Writer
	if config.Global.Log.Stdout && len(config.Global.Log.LogPath) > 0 {
		lumberjackLogrotate := &lumberjack.Logger{
			Filename:   config.Global.Log.LogPath,
			MaxSize:    config.Global.Log.MaxSize,  // Max megabytes before log is rotated
			MaxBackups: config.Global.Log.MaxBackups, // Max number of old log files to keep
			MaxAge:     config.Global.Log.MaxAge, // Max number of days to retain log files
			Compress:   config.Global.Log.Compress,
		}
		logMultiWriter = io.MultiWriter(os.Stdout, lumberjackLogrotate)
	} else if config.Global.Log.Stdout {
		logMultiWriter = io.MultiWriter(os.Stdout)
	} else if len(config.Global.Log.LogPath) > 0 {
		lumberjackLogrotate := &lumberjack.Logger{
			Filename:   config.Global.Log.LogPath,
			MaxSize:    config.Global.Log.MaxSize,  // Max megabytes before log is rotated
			MaxBackups: config.Global.Log.MaxBackups, // Max number of old log files to keep
			MaxAge:     config.Global.Log.MaxAge, // Max number of days to retain log files
			Compress:   config.Global.Log.Compress,
		}
		logMultiWriter = io.MultiWriter(lumberjackLogrotate)
	} else {
		return fmt.Errorf("Both of LogPath and Stdout is not configured. Exit code = %d", 1)
	}

	logger.SetOutput(logMultiWriter)

	level, err := logger.ParseLevel(config.Global.Log.LogLevel)
	if err != nil {
		logger.WithField("error", err).Fatalln("Wrong log level")
	}

	logger.SetLevel(level)

	return nil
}