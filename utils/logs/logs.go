package logs

import (
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"syscall"

	"github.com/mohprilaksono/url-shortener/utils"
)

func initLog(prefix, message string) error {
	file, err := os.OpenFile(filepath.Join("storage", "log", "app.log"), os.O_RDWR|os.O_CREATE, fs.ModePerm)
	if err != nil {
		return err
	}

	defer utils.CloseFile(file)

	err = syscall.Flock(int(file.Fd()), syscall.LOCK_EX)
	if err != nil {
		return err
	}

	logger := log.New(file, prefix, log.Ldate|log.Ltime|log.Lmicroseconds|log.Llongfile|log.LUTC|log.Lmsgprefix)

	logger.Println(message)
	
	return nil
}

func Info(message string) error {
	return initLog("INFO", message)
}

func Debug(message string) error {
	return initLog("DEBUG", message)
}

func Warn(message string) error {
	return initLog("WARN", message)
}

func Trace(message string) error {
	return initLog("TRACE", message)
}

func Err(message string) error {
	return initLog("ERROR", message)
}

func Fatal(message string) error {
	return initLog("FATAL", message)
}