# Zombie 

> 一个无聊的弱口令爆破工具


## 设计初衷
* 日常红队过程中普遍使用超级弱口令检测工具和hydra

* 与之相比的区别:
    1. 命令行工具,且是单文件版本
    2. 由golang编写,可以支持多平台
    3. 体积较小,可以上传到目标上使用(很多时候代理是真的傻逼)

* 在爆破到密码后会停止,爆破,然后用现在的用户名进行一次随机密码的爆破,如果依然成功,则会判断为低交互蜜罐

## 基础用法

###  Brute 模块
> 即基础爆破模块

* 参数
    * ip:
        1. 例如:127.0.0.1:3306,支持以逗号分隔传入多个ip
    
    * username or u
        1. 支持以逗号输入传入多个要爆破的用户名
    
    * password or p
        1. 支持以逗号输入传入多个要爆破的用户名
    
    * server or s 
        * 传入要爆破的服务(暂时一轮只能爆破一种)

* example
    在使用默认端口的时候可以不写端口或者不写服务名
    * `Brute -u admin,admin321,admin888,root,postgres -p aaaa,bbbb,ccc,cccd,ddd,321 -ip 127.0.0.1:6379 -s redis`
    * `Brute -u admin,admin321,admin888,root,postgres -p aaaa,bbbb,ccc,cccd,ddd,321 -ip 127.0.0.1:6379 `
    * `Brute -u admin,admin321,admin888,root,postgres -p aaaa,bbbb,ccc,cccd,ddd,321 -ip 127.0.0.1 -s redis`
    * -U,-P,-IP 则是读取文件模式,(使用绝对路径)
    * `Brute -U user.dic -P pass.dic -f log.txt -s tomcat -IP ip.dic`
* 目前支持的协议
  {
  21:    "FTP",
  22:    "SSH",
  445:   "SMB",
  1433:  "MSSQL",
  3306:  "MYSQL",
  5432:  "POSTGRESQL",
  6379:  "REDIS",
  9200:  "ELASTICSEARCH",
  27017: "MONGO",
  5900:  "VNC",
  8080: "TOMCAT"
  }

### Exec 模块

* 参数
    * ip:
        1. 例如:127.0.0.1:3306
    
    * username or u
      
    * password or p
      
    * server or s 
        * 传入要执行命令的服务 :仅支持 mysql 
    
*  example
    Exec -u root -p test -ip 127.0.0.1:3306 -q "show tables"



   



