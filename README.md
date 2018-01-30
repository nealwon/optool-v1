# optool
A tool to execute commands,transfer files on multiple remote hosts

------------
### Usage:
```bash
Usage:
  -V	print sample configure
  -config string
    	set config file path (default "/optool.yml")
  -encrypt
    	encrypt a password/phrase
  -g string
    	set default group name for hosts
  -get string
    	get a file from remote host
  -gz
    	enable gzip for transfer./usr/bin/gzip must be executable at remote host
  -host string
    	set run host
  -key string
    	set private key
  -nh int
    	(1)1<<0=no header,(2)1<<1=no server ip,3=none
  -o string
    	set output file (default "-")
  -override
    	Override remote file if exists
  -path string
    	set path.if get is set this is local path,if put is set this is remote path
  -port int
    	set default ssh port
  -put string
    	put a file to remote host
  -s string
    	read commands from script
  -t string
    	set tagged command
  -ta string
    	append tagged command parameters, overflow params will be dropped, separated by comma(,).
	 to replace in tags use string: _REPLACE_
  -tl
    	list all tags
  -tp
    	print tag line
  -u string
    	set ssh auth user
  -v	verbose all configs
  -version
    	print version and exit
  -x string
    	execute command directly
```

### Sample configure:
```yaml
server:
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
# transfer_max_size: 1099511627776 #100MB
```
