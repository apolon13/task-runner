package ssh

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"log"
	"net"
	"os"
	"path"
	"time"
)

type Params struct {
	Username           string
	Host               string
	Port               int
	PrivateKey string
	Password           string
}

type Client struct {
	Params     Params
	Connection *ssh.Client
}

func (p Params) getPrivateKey() string {
	if p.PrivateKey != "" {
		return p.PrivateKey
	}
	home := os.Getenv("HOME")
	if len(home) > 0 {
		return path.Join(home, ".ssh/id_rsa")
	}
	return ""
}

func (p Params) readKeyFile() ([]byte, error) {
	key, err := ioutil.ReadFile(p.getPrivateKey())
	if err != nil {
		return nil, err
	}
	return key, nil
}

func (client Client) authMethod() ([]ssh.AuthMethod, error) {
	if client.Params.Password != "" {
		return []ssh.AuthMethod{ssh.Password(client.Params.Password)}, nil
	}
	keyFile, err := client.Params.readKeyFile()
	if err != nil {
		return nil, err
	}
	signer, err := ssh.ParsePrivateKey(keyFile)
	if err != nil {
		return nil, err
	}
	return []ssh.AuthMethod{ssh.PublicKeys(signer)}, nil
}

func (client *Client) Connect() {
	authMethod, err := client.authMethod()
	if err != nil {
		log.Fatal(fmt.Errorf("auth method parsing error: %s", err))
	}
	config := &ssh.ClientConfig{
		User: client.Params.Username,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
		Auth:    authMethod,
		Timeout: 60 * time.Second,
	}
	addr := fmt.Sprintf("%s:%d", client.Params.Host, client.Params.Port)
	client.Connection, err = ssh.Dial("tcp", addr, config)
	if err != nil {
		log.Fatal(fmt.Errorf("connection error: %s", err))
	}
}
