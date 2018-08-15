package main

import "fmt"
import "net"
import "log"
import "time"
import "strings"
import (
	"strconv"
)

func main() {
	fmt.Println("vim-go")
	newSocket()
	time.Sleep(time.Second * 1000)
}

func newSocket() {
	//serverAddr := fmt.Sprintf("%s:%d", masterHost, masterPort)
	tcpAddr, _ := net.ResolveTCPAddr("tcp", "127.0.0.1:1008")
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		log.Fatalf("Failed to connect to the Locust master: %s %s", tcpAddr, err)
	}
	conn.SetNoDelay(true)
	go func() {
		for z := 0; z < 10; z++ {
			x := "GET / HTTP/1.1\r\nHost: 127.0.0.1:1008\r\n\r\n"
			fmt.Println(">>>>>>>> FanPrint[0].newSocket", conn)
			conn.Write([]byte(x))
			read(conn)
		}
		time.Sleep(time.Second * 1000)
	}()
}

func read(c *net.TCPConn) {
	rH := false
	rR := false
	b := make([]byte, 1024)
	h := map[string]string{}
	buffer := make([]byte, 0)
	clength := 0
	for {
		count, err := c.Read(b)
		if err != nil {
			fmt.Println(">>>>>>>> FanPrint[0].read", err)
			break
		}
		buffer = append(buffer, b[:count]...)
		if !rH {
			first := 0
			for j := 0; j < len(buffer); j++ {
				x := buffer[j]
				if x != byte('\r') {
					continue
				}
				if j == len(buffer) {
					break
				}
				if buffer[j+1] == byte('\n') {
					head := string(buffer[first:j])
					first = j + 2
					if !rR {
						rR = true
						fmt.Println(">>>>>>>> FanPrint[3].read", head)
					} else if head == "" {
						rH = true
						break
					} else {
						shead := string(head)
						shead = strings.Replace(shead, " ", "", -1)
						shead = strings.ToLower(shead)
						splith := strings.Split(shead, ":")
						h[splith[0]] = splith[1]
					}
				}
			}
			if first <= len(buffer) {
				buffer = buffer[first:]
			}
		}
		if clength == 0 && rH {
			clength, err = strconv.Atoi(h["content-length"])
		}
		if rH && len(buffer) == clength {
			break
		}
	}
}
