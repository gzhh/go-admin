package logger

import (
	"fmt"
)

var (
	LogSavePath string
	LogSaveName string
	LogFileExt  string
)

func getLogFilePath() string {
	return LogSavePath
}

func getLogFileFullPath() string {
	prefixPath := getLogFilePath()
	suffixPath := fmt.Sprintf("%s.%s", LogSaveName, LogFileExt)
	return fmt.Sprintf("%s%s", prefixPath, suffixPath)
}
