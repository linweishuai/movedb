# movedb —— MySQL 数据库迁移工具

movedb 是一个基于 Web 可视化界面的 MySQL 数据迁移工具。它把「从哪些源表取数据、映射到哪些目标表、按什么规则过滤和转换」做成了一张可拖拽的画布：添加数据源表、用连线配置字段映射和过滤条件，然后一键执行迁移，并在页面上实时查看进度。

## 功能特性

- 纯 Go 实现，提供一个本地 Web 界面（默认 8888 端口），不需要写代码即可配置和启动迁移
- 字段级映射：`Default`（原样复制）、`OnetoOne`（枚举值一一对应）
- 表级过滤：用类 SQL 条件（如 `$account_type eq 'money'`）只迁移命中的行
- 并发执行：每张导出表独立 goroutine 导出；导入最多同时 `2×CPU` 个 goroutine 写入
- 分批处理：导出每 5 万行查一次，导入每 5 千行生成一条批量 INSERT
- 通过 WebSocket 向页面实时推送迁移日志

## 项目结构

```
movedb
├── main.go            # 入口：启动 HTTP 服务（8888 端口）、配置日志、打开浏览器
├── config.json        # 迁移配置（运行时会被 Web 界面覆盖生成）
├── configsimple.json  # 配置示例
├── controller/        # HTTP 处理器：页面、读取表字段、执行迁移、WebSocket
├── config/            # 配置解析与初始化（构建导出/导入实例、规则表、数据库连接）
├── dbconfig/          # MySQL 连接
├── exporter/          # 从源表分批导出数据
├── importer/          # 批量插入到目标表
├── tranfer/           # 字段值转换规则 Default / OnetoOne
├── sqlmaker/          # 拼接 SELECT / INSERT 语句
├── ws/                # WebSocket 客户端管理，推送迁移日志
├── common/            # 通用 JSON 响应
└── template/          # 前端页面与静态资源（jsPlumb 画布）
```

## 编译运行

前置条件：Go 环境，以及可访问的 MySQL 源库 / 目标库。

> 本项目是 2018 年的 GOPATH 时代项目，没有 `go.mod` / `go.sum`，代码以 `movedb/...` 作为 import 路径，需要把整个仓库放在 `$GOPATH/src/movedb` 下才能编译。

```bash
cd $GOPATH/src/movedb
go get ./...
go build -o movedb.exe .
./movedb.exe
```

程序启动后默认监听 `http://127.0.0.1:8888`。在 Windows 上会自动调用系统命令打开默认浏览器；其他平台请手动在浏览器打开该地址。

## 使用步骤

1. **添加数据源**：在页面上方填写 `alias`、`host`、`username`、`passwd`、`dbname`、`tablename`，选择「导出表」或「导入表」后点「添加」。程序会连接数据库，并通过 `information_schema.COLUMNS` 读取该表的字段列表。
2. **布置画布**：把左侧菜单中的表拖到画布上。红色的是导出表（数据源），蓝色的是导入表（目标）。
3. **配置迁移规则**：
   - 字段连线：把一个字段端点连到另一个字段端点，在弹出的对话框里选择 `Default` 或 `OnetoOne`（OnetoOne 需填写值映射）。
   - 表级连线：把表顶部的表级锚点连到另一张表，在弹出的对话框里填写过滤条件（如 `$account_type eq 'money'`）。
4. **开始迁移**：点击「开始导出」。程序会把当前的连线配置生成为 `config.json` 并立即执行迁移。
5. **查看进度**：页面下方的黑色区域实时显示迁移日志（通过 WebSocket 推送）。

## 配置文件说明

迁移规则本质上是一个 JSON 配置文件（`config.json`）。可以通过 Web 界面生成，也可以手工编写。

```json
{
  "ExportDb": {
    "a": {"host": "127.0.0.1", "port": "3306", "username": "root", "passwd": "root", "dbname": "db_old", "tablename": "admin"}
  },
  "ImportDb": {
    "b": {"host": "127.0.0.1", "port": "3306", "username": "root", "passwd": "root", "dbname": "db_new", "tablename": "admin"}
  },
  "Fieldrule": {
    "b.id": ["a.id", "Default"],
    "b.status": ["a.status", "OnetoOne", "1,2:0,1"]
  },
  "RuleField": {
    "a": ["account_type"]
  },
  "Tablerule": {
    "ab": "$account_type eq 'money'"
  }
}
```

- `ExportDb`：导出数据源表，key 是别名（不能重复），value 里的字段名固定。
- `ImportDb`：导入目标表，key 是别名（不能重复），value 里的字段名固定。
- `Fieldrule`：字段级映射规则，格式为 `"目标表别名.目标字段": ["源表别名.源字段", "规则", "附加参数"]`。规则说明见下节。
- `RuleField`：需要额外从源表读取的字段。这些字段本身不参与映射，但会出现在表级过滤条件里（例如 `$account_type`），因此必须被 SQL 查出来才能参与判断。
- `Tablerule`：表级过滤条件，key 是「源表别名 + 目标表别名」拼接（如 `ab` 表示 a 表 → b 表）。

### 字段映射规则

| 规则 | 说明 |
|------|------|
| `Default` | 原样复制：源字段的值直接写入目标字段 |
| `OnetoOne` | 值一一对应：第三项用 `导出值1,导出值2,...:导入值1,导入值2,...` 指定映射表。例如 `"1,2:0,1"` 表示导出数据中的 `1`、`2` 在导入时分别写成 `0`、`1` |

### 表级过滤条件

`Tablerule` 里的条件用 `$字段名` 加比较运算符书写，例如 `$account_type eq 'money' and $id egt 0`。程序会把：

- `$字段名` 替换为当前行的实际取值
- `eq`→`==`、`neq`→`!=`、`gt`→`>`、`lt`→`<`、`egt`→`>=`、`elt`→`<=`、`and`→`&&`、`or`→`||`

翻译成一段 JavaScript 表达式，用内置 JS 引擎（otto）求值，表达式为真的行才会被迁移。

## 处理流程

迁移在内存中分阶段进行，全程并发：

1. 读取 `./config.json`，解析出导出表、导入表、字段规则、表规则，并为每张表建立 MySQL 连接。
2. **导出**：每张导出表一个 goroutine，先 `SELECT count(*)` 得到总行数，再按 `LIMIT 50000` 分批查询，把结果通过 channel 交给主流程。
3. **转换**：对每一批数据逐行套用表级过滤条件过滤；命中的行再按字段映射规则和目标表的字段列表，生成待插入的数据。
4. **导入**：每满 5000 行生成一条批量 `INSERT`，最多同时 `2×CPU` 个导入 goroutine 并发写入。
5. 迁移进度通过 WebSocket 实时推送到页面。

## 注意事项

- MySQL 连接固定使用 3306 端口；配置文件里的 `port` 字段暂未生效。
- 一个目标字段只能连一条来源（目标字段唯一）；一张源表可以拆分映射到多张目标表。
- 大批量插入会生成很大的 SQL 语句，可能需要调大 MySQL 的 `max_allowed_packets`。
- 程序从当前目录读取 `./config.json`，配置文件缺失或解析出错会报错。
