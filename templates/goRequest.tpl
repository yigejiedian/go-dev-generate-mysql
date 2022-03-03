package dto

type List{{.ClassName}}Request struct {
	Id string `json:"id"`

	PageNum  int `json:"page_num"`
	PageSize int `json:"page_size"`
}


type Get{{.ClassName}}Request struct {
	Id string `json:"id"`
}



type Create{{.ClassName}}Request struct {
{{ range $index, $element := .DbResult }}{{ if and (ne $element.ClassName "Id") (ne $element.ClassName "CreateUserId") (ne $element.ClassName "CreateTime") (ne $element.ClassName "UpdateUserId") (ne $element.ClassName "UpdateTime") (ne $element.ClassName "Deleted") (ne $element.ClassName "Version") }}{{"\t"}}{{$element.ClassName}}{{"\t"}}{{$element.GoType}}{{"\t"}}`{{$element.TagByJson}}`{{"\t\t"}}// {{$element.Comment}}{{"\n"}}{{end}}{{end}}
}



type Update{{.ClassName}}Request struct {
    Id          string  `json:"id" binding:"required"`      // ID

{{ range $index, $element := .DbResult }}{{ if and (ne $element.ClassName "Id") (ne $element.ClassName "CreateUserId") (ne $element.ClassName "CreateTime") (ne $element.ClassName "UpdateUserId") (ne $element.ClassName "UpdateTime") (ne $element.ClassName "Deleted") (ne $element.ClassName "Version") }}{{"\t"}}{{$element.ClassName}}{{"\t"}}{{$element.GoType}}{{"\t"}}`{{$element.TagByJson}}`{{"\t\t"}}// {{$element.Comment}}{{"\n"}}{{end}}{{end}}

    Version      int       `json:"version"`        // 版本号
}



type Delete{{.ClassName}}Request struct {
	Ids []string `json:"ids" binding:"required"`
}
