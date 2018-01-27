# deployer
A tool to execute commands on multiple remote hosts

------------
### Usage:
```bash
Usage:
  -V    print sample configure
  -config string
        set config file path (default "/deployer.yml")
  -g string
        set default group name for hosts
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
  -port int
        set default ssh port
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
  -v    verbose all configs
  -x string
        execute command directly
```

### Sample configure:
```yaml
server:
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
```