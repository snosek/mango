package utils

import (
	"context"
	"os"
	"path/filepath"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

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

func IsValidCtrlRequest(r string) bool {
	return r == "pause" || r == "resume" || r == "next" || r == "previous" || r == "changePosition" || r == "playTrack"
}

func FirstOrEmpty(s []string) string {
	if len(s) > 0 {
		return s[0]
	}
	return ""
}

func FirstOrFallback(primary, fallback []string) []string {
	if len(primary) > 0 {
		return primary
	}
	return fallback
}

func HashTitle(title string) string {
	h := fnv32a(title)
	return string(h)
}

func fnv32a(s string) []byte {
	hash := uint32(2166136261)
	for i := 0; i < len(s); i++ {
		hash ^= uint32(s[i])
		hash *= 16777619
	}
	b := make([]byte, 4)
	b[0] = byte(hash >> 24)
	b[1] = byte(hash >> 16)
	b[2] = byte(hash >> 8)
	b[3] = byte(hash)
	return b
}
