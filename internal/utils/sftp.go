package utils

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/knownhosts"
)

type SftpClient struct {
}

func (s *SftpClient) UseWithPwd() {
	host := "sftp.example.com"
	port := 22
	user := "your_username"
	password := "your_password"

	// Create SSH client configuration
	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	conn, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", host, port), config)
	if err != nil {
		fmt.Println("Failed to connect to SSH server:", err)
		return
	}
	defer conn.Close()

	// Open SFTP session
	sftpClient, err := sftp.NewClient(conn)
	if err != nil {
		fmt.Println("Failed to open SFTP session:", err)
		return
	}
	defer sftpClient.Close()

	remoteFile, err := sftpClient.Open("/path/to/remote/file.txt")
	if err != nil {
		fmt.Println("Failed to open remote file:", err)
		return
	}
	defer remoteFile.Close()

	localFile, err := os.Create("/path/to/local/file.txt")
	if err != nil {
		fmt.Println("Failed to create local file:", err)
		return
	}
	defer localFile.Close()

	_, err = io.Copy(localFile, remoteFile)
	if err != nil {
		fmt.Println("Failed to download file:", err)
		return
	}
}

func (s *SftpClient) UseWithCert() {
	host := "example.com:22"
	user := "myuser"
	keyPath := "/home/me/.ssh/id_ed25519"
	knownHostsPath := "/home/me/.ssh/known_hosts"

	// ----- Load the private key file -----
	key, err := os.ReadFile(keyPath)
	if err != nil {
		log.Fatal("Unable to read private key:", err)
	}

	// ----- Parse the key (supports passphrase or no passphrase) -----
	signer, err := ssh.ParsePrivateKey(key)
	// If your key is encrypted:
	// signer, err := ssh.ParsePrivateKeyWithPassphrase(key, []byte("myPassphrase"))
	if err != nil {
		log.Fatal("Unable to parse private key:", err)
	}

	// ----- Known-hosts verification (recommended) -----
	hostKeyCallback, err := knownhosts.New(knownHostsPath)
	if err != nil {
		log.Fatal("Could not load known_hosts:", err)
	}

	// ----- SSH client configuration -----
	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: hostKeyCallback,
		// If using old servers that don't support modern KEX:
		// HostKeyAlgorithms: []string{ssh.KeyAlgoRSA, ssh.KeyAlgoED25519},
	}

	// ----- Establish the SSH connection -----
	conn, err := ssh.Dial("tcp", host, config)
	if err != nil {
		log.Fatal("Failed to dial:", err)
	}
	defer conn.Close()

	// ----- Open the SFTP session -----
	client, err := sftp.NewClient(conn)
	if err != nil {
		log.Fatal("Failed to create SFTP client:", err)
	}
	defer client.Close()

	// ----- Test: list remote directory -----
	files, err := client.ReadDir(".")
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range files {
		fmt.Println(f.Name())
	}
}
