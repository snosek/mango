package utils

import (
	"context"
	"hash/fnv"
	"os"
	"path/filepath"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

func Hash(s string) string {
	hash := fnv.New32a()
	hash.Write([]byte(s))
	return string(hash.Sum([]byte{}))
}

func GetDirPath(ctx context.Context) (string, error) {
	return runtime.OpenDirectoryDialog(ctx, runtime.OpenDialogOptions{})
}

func FetchDirectories(fp string) ([]string, error) {
	entries, err := os.ReadDir(fp)
	if err != nil {
		return nil, err
	}
	var dirs []string
	for _, f := range entries {
		if f.IsDir() {
			dirs = append(dirs, filepath.Join(fp, f.Name()))
		}
	}
	return dirs, nil
}

func FetchAudioFiles(dirPath string) ([]string, error) {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}
	var audioFiles []string
	for _, f := range entries {
		if isAudioFile(f.Name()) {
			audioFiles = append(audioFiles, filepath.Join(dirPath, f.Name()))
		}
	}
	return audioFiles, nil
}

func isAudioFile(fileName string) bool {
	ext := filepath.Ext(fileName)
	return ext == ".flac" || ext == ".mp3" || ext == ".wav" || ext == ".ogg"
}

func IsValidCtrlRequest(request string) bool {
	return request == "pause" || request == "resume" || request == "next" || request == "previous"
}
