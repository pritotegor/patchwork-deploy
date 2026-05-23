package ssh

import (
	"bytes"
	"fmt"
	"os"
	"time"

	"golang.org/x/crypto/ssh"
)

// Config holds SSH connection parameters.
type Config struct {
	Host       string
	Port       int
	User       string
	KeyPath    string
	Timeout    time.Duration
}

// Client wraps an SSH client connection.
type Client struct {
	conn *ssh.Client
}

// Connect establishes an SSH connection using the provided config.
func Connect(cfg Config) (*Client, error) {
	key, err := os.ReadFile(cfg.KeyPath)
	if err != nil {
		return nil, fmt.Errorf("reading private key %q: %w", cfg.KeyPath, err)
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return nil, fmt.Errorf("parsing private key: %w", err)
	}

	timeout := cfg.Timeout
	if timeout == 0 {
		timeout = 30 * time.Second
	}

	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	conn, err := ssh.Dial("tcp", addr, &ssh.ClientConfig{
		User:            cfg.User,
		Auth:            []ssh.AuthMethod{ssh.PublicKeys(signer)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // TODO: replace with known_hosts
		Timeout:         timeout,
	})
	if err != nil {
		return nil, fmt.Errorf("dialing %s: %w", addr, err)
	}

	return &Client{conn: conn}, nil
}

// RunCommand executes a shell command on the remote host and returns combined output.
func (c *Client) RunCommand(cmd string) (string, error) {
	session, err := c.conn.NewSession()
	if err != nil {
		return "", fmt.Errorf("creating session: %w", err)
	}
	defer session.Close()

	out, err := session.CombinedOutput(cmd)
	if err != nil {
		return string(out), fmt.Errorf("running command: %w", err)
	}
	return string(out), nil
}

// RunCommandSeparate executes a shell command on the remote host and returns
// stdout and stderr as separate strings.
func (c *Client) RunCommandSeparate(cmd string) (stdout, stderr string, err error) {
	session, err := c.conn.NewSession()
	if err != nil {
		return "", "", fmt.Errorf("creating session: %w", err)
	}
	defer session.Close()

	var outBuf, errBuf bytes.Buffer
	session.Stdout = &outBuf
	session.Stderr = &errBuf

	runErr := session.Run(cmd)
	if runErr != nil {
		return outBuf.String(), errBuf.String(), fmt.Errorf("running command: %w", runErr)
	}
	return outBuf.String(), errBuf.String(), nil
}

// Close closes the underlying SSH connection.
func (c *Client) Close() error {
	return c.conn.Close()
}
