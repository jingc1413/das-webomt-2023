package sftp

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"path"
	"time"

	"github.com/pkg/errors"
	"github.com/pkg/sftp"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
)

func connect(username, password, host string, port int) (*sftp.Client, error) {
	var (
		auth         []ssh.AuthMethod
		addr         string
		clientConfig *ssh.ClientConfig
		sshClient    *ssh.Client
		sftpClient   *sftp.Client
		err          error
	)
	auth = make([]ssh.AuthMethod, 0)
	auth = append(auth, ssh.Password(password))
	clientConfig = &ssh.ClientConfig{
		User:    username,
		Auth:    auth,
		Timeout: 30 * time.Second,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}

	addr = fmt.Sprintf("%s:%d", host, port)
	if sshClient, err = ssh.Dial("tcp", addr, clientConfig); err != nil {
		return nil, err
	}
	if sftpClient, err = sftp.NewClient(sshClient); err != nil {
		return nil, err
	}

	return sftpClient, nil
}

func DownloadNone(addr string, port uint16, username string, password string, dir string, filename string) error {
	buf := new(bytes.Buffer)
	err := Download(addr, port, username, password, dir, filename, buf)
	if err != nil {
		return err
	}
	return nil
}

func Download(addr string, port uint16, username string, password string, dir string, filename string, w io.Writer) error {
	cli, err := connect(username, password, addr, int(port))
	if err != nil {
		return errors.Wrap(err, "ssh connect")
	}
	defer cli.Close()

	remotePath := path.Join(dir, filename)
	srcFile, err := cli.Open(remotePath)
	if err != nil {
		return errors.Wrap(err, "open remote file")
	}
	defer srcFile.Close()

	n, err := srcFile.WriteTo(w)
	if err != nil {
		return errors.Wrap(err, "write to buffer")
	}
	log.Tracef("download file, path=%s, filename=%s, size=%d", dir, filename, n)
	return nil
}

func Upload(addr string, port uint16, username string, password string, dir string, filename string, r io.Reader) error {
	cli, err := connect(username, password, addr, int(port))
	if err != nil {
		return errors.Wrap(err, "ssh connect")
	}
	defer cli.Close()

	remotePath := path.Join(dir, filename)
	dstFile, err := cli.Create(remotePath)
	if err != nil {
		return errors.Wrap(err, "create remote file")
	}
	defer dstFile.Close()

	n, err := dstFile.ReadFrom(r)
	if err != nil {
		return errors.Wrap(err, "write to buffer")
	}
	log.Tracef("upload file, path=%s, filename=%s, size=%d", dir, filename, n)
	return nil
}
