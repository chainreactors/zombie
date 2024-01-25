# Zombie 

一个轻量级的服务口令爆破工具, 继承了hydra的命令行设计, hashcat的字典生成, 以及红队向的功能设计. 

## QuickStart

参考了hydra的命令行设计, 小写为命令行输出, 大写为文件输入, 留空为使用默认值.

使用默认字典爆破ssh的root用户口令

`zombie -i 1.1.1.1 -u root -s ssh`

使用指定的密码批量喷洒ssh口令

`zombie -I targets.txt -u root -p password -s ssh`

targets.txt
```
1.1.1.1
2.2.2.2
3.3.3.3
...
```

从文件中自动解析输入

`zombie -I targets.txt`

targets.txt:
```
mysql://user:pass@1.1.1.1:3307  # 指定了用户与密码以及端口, 尝试登录mysql
ssh://user@2.2.2.2              # 自动解析ssh默认端口22, 使用默认密码爆破指定user的ssh
mssql://3.3.3.3:1433            # 未指定user与pass, 自动选用默认的用户与密码字典
```

使用已知的所有用户与密码,  进行笛卡尔积的方式对服务进行最大可能的爆破.

`zombie -I targets.txt -U user.txt -P pass.txt`

targets.txt:
```
mysql://1.1.1.1
ssh://2.2.2.2
mssql://3.3.3.3
```

从gogo结果开始扫描

`zombie --gogo 1.dat`

从json开始扫描

`zombie -j 1.json`


简单配置自定义密码生成器

`zombie -l 1.txt -p google --weakpass`

将会根据google关键字生成常见的密码组合, 以google为例， 将会生成以下密码

```
google
Google
GOOGLE
gOOGLE
google1
google2
google3
google4
google5
google6
google7
google8
google9
google0
google123
google1234
google12345
google123456
google2018
google2019
google2020
google2021
google2022
...
```

`--weakpass` 的规则位于 https://github.com/chainreactors/templates/blob/master/zombie/rule/weakpass.rule , 欢迎提供新规则



## Usage

当前支持的协议

* FTP
* SSH
* SMB
* MSSQL
* MYSQL
* Mongo
* POSTGRESQL
* REDIS
* ORACLE
* VNC
* LDAP
* SNMP
* RDP 
* HTTP/HTTPS 401 
  * kibana
  * tomcat
* SOCKS5
* TELNET
* RSYNC
* telnet
* HTTP 401
* POP3
* SOCKS5


### TODO

* 支持基本信息收集
* 支持基本的后利用(希望能像cme一样)
* 支持更多协议
  * RSTP
  * HTTP PROXY
  * rlogin
  * zookeeper
  * memcache
  * amqp
  * mqtt
  * http
    * rabbitmq management
    * solr
    * webdav
    * ...
  * web
    * weblogic
    * websphere
    * jenkins
    * grafana
    * zabbix
    * ...
* 支持neutron引擎, 允许通过模板配置插件
* 新增密码策略限制的功能, 减少爆破次数
* 新增爆破限制的功能, 防止被封禁
## Make

```bash
# download
git clone --recurse-submodules https://github.com/chainreactors/zombie
cd zombie

# sync dependency
go mod tidy   

# generate template.go
go generate  

# build 
go build .
```