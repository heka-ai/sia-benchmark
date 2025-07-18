package ssh

import (
	"net"

	log "github.com/heka-ai/benchmark-cli/internal/logs"
	"github.com/melbahja/goph"
	"golang.org/x/crypto/ssh"

	"crypto/rand"
	"crypto/rsa"
	"encoding/pem"
	"fmt"
)

var logger = log.GetLogger("ssh")

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

func (c *SSHClient) Run(command string) error {
	auth, err := goph.Key(c.KeyPath, "")

	if err != nil {
		return err
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
		return err
	}

	// Defer closing the network connection.
	defer client.Close()

	// Execute your command.
	out, err := client.Run(command)

	logger.Debug().Str("command", command).Msgf("Executed command %s", out)

	if err != nil {
		return err
	}

	return nil
}

// GenerateSSHKey generates an SSH key pair and saves them to files
func GenerateSSHKey(bits int) (string, string, error) {
	// Generate RSA private key
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate private key: %v", err)
	}

	pemBlock, err := ssh.MarshalPrivateKey(privateKey, "")
	if err != nil {
		return "", "", fmt.Errorf("failed to generate private key PEM Block: %v", err)
	}

	// Encode private key to PEM format
	privateKeyPEM := pem.EncodeToMemory(pemBlock)

	// Generate public key from private key
	publicKey, err := ssh.NewPublicKey(&privateKey.PublicKey)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate public key: %v", err)
	}

	// Get OpenSSH format for the public key
	publicKeyBytes := ssh.MarshalAuthorizedKey(publicKey)

	return string(publicKeyBytes), string(privateKeyPEM), nil
}
