package tranfer

import (
	"movedb/config"
	"fmt"
	"reflect"
)

type Transfer struct {

}
func (H Transfer) Default(Fieldname string, value map[string]string,rule config.Fieldrule) string  {
	return value[Fieldname]
}
func (H Transfer) OnetoOne(Fieldname string, value map[string]string,rule config.Fieldrule) string  {
	index :=0
	for nowIndex,content:=range rule.ExtraData[0]{
		if(content==value[Fieldname]){
			index=nowIndex
		}
	}
   return rule.ExtraData[1][index]
}

func DoTransfer(Rule string,args...interface{}) string {
	inputs := make([]reflect.Value, len(args))
	for i, _ := range args {
		inputs[i] = reflect.ValueOf(args[i])
	}
	return fmt.Sprintf("%q",reflect.ValueOf(new(Transfer)).MethodByName(Rule).Call(inputs)[0].String())
}


