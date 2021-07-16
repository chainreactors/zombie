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
    
    * UA 参数
        1. 支持以键值对文本的方式读入用户名密码
    
    * server or s 
        * 传入要爆破的服务(暂时一轮只能爆破一种)
      
    * thread or t
        * 爆破的协程池大小
      
    * timeout 
        * 超时时间 记得是 --timeout 默认是2
    * *simple or e*
        * 在爆用户名密码对比较少,目标少且字典大的情况,加上-e参数（即新增是默认开启喷射模式）
    * proc
        * 进度条现在默认关闭，如果要开始则加上-proc 数字（多少次显示一次）
    * file or f
        * 输出的文件，默认为当前目录下的res.log，如果不想生成文件则输入-f null
    * type
        * 调整输出文件的格式，默认为raw，若要json，则输入-type json
    * proc 
        * 默认不显示进度条，可以--proc int 来设置爆破多少次显示一次
    * uppair or UP
        * 支持出入键值对的用户名密码
    

* example
    在使用默认端口的时候可以不写端口或者不写服务名
    * `Brute -u admin,admin321,admin888,root,postgres -p aaaa,bbbb,ccc,cccd,ddd,321 -ip 127.0.0.1:6379 -s redis`
    * `Brute -u admin,admin321,admin888,root,postgres -p aaaa,bbbb,ccc,cccd,ddd,321 -ip 127.0.0.1:6379 `
    * `Brute -u admin,admin321,admin888,root,postgres -p aaaa,bbbb,ccc,cccd,ddd,321 -ip 127.0.0.1 -s redis`
    * -U,-P,-IP 则是读取文件模式,(使用绝对路径)
    * `Brute -U user.dic -P pass.dic -f log.txt -s tomcat -IP ip.dic`
* 目前支持的协议
    * 21:    "FTP",
    * 22:    "SSH",
    * 445:   "SMB",
    * 1433:  "MSSQL",
    * 3306:  "MYSQL",
    * 5432:  "POSTGRESQL",
    * 6379:  "REDIS",  
    * 9200:  "ELASTICSEARCH",
    * 27017: "MONGO",
    * 5900:  "VNC",
    * 8080: "TOMCAT"
    
    

* redis支持爆破成功后系统检测和linux 是否为root权限

* 加入特殊解析如果password中带有%user%的字符串,则会替换为用户名,便于减少密码字典的数量

* SMB支持哈希传递,密码字段加上前缀 hash: , 且无论是否成功都会返回winodws版本

* 支持IP段扫描,如果端口自定义则为 127.0.0.1/24:445 的形式

* 新增内置字典，只需要输入IP和爆破类型就会使用内置字典

TODO：

    [+] ssh批量执行命令，ssh私钥喷射
    [+] 自定义get，post or yaml解析
    
    

### Query 模块

* 参数 
  
    * --username value, -u value   
    * --password value, -p value   
    * --ip value                   
  * --input value, -i value      
  * --InputFile value, -F value  
  * --server value, -s value     
  * --auto, -a                   (default: false)
  * --help, -h                   show help (default: false)

    
*  example
    Query -u root -p test -ip 127.0.0.1:3306 -q "show tables"


* 目前auto支持postgresql,mysql,mssql,(收集操作系统,数据总量,和一些敏感信息)且可以通过-F参数将zombie爆破结果至今放入自动化收集

   
* auto模式下会对扫描结果进行去重

### Decrypt 模块

* 目前支持Xshell，Xftp，低于7版本的解密
* 支持Navicat全版本解密（最新版测试通过）
  * all, a
    * 同时运行两款解密
  * Navicat, N
    > 可以手动输入来解密，如果不输入，则自己从注册表读取数据
    * --cipher value, -c value
    * --OutputFile value, -f value  (default: "./DeRes.log")

  * Xshell, X
    > 默认从用户目录下读取判断版本，也可以自己输入数据来解密
    * --cipher value, -c value
    * --username value, -u value
    * --sid value, -s value
    * --version value, -v value     (default: 0)
    * --OutputFile value, -f value  (default: "./DeRes.log")
    

  

