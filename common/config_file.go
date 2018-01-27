// +build: darwin
package common

import "runtime"

var ConfigFileList = []string{
	"./deployer.yml",
	"~/deployer.yml",
	"/etc/deployer.yml",
	"/tmp/deployer.yml",
}

func init() {
	if runtime.GOOS == "windows" {
		ConfigFileList = []string{
			"./deployer.yml",
			"%USERPROFILE%/deployer.yml",
			"C:/deployer.yml",
			"D:/deployer.yml",
		}
	}
}
