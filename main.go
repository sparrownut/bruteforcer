package main

import (
	"bruteforcer/Global"
	"bruteforcer/brute"
	"bruteforcer/exp"
	"fmt"
	"github.com/urfave/cli/v2"
	"os"
)

func main() {
	port := "8080"
	protocol := "http"
	app := &cli.App{
		Name:      "bruteforcer",
		Usage:     "bruteforcer 虽然这么叫 但是是一个exp集合\nssh\nyonyou\nredis \n仅供授权的渗透测试使用 请遵守法律!", // 这里写协议
		UsageText: "lazy to write...",
		Version:   "0.4.5",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "port", Aliases: []string{"p"}, Destination: &port, Value: "22", Usage: "port", Required: true},
			&cli.StringFlag{Name: "protocol", Aliases: []string{"P"}, Destination: &protocol, Value: "ssh", Usage: "protocol", Required: true},
			&cli.BoolFlag{Name: "DBG", Aliases: []string{"D"}, Destination: &Global.DBG, Value: false, Usage: "DBG MOD", Required: false},
			&cli.StringFlag{Name: "cmd", Aliases: []string{"C"}, Destination: &Global.CMD, Value: "whoami", Usage: "shell command", Required: false},
			&cli.StringFlag{Name: "pwd", Aliases: []string{"W"}, Destination: &Global.PWD, Value: "", Usage: "pwd", Required: false},
			&cli.StringFlag{Name: "user", Aliases: []string{"U"}, Destination: &Global.USR, Value: "", Usage: "user", Required: false},
		},
		HideHelpCommand: true,
		Action: func(c *cli.Context) error {
			err := do(port, protocol)
			if err != nil {

			}
			return nil
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		//panic(err)
	}

	//fmt.Printf(os.Args[1])
}
func do(port string, protocol string) error {
	outfile, _ := os.OpenFile("bruteoutput.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	defer func(file *os.File) {
		_ = file.Close()
	}(outfile)
start: // 在这里循环
	host := ""

	//if !terminal.IsTerminal(0) {
	//	// 如果是管道模式
	//} else {
	_, _ = fmt.Scanln(&host) // 管道或输入
	//}
	dofunc(host, port, outfile, protocol) // 执行
	goto start
}
func dofunc(host string, port string, file *os.File, protocol string) {

	if protocol == "ssh" { //ssh爆破
		brute.SSHBrute(host, port, file)
	} else if protocol == "yonyou" {
		exp.YonYouNCEXP(host, Global.CMD)
	} else if protocol == "redis" {
		exp.REDISUnauthorizedEXP(host, port, Global.CMD)
	} else {
		println("exp无此协议")
		os.Exit(0)
	}
	return
}
