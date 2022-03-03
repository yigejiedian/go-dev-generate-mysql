package main

import (
	"context"
	"fmt"
	generate "go-dev-generate-mysql/internal/pkg"
	"go-dev-generate-mysql/internal/pkg/db"
	"golang.org/x/sync/errgroup"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
	"sync"
	"sync/atomic"
	"text/template"
)

//var tableName ="sys_user"
//var tablePrefix ="sys_"
//var tableComment ="用户表"
var rootPath string

var wait sync.WaitGroup
var count int64 = 0

func main() {
	Configs := make(map[string]interface{})
	TableConfigs := make(map[string]interface{})

	data, _ := ioutil.ReadFile("./conf/conf.yaml")
	err := yaml.Unmarshal(data, &Configs)
	if err != nil {
		log.Fatal(err)
	}

	table, _ := ioutil.ReadFile("./conf/table.yaml")
	err = yaml.Unmarshal(table, &TableConfigs)
	if err != nil {
		log.Fatal(err)
	}

	// 数据库初始化
	mysqlConfs := Configs["mysql"].(map[string]interface{})
	db.InitDB(mysqlConfs)

	//配置生成表相关信息
	tablePrefix, tableName, tableComment := SetTableConfigs(TableConfigs)

	var templateConfs []*TemplateConf

	rootPath = Configs["fileRootPath"].(string)

	templates := Configs["templates"].([]interface{})
	for _, tp := range templates {
		tpName := tp.(string)

		confMaps := Configs[tpName].(map[string]interface{})

		templateConf := &TemplateConf{
			TemplateName: confMaps["templateName"].(string),
			BuildPath:    confMaps["buildPath"].(string),
			FileName:     confMaps["fileName"].(string),
		}
		templateConfs = append(templateConfs, templateConf)
	}

	withContext, ctx := errgroup.WithContext(context.Background())

	rootGenerate, err := NewGenerate(tablePrefix, tableName, tableComment)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%v\n", rootGenerate)
	fmt.Printf("%v\n", rootGenerate.DbResult)

	for i := 0; i < len(templateConfs); i++ {
		wait.Add(1)
		//go createCodeFile(templateConfs[i], rootGenerate)
		conf := templateConfs[i]
		withContext.Go(func() error {
			createCodeFile(conf, rootGenerate)
			return nil
		})
	}

	err = withContext.Wait()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(ctx.Err())

	wait.Wait()
	fmt.Println("构建成功=", atomic.LoadInt64(&count))
}

//配置生成表相关信息
func SetTableConfigs(tableConfigs map[string]interface{}) (tablePrefix string, tableName string, tableComment string) {
	tableConfs, ok := tableConfigs["table"]
	if !ok {
		log.Fatal("table not find")
	}
	tableMaps := tableConfs.(map[string]interface{})
	tablePrefixConf, ok := tableMaps["prefix"]
	if !ok {
		log.Fatal("table.prefix not find")
	}
	tablePrefix = tablePrefixConf.(string)

	nameConf, ok := tableMaps["name"]
	if !ok {
		log.Fatal("table.name not find")
	}
	tableName = nameConf.(string)

	commentConf, ok := tableMaps["comment"]
	if !ok {
		log.Fatal("table.comment not find")
	}
	tableComment = commentConf.(string)
	return
}

type TemplateConf struct {
	TemplateName string
	BuildPath    string
	FileName     string
}

func createCodeFile(templateConf *TemplateConf, rootGenerate *Generate) {
	defer wait.Done()
	atomic.AddInt64(&count, 1)

	buildFileName := strings.ReplaceAll(templateConf.FileName, "{{ClassName}}", rootGenerate.ClassName)
	buildFileName = strings.ReplaceAll(buildFileName, "{{GenerateName}}", rootGenerate.GenerateName)
	buildFile("templates/"+templateConf.TemplateName, templateConf.BuildPath, buildFileName, rootGenerate)
}

func buildFile(templatePath string, buildPath string, buildFile string, generate *Generate) {
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		log.Fatalln(err)
	}

	tmpPath := rootPath + buildPath
	// buildPath 是否存在,不存在则创建
	if _, err := os.Stat(tmpPath); err != nil && os.IsNotExist(err) {
		os.MkdirAll(tmpPath, os.ModePerm)
	}

	tmpPath = tmpPath + "/" + buildFile
	if _, err := os.Stat(tmpPath); err != nil && os.IsNotExist(err) {
		//f, err := os.OpenFile(buildPath, os.O_WRONLY|os.O_CREATE, 0644)
		f, err := os.Create(tmpPath)
		defer f.Close()
		if err != nil {
			log.Fatalln(err)
		}
		if err := tmpl.Execute(f, generate); err != nil {
			log.Fatalln(err)
		}
	}

}

func NewGenerate(tablePrefix string, tableName string, tableComment string) (*Generate, error) {
	var dbResult []DbResult
	err := db.DB.Raw("SHOW FULL COLUMNS FROM " + tableName).Scan(&dbResult).Error
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	rootVal := &Generate{
		TableName:    tableName,
		Comment:      tableComment,
		GenerateName: strings.TrimPrefix(tableName, tablePrefix),
	}
	rootVal.ClassName = ClassNameFormat(rootVal.GenerateName)
	rootVal.VariableName = VariableNameFormat(rootVal.ClassName)

	for i := 0; i < len(dbResult); i++ {
		v := &dbResult[i]
		v.init()
	}

	rootVal.DbResult = &dbResult

	return rootVal, nil
}

type Generate struct {
	TableName    string
	Comment      string
	ClassName    string
	VariableName string
	GenerateName string
	DbResult     *[]DbResult
}

type DbResult struct {
	Field      string
	Type       string
	Collation  string
	Null       string
	Key        string
	Default    string
	Extra      string
	Privileges string
	Comment    string

	ClassName    string
	VariableName string
	GoType       string
	TagByJson    string
	TagByGorm    string
}

func (d *DbResult) init() {
	d.Field = strings.ToLower(d.Field)
	d.setClassName()
	d.setVariableName()
	d.setGoType()
	d.setTagByJson()
	d.setTagByGorm()
}

func (d *DbResult) setClassName() {
	d.ClassName = ClassNameFormat(d.Field)
}
func (d *DbResult) setVariableName() {
	d.VariableName = VariableNameFormat(d.ClassName)
}

func (d *DbResult) setGoType() {
	d.GoType = getTypeDef(d.Type)
}

func (d *DbResult) setTagByJson() {
	d.TagByJson = fmt.Sprintf("json:\"%s\"", d.Field)
}

func (d *DbResult) setTagByGorm() {
	str := ""
	if d.Key != "" {
		str += d.Key + ";"
	}
	str += "column:" + d.Field + ";"
	str += "type:" + d.GoType + ";"
	if d.Null == "NO" {
		str += "not null;"
	}

	// gorm:"primaryKey;column:id;type:bigint;not null"
	d.TagByGorm = fmt.Sprintf("gorm:\"%s\"", str)
}

//lowercase

// 类格式命名(驼峰)
func ClassNameFormat(name string) string {
	split := strings.Split(name, "_")
	str := ""
	for _, v := range split {
		runes := []rune(v)
		if len(runes) > 0 {
			if runes[0] >= 'a' && runes[0] <= 'z' {
				runes[0] -= 32
			}
		}
		str += string(runes)
	}
	return str
}

// 变量格式命名(首字母小写)
func VariableNameFormat(name string) string {
	runes := []rune(name)
	if len(runes) > 0 {
		if runes[0] >= 'A' && runes[0] <= 'Z' {
			runes[0] += 32
		}
	}
	return string(runes)
}

func getTypeDef(mysqlType string) string {
	mp := generate.TypeMysqlDicMp
	val, ok := mp[mysqlType]
	if ok {
		return val
	}

	//index := strings.Index(mysqlType, "(")
	//if index != 0 {
	//	bytes := []byte(mysqlType)
	//	newBytes := bytes[0:index]
	//	return getTypeDef(string(newBytes))
	//}

	for _, l := range generate.TypeMysqlMatchList {
		if ok, _ := regexp.MatchString(l.Key, mysqlType); ok {
			return l.Value
		}
	}
	return ""
}
