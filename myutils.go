package main

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func GetCurrentDirectory() string {
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	return strings.Replace(dir, "\\", "/", -1)
}

func get_appname() string {
	full := strings.Replace(os.Args[0], "\\", "/", -1)
	splits := strings.Split(full, "/")
	if len(splits) >= 1 {
		name := splits[len(splits)-1]
		name = strings.TrimSuffix(name, ".exe")
		return name
	}
	return ""
}

func GetLogFilePath() string {
	return filepath.Join(GetCurrentDirectory(), "log", get_appname()+".log")
}

func GetConfigFilePath() string {
	return GetYAMLConfigFilePath()
}

func GetConfigFilePathWithName(name string) string {
	return filepath.Join(GetCurrentDirectory(), name)
}

func GetConfigFilePathWithType(ext string) string {
	return filepath.Join(GetCurrentDirectory(), get_appname()+"."+ext)
}

func GetYAMLConfigFilePath() string {
	return GetConfigFilePathWithType("yaml")
}

func GetJSONConfigFilePath() string {
	return GetConfigFilePathWithType("json")
}

func GetTOMLConfigFilePath() string {
	return GetConfigFilePathWithType("toml")
}

func CreateFolder(path string) error {
	var err error
	if _, err = os.Stat(path); errors.Is(err, fs.ErrNotExist) {
		err = os.MkdirAll(path, 0755)
	}
	return err
}
