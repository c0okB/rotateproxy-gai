package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

var crawlDone = make(chan struct{})

type InfoMation struct {
	Error   bool       `json:"error"`
	Mode    string     `json:"mode"`
	Page    int        `json:"page"`
	Query   string     `json:"query"`
	Results [][]string `json:"results"`
	Size    int        `json:"size"`
}



func RunCrawler(fofaApiKey, fofaEmail, rule string, limit int) {
	var Info InfoMation
	//transCfg := &http.Transport{
	//	TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // ignore expired SSL certificates
	//}
	//
	//proxy := "http://127.0.0.1:8080"
	//if proxy != "" {
	//	proxyUrl, err := url.Parse(proxy)
	//	if err == nil { // 使用传入代理
	//		transCfg.Proxy = http.ProxyURL(proxyUrl)
	//	}
	//}
	//
	//httpClient := &http.Client{
	//	Transport: transCfg,
	//}

	httpClient := &http.Client{}

	url := fmt.Sprintf("https://fofa.so/api/v1/search/all?limit=%d&",limit)
	req, err1 := http.NewRequest("GET", url, nil)
	if err1 != nil {
		fmt.Printf("[-] %s\n",err1)
	}
	rule = base64.StdEncoding.EncodeToString([]byte(rule))
	q := req.URL.Query()
	q.Add("email", fofaEmail)
	q.Add("key", fofaApiKey)
	q.Add("qbase64", rule)
	req.URL.RawQuery = q.Encode()
	resp, err2 := httpClient.Do(req)
	if err2 != nil {
		fmt.Printf("[-] %s\n",err2)
	}
	fmt.Printf("start to parse proxy url from response\n")
	defer resp.Body.Close()
	body,_ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body,&Info)

	fmt.Printf("Get %d hosts\n", len(Info.Results))
	for _, value := range Info.Results {
		host := value[0]
		//fmt.Println(host)
		addProxyURL(fmt.Sprintf("socks5://%s", host))
	}
}

func StartRunCrawler(fofaApiKey, fofaEmail, rule string,CountLimit int) {

	RunCrawler(fofaApiKey,fofaEmail,rule,CountLimit)

}
