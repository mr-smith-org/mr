package functions

import (
	"io"
	"strings"

	"github.com/spf13/afero"
)

func GetFileContent(filePath string) string {
	fs := afero.NewOsFs()
	file, err := fs.Open(filePath)
	if err != nil {
		return ""
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		return ""
	}

	return string(content)
}

func GetFilesList(path string) []string {
	fs := afero.NewOsFs()
	files, err := afero.ReadDir(fs, path)
	if err != nil {
		return []string{}
	}
	var fileNames []string
	for _, file := range files {
		fileNames = append(fileNames, file.Name())
	}
	return fileNames
}

func GetFileExtension(filePath string) string {
	splitFilePath := strings.Split(filePath, ".")
	if len(splitFilePath) > 1 {
		return splitFilePath[len(splitFilePath)-1]
	}
	return ""
}

func GetFileName(filePath string) string {
	splitFilePath := strings.Split(filePath, "/")
	return splitFilePath[len(splitFilePath)-1]
}

func GetFilePath(filePath string) string {
	splitFilePath := strings.Split(filePath, "/")
	return strings.Join(splitFilePath[:len(splitFilePath)-1], "/")
}

func FileExists(filePath string) bool {
	fs := afero.NewOsFs()
	exists, err := afero.Exists(fs, filePath)
	if err != nil {
		return false
	}
	return exists
}

func IsDirectory(filePath string) bool {
	fs := afero.NewOsFs()
	fileInfo, err := fs.Stat(filePath)
	if err != nil {
		return false
	}
	return fileInfo.IsDir()
}

func IsFile(filePath string) bool {
	fs := afero.NewOsFs()
	fileInfo, err := fs.Stat(filePath)
	if err != nil {
		return false
	}
	return !fileInfo.IsDir()
}

func GetFileSize(filePath string) int64 {
	fs := afero.NewOsFs()
	fileInfo, err := fs.Stat(filePath)
	if err != nil {
		return 0
	}
	return fileInfo.Size()
}
