package catalog

import (
	"os"

	"github.com/labstack/gommon/log"
)

func FetchDirectories(fp string) ([]string, error) {
	entries, err := os.ReadDir(fp)
	if err != nil {
		log.Info(err.Error())
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
