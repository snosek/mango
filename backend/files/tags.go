package files

import (
	"fmt"

	"go.senan.xyz/taglib"
)

func ReadTags(fp string) (map[string][]string, error) {
	tags, err := taglib.ReadTags(fp)
	if err != nil {
		return nil, fmt.Errorf("error reading tags: %v", err)
	}
	return tags, nil
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
