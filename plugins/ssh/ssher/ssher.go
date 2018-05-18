package ssher

import (
	"net"

	"golang.org/x/crypto/ssh"
)

// SSHer main struct of package, stored information about connection and client variable
type SSHer struct {
	Host string
	Port string
	User string
	Pass string

	client *ssh.Client
}

// New function return SSHer pointer for next use
func New(host, port, user, pass string) *SSHer {
	return &SSHer{
		Host: host,
		Port: port,
		User: user,
		Pass: pass,
	}
}

// Connect function connect to server
func (s *SSHer) Connect() (err error) {
	addr := net.JoinHostPort(s.Host, s.Port)
	config := &ssh.ClientConfig{
		User: s.User,
		Auth: []ssh.AuthMethod{
			ssh.Password(s.Pass),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	if s.client, err = ssh.Dial("tcp", addr, config); err != nil {
		return
	}

	return
}

// Close function close connection with server
func (s *SSHer) Close() error {
	return s.client.Close()
}

// Run function run one command over connection
func (s *SSHer) Run(cmd string) (output []byte, err error) {
	var session *ssh.Session
	if session, err = s.client.NewSession(); err != nil {
		return
	}
	defer session.Close()

	if output, err = session.CombinedOutput(cmd); err != nil {
		return
	}
	return
}
