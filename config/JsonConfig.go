package config

type Conifg struct {//config.json 结构体
	ExportDb map[string]map[string]string
	ImportDb map[string]map[string]string
	Fieldrule map[string][]string
	Tablerule map[string]string
	RuleField map[string][]string
}
