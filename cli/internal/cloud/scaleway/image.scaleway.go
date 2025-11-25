package scaleway

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/heka-ai/benchmark-cli/internal/ssh"
	iam "github.com/scaleway/scaleway-sdk-go/api/iam/v1alpha1"
	"github.com/scaleway/scaleway-sdk-go/api/instance/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

func (c *ScalewayClient) CreateBenchInstance() error {
	logger.Info().Msg("Starting the creation of the CPU instance image")

	err := c.CreateTemplateInstance("cpu")

	return err
}

func (c *ScalewayClient) CreateLLMInstance() error {
	logger.Info().Msg("Starting the creation of the LLM instance image")

	err := c.CreateTemplateInstance("llm")

	return err
}

func (c *ScalewayClient) CreateTemplateInstance(installType string) error {
	keyPairName := fmt.Sprintf("benchmark-key-pair-%s", uuid.New().String())
	logger.Debug().Str("keyName", keyPairName).Msg("Generating a key pair")

	PublicKey, PrivateKey, err := ssh.GenerateSSHKey(2048)
	if err != nil {
		return err
	}

	keyPair, err := c.iam.CreateSSHKey(&iam.CreateSSHKeyRequest{
		Name:      keyPairName,
		PublicKey: PublicKey,
	})
	if err != nil {
		return err
	}

	defer func() {
		logger.Debug().Str("keyName", keyPairName).Msg("Deleting the key pair")
		c.iam.DeleteSSHKey(&iam.DeleteSSHKeyRequest{
			SSHKeyID: keyPair.ID,
		})
		logger.Debug().Str("keyName", keyPairName).Msg("Key pair deleted")
	}()

	// write the key to a temp file
	keyFile, err := os.CreateTemp("/tmp", "tmp-key-ssh-benchmark-*.pem")
	if err != nil {
		return err
	}

	defer func() {
		logger.Debug().Str("keyName", keyPairName).Msg("Deleting the key pair file")
		os.Remove(keyFile.Name())
	}()

	// write the key to the file
	_, err = keyFile.WriteString(PrivateKey)
	if err != nil {
		return err
	}

	logger.Debug().Str("keyName", keyPairName).Msg("Key pair generated")
	logger.Debug().Str("CommercialType ", "DEV1-S").Msg("Creating the template server")

	server, err := c.createServer(
		c.config.ScalewayConfig.BaseCommercialType,
		scw.StringPtr(c.config.ScalewayConfig.BaseImage),
		scw.StringPtr(keyPair.ID),
	)

	defer func() error {
		err = c.deleteServer(*server)
		if err != nil {
			return err
		}
		return nil
	}()

	if err != nil {
		return err
	}

	serverIP := server.PublicIPs[0]
	logger.Debug().Str("publicIP", string(serverIP.Address)).Str("keyFile", keyFile.Name()).Msgf("Creating the SSH client")

	sshClient := ssh.NewSSHClient(keyFile.Name(), string(serverIP.Address), "ubuntu")

	err = c.runSetup(sshClient, installType)
	if err != nil {
		return err
	}

	logger.Info().Msg("Setup complete, creating the Image")

	amiName := fmt.Sprintf("benchmark-%s-%s", installType, time.Now().Format("2006-01-02-15-04-05"))

	var firstVolume *instance.VolumeServer
	for _, volume := range server.Volumes {
		firstVolume = volume
		break
	}

	// Now `firstVolume` contains the first volume found in the map.
	if firstVolume == nil {
		return errors.New("no volumes found for the server")
	}

	createSnapshotResponse, err := c.instance.CreateSnapshot(&instance.CreateSnapshotRequest{
		VolumeID: scw.StringPtr(firstVolume.ID),
	})
	if err != nil {
		return err
	}

	createImageResponse, err := c.instance.CreateImage(&instance.CreateImageRequest{
		Name:       amiName,
		RootVolume: createSnapshotResponse.Snapshot.ID,
		Public:     scw.BoolPtr(false),
	})
	if err != nil {
		return err
	}

	logger.Info().Str("imageID", createImageResponse.Image.ID).Str("image", createImageResponse.Image.Name).Msg("Image created")

	return nil
}

func (c *ScalewayClient) runSetup(sshClient *ssh.SSHClient, installType string) error {
	// clone the benchmark-cli repo
	logger.Info().Msg("Installing the benchmark-cli repo")
	err := sshClient.Run("git clone https://github.com/heka-ai/sia-benchmark.git /home/ubuntu/benchmark-cli")
	if err != nil {
		return err
	}

	err = sshClient.Run("mv /home/ubuntu/benchmark-cli/instance-builder/aws/ec2 /home/ubuntu/")
	if err != nil {
		return err
	}

	sshClient.Run("ls /home/ubuntu")

	logger.Info().Msg("Installing the control API")
	err = sshClient.Run("sudo bash /home/ubuntu/ec2/api/install.sh")
	if err != nil {
		return err
	}

	if installType == "llm" {
		logger.Info().Msg("Installing the LLM dependencies")
		err = sshClient.Run("sudo bash /home/ubuntu/ec2/gpu/install.sh")
	} else {
		logger.Info().Msg("Installing the CPU dependencies")
		err = sshClient.Run("bash /home/ubuntu/ec2/cpu/install.sh")
	}

	return err
}
