package main

import (
	"log"
	"net"
	"net/http"
	"net/rpc"
	"path/filepath"

	config "github.com/panjjo/go-config"
	"github.com/panjjo/ppp"
)

var configPath string

func main() {
	//config
	configPath = "/opt/workplace/gopath/src/github.com/panjjo/ppp"
	err := config.ReadConfigFile(filepath.Join(configPath, "/config.yml"))
	if err != nil {
		return
	}
	user := new(ppp.Account)
	rpc.Register(user)
	config.Mod = "alipay"
	if ok, err := config.GetBool("status"); ok {
		initAliPay()
		ali := new(ppp.AliPay)
		rpc.Register(ali)
	} else {
		log.Fatal(err)
	}

	rpc.HandleHTTP()
	l, e := net.Listen("tcp", ":1234")
	if e != nil {
		log.Fatal("Listen error:", e)
	}
	http.Serve(l, nil)
}
func initAliPay() {
	ali := ppp.AliPayInit{
		ServiceProviderId: config.GetStringDefault("serviceid", ""),
		ConfigPath:        configPath,
	}
	var err error
	if ali.AppId, err = config.GetString("appid"); err != nil {
		log.Fatal("Init Error:Not Found alipay:appid")
	}
	if ali.Url, err = config.GetString("url"); err != nil {
		log.Fatal("Init Error:Not Found alipay:url")
	}
	ali.Init()
}