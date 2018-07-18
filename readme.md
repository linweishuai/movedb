
```
{
"ExportDb": {//导出数据库信息(字段名称不能动)
        "a": {//别名 不能重复
        "host": "127.0.0.1",
        "port": "3306",
        "username": "root",
        "passwd": "root",
        "dbname": "linweishuai_naoshijie",
        "tablename": "admin"
    }
},
"ImportDb": {//导出数据库信息(字段名称不能动)
        "b": {//别名 不能重复
        "host": "127.0.0.1",
        "port": "3306",
        "username": "root",
        "passwd": "root",
        "dbname": "linweishuai_naoshijiev3",
        "tablename": "admin"
    }
},
"Fieldrule": {(字段名称不能动)
        "b.id": [ //导入目标字段
            "a.id",//导出来源字段
            "Default"//两个字段之间的关系类型
        ],
        "b.account": [
            "a.account",
            "Default"
        ],
        "b.status": [
            "a.status",
            "OnetoOne",
            "1,2:0,1"
        ]
    }
}
```
# 使用 
main.exe |-c string "配置文件位置(./config.json)" | -n float 每个协程导入数据数量 (默认5000.00)
## 示例
个人设置
main.exe -c "./config.json" -n 5000.00
# mysql配置修改
导入max_allowed_packets(每次导入数据量*每条数据大小/1024 一般8M就够够的了)
# makejson.html
1.生成config.json 文件 (导入字段唯一)
# movedb
1. Default(原样输出)
2. OnetoOne(一一对应)例如以上例子中的 1,2(导出数据可能存在值):0,1(导入数据对应值)


