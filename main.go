package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"

	"github.com/go-yaml/yaml"
	"github.com/nealwon/deployer/common"
)

// REPLACEMENT variable to replace by args
const REPLACEMENT = "_REPLACE_"

var (
	pConfigFile   = flag.String("config", "/deployer.yml", "set config file path")
	pTag          = flag.String("t", "", "set tagged command")
	pTagArgs      = flag.String("ta", "", "append tagged command parameters, overflow params will dropped,params separated by comma(,).\n\t to replace in tags use string: _REPLACE_")
	pTagPrint     = flag.Bool("tp", false, "print tag line")
	pTagList      = flag.Bool("tl", false, "list all tags")
	pGzip         = flag.Bool("gz", false, "enable gzip for transfer.remote host must has executable: /usr/bin/gzip")
	pGroup        = flag.String("g", "", "set hosts group name")
	pUser         = flag.String("u", "", "set ssh auth user")
	pOutput       = flag.String("o", "-", "set output file")
	pCommand      = flag.String("c", "", "direct run command")
	pNoHeader     = flag.Bool("nh", false, "no header output")
	pNoHost       = flag.Bool("ns", false, "no server ip output")
	pHost         = flag.String("host", "", "set run host")
	pPort         = flag.Int("port", 0, "set default ssh port")
	pPrivateKey   = flag.String("key", "", "set private key")
	pVerbose      = flag.Bool("v", false, "verbose all configures")
	pSampleConfig = flag.Bool("V", false, "print sample configure")
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	flag.Parse()
	if *pSampleConfig {
		printSample()
		os.Exit(0)
	}
	var err error
	if _, err = os.Stat(*pConfigFile); err != nil {
		for _, cf := range common.ConfigFileList {
			if _, err = os.Stat(cf); err == nil {
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
	if cmd == "" {
		if *pCommand != "" {
			cmd = *pCommand
		}
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
	rc.PrettyPrint(wo, os.Stderr, !*pNoHeader, !*pNoHost)
}

func printSample() {
	fmt.Print(`server:
  default_group: web
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
