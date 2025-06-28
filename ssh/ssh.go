package ssh

import (
	"context"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/ssh"
	"log"
	"net"
	"time"
)

type Dialer struct {
	client *ssh.Client
}

type Ssh struct {
	Host     string `json:"host"`
	User     string `json:"user"`
	Port     int    `json:"port"`
	Password string `json:"password"`
}

func (v *Dialer) Dial(address string) (net.Conn, error) {
	return v.client.Dial("tcp", address)
}
func (s *Ssh) dialWithPassword() (*ssh.Client, error) {
	address := fmt.Sprintf("%s:%d", s.Host, s.Port)
	config := &ssh.ClientConfig{
		User: s.User,
		Auth: []ssh.AuthMethod{
			ssh.Password(s.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         time.Second * 5, // 5秒超时
	}
	return ssh.Dial("tcp", address, config)
}

func RegisterSsh(client *Ssh) {
	dial, err := client.dialWithPassword()
	if err != nil {
		//log.Fatalf("ssh connect error: %s", err.Error())
		log.Printf("ssh connect error: %s\n", err.Error())
		return
	}
	//defer dial.Close()
	// 注册ssh代理
	mysql.RegisterDialContext("mysql+tcp", func(ctx context.Context, addr string) (net.Conn, error) {
		return (&Dialer{client: dial}).Dial(addr)
	})
	return
}
