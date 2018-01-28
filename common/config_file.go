package common

import (
	"os"
	"os/exec"
	"os/user"
	"runtime"
	"strings"
)

// ConfigFileList read config file one by one if not exists
var ConfigFileList = []string{
	"./optool.yml",
	homeDir() + "/optool.yml",
	"/etc/optool.yml",
	"/tmp/optool.yml",
}

func init() {
	if runtime.GOOS == "windows" {
		ConfigFileList = []string{
			"./optool.yml",
			homeDir() + "/optool.yml",
			"C:/optool.yml",
			"D:/optool.yml",
		}
	}
}

// homeDir get current user's home dir
func homeDir() string {
	user, err := user.Current()
	if err == nil {
		return user.HomeDir
	}
	if runtime.GOOS == "windows" {
		drive := os.Getenv("HOMEDRIVE")
		path := os.Getenv("HOMEPATH")
		home := drive + path
		if drive == "" || path == "" {
			home = os.Getenv("USERPROFILE")
		}
		return home
	}
	if home := os.Getenv("HOME"); home != "" {
		return home
	}
	output, err := exec.Command("sh", "-c", "eval echo ~$USER").Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(output))
}
