package common

import "runtime"

// ConfigFileList read config file one by one if not exists
var ConfigFileList = []string{
	"./optool.yml",
	"~/optool.yml",
	"/etc/optool.yml",
	"/tmp/optool.yml",
}

func init() {
	if runtime.GOOS == "windows" {
		ConfigFileList = []string{
			"./optool.yml",
			"%USERPROFILE%/optool.yml",
			"C:/optool.yml",
			"D:/optool.yml",
		}
	}
}
