package utils

import (
	"context"
	"errors"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

func GetDirPath(ctx context.Context) (string, error) {
	dirPath, err := runtime.OpenDirectoryDialog(ctx, runtime.OpenDialogOptions{})
	if err != nil {
		log.Print("Problem reading directory.")
		return "", err
	}
	return dirPath, nil
}

func FetchDirectories(fp string) ([]string, error) {
	entries, err := os.ReadDir(fp)
	if err != nil {
		log.Print(err.Error())
		return nil, err
	}
	result := make([]string, 0)
	for _, f := range entries {
		if f.IsDir() {
			result = append(result)
		}
	}
	return result, nil
}

func FetchAudioFiles(dirPath string) ([]string, error) {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		log.Print(err.Error())
		return nil, err
	}
	result := make([]string, 0)
	err = nil
	for _, f := range entries {
		if !isAudioFile(f.Name()) {
			err = errors.New("Unsupported file type")
			continue
		}
		result = append(result, filepath.Join(dirPath, f.Name()))
	}
	return result, err
}

func isAudioFile(fileName string) bool {
	switch {
	case strings.HasSuffix(fileName, ".flac"):
		fallthrough
	case strings.HasSuffix(fileName, ".mp3"):
		fallthrough
	case strings.HasSuffix(fileName, ".wav"):
		fallthrough
	case strings.HasSuffix(fileName, ".ogg"):
		return true
	default:
		return false
	}
}
