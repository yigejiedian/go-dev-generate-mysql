# GO-DEV-GENERATE-MYSQL
#项目简介
一个基于MYSQL数据库开发的，轻量级代码生成工具。
可帮助开发者快速开发golang相关代码

#使用说明
## conf/conf.yaml
用于配置mysql和自定义模板的相关信息
```yaml
# mysql数据库相关配置
mysql:
  url: "tcp(127.0.0.1:3306)/ht-admin?charset=utf8mb4&parseTime=True&loc=Local"
  username: "root"
  password: "admin"

# 生成的文件放置的位置
fileRootPath: "/Users/haitao/Work/gitdev/go-dev-generate-mysql/_tmp"

# 生成的模板信息
templates: ["goDo","goRequest"]

goDo:
  # 对应templates文件夹中的模板文件名称
  templateName: "goDo.tpl"
  # 生成的相对路径
  buildPath: "/biz/do"
  # 生成的文件名称
  fileName: "{{ClassName}}Do.go"

goRequest:
  templateName: "goRequest.tpl"
  buildPath: "/api/dto"
  fileName: "{{ClassName}}Request.go"

```

## conf/table.yaml
用于配置基于哪张表生成代码的相关信息
```yaml
table:
  # 删除前缀
  prefix: "sys_"
  # 表名
  name: "sys_user"
  # 功能描述
  comment: "用户表"
```

## cmd/app/go-dev-generate-mysql.do
Main
