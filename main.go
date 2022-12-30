package main

import (
	"bruteforcer/Global"
	"bruteforcer/brute"
	"fmt"
	"github.com/urfave/cli/v2"
	"os"
)

func main() {
	port := "8080"
	protocol := "http"
	app := &cli.App{
		Name:      "protocaldetect",
		Usage:     "judg protocol\n protocol:\nssh\nmysql", // 这里写协议
		UsageText: "lazy to write...",
		Version:   "0.1.1",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "port", Aliases: []string{"p"}, Destination: &port, Value: "8080", Usage: "port", Required: true},
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
	_, _ = fmt.Scanln(&host)
	dofunc(host, port, outfile, protocol) // 执行
	goto start
}
func dofunc(host string, port string, file *os.File, protocol string) {

	if protocol == "ssh" { //ssh爆破
		brute.SSHBrute(host, port, file)
	} else {
		println("爆破机无此协议")
		os.Exit(0)
	}
	return
}
