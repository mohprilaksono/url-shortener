package utils

import (
	"os"
	"path/filepath"
	"syscall"
)

func LoadFile() (*os.File, error) {
	file, err := os.OpenFile(filepath.Join("database", "db"), os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		return nil, err
	}

	err = syscall.Flock(int(file.Fd()), syscall.LOCK_EX)
	if err != nil {
		return nil, err
	}

	return file, nil
}

func CloseFile(file *os.File) error {
	err := syscall.Flock(int(file.Fd()), syscall.LOCK_UN)
	if err != nil {
		return err
	}
	
	err = file.Close()
	if err != nil {
		return err
	}

	return nil
}