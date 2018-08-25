package controller

import (
	"net/http"
	"movedb/common"
	"movedb/dbconfig"
	"fmt"
)

func GetDbInfoHandler(w http.ResponseWriter,r *http.Request)  {
	w.Header().Set("content-type", "application/json")
	err:=r.ParseForm()
	if err != nil {
		common.OutputJson(w, 0, "参数错误")
		return
	}

	host:=r.FormValue("host");
	username:=r.FormValue("username");
	password:=r.FormValue("password");
	dbname:=r.FormValue("dbname");
	tablename:=r.FormValue("tablename");
	dbconfg:=dbconfig.DbConfig{Host:host,Username:username,Passwd:password,Dbname:dbname}
	connction,err:=dbconfg.GetDbInstance()
	if err!=nil{
		common.OutputJson(w, 0, dbname+"链接错误"+err.Error())
		return
	}
	sql:="select COLUMN_NAME from information_schema.COLUMNS where table_name = '"+tablename+"' and table_schema = '"+dbname+"';"
	exportRs, err := connction.Query(sql)
	if err!=nil{
		common.OutputJson(w, 0, tablename+"查询字段出错"+err.Error())
		return
	}

	result:=make([]string,0)
	for exportRs.Next(){
		var COLUMN_NAME string
		err = exportRs.Scan(&COLUMN_NAME)
		if err!=nil{
			common.OutputJson(w, 0, fmt.Sprintf("查询字段出错%v",err))
			return
		}
		result=append(result,COLUMN_NAME)
	}
	common.OutputJson(w, 1, result)

}