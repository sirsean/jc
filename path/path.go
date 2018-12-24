package path

import (
    "path"
    "io/ioutil"
    "os"
)

var basePath = ".jc"

func RequestDirPath(id string) string {
    return path.Join(basePath, id)
}

func RequestPath(id string) string {
	return path.Join(basePath, id, "request.json")
}

func RequestBodyPath(id string) string {
	return path.Join(basePath, id, "body.json")
}

func ResponsePath(id string) string {
	return path.Join(basePath, id, "response.json")
}

func ListFiles() ([]os.FileInfo, error) {
    return ioutil.ReadDir(basePath)
}

func MakeBasePath() error {
	return os.MkdirAll(basePath, os.ModePerm)
}
