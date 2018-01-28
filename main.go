package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"strings"

	"github.com/go-yaml/yaml"
	"github.com/nealwon/optool/common"
)

// REPLACEMENT variable to replace by args
const REPLACEMENT = "_REPLACE_"
const (
	// NoHeader no output/error header
	NoHeader = 1 << 0
	//NoServer no server ip line
	NoServer = 1 << 1
)

// OptoolVersion define current version
const OptoolVersion = "v0.2"

var (
	pConfigFile   = flag.String("config", "/optool.yml", "set config file path")
	pTag          = flag.String("t", "", "set tagged command")
	pTagArgs      = flag.String("ta", "", "append tagged command parameters, overflow params will be dropped, separated by comma(,).\n\t to replace in tags use string: _REPLACE_")
	pTagPrint     = flag.Bool("tp", false, "print tag line")
	pTagList      = flag.Bool("tl", false, "list all tags")
	pGzip         = flag.Bool("gz", false, "enable gzip for transfer./usr/bin/gzip must be executable at remote host")
	pGroup        = flag.String("g", "", "set default group name for hosts")
	pUser         = flag.String("u", "", "set ssh auth user")
	pOutput       = flag.String("o", "-", "set output file")
	pCommand      = flag.String("x", "", "execute command directly")
	pScript       = flag.String("s", "", "read commands from script")
	pNoHeader     = flag.Int("nh", 0, "(1)1<<0=no header,(2)1<<1=no server ip,3=none")
	pHost         = flag.String("host", "", "set run host")
	pPort         = flag.Int("port", 0, "set default ssh port")
	pPrivateKey   = flag.String("key", "", "set private key")
	pVerbose      = flag.Bool("v", false, "verbose all configs")
	pSampleConfig = flag.Bool("V", false, "print sample configure")
	pVersion      = flag.Bool("version", false, "print version and exit")
	pEncrypt      = flag.Bool("encrypt", false, "encrypt a password/phrase")
)

func main() {
	flag.Parse()
	if *pVersion {
		fmt.Println("Opstool", OptoolVersion)
		os.Exit(0)
	}
	if *pEncrypt {
		doEncryption()
		os.Exit(0)
	}
	runtime.GOMAXPROCS(runtime.NumCPU())
	if *pSampleConfig {
		printSample()
		os.Exit(0)
	}
	var err error
	if _, err = os.Stat(*pConfigFile); err != nil {
		for _, cf := range common.ConfigFileList {
			_, err = os.Stat(cf)
			if *pVerbose {
				fmt.Println(cf, err)
			}
			if err == nil {
				*pConfigFile = cf
				break
			}
		}
	}

	if err = common.ParseConfig(*pConfigFile); err != nil {
		log.Fatalln("ParseConfig: ", err)
	}
	// tag list,print,arg parse
	if *pTagList {
		common.TagList() // exit
	}
	if *pTagPrint && *pTag != "" {
		common.TagPrint(*pTag) // exit
	}
	var tagArgs []string
	if *pTagArgs != "" {
		tagArgs = strings.Split(*pTagArgs, ",")
	}
	// output handle
	wo := os.Stdout
	if *pOutput != "-" {
		wo, err = os.OpenFile(*pOutput, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0755)
		if err != nil {
			log.Fatalln("Open output: ", err)
		}
		defer wo.Close()
	}
	// gzip or not
	common.C.Gzip = *pGzip
	// user
	if *pUser != "" {
		common.C.Auth.User = *pUser
	}
	// hosts
	var hosts []string
	if *pHost != "" {
		hosts = []string{*pHost}
	} else {
		var ok bool
		if *pGroup != "" {
			common.C.Server.DefaultGroup = *pGroup
		}
		if hosts, ok = common.C.Server.Hosts[common.C.Server.DefaultGroup]; !ok {
			log.Fatalln("Host group not found. Group: ", common.C.Server.DefaultGroup)
		}
	}
	// port
	if *pPort > 0 && *pPort < 65536 {
		common.C.Server.DefaultPort = *pPort
	}
	// private key
	if *pPrivateKey != "" {
		common.C.Auth.PrivateKey = *pPrivateKey
		common.C.Auth.PrivateKeyPhrase = ""
	}
	// command
	var cmd string
	if *pTag != "" {
		for t, tcmd := range common.C.Tags {
			if t == *pTag {
				cmd = tcmd
				break
			}
		}
	}
	// direct command
	if cmd == "" {
		if *pCommand != "" {
			cmd = *pCommand
		}
	}

	// script
	if *pScript != "" {
		script, err := ioutil.ReadFile(*pScript)
		if err != nil {
			log.Fatalln("Script: ", err)
		}
		cmd = string(script)
	}
	if cmd == "" {
		log.Fatal("Command cannot be empty")
	}
	toReplaceCount := strings.Count(cmd, REPLACEMENT)
	if len(tagArgs) < toReplaceCount {
		log.Fatalln("Parameter is not enough. Required is", toReplaceCount)
	}
	for i := 0; i < toReplaceCount; i++ {
		cmd = strings.Replace(cmd, REPLACEMENT, tagArgs[i], 1)
	}
	if *pVerbose {
		fmt.Println("Config file: ", *pConfigFile)
		fmt.Println("================================ Config ===================================")
		ox, _ := yaml.Marshal(common.C)
		os.Stdout.Write(ox)
		os.Exit(0)
	}
	// run
	//cmd := "/bin/cat /data/tmp/phalcon-cli.log"
	rc := common.NewRemoteCommand(hosts, cmd)
	if err := rc.Start(); err != nil {
		log.Fatalln(err)
	}
	rc.PrettyPrint(wo, os.Stderr, (*pNoHeader&NoHeader) > 0, (*pNoHeader&NoServer) > 0)
}

func printSample() {
	fmt.Print(`server:
  default_group: vm
  default_port: 22
  hosts:
	vm:
	  - 172.16.80.129
	router:
	  - 192.168.11.1
auth:
  user: root
  password: {my password}
  # higher priority than password
  private_key: {/path/to/my/private/key.pem}
  # not used
  #private_key_content: ""
  private_key_phrase: ""
  plain_password: true
tags:
  ps: "/bin/ps"
  netstat: "/bin/netstat -lntpu"
  err: "/bin/grep ERROR /var/log/nginx/error.log_REPLACE_"
`)
}

func doEncryption() {
	var str string
	var restr string
	fmt.Printf("   Input string:")
	fmt.Scanln(&str)
	fmt.Printf("Re-input string:")
	fmt.Scanln(&restr)
	if str != restr {
		fmt.Println("Your input mismatch.")
		os.Exit(1)
	}
	fmt.Println(string(common.Encrypt(str)))
}
