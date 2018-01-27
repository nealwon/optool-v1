package common

import (
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"golang.org/x/crypto/ssh"
)

// RemoteCommand remote command structure
type RemoteCommand struct {
	lock     sync.Mutex
	wg       *sync.WaitGroup
	Hosts    []string
	Cmd      string
	PipeMode bool

	PipeChan  chan bool
	PipeIn    map[string]io.WriteCloser
	PipeOut   map[string]io.Reader
	PipeError map[string]io.Reader

	Output  map[string]string
	Error   map[string]string
	Running map[string]*ssh.Session
}

// NewRemoteCommand prepare a remote execution
func NewRemoteCommand(hosts []string, cmd string) *RemoteCommand {
	if C.Gzip {
		cmd = cmd + " | /usr/bin/gzip -f"
	}
	return &RemoteCommand{
		lock:      sync.Mutex{},
		wg:        &sync.WaitGroup{},
		Hosts:     hosts,
		Cmd:       cmd,
		Output:    make(map[string]string),
		Error:     make(map[string]string),
		Running:   make(map[string]*ssh.Session),
		PipeIn:    make(map[string]io.WriteCloser),
		PipeOut:   make(map[string]io.Reader),
		PipeError: make(map[string]io.Reader),
		PipeChan:  make(chan bool),
	}
}

// Start run remote command
func (rc *RemoteCommand) Start() error {
	cfg := &ssh.ClientConfig{
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         time.Second * 10,
	}
	if C.Auth.User != "" {
		cfg.User = C.Auth.User
		if _, err := os.Stat(C.Auth.PrivateKey); err != nil {
			return err
		}
		key, err := ioutil.ReadFile(C.Auth.PrivateKey)
		if err != nil {
			return err
		}
		var signer ssh.Signer
		if C.Auth.PrivateKeyPhrase == "" {
			signer, err = ssh.ParsePrivateKey(key)
		} else {
			signer, err = ssh.ParsePrivateKeyWithPassphrase(key, []byte(C.Auth.PrivateKeyPhrase))
		}
		if err != nil {
			return err
		}
		cfg.Auth = []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		}
	}
	for _, host := range rc.Hosts {
		rc.wg.Add(1)
		//L.Info("host=", host)
		go rc.execute(host, cfg)
	}
	if rc.PipeMode {
		rc.PipeChan <- true
	}
	rc.wg.Wait()
	return nil
}

// execute execute command at host
func (rc *RemoteCommand) execute(host string, cfg *ssh.ClientConfig) {
	ohost := host
	if strings.Index(host, ":") < 0 {
		host = host + ":" + strconv.Itoa(C.Server.DefaultPort)
	}
	client, err := ssh.Dial("tcp", host, cfg)
	if err != nil {
		rc.lock.Lock()
		rc.Error[ohost] = err.Error()
		rc.lock.Unlock()
		rc.wg.Done()
		return
	}
	defer client.Close()
	sess, err := client.NewSession()
	if err != nil {
		rc.lock.Lock()
		rc.Error[ohost] = err.Error()
		rc.lock.Unlock()
		return
	}
	defer sess.Close()
	var o []byte
	var e error
	// @todo std pipes
	if rc.PipeMode {
		rc.Running[ohost] = sess
		//rc.PipeIn[ohost], e = sess.StdinPipe()
		rc.PipeOut[ohost], e = sess.StdoutPipe()
		rc.PipeError[ohost], e = sess.StderrPipe()
		e = sess.Start(rc.Cmd)
		e = sess.Wait()
		rc.wg.Done()
		return
	}
	o, e = sess.Output(rc.Cmd)
	//L.Debugf("RemoteCommand: [%s] cmd=%s, output=%s, error=%s\n", ohost, rc.Cmd, string(o), e)
	rc.lock.Lock()
	rc.Output[ohost] = string(o)
	if e != nil {
		rc.Error[ohost] = e.Error()
	}
	rc.lock.Unlock()
	rc.wg.Done()
}

// ClosePipe close ssh sessions
func (rc *RemoteCommand) ClosePipe() {
	for _, sess := range rc.Running {
		sess.Signal(ssh.SIGTERM)
		sess.Close()
	}
}

// PrettyPrint print output and errors
func (rc *RemoteCommand) PrettyPrint(wo io.Writer, we io.Writer, withHeader bool, withHost bool) {
	if len(rc.Error) > 0 && withHost {
		if withHeader {
			we.Write([]byte("================================= ERROR ================================="))
		}
		for h, e := range rc.Error {
			fmt.Fprintln(we, h, "\n", e)
		}
	}
	if len(rc.Output) > 0 {
		if withHeader {
			fmt.Fprintln(wo, "================================= OUTPUT =================================")
		}
		for h, o := range rc.Output {
			if C.Gzip {
				gr, err := gzip.NewReader(strings.NewReader(o))
				if err != nil {
					log.Println(err)
					continue
				}
				defer gr.Close()
				data, err := ioutil.ReadAll(gr)
				if err != nil {
					log.Println(err)
				}
				if withHost {
					wo.Write([]byte(h + ": \n"))
				}
				wo.Write(data)
				wo.Write([]byte("\n"))
				continue
			}
			if withHost {
				wo.Write([]byte(h + ": \n"))
			}
			wo.Write([]byte(o))
			wo.Write([]byte("\n"))
		}
	}
}
