package main

import (
	"net"
	"time"
	"sync"
	"encoding/hex"
	"flag"
	"strings"
	"regexp"
)

const ADDR  = "openbarrage.douyutv.com:8601"

func main() {
	roomid := flag.String("room","9999","roomid")
	flag.Parse()
	println("start connecting server .....")
	conn,err := net.Dial("tcp",ADDR)
	if err != nil {
		println("connect to server failed!")
		return
	}
	println("connect to server success!")
	println("start into room.....")
	conn.Write(buildRequest("type@=loginreq/roomid@=" + *roomid + "/"))
	conn.Write(buildRequest("type@=joingroup/rid@=" + *roomid + "/gid@=-9999/"))

	go reciveMessage(conn)
	go heartBeat(conn)

	var wg sync.WaitGroup
	wg.Add(1)
	wg.Wait()
}

func reciveMessage(conn net.Conn) {
	danmuReg  := regexp.MustCompile("type@=chatmsg/.*rid@=(\\d*?)/.*uid@=(\\d*).*nn@=(.*?)/txt@=(.*?)/(.*)/")
	for true {
		buf := make([]byte, 1024 * 80)
		n,_ := conn.Read(buf)
		data := hex.EncodeToString(buf[:n])
		messages := strings.Split(data,"b2020000")
		for _, m := range messages {
			if strings.Contains(m,"00"){
				end := strings.Index(m,"00")
				m = string([]rune(m)[:end])
				mb, _ := hex.DecodeString(m)
				dm := string(mb)
				match := danmuReg.FindStringSubmatch(dm)
				if len(match) >0 {
					println(getFormatTime() +"\t" +match[3]+"\t:" +match[4])
				}
			}
		}
	}
}


func getFormatTime() string{
	return time.Now().Format("2006-01-02 15:04:05")
}
func heartBeat(conn net.Conn) {
	for true {
		conn.Write(buildRequest("type@=mrkl/"))
		time.Sleep(1e10)
	}
}

func buildRequest(str string) []byte{
	var data []byte
	length := byte(len(str)+9)
	data = append(data,length,0,0,0)//length
	data = append(data,length,0,0,0)//code
	data = append(data, 0xb1, 0x02, 0x00, 0x00 )//magic
	data = append(data, []byte(str)...)//content
	data = append(data, 0)//end
	return data
}