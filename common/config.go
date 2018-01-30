package common

import (
	"io/ioutil"
	"os"

	"golang.org/x/crypto/ssh"

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
	TransferMaxSize int64 `yaml:"transfer_max_size"`
}

// Server server groups and default port/group config
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

// GetAuth get auth method list from configs
func GetAuth() (auth []ssh.AuthMethod, err error) {
	password := C.Auth.Password
	if !C.Auth.PlainPassword {
		password = string(Decrypt(C.Auth.Password))
	}
	if C.Auth.PrivateKey != "" {
		if _, err := os.Stat(C.Auth.PrivateKey); err != nil {
			return nil, err
		}
		key, err := ioutil.ReadFile(C.Auth.PrivateKey)
		if err != nil {
			return nil, err
		}
		var signer ssh.Signer
		if C.Auth.PrivateKeyPhrase == "" {
			signer, err = ssh.ParsePrivateKey(key)
		} else {
			passphrase := []byte(C.Auth.PrivateKeyPhrase)
			if !C.Auth.PlainPassword {
				passphrase = Decrypt(C.Auth.PrivateKeyPhrase)
			}
			signer, err = ssh.ParsePrivateKeyWithPassphrase(key, passphrase)
		}
		if err != nil {
			return nil, err
		}
		auth = []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		}
		if password != "" {
			auth = append(auth, ssh.Password(password))
		}
	} else {
		auth = []ssh.AuthMethod{
			ssh.Password(password),
		}
	}
	return
}
