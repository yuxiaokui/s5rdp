package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"sync"
  
  "golang.org/x/net/proxy"
)

var lock sync.Mutex

var socks5 string
var target string
var lport string

var help = func() {
	fmt.Println("Usage for s5rdp tools, coded by xi4okv")
	fmt.Println("==========================================================")
	fmt.Println("S5rdp Socks5Ip:Socks5Port RdpIp:RdpPort ListenPort")
	fmt.Println("==========================================================")
}

func main() {
	args := os.Args

	if len(args) != 4 || args == nil {
		help()
		os.Exit(0)
	}else
	{
		socks5 = args[1]
		target = args[2]
		lport = args[3]
		server(args[2])
	}
}

func server(target string) {
	fmt.Printf("Listening %s\n", lport)
	lis, err := net.Listen("tcp", "0.0.0.0:" + lport)
	if err != nil {
		return
	}
	defer lis.Close()

	for {
		conn, err := lis.Accept()
		if err != nil {
			continue
		}
		go handle(conn, target)
	}
}

func handle(sconn net.Conn, target string) {
	defer sconn.Close()
	ip := target
	dialer, err := proxy.SOCKS5("tcp", socks5, nil, proxy.Direct)
	dconn, err := dialer.Dial("tcp", ip)
	if err != nil {
		fmt.Printf("连接%v失败:%v\n", ip, err)
		return
	}
	ExitChan := make(chan bool, 1)
	go func(sconn net.Conn, dconn net.Conn, Exit chan bool) {
		io.Copy(dconn, sconn)
		ExitChan <- true
	}(sconn, dconn, ExitChan)

	go func(sconn net.Conn, dconn net.Conn, Exit chan bool) {
		io.Copy(sconn, dconn)
		ExitChan <- true
	}(sconn, dconn, ExitChan)
	<-ExitChan
	dconn.Close()
}
