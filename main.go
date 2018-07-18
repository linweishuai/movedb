package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"encoding/json"
	"strings"
	"movedb/dbconfig"//数据哭连接
	"movedb/exporter"//导出
	"movedb/sqlmaker"//sql生成
	"movedb/importer"//导入
	"sync"
	"math"
	"time"
	"github.com/cihub/seelog"
	"database/sql"
	"os"
	"runtime"
)

func main() {
	SetLogger("logConfig.xml")
	cpunumber:=runtime.NumCPU()//导入进程做好是cpu的两倍
	config_path := flag.String("c","./config.json","config_file path")
	goruntimeNumber:=flag.Float64("n",5000.00,"每个导入携程处理数量")
	flag.Parse()
	fmt.Printf("配置文件位置 : %s\n",*config_path)
	content,err:=ioutil.ReadFile(*config_path)
	 if err!=nil{
		seelog.Infof("读取配置文件出错")
	}
	type Conifg struct {
		ExportDb map[string]map[string]string
		ImportDb map[string]map[string]string
		Fieldrule map[string][]string
	}
	var inter Conifg
	jsonerr := json.Unmarshal(content, &inter)
	if jsonerr != nil {
		seelog.Infof("error in translating,", err.Error())
		return
	}
	//解析出要导出的所有字段
	var exportField=make(map[string][]string)
	//var ExportResult sync.Map
	//fmt.Print(inter.Fieldrule)
	var ImportTableFieldRelation=make(map[string][]string)//导入数据库的字段
	var ExImRelation=make(map[string][]string)//标注我导出的数据导入到那些表里面去
	var tempmap=make(map[string]struct{})//临时表
	var exportFieldtemp=make(map[string]struct{})
	type Fieldrule struct {
		TransferRule string
		ExtraData [][]string
	}
	var NewFeildmap=make(map[string]Fieldrule)
	for Importtarget,exportsource:=range inter.Fieldrule{
			//TODO importtarget not ready
			ImporttableFieldslice:=strings.Split(Importtarget,".")
			ImportTableFieldRelation[ImporttableFieldslice[0]]=append(ImportTableFieldRelation[ImporttableFieldslice[0]],ImporttableFieldslice[1])
			tablefield:=exportsource[0]
			ExportableFieldslice:=strings.Split(tablefield,".")
			if _,ok:=exportFieldtemp[ExportableFieldslice[0]+ExportableFieldslice[1]];!ok{
				exportField[ExportableFieldslice[0]]=append(exportField[ExportableFieldslice[0]],ExportableFieldslice[1])
				exportFieldtemp[ExportableFieldslice[0]+ExportableFieldslice[1]]= struct{}{}
			}
		    _,ok:=tempmap[ImporttableFieldslice[0]]
		    if!ok{//如果读取不到那么就说明不存在这导入表
				tempmap[ImporttableFieldslice[0]]=struct{}{}
				ExImRelation[ExportableFieldslice[0]]=append(ExImRelation[ExportableFieldslice[0]],ImporttableFieldslice[0])
			}
			//如果是onetoone 模式 那么就需要 两个一一对应的参数
			var ExtraData [][]string
			if(exportsource[1]=="OnetoOne"){
				rulemap:=strings.Split(exportsource[2],":")//分隔出来然后
				for _,v:=range rulemap{
					ExtraData=append(ExtraData,strings.Split(v,","))
				}
			}
			NewFeildmap[Importtarget+":"+tablefield]=Fieldrule{
				TransferRule:exportsource[1],
				ExtraData:ExtraData,
			}
			//parse condition from neq eq egt elt
	}
		//fmt.Println(ImportTableFieldRelation)
		//os.Exit(1)
		var Transferchan=make(chan exporter.Tranferdata,len(inter.ExportDb))
		//通道的方式传递数据
		var ExporterMap=make(map[string]exporter.Exporter);
		for tableName,FieldSlice:=range exportField {
			exportDB := dbconfig.DbConfig{
				Host:     inter.ExportDb[tableName]["host"],
				Username: inter.ExportDb[tableName]["username"],
				Passwd:   inter.ExportDb[tableName]["passwd"],
				Dbname:   inter.ExportDb[tableName]["dbname"],
			}
			var Dbconnection = exportDB.GetDbInstance();
			Exporter := exporter.Exporter{
				Dbconnction: Dbconnection,
				Selectsqlmaker: sqlmaker.Selectsqlmaker{
					Sqlmaker: sqlmaker.Sqlmaker{
						Tablename: inter.ExportDb[tableName]["tablename"],
						Field:     FieldSlice,
					},
				},
			}
			ExporterMap[tableName]=Exporter
		}
		go func() {
			var wg sync.WaitGroup
			for Tname,Export:=range ExporterMap{
				wg.Add(1)
				go Export.Export(Transferchan,Tname,&wg)
			}
			wg.Wait()
			close(Transferchan)
			seelog.Infof("所有表数据到处完成")
		}()
		//从导出的数据入手 导出的数据多有几个那么就遍历这些数据
		beginTime := time.Now().Unix()
		var DbInstance=make(map[string]*sql.DB)
		for transferdata:=range Transferchan {
			//每次导入5000数据
			key:=transferdata.TableName
			value:=transferdata.Data
			seelog.Infof("处理导出%s导出数据共%d条数据",key,len(value))
			var ig sync.WaitGroup
			ProcessChan := make(chan struct{}, cpunumber*2)//导入进程是cpu*2
			goroutineNumber:=*goruntimeNumber
			GoroutineNumber := int(goroutineNumber)
			goNumber := math.Ceil(float64(len(value)) / goroutineNumber)
			for i := 0; i < int(goNumber); i++ {
				start := i * GoroutineNumber
				end := (i + 1) * GoroutineNumber
				var tempSlice []map[string]string
				if i == int(goNumber)-1 {
					tempSlice = value[start:]
				} else {
					tempSlice = value[start:end]
				}
				//fmt.Println(tempSlice)
				//os.Exit(1)
				Importslice:=make(map[string][]map[string]string)
				seelog.Infof("处理导出%s导出数 据第%d到%d数据",key,start,end)
				for _,values:=range tempSlice{
					//seelog.Infof(values)
					for _,Importable:=range ExImRelation[key]{
						var rowdata=make(map[string]string)
						//seelog.Infof(key.(string)+Importable)
						for _,exportfield:= range exportField[key]  {
							//seelog.Infof(exportfield)
							for _,importfield:=range ImportTableFieldRelation[Importable]{
								//seelog.Infof(importfield)
								//seelog.Infof(Importable+"."+importfield+":"+key.(string)+"."+exportfield)
								rule,ok:=NewFeildmap[Importable+"."+importfield+":"+key+"."+exportfield]
								if ok{
									switch (rule.TransferRule){
									case "Default":
										rowdata[importfield]=fmt.Sprintf("%q", values[exportfield])
									case "OnetoOne":
										index :=0
										for nowIndex,content:=range rule.ExtraData[0]{
											if(content==values[exportfield]){
												index=nowIndex
											}
										}
										rowdata[importfield]=fmt.Sprintf("%q", rule.ExtraData[1][index])
									default:
										rowdata[importfield]=fmt.Sprintf("%q", values[exportfield])
									}
								}else{
									continue
								}
							}
						}
						Importslice[Importable]=append(Importslice[Importable],rowdata)
					}
				}
				//fmt.Println(Importslice)
				//fmt.Println(ImportTableFieldRelation)
				//os.Exit(1)
				for tablealias,importtempslice:=range Importslice{
					ig.Add(1)
						var insertsqlmaker=sqlmaker.Insertsqlmaker{
								Sqlmaker:sqlmaker.Sqlmaker{
									Tablename:inter.ImportDb[tablealias]["tablename"],
									Field:ImportTableFieldRelation[tablealias],
								},
								Dataslice:importtempslice,
						}
						//fmt.Println(insertsqlmaker)
						var ImportDbconfig=dbconfig.DbConfig{
							Host:inter.ImportDb[tablealias]["host"],
							Username:inter.ImportDb[tablealias]["username"],
							Passwd:inter.ImportDb[tablealias]["passwd"],
							Dbname:inter.ImportDb[tablealias]["dbname"],
						}
						_,ok:=DbInstance[tablealias]
						if !ok{
							DbInstance[tablealias]=ImportDbconfig.GetDbInstance()
						}
						var Importer=importer.Importer{
							Dbconnction:DbInstance[tablealias],
							Insertsqlmaker:insertsqlmaker,
						}
						ProcessChan <- struct{}{}
						go func(i int) {
							seelog.Infof("导入进程%v开启", i+1)
							Importer.Import(ProcessChan)
							defer ig.Done()
						}(i)
				}
				seelog.Infof("处理导出%s导出数据第%d到%d数据完成",key,start,end)
			}
			ig.Wait()
		}
		finishTime := time.Now().Unix()
		defer func() {
			seelog.Infof("导入完成实际消耗时间为：%v秒", finishTime-beginTime)
			seelog.Flush()
		}()

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
