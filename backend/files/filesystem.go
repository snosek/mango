package files

import (
	"image"
	"image/jpeg"
	"os"
	"path/filepath"
)

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

func GetModificationTime(fp string) string {
	albumDir, err := os.Open(fp)
	if err != nil {
		return ""
	}
	defer albumDir.Close()
	albumDirStat, err := albumDir.Stat()
	if err != nil {
		return ""
	}
	albumModTime := albumDirStat.ModTime()
	return albumModTime.String()
}

func IsSystemFile(fp string) bool {
	fileBase := filepath.Base(fp)
	for _, pattern := range systemFilePatterns {
		matched, _ := filepath.Match(pattern, fileBase)
		if matched {
			return true
		}
	}
	return false
}

var systemFilePatterns = []string{
	".DS_Store",
	"._*",
	".Trash*",
	".fseventsd",
	".Spotlight-V100",
	".TemporaryItems",
	".apdisk",
	"Thumbs.db",
	"desktop.ini",
	"$RECYCLE.BIN",
	".Trash-1000",
	".nfs*",
}

func ReadAlbumCover(fp string) (image.Image, error) {
	file, err := os.Open(filepath.Join(fp, "folder.jpg"))
	if err != nil {
		return nil, err
	}
	defer file.Close()
	cover, err := jpeg.Decode(file)
	if err != nil {
		return nil, err
	}
	return cover, nil
}
