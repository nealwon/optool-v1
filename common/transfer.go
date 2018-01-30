package common

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

const (
	// TransferGet get file from remote servers
	TransferGet = "GET"
	// TransferPut put file to remote servers
	TransferPut = "PUT"
)

// Transfer transfer files via ssh
type Transfer struct {
	Inited     bool
	Method     string // GET,PUT
	LocalPath  string
	RemotePath string
	Recursive  bool
	Hosts      []string
	Clients    map[string]*ssh.Client
	SftpClient map[string]*sftp.Client
	Override   bool // override remote existed file?
}

// NewTransfer get file transfer instance
func NewTransfer(method, localPath, remotePath string, hosts []string) *Transfer {
	return &Transfer{
		Inited:     true,
		Method:     method,
		LocalPath:  localPath,
		RemotePath: remotePath,
		Recursive:  false,
		Clients:    make(map[string]*ssh.Client),
		SftpClient: make(map[string]*sftp.Client),
		Hosts:      hosts,
		Override:   false,
	}
}

// Start start file transfer
func (t *Transfer) Start() (err error) {
	if err = t.initClient(); err != nil {
		return
	}
	// close connections
	defer func() {
		for _, sc := range t.SftpClient {
			sc.Close()
		}
		for _, c := range t.Clients {
			c.Close()
		}
	}()
	if t.Method == TransferGet {
		return t.batchGet()
	}
	if t.Method == TransferPut {
		return t.batchPut()
	}
	return nil
}

func (t *Transfer) batchGet() (err error) {
	fi, err := os.Stat(t.LocalPath)
	if err != nil {
		err = os.MkdirAll(t.LocalPath, 0755)
		if err != nil {
			return
		}
	} else {
		if !fi.IsDir() {
			log.Fatalln("Local path cannot be a file")
		}
	}
	wg := sync.WaitGroup{}
	for h, sc := range t.SftpClient {
		c := t.Clients[h]
		wg.Add(1)
		go func(sc *sftp.Client, c *ssh.Client) {
			defer wg.Done()
			t.get(sc, c, t.RemotePath, t.LocalPath)
		}(sc, c)
	}
	wg.Wait()
	return
}

func (t *Transfer) batchPut() (err error) {
	fi, err := os.Stat(t.LocalPath)
	if err != nil {
		return
	}
	if fi.IsDir() {
		return errors.New("Local is dir,recursive transfer not supported now")
	}
	wg := sync.WaitGroup{}
	for _, sc := range t.SftpClient {
		wg.Add(1)
		go func(sc *sftp.Client) {
			defer wg.Done()
			err := t.put(sc, t.LocalPath, t.RemotePath)
			fmt.Println(err)
		}(sc)
	}
	return
}

func (t *Transfer) get(sc *sftp.Client, c *ssh.Client, remotePath, localPath string) (err error) {
	fi, err := sc.Stat(remotePath)
	if err != nil {
		return
	}
	if fi.IsDir() {
		return errors.New("Remote dir get is not supported")
	}
	basename := path.Base(fi.Name())
	srcFile, err := sc.Open(remotePath)
	if err != nil {
		return
	}
	defer srcFile.Close()
	addr := c.Conn.RemoteAddr().String()
	xaddr := strings.Split(addr, ":")
	addr = strings.Replace(xaddr[0], ".", "-", -1)
	exp := strings.Split(basename, ".")
	var ext, prefName string
	lenth := len(exp)
	if lenth > 1 {
		ext = exp[lenth-1]
		prefName = strings.Join(exp[0:lenth-1], ".")
	} else {
		prefName = basename
	}
	dstFile, err := os.OpenFile(path.Join(localPath, prefName+"-"+addr+"."+ext), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
	if err != nil {
		return
	}
	defer dstFile.Close()
	buf := make([]byte, 1024)
	for {
		n, _ := srcFile.Read(buf)
		if n < 1 {
			break
		}
		dstFile.Write(buf[0:n])
	}
	return
}
func (t *Transfer) put(sc *sftp.Client, localPath, remotePath string) (err error) {
	// remote path is dir
	if strings.HasSuffix(remotePath, "/") {
		basename := path.Base(localPath)
		remotePath = path.Join(remotePath, basename)
	}
	_, e := sc.Stat(remotePath)
	if e == nil {
		if !t.Override {
			fmt.Println("Remote file exists")
			return errors.New("Remote file exists")
		}
	}
	srcFile, err := os.OpenFile(localPath, os.O_RDONLY, 0755)
	if err != nil {
		return
	}
	defer srcFile.Close()
	dstFile, err := sc.OpenFile(remotePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC)
	if err != nil {
		return
	}
	defer dstFile.Close()
	buf := make([]byte, 1024)
	for {
		n, _ := srcFile.Read(buf)
		if n < 1 {
			break
		}
		dstFile.Write(buf[0:n])
	}
	return
}

func (t *Transfer) initClient() error {
	auth, err := GetAuth()
	if err != nil {
		log.Fatalln(err)
	}
	clientConfig := &ssh.ClientConfig{
		User:            C.Auth.User,
		Auth:            auth,
		Timeout:         30 * time.Second,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	for _, h := range t.Hosts {
		if strings.Index(h, ":") < 0 {
			h = h + ":" + strconv.Itoa(C.Server.DefaultPort)
		}
		client, err := ssh.Dial("tcp", h, clientConfig)
		if err != nil {
			return err
		}
		t.Clients[h] = client
		t.SftpClient[h], err = sftp.NewClient(client, sftp.MaxPacket(33788))
		if err != nil {
			return err
		}
	}
	return nil
}

// PrettyPrint print transfer result
func (t *Transfer) PrettyPrint() {}
