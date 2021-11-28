package main

import (
	"flag"
	"fmt"
)

var baseCfg BaseConfig


func main() {
	fofaEmail := flag.String("email","","fofa email")
	fofaApi := flag.String("token","","fofa api token")
	CountLimit := flag.Int("limit",3500,"the Count of you want")
	rule := `protocol=="socks5" && "Version:5 Method:No Authentication(0x00)" && after="2021-01-01" && country="CN"`
	a := flag.String("lport","8899","local listen address")
	Type := flag.String("type","","type : 1 | 2")
	fmt.Println(`[*] First : when you want to Catch the socks5 url from fofa`)
	fmt.Println(`[*]   Please input: ./rotateproxy -type 1  -email xxxxxxxxx@gmail.com -token xxxxxxxxxxxx (-rule) (-limit)`)
	fmt.Println(`[*] Second	: when you want to Use the socks5 url from fofa`)
	fmt.Println(`[*]   Please input: ./rotateproxy -type 2  (-lport)`)
	fmt.Println("[*] By default lport is 8899 and limit is 3500. enjoy ~ ~ ~")
	flag.Parse()
	if *Type == "1" {
		StartRunCrawler(*fofaApi, *fofaEmail, rule, *CountLimit)
	}else if *Type == "2" {
		baseCfg.ListenAddr = "127.0.0.1:"+*a
		StartCheckProxyAlive()
		c := NewRedirectClient(WithConfig(&baseCfg))
		c.Serve()
		select {}
	}else {
		fmt.Println("[-] Check your type")
	}

}
