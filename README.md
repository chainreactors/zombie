# Zombie 

一个轻量级的服务口令爆破工具, 继承了hydra的命令行设计, hashcat的字典生成, 以及红队向的功能设计. 

## QuickStart

参考了hydra的命令行设计, 小写为命令行输出, 大写为文件输入, 留空为使用默认值.

使用默认字典爆破ssh口令

`zombie -I targets.txt -u root -s ssh`

打开debug,判断是否存在漏报

`zombie -I targets.txt -u root -p password123 -s ssh --debug`

从gogo中导入

`gogo -F 1.dat -o zombie -f zombie.json`

`zombie --go zombie.json`


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