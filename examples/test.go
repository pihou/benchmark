package main

import "fmt"
import "net"
import "log"

//import "time"
import "strings"
import (
	"strconv"
)

var sPool = make(chan *net.TCPConn, 256)

//func main() {
//	for i := 0; i < 3; i++ {
//		c := newSocket()
//		sPool <- c
//	}
//	for i := 0; i < 10; i++ {
//		go func(index int) {
//			for {
//				c := <-sPool
//				succ := read(c, index)
//				if succ {
//					sPool <- c
//				} else {
//					sPool <- newSocket()
//				}
//			}
//		}(i)
//	}
//	time.Sleep(time.Second * 1000)
//}

func Init() {
	for i := 0; i < 256; i++ {
		c := newSocket()
		sPool <- c
	}
}

func newSocket() *net.TCPConn {
	tcpAddr, _ := net.ResolveTCPAddr("tcp", "k8s-node1.shoupihou.site:1008")
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		log.Fatalf("Failed to connect to the Locust master: %s %s", tcpAddr, err)
	}
	conn.SetNoDelay(true)
	return conn
}

func read(c *net.TCPConn, index int) bool {
	x := "GET /app/benchmark/ HTTP/1.1\r\nHost: k8s-node1.shoupihou.site:1008\r\n\r\n"
	length := len(x)
	for {
		count, err := c.Write([]byte(x))
		if err != nil {
			fmt.Println(">>>>>>>> FanPrint[2].write", err)
			return false
		}
		x = x[count:]
		length -= count
		if length == 0 {
			break
		}
	}

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
			return false
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
						fmt.Println(">>>>>>>> FanPrint[3].read", index, head)
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
	return true
}
