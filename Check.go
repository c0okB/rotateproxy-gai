package main

import "fmt"



func checkAlive() {
	proxies, err := QueryProxyURL()
	if err != nil {
		fmt.Printf("[!] query db error: %v\n", err)
	}
	for i := range proxies {
		proxy := proxies[i]
		//测试代理是否可用
		//按需求修改 加入gorouting
		//fmt.Printf("No %d ==> %s\n",i,proxy)
		go func() {
			respBody, avail := CheckProxyAlive(proxy.URL)
			//fmt.Printf("No %d ==> %s\n",i,proxy)
			if avail {
				fmt.Printf("%v 可用\n", proxy.URL)
				SetProxyURLAvail(proxy.URL, CanBypassGFW(respBody))
			} else {
				AddProxyURLRetry(proxy.URL)
			}
		}()
	}
}



func StartCheckProxyAlive() {
	checkAlive()
}