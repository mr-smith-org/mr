package functions

import (
	"html/template"

	"github.com/go-sprout/sprout/sprigin"
)

func GetFuncMap() template.FuncMap {
	fnMap := sprigin.TxtFuncMap()
	fnMap["toYaml"] = ToYaml
	fnMap["getRefFrom"] = GetRefFrom
	fnMap["getPathsByTag"] = GetPathsByTag
	fnMap["getParamsByType"] = GetParamsByType
	fnMap["groupByKey"] = GroupByKey
	fnMap["getFileContent"] = GetFileContent
	fnMap["getFilesList"] = GetFilesList
	fnMap["getFileExtension"] = GetFileExtension
	fnMap["getFileName"] = GetFileName
	fnMap["getFilePath"] = GetFilePath
	fnMap["fileExists"] = FileExists
	fnMap["isDirectory"] = IsDirectory
	fnMap["isFile"] = IsFile
	fnMap["getFileSize"] = GetFileSize
	return fnMap
}
