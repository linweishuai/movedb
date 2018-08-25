package main

import (
	"github.com/cihub/seelog"
	"os"
	"net/http"
	"movedb/controller"
)

func main() {
	SetLogger("logConfig.xml")
	//开启http 服务
	go Registerwebservice();
	//开启ws服务日志输出
	go Registerwsservice()

	select {}
}
func SetLogger(fileName string) {
	if _, err := os.Stat(fileName); err == nil {
		logger, err := seelog.LoggerFromConfigAsFile(fileName)
		if err != nil {
			panic(err)
		}
		seelog.ReplaceLogger(logger)
	} else {
		configString := `<seelog>
                        <outputs formatid="main">
                            <filter levels="info,error,critical">
                                <rollingfile type="date" filename="log/AppLog.log" namemode="prefix" datepattern="02.01.2006"/>
                            </filter>
                            <console/>
                        </outputs>
                        <formats>
                            <format id="main" format="%Date %Time [%LEVEL] %Msg%n"/>
                        </formats>
                        </seelog>`
		logger, err := seelog.LoggerFromConfigAsString(configString)
		if err != nil {
			panic(err)
		}
		seelog.ReplaceLogger(logger)
	}
}
func Registerwebservice(){
	http.Handle("/css/", http.FileServer(http.Dir("template")))
	http.Handle("/js/", http.FileServer(http.Dir("template")))
	//注册几个
	http.HandleFunc("/getDbinfo", controller.GetDbInfoHandler)
	http.HandleFunc("/exec",controller.ExecHandler)
	http.HandleFunc("/",controller.IndexHandler)
	http.ListenAndServe(":8888", nil)
}
func Registerwsservice()  {
	go controller.Manager.Start()
	http.HandleFunc("/ws", controller.WsPage)
	http.ListenAndServe(":8888", nil)
}