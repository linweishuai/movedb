
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
# movedb
Json 必须严格按照此模式来书写
可以现在php里面写好然后json_encode一下
- ExportDb ImportDb Fieldrule 三个字段写不动
- 导入目标字段在Fieldrule中必须唯一,所以现在仅支持一个表拆分成两个表 单表对单表
- 5000记录做一次插入操作,同时保证5个goruntime持续写入数据,根据自己数据库更改,当插入数据较大是需要改变数据库配置 max_allowed_packets
- 字段对应关系现在只有两种
1. Default(原样输出)
2. OnetoOne(一一对应)例如以上例子中的 1,2(导出数据可能存在值):0,1(导入数据对应值)

1,2 | 0,1
---|---
导出数据可能存在值 | 导入数据对应值


