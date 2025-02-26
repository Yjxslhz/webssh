package core

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"golang.org/x/crypto/ssh"
	"log"
	"net"
	"strconv"
	"strings"
	"time"
)

// DecodedMsgToSSHClient 字符串信息解析为ssh客户端
func DecodedMsgToSSHClient(sshInfo string) (SSHClient, error) {
	client := NewSSHClient()
	decoded, err := base64.StdEncoding.DecodeString(sshInfo)
	if err != nil {
		return client, err
	}
	err = json.Unmarshal(decoded, &client)
	if err != nil {
		return client, err
	}
	if strings.Contains(client.IPAddress, ":") && string(client.IPAddress[0]) != "[" {
		client.IPAddress = "[" + client.IPAddress + "]"
	}
	return client, nil
}

// GenerateClient 创建ssh客户端
func (sclient *SSHClient) GenerateClient() error {
	var (
		auth         []ssh.AuthMethod
		addr         string
		clientConfig *ssh.ClientConfig
		client       *ssh.Client
		config       ssh.Config
		err          error
	)
	auth = make([]ssh.AuthMethod, 0)
	
	// 根据登录类型选择认证方式
	if sclient.LoginType == 0 {
		// 密码认证
		auth = append(auth, ssh.Password(sclient.Password))
	} else {
		// 私钥认证
		var signer ssh.Signer
		
		// 检查是否有密钥密码
		log.Printf("KeyPassphrase: '%s', LoginType: %d", sclient.KeyPassphrase, sclient.LoginType)
		
		// 尝试先用密码解析私钥
		if sclient.KeyPassphrase != "" {
			// 使用密码解析私钥
			log.Println("Attempting to parse private key with passphrase")
			signer, err = ssh.ParsePrivateKeyWithPassphrase([]byte(sclient.Password), []byte(sclient.KeyPassphrase))
			if err != nil {
				log.Printf("Failed to parse private key with passphrase: %v", err)
				return fmt.Errorf("failed to parse private key with passphrase: %v", err)
			}
		} else {
			// 尝试解析没有密码的私钥
			log.Println("Attempting to parse private key without passphrase")
			signer, err = ssh.ParsePrivateKey([]byte(sclient.Password))
			if err != nil {
				// 如果解析失败，检查是否是因为私钥有密码保护
				if strings.Contains(err.Error(), "private key is passphrase protected") || 
				   strings.Contains(err.Error(), "cannot decode encrypted private keys") {
					log.Println("Private key is passphrase protected, but no passphrase provided")
					return fmt.Errorf("private key is passphrase protected, please provide the passphrase")
				}
				log.Printf("Failed to parse private key: %v", err)
				return fmt.Errorf("failed to parse private key: %v", err)
			}
		}
		
		auth = append(auth, ssh.PublicKeys(signer))
	}
	
	config = ssh.Config{
		Ciphers: []string{"aes128-ctr", "aes192-ctr", "aes256-ctr", "aes128-gcm@openssh.com", "arcfour256", "arcfour128", "aes128-cbc", "3des-cbc", "aes192-cbc", "aes256-cbc"},
	}
	
	clientConfig = &ssh.ClientConfig{
		User:    sclient.Username,
		Auth:    auth,
		Timeout: 5 * time.Second,
		Config:  config,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}
	
	addr = fmt.Sprintf("%s:%d", sclient.IPAddress, sclient.Port)
	if client, err = ssh.Dial("tcp", addr, clientConfig); err != nil {
		return err
	}
	sclient.Client = client
	return nil
}

// InitTerminal 初始化终端
func (sclient *SSHClient) InitTerminal(ws *websocket.Conn, rows, cols int) *SSHClient {
	sshSession, err := sclient.Client.NewSession()
	if err != nil {
		log.Println(err)
		return nil
	}
	sclient.Session = sshSession
	sclient.StdinPipe, _ = sshSession.StdinPipe()
	wsOutput := new(wsOutput)
	//ssh.stdout and stderr will write output into comboWriter
	sshSession.Stdout = wsOutput
	sshSession.Stderr = wsOutput
	wsOutput.ws = ws
	modes := ssh.TerminalModes{
		ssh.ECHO:          1,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}

	if err := sshSession.RequestPty("xterm", rows, cols, modes); err != nil {
		return nil
	}
	if err := sshSession.Shell(); err != nil {
		return nil
	}
	return sclient
}

// Connect ws连接
func (sclient *SSHClient) Connect(ws *websocket.Conn, timeout time.Duration, closeTip string) {
	stopCh := make(chan struct{})
	//这里第一个协程获取用户的输入
	go func() {
		for {
			// p为用户输入
			_, p, err := ws.ReadMessage()
			if err != nil {
				close(stopCh)
				return
			}
			if string(p) == "ping" {
				continue
			}
			if strings.Contains(string(p), "resize") {
				resizeSlice := strings.Split(string(p), ":")
				rows, _ := strconv.Atoi(resizeSlice[1])
				cols, _ := strconv.Atoi(resizeSlice[2])
				err := sclient.Session.WindowChange(rows, cols)
				if err != nil {
					log.Println(err)
					close(stopCh)
					return
				}
				continue
			}
			_, err = sclient.StdinPipe.Write(p)
			if err != nil {
				close(stopCh)
				return
			}
		}
	}()

	defer func() {
		ws.Close()
		sclient.Close()

		if err := recover(); err != nil {
			log.Println(err)
		}
	}()
	// 设置ws超时时间timer
	stopTimer := time.NewTimer(timeout)
	defer stopTimer.Stop()
	// 主循环
	for {
		select {
		case <-stopCh:
			return
		case <-stopTimer.C:
			ws.WriteMessage(1, []byte(fmt.Sprintf("\u001B[33m%s\u001B[0m", closeTip)))
			return
		}
	}
}
