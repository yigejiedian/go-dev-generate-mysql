package do

import (
	"ht-admin-server/internal/pkg/code/global"
)

// {{.ClassName}} {{.Comment}}
type {{.ClassName}}DO struct {
	global.BaseModel{{"\n"}}{{ range $index, $element := .DbResult }}{{"\t"}}{{$element.ClassName}}{{"\t"}}{{$element.GoType}}{{"\t"}}`{{$element.TagByGorm}} {{$element.TagByJson}}`{{"\t\t"}}//{{$element.Comment}}{{"\n"}}{{end}}
}

func ({{.ClassName}}DO) TableName() string{
	return "{{.TableName}}"
}


// {{.ClassName}}Columns get sql column name.获取数据库列名
var {{.ClassName}}Columns = struct {
{{ range $index, $element := .DbResult }}{{"\t"}}{{$element.ClassName}}{{"\t\t\t"}}string{{"\n"}}{{end}}
}{
{{ range $index, $element := .DbResult }}{{"\t"}}{{$element.ClassName}}:{{"\t\t\t"}}"{{$element.Field}}",{{"\n"}}{{end}}
}
