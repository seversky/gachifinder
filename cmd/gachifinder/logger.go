package main

import (
	"io"
	"os"
	"time"

	"github.com/natefinch/lumberjack"
	logger "github.com/sirupsen/logrus"

	"github.com/seversky/gachifinder"
)

func setLogger(config *gachifinder.Config) {
	var timeFormat string
	if len(config.Global.Log.GoTimeFormat) < 1 {
		timeFormat = time.RFC3339
	} else {
		timeFormat = config.Global.Log.GoTimeFormat
	}

	if config.Global.Log.Format == gachifinder.TEXT {
		logger.SetFormatter(&logger.TextFormatter{
			ForceColors: true,
			FullTimestamp: true,
			TimestampFormat: timeFormat})
	} else {
		logger.SetFormatter(&logger.JSONFormatter{TimestampFormat: timeFormat})
	}

	lumberjackLogrotate := &lumberjack.Logger{
		Filename:   config.Global.Log.LogPath,
		MaxSize:    5,  // Max megabytes before log is rotated
		MaxBackups: 90, // Max number of old log files to keep
		MaxAge:     60, // Max number of days to retain log files
		Compress:   true,
	}
	logMultiWriter := io.MultiWriter(os.Stdout, lumberjackLogrotate)
	logger.SetOutput(logMultiWriter)

	level, err := logger.ParseLevel(config.Global.Log.LogLevel)
	if err != nil {
		logger.Fatalln("E! error:", err)
	}

	logger.SetLevel(level)
}