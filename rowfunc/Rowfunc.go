package fieldrule

type RowFunc func(map[string]*[]byte,string) string
type FieldRule interface {
	NumbertoNumber([]string,[]string) RowFunc
	Default() RowFunc
}