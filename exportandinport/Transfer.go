package exportandinport

import (
	"database/sql"
	"sync"
	"movedb/rowfunc"
)

type Transfer struct {
	ImportDb    *sql.DB
	ExportDb    *sql.DB
	ExportField []string
	ImportField []string
	sync.Mutex
	ImportAllData   []map[string]string //查询出的所有结果
	GoroutineNumber float64             //每个goroutine导入的数量
	GoroutineExportNumber int             //每个goroutine读取的数量
	TotalNumber int
	ReadyChan chan struct{} //是否准备进行插入操作
	ActualInsert  int64
	ProcessChan chan struct{} //两个进程可以同时插入,使用channel控制插入频率
	ExportChan chan struct{}
	RowDefaultData map[string]fieldrule.RowFunc
}

func (d *Transfer) InitFinish() {
	d.ProcessChan = make(chan struct{}, 5)
	d.ExportChan = make(chan struct{}, 5)
	d.ImportField = make([]string, 0, len(d.ExportField))

	for _, field := range d.ExportField {
		d.ImportField = append(d.ImportField, field)
	}
}