# Zombie 

一个轻量级的服务口令爆破工具, 继承了hydra的命令行设计, hashcat的字典生成, 以及红队向的功能设计. 

## QuickStart

完整文档位于: https://chainreactors.github.io/wiki/zombie/

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

`zombie -l 1.txt -p admin --weakpass`

将会根据admin关键字生成常见的密码组合, 以admin为例， 将会生成以下密码

<details>
  <summary>--weakpass 生成的密码</summary>

```
admin
Admin
ADMIN
aDMIN
admin1
admin2
admin3
admin4
admin5
admin6
admin7
admin8
admin9
admin0
admin123
admin1234
admin12345
admin123456
admin2018
admin2019
admin2020
admin2021
admin2022
admin01
admin02
admin03
admin04
admin05
admin06
admin07
admin08
admin09
admin10
admin11
admin12
admin13
admin14
admin15
admin16
admin17
admin18
admin19
admin20
admin21
admin22
admin23
admin24
admin25
admin26
admin27
admin28
admin29
admin30
admin31
admin!
admin@
admin#
admin$
admin!@
admin!@#
admin!@#$
admin123!
admin!123
admin1@
admin2018!
admin2019!
admin2020!
admin2021!
admin2022!
admin!2018
admin!2019
admin!2020
admin!2021
admin!2022
admin2018!@#
admin2019!@#
admin2020!@#
admin2021!@#
admin2022!@#
admin01!
admin02!
admin03!
admin04!
admin05!
admin06!
admin07!
admin08!
admin09!
admin10!
admin11!
admin12!
admin13!
admin14!
admin15!
admin16!
admin17!
admin18!
admin19!
admin20!
admin21!
admin22!
admin23!
admin24!
admin25!
admin26!
admin27!
admin28!
admin29!
admin30!
admin31!
Admin1
Admin2
Admin3
Admin4
Admin5
Admin6
Admin7
Admin8
Admin9
Admin0
Admin123
Admin1234
Admin12345
Admin123456
Admin2018
Admin2019
Admin2020
Admin2021
Admin2022
Admin!
Admin@
Admin#
Admin$
Admin!@
Admin!@#
Admin!@#$
Admin123!
Admin!123
Admin1@
Admin2018!
Admin2019!
Admin2020!
Admin2021!
Admin2022!
Admin!2018
Admin!2019
Admin!2020
Admin!2021
Admin!2022
Admin2018!@#
Admin2019!@#
Admin2020!@#
Admin2021!@#
Admin2022!@#
Admin01!
Admin02!
Admin03!
Admin04!
Admin05!
Admin06!
Admin07!
Admin08!
Admin09!
Admin10!
Admin11!
Admin12!
Admin13!
Admin14!
Admin15!
Admin16!
Admin17!
Admin18!
Admin19!
Admin20!
Admin21!
Admin22!
Admin23!
Admin24!
Admin25!
Admin26!
Admin27!
Admin28!
Admin29!
Admin30!
Admin31!
Admin01
Admin02
Admin03
Admin04
Admin05
Admin06
Admin07
Admin08
Admin09
Admin10
Admin11
Admin12
Admin13
Admin14
Admin15
Admin16
Admin17
Admin18
Admin19
Admin20
Admin21
Admin22
Admin23
Admin24
Admin25
Admin26
Admin27
Admin28
Admin29
Admin30
Admin31
```

</details>


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
* SOCKS5
* HTTP 401
* POP3
* SOCKS5

### 通过neutron template支持的服务

* apollo
* canal
* hikvision_camer
* rabbitmq
* ruijie_ap
* zte-epon
* h3c_router
* minio
* snmp_manage
* nacos
* gitlab
* huawei_ibmc
* weblogic
* activemq
* boda
* dubbo
* grafana
* jenkins
* tomcat
* apisix
* druid
* nexus
* xxl-job



### TODO

- [ ] 支持基本信息收集
- [ ] 支持基本的后利用(希望能像cme一样)
- [ ] 支持更多协议
  * RSTP
  * HTTP PROXY
  * rlogin
  * RSYNC
  * zookeeper
  * memcache
  * amqp
  * mqtt
  * http
- [x] 支持neutron引擎, 允许通过模板配置插件
- [ ] 新增密码策略限制的功能, 减少爆破次数
- [ ] 新增爆破限制的功能, 防止被封禁

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