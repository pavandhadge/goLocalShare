package utils

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

func SecurePath(base, requested string) (string, error) {
	if strings.Contains(requested, "../") || strings.Contains(requested, "~/") ||
		strings.Contains(requested, "..\\") || strings.Contains(requested, "\\") {
		return "", errors.New("path traversal attempt")
	}

	joined := filepath.Join(base, requested)
	absPath, err := filepath.Abs(joined)
	if err != nil {
		return "", err
	}

	if !strings.HasPrefix(absPath, base) {
		return "", errors.New("path outside base directory")
	}

	realPath, err := filepath.EvalSymlinks(absPath)
	if err != nil {
		return "", err
	}

	if !strings.HasPrefix(realPath, base) {
		return "", errors.New("symlink points outside base directory")
	}

	return realPath, nil
}

func SecureStat(path string) (os.FileInfo, error) {
	fi, err := os.Lstat(path)
	if err != nil {
		return nil, err
	}

	if !fi.Mode().IsRegular() && !fi.Mode().IsDir() {
		return nil, errors.New("special files not allowed")
	}

	return fi, nil
}

func SecureOpen(path string) (*os.File, error) {
	file, err := os.OpenFile(path, os.O_RDONLY|0x20000, 0)
	if err != nil {
		return nil, err
	}

	fi, err := file.Stat()
	if err != nil {
		file.Close()
		return nil, err
	}

	if fi.Mode()&os.ModeSymlink != 0 {
		file.Close()
		return nil, errors.New("symlinks not allowed")
	}

	return file, nil
} 