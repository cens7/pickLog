package ot

import (
	"fmt"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"io"
	"log"
	"net"
	"os"
	"path"
	"time"
)

func connect(user, password, host string, port int) (*sftp.Client, error) {
	var (
		auth []ssh.AuthMethod
		addr string
		clientConfig *ssh.ClientConfig
		sshClient *ssh.Client
		sftpClient *sftp.Client
		err error
	)
	// get auth method
	auth = make([]ssh.AuthMethod, 0)
	auth = append(auth, ssh.Password(password))

	clientConfig = &ssh.ClientConfig{
		User: user,
		Auth: auth,
		Timeout: 30 * time.Second,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}

	// connet to ssh
	addr = fmt.Sprintf("%s:%d", host, port)

	if sshClient, err = ssh.Dial("tcp", addr, clientConfig); err != nil {
		return nil, err
	}

	// create sftp client
	if sftpClient, err = sftp.NewClient(sshClient); err != nil {
		return nil, err
	}

	return sftpClient, nil
}
func Run(se *ServeInfo) {

	if se.username == "" {
		fmt.Println("获取数据库服务器信息错误")
		return
	}



	var (
		err error
		sftpClient *sftp.Client
	)

	sftpClient, err = connect(se.username, se.passoword, se.hostIp, 22)
	if err != nil {
		log.Fatal(err)
	}
	defer sftpClient.Close()

	fmt.Println("\n---------->连接远程服务器成功....")

	localDir, e := os.Getwd()

	//localDir = localDir+"/"

	if e != nil {
		panic(e)
	}

	var remoteFilePath = "/app/app-log/bss-order-shopcart-dubbo/"+se.appName+".log"

	srcFile, err := sftpClient.Open(remoteFilePath)
	if err != nil {
		log.Fatal(err)
	}

	//fmt.Println("当前文件名：",srcFile.Name())

	defer srcFile.Close()

	var localFileName = path.Base(remoteFilePath)
	dstFile, err := os.Create(path.Join(localDir, localFileName))
	if err != nil {
		log.Fatal(err)
	}
	defer dstFile.Close()

	//if _, err = srcFile.WriteTo(dstFile); err != nil {
	//	log.Fatal(err)
	//}
	buf := make([]byte, 1024)
	for {
		// read a chunk
		n, err := srcFile.Read(buf)
		if err != nil && err != io.EOF {
			fmt.Println("下载日志失败")
		}
		if n == 0 {
			break
		}

		// write a chunk
		if _, err := dstFile.Write(buf[:n]); err != nil {
			fmt.Println("下载日志失败")
		}
	}

	//todo 处理日志
	fmt.Println("\n---------->复制服务器日志文件完成!")
	fmt.Println("\n---------->文件所在目录：",localDir,"， 文件名：",localFileName)


}