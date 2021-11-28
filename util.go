package main

import "fmt"

func IsProxyURLBlank() bool {
	proxies, err := QueryAvailProxyURL()
	if err != nil {
		fmt.Printf("[!] Error: %v\n", err)
		return false
	}
	return len(proxies) == 0
}

