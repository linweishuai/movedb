package controller

import (
	"net/http"
	"movedb/common"
	"io/ioutil"
)

func ExecHandler(w http.ResponseWriter,r *http.Request)  {
	err:=r.ParseForm()
	if err != nil {
		common.OutputJson(w, 0, "参数错误")
		return
	}
	jsonstr:=r.FormValue("jsonstr");
	//写入文件
	err=ioutil.WriteFile("configbak.json",[]byte(jsonstr),0755)
	if err!=nil{
		common.OutputJson(w, 0, "写入配置文件错误"+err.Error())
		return
	}
	//开始导入
	Doexport();
}