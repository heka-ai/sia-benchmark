package ssh

import (
	"net"

	"github.com/melbahja/goph"
	"golang.org/x/crypto/ssh"
)

type SSHClient struct {
	Host    string
	User    string
	KeyPath string
}

func NewSSHClient(keyPath string, host string, user string) *SSHClient {
	return &SSHClient{
		Host:    host,
		User:    user,
		KeyPath: keyPath,
	}
}

func (c *SSHClient) Run(command string) (string, error) {
	auth, err := goph.Key(c.KeyPath, "")

	if err != nil {
		return "", err
	}

	client, err := goph.NewConn(&goph.Config{
		User: c.User,
		Auth: auth,
		Addr: c.Host,
		Port: 22,
		Callback: func(host string, remote net.Addr, key ssh.PublicKey) error {
			return goph.AddKnownHost(host, remote, key, "")
		},
	})

	if err != nil {
		return "", err
	}

	// Defer closing the network connection.
	defer client.Close()

	// Execute your command.
	out, err := client.Run(command)

	if err != nil {
		return "", err
	}

	return string(out), nil
}
