package utility

import (
	"bytes"
	"fmt"
	"os"

	"golang.org/x/crypto/ssh"
)

type SSHClient struct {
	client *ssh.Client
}

func (client *SSHClient) Connect(ip string, port string, userName string, password string) error {
	authPassword := []ssh.AuthMethod{ssh.Password(password)}
	conf := ssh.ClientConfig{User: userName, Auth: authPassword,
		HostKeyCallback: ssh.InsecureIgnoreHostKey()}
	var err error

	client.client, err = ssh.Dial("tcp", ip+":"+port, &conf)
	return err
}

func (client *SSHClient) Disconnect() {
	client.client.Close()
}

func (client *SSHClient) Run(command string) string {
	var outbuf bytes.Buffer
	if session, err := client.client.NewSession(); err == nil {
		session.Stdout = &outbuf
		session.Stderr = os.Stderr
		err := session.Run(command)
		if err != nil {
			fmt.Println(err)
			return ""
		}
		session.Close()
	}
	return outbuf.String()
}
