package utils

import (
	"context"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

func GetDirPath(ctx context.Context) (string, error) {
	return runtime.OpenDirectoryDialog(ctx, runtime.OpenDialogOptions{})
}

func IsValidCtrlRequest(r string) bool {
	return r == "pause" || r == "resume" || r == "next" || r == "previous" || r == "changePosition" || r == "playTrack"
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
