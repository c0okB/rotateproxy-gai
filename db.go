package main

import (
	"crypto/tls"
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var DB *gorm.DB

type ProxyURL struct {
	gorm.Model
	URL          string `gorm:"uniqueIndex;column:url"`
	Retry        int    `gorm:"column:retry"`
	Available    bool   `gorm:"column:available"`
	CanBypassGFW bool   `gorm:"column:can_bypass_gfw"`
}

func (ProxyURL) TableName() string {
	return "proxy_urls"
}


func init() {
	var err error
	DB, err = gorm.Open(sqlite.Open("db.db"), &gorm.Config{
		Logger: logger.Discard,
	})
	checkErr(err)
	DB.AutoMigrate(&ProxyURL{})
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func CreateProxyURL(url string) error {
	tx := DB.Create(&ProxyURL{
		URL:       url,
		Retry:     0,
		Available: false,
	})
	return tx.Error
}

func addProxyURL(url string) {
	CreateProxyURL(url)
}


func CheckProxyAlive(proxyURL string) (respBody string, avail bool) {
	proxy, _ := url.Parse(proxyURL)
	httpclient := &http.Client{
		Transport: &http.Transport{
			Proxy:           http.ProxyURL(proxy),
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		Timeout: 20 * time.Second,
	}
	resp, err := httpclient.Get("http://cip.cc/")
	if err != nil {
		return "", false
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", false
	}
	if !strings.Contains(string(body), "地址") {
		return "", false
	}
	return string(body), true
}

func AddProxyURLRetry(url string) error {
	tx := DB.Model(&ProxyURL{}).Where("url = ?", url).Update("retry", gorm.Expr("retry + 1"))
	return tx.Error
}


func SetProxyURLAvail(url string, canBypassGFW bool) error {
	tx := DB.Model(&ProxyURL{}).Where("url = ?", url).Updates(ProxyURL{Retry: 0, Available: true, CanBypassGFW: canBypassGFW})
	return tx.Error
}

func QueryProxyURL() (proxyURLs []ProxyURL, err error) {
	tx := DB.Find(&proxyURLs)
	err = tx.Error
	return
}

func CanBypassGFW(respBody string) bool {
	return strings.Contains(respBody, "香港") ||
		strings.Contains(respBody, "台湾") ||
		strings.Contains(respBody, "澳门") || !strings.Contains(respBody, "中国")
}


func QueryAvailProxyURL() (proxyURLs []ProxyURL, err error) {
	tx := DB.Where("available = ?", true).Find(&proxyURLs)
	err = tx.Error
	return
}

func RandomProxyURL(regionFlag int) (string, error) {
	var proxyURL ProxyURL
	var tx *gorm.DB
	switch regionFlag {
	case 1:
		tx = DB.Raw(fmt.Sprintf("SELECT * FROM %s WHERE available = ? AND can_bypass_gfw = ? ORDER BY RANDOM() LIMIT 1;", proxyURL.TableName()), true, false).Scan(&proxyURL)
	case 2:
		tx = DB.Raw(fmt.Sprintf("SELECT * FROM %s WHERE available = ? AND can_bypass_gfw = ? ORDER BY RANDOM() LIMIT 1;", proxyURL.TableName()), true, true).Scan(&proxyURL)
	default:
		tx = DB.Raw(fmt.Sprintf("SELECT * FROM %s WHERE available = 1 ORDER BY RANDOM() LIMIT 1;", proxyURL.TableName())).Scan(&proxyURL)
	}
	return proxyURL.URL, tx.Error
}
