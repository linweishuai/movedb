package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"encoding/json"
	"movedb/exporter"//导出
	"movedb/sqlmaker"//sql生成
	"movedb/importer"//导入
	"sync"
	"math"
	"time"
	"github.com/cihub/seelog"
	"os"
	"github.com/robertkrimen/otto"
	"movedb/config"
	"regexp"
	"strings"
	"movedb/tranfer"
)

func main() {
	SetLogger("logConfig.xml")
	config_path := flag.String("c","./config.json","config_file path")
	goruntimeNumber:=flag.Float64("n",5000.00,"每个导入携程处理数量")
	flag.Parse()
	fmt.Printf("配置文件位置 : %s\n",*config_path)
	content,err:=ioutil.ReadFile(*config_path)
	if err!=nil{
		seelog.Infof("读取配置文件出错")
	}
	var Allconfig config.Allconfig;
	var config config.Conifg
	jsonerr := json.Unmarshal(content, &config)
	if jsonerr != nil {
		seelog.Infof("error in translating,", err.Error())
		return
	}
	Allconfig.Initrule(config)
	Allconfig.InitexportdbInstance(config)
	var Transferchan=make(chan exporter.Tranferdata,len(config.ExportDb))
	go func() {
		var wg sync.WaitGroup
		for Tname,Export:=range Allconfig.Exporter{
			wg.Add(1)
			go Export.Export(Transferchan,Tname,&wg)
		}
		wg.Wait()
		close(Transferchan)
		seelog.Infof("所有表数据到处完成")
	}()
	//从导出的数据入手 导出的数据多有几个那么就遍历这些数据
	beginTime := time.Now().Unix()
	for transferdata:=range Transferchan {
		//每次导入5000数据
		key:=transferdata.TableName
		value:=transferdata.Data
		seelog.Infof("处理导出%s导出数据共%d条数据",key,len(value))
		var ig sync.WaitGroup
		ProcessChan := make(chan struct{}, Allconfig.Cpunumber*2)//导入进程是cpu*2
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
			vm := otto.New()
			for _,values:=range tempSlice{
				//seelog.Infof(values)
				for _,Importable:=range Allconfig.ExImRelation[key]{
					//如果有条件限制 就取出条件限制
					//example $account_type eq 'money' and id egt 5000
					//transfer $account_type == 'money' && id >= 5000
					//replace 'coupon'=='money' && '20000'>=5000 `javascript`
					if conditon,ok:=Allconfig.Tablerule[key+Importable];ok{
						var reg=regexp.MustCompile("\\$\\w+")
						javascript:=string(reg.ReplaceAllFunc([]byte(conditon), func(i []byte) []byte {
							return []byte(fmt.Sprintf("%q",values[strings.Replace(string(i),"$","",-1)]))
						}))
						jsres,_:=vm.Run(javascript)
						res,_:=jsres.ToBoolean()
						if !res{
							continue;
						}
					}
					var rowdata=make(map[string]string)
					//seelog.Infof(key.(string)+Importable)
					for _,exportfield:= range Allconfig.ExportTableField[key]  {
						//seelog.Infof(exportfield)
						for _,importfield:=range Allconfig.ImportTableField[Importable]{
							//seelog.Infof(importfield)
							//seelog.Infof(Importable+"."+importfield+":"+key.(string)+"."+exportfield)
							rule,ok:=Allconfig.Fieldrule[Importable+"."+importfield+":"+key+"."+exportfield]
							if ok{
								rowdata[importfield]=tranfer.DoTransfer(rule.TransferRule,exportfield,values,rule)
								//switch (rule.TransferRule){
								//case "Default":
								//	rowdata[importfield]=fmt.Sprintf("%q", values[exportfield])
								//case "OnetoOne":
								//	index :=0
								//	for nowIndex,content:=range rule.ExtraData[0]{
								//		if(content==values[exportfield]){
								//			index=nowIndex
								//		}
								//	}
								//	rowdata[importfield]=fmt.Sprintf("%q", rule.ExtraData[1][index])
								//default:
								//	rowdata[importfield]=fmt.Sprintf("%q", values[exportfield])
								//}
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
						Tablename:config.ImportDb[tablealias]["tablename"],
						Field:Allconfig.ImportTableField[tablealias],
					},
					Dataslice:importtempslice,
				}
				var Importer=importer.Importer{
					Dbconnction:Allconfig.Importer[tablealias],
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
