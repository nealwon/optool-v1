package common

import (
	"io/ioutil"

	"github.com/go-yaml/yaml"
)

// AuthConfig configures for host authorization
type AuthConfig struct {
	User              string `yaml:"user"`
	Password          string `yaml:"password"`
	PrivateKey        string `yaml:"private_key"`
	PrivateKeyContent string `yaml:"private_key_content"`
	PrivateKeyPhrase  string `yaml:"private_key_phrase"`
	PlainPassword     bool   `yaml:"plain_password"` // 是否是明文的密码(通用password和phrase)
}

// Configure global configure
type Configure struct {
	Server Server `yaml:"server"`
	//Hosts        map[string][]string `yaml:"hosts"` // structure: group => host list
	Auth AuthConfig        `yaml:"auth"`
	Tags map[string]string `yaml:"tags"` // shortcut for frequently used commands
	Gzip bool              `yaml:"-"`    // enable gzip transfer
	//DefaultGroup string              `yaml:"default_group"` // set default host group
}
type Server struct {
	DefaultGroup string              `yaml:"default_group"`
	DefaultPort  int                 `yaml:"default_port"`
	Hosts        map[string][]string `yaml:"hosts"`
}

// C exported parsed configure
var C *Configure

func init() {
	C = &Configure{}
}

// ParseConfig parse configure file
func ParseConfig(f string) error {
	s, err := ioutil.ReadFile(f)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(s, C)
	if err != nil {
		return err
	}
	return nil
}

/**
server:
    default_group: web
    default_port: 22
    hosts:
        web:
            - 172.16.80.129
        ow:
            - 192.168.11.1
        adm:
            - 106.75.60.188
auth:
    user: root
    password: pwds
    private_key: /opt/keys/pk.pem
    private_key_content: ""
    private_key_phrase: ""
    plain_password: true
tags:
    ps: "/bin/ps"
    netstat: "/bin/netstat -lntpu"
    err: "/bin/grep ERROR /data/tmp/phalcon-admin.log_REPLACE_"
*/
