## logger 简介 

当前版本只支持 linux 系统使用 , 后续可能会支持Windows 

支持追加方式写入文件 , 保持到磁盘

提供日志处理接口用于: `上传到云`或者`存储到数据库`

在`goland` IDE 中可快速定为到打印的行数

部分代码参考了其他[Github项目](<https://github.com/wonderivan/logger>)

## 安装

```shel
go get github.com/eliot-jay/logger
```

####	如果get发生错误:

```shell
go get: github.com/eliot-jay/logger@v1.0.0: parsing go.mod:
        module declares its path as: logger
                but was required as: github.com/eliot-jay/logger
```

##### 错误解决方法:

在go.mod 里使用`replace`关键字添加包的`绝对路径`

路径在你 **GOPATH** 里

`gopath/pkg/mod/github.com/eliot-jay/logger@v1.0.0`

```
replace github.com/eliot-jay/logger v1.0.0 => /home/jay/gosrc/pkg/mod/github.com/eliot-jay/logger@v1.0.0
```

没有错误请忽略

##	日志的等级:

一共有五个等级 , `打印时的等级要等于或者小于配置的等级时`才会输出消息

| 等级 | 配置 | 释义                                             | 控制台颜色 |
| ---- | ---- | ------------------------------------------------ | :--------: |
| 0    | SERIOUS |可能有危险的严重错误,如:初始化,数据库连接错误等 |红色底|
| 1    | ERROR |普通错误,断言失败,类型转换失败等   						 |红色|
| 2    | WARN | 普通警告，比如权限出错，访问异常等               |紫色底|
| 3    | INFO | 重要消息                   									 |蓝色 |
| 4    | DEBUG | 调试消息                                    |绿色|

####	Goland IDE的效果

![1584630606706](assets/1584630606706.png)

**如:** 在`goland` IDE 中点击 `main.go:10` 可快速`跳转`到打印消息的行数



#### Linux 终端的效果

终端的效果是不同的 , 因为配色方案的结果 , 终端代码的行数无法进行跳转 , 默认颜色为`浅蓝色`**+** `下划线 `

![1584632774672](assets/1584632774672.png)



##	简单使用:

```go
package test
import "github.com/eliot-jay/logger"

func main () {
  
    logger:=logger.DefaultLogger(false) //true为开启文件记录
	logger.DEBUG("This's debug message")
	logger.INFO("This's info message")
	logger.WARN("This's warn message")
	logger.ERROR("this's error message")
	logger.SERIOUS("this's serious message")
  
}
```

####	json文件中获取配置:

```go
package main

import "github.com/eliot-jay/logger"

func main() {
	logger,err:=logger.NewLogByJsonFile("./.json")
	if err!=nil{
		panic(err)
	}
	logger.DEBUG("This's debug message")
	logger.INFO("This's info message")
	logger.WARN("This's warn message")
	logger.ERROR("this's error message")
	logger.SERIOUS("this's serious message")
  
}
```

#### json文件的配置

```json
{
  "logger": {
    "file_name": "app.",   //保存日志到磁盘的文件名
    "file_cording": true,  //是否开启文件记录
    "level": "DEBUG",	   //日志的等级
    "identifier": "$",     //打印消息的标示符
    "time_format": "2006-01-02 15:04:05"   //你喜欢的时间格式
  }
}

```

## 接收打印时的消息

```go
package main

import "github.com/eliot-jay/logger"
import "fmt"

func main()  {
	logger := logger.DefaultLogger(false)
  //上传到 云 或者数据库的 接口
	logger.ReceiveLog(func( string) {
		fmt.Println("receive: ",)
	})
	
	logger.DEBUG("hello debug")
	logger.INFO("hello info")
	logger.ERROR("hello error")
	logger.WARN("hello warn")
	logger.SERIOUS("hello serious")
}

```

##	支持的时间格式

```
ANSIC           "Mon Jan _2 15:04:05 2006"
UnixDate      	"Mon Jan _2 15:04:05 MST 2006"
RubyDate     		"Mon Jan 02 15:04:05 -0700 2006"
RFC822     		  "02 Jan 06 15:04 MST"
RFC822Z         "02 Jan 06 15:04 -0700"
RFC850      		"Monday, 02-Jan-06 15:04:05 MST"
RFC1123         "Mon, 02 Jan 2006 15:04:05 MST"
RFC1123Z        "Mon, 02 Jan 2006 15:04:05 -0700"
RFC3339         "2006-01-02T15:04:05Z07:00"
RFC3339Nano     "2006-01-02T15:04:05.999999999Z07:00"
Kitchen         "3:04PM"
Stamp        		"Jan _2 15:04:05"
StampMilli    	"Jan _2 15:04:05.000"
StampMicro    	"Jan _2 15:04:05.000000"
StampNano     	"Jan _2 15:04:05.000000000"
RFC3339Nano1   	"2006-01-02 15:04:05.999999999 -0700 MST"
DEFAULT         "2006-01-02 15:04:05"
```
