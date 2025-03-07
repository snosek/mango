package utils

import (
	"context"
	"log"
	"os"

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

func GetDirectories(fp string) ([]string, error) {
	entries, err := os.ReadDir(fp)
	if err != nil {
		log.Print(err.Error())
		return nil, err
	}
	result := make([]string, 0)
	for _, f := range entries {
		if f.IsDir() {
			result = append(result, f.Name())
		}
	}
	return result, nil
}
