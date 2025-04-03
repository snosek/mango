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

func Hash(title string) string {
	h := fnv64a(title)
	return string(h)
}

func fnv64a(s string) []byte {
	hash := uint64(14695981039346656037)
	const prime uint64 = 1099511628211
	for i := 0; i < len(s); i++ {
		hash ^= uint64(s[i])
		hash *= prime
	}
	b := make([]byte, 8)
	b[0] = byte(hash >> 56)
	b[1] = byte(hash >> 48)
	b[2] = byte(hash >> 40)
	b[3] = byte(hash >> 32)
	b[4] = byte(hash >> 24)
	b[5] = byte(hash >> 16)
	b[6] = byte(hash >> 8)
	b[7] = byte(hash)
	return b
}
