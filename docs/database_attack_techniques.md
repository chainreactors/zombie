# 数据库攻击手法综合参考文档

> 本文档基于对15+款主流数据库攻击工具的深度调研，系统整理了所有已武器化的数据库攻击手法。
> 
> 调研工具：MDUT、MDUT-Extend、ODAT、MSDAT、PowerUpSQL、SQLRecon、NoSQLMap、RedisEXP、Databasetools、sqlmap、SharpSQLTools、enumdb、Sylas
> 
> 覆盖数据库：MySQL、Microsoft SQL Server、Oracle、PostgreSQL、Redis、MongoDB、CouchDB、Cassandra
> 
> 文档生成时间：2025年7月

---

## 摘要

| 数据库类型 | 已武器化攻击手法数量 |
|-----------|-------------------|
| MySQL | 8 |
| Microsoft SQL Server | 48+ |
| Oracle | 37+ |
| PostgreSQL | 6 |
| Redis | 16 |
| MongoDB | 7 |
| CouchDB | 4 |
| Cassandra | 0（计划中） |
| 多数据库/通用 | 20+ |

**总计：140+ 种已武器化攻击手法**

---

## 一、MySQL 攻击手法

> 涉及工具：MDUT、MDUT-Extend、sqlmap、enumdb

### 1.1 UDF (User-Defined Function) 提权

| 属性 | 描述 |
|------|------|
| **描述** | 通过向MySQL插件目录写入自定义UDF库文件，创建`sys_eval`函数来执行系统命令。支持Windows 32/64位和Linux 32/64位平台，根据目标系统自动选择对应的UDF库文件。 |
| **目标数据库** | MySQL (全版本) |
| **目标平台** | Windows 32/64、Linux 32/64 |
| **涉及工具** | MDUT |
| **武器化状态** | **已实现** |
| **相关SQL** | `SELECT [hex] INTO DUMPFILE '[plugin_path]'; CREATE FUNCTION sys_eval RETURNS STRING SONAME '[udf_file]';` |
| **出处** | https://github.com/SafeGroceryStore/MDUT/blob/main/MDAT-DEV/src/main/java/Dao/MysqlDao.java |

### 1.2 Windows 反弹 Shell (backShell)

| 属性 | 描述 |
|------|------|
| **描述** | 通过加载专门的反弹Shell UDF插件(`udf_win_ex_hex.txt`)，创建`backshell`函数，实现从目标MySQL服务器向攻击者指定IP和端口反弹Windows Shell。 |
| **目标数据库** | MySQL |
| **目标平台** | Windows（仅Windows平台支持） |
| **涉及工具** | MDUT |
| **武器化状态** | **已实现** |
| **相关SQL** | `SELECT backshell('%s','%s') AS s;` |
| **出处** | https://github.com/SafeGroceryStore/MDUT/blob/main/MDAT-DEV/src/main/java/Dao/MysqlDao.java |

### 1.3 NTFS 提权 (ADS Alternate Data Streams)

| 属性 | 描述 |
|------|------|
| **描述** | 利用NTFS文件系统的备用数据流(Alternate Data Streams)特性创建目录，绕过某些安全限制。通过向`::$INDEX_ALLOCATION`流写入数据来创建特殊目录。 |
| **目标数据库** | MySQL |
| **目标平台** | Windows（依赖NTFS文件系统） |
| **涉及工具** | MDUT |
| **武器化状态** | **已实现** |
| **相关SQL** | `SELECT '1' INTO DUMPFILE '%s::$INDEX_ALLOCATION'` |
| **出处** | https://github.com/SafeGroceryStore/MDUT/blob/main/MDAT-DEV/src/main/java/Dao/MysqlDao.java |

### 1.4 系统命令执行 (UDF Eval)

| 属性 | 描述 |
|------|------|
| **描述** | 通过已创建的UDF函数`sys_eval`执行任意系统命令，支持UTF-8、GB2312、GBK编码选择。sqlmap通过注入自定义UDF `sys_exec()`/`sys_eval()` 执行命令。 |
| **目标数据库** | MySQL |
| **涉及工具** | MDUT、sqlmap |
| **武器化状态** | **已实现** |
| **出处(MDUT)** | https://github.com/SafeGroceryStore/MDUT/blob/main/MDAT-DEV/src/main/java/Dao/MysqlDao.java |
| **出处(sqlmap)** | https://github.com/sqlmapproject/sqlmap/wiki/Features |

### 1.5 HTTP 隧道连接

| 属性 | 描述 |
|------|------|
| **描述** | 支持通过HTTP隧道连接MySQL数据库，用于绕过网络限制。将所有数据库操作封装为HTTP请求，通过中间HTTP隧道代理与目标数据库通信。 |
| **目标数据库** | MySQL、MSSQL、Oracle、PostgreSQL（Redis不支持） |
| **涉及工具** | MDUT |
| **武器化状态** | **已实现** |
| **出处** | https://github.com/SafeGroceryStore/MDUT/blob/main/MDAT-DEV/src/main/java/Dao/MysqlHttpDao.java |

### 1.6 SQL注入检测与利用 (sqlmap)

| 属性 | 描述 |
|------|------|
| **描述** | sqlmap支持对MySQL的完整SQL注入检测与利用，包括布尔盲注、时间盲注、报错注入、UNION注入、堆叠查询等5种核心注入技术；支持数据库指纹识别、用户/密码哈希/权限/角色枚举、全库数据Dump、文件读写、OS命令执行等。 |
| **目标数据库** | MySQL (全版本) |
| **涉及工具** | sqlmap |
| **武器化状态** | **已实现** |
| **出处** | https://github.com/sqlmapproject/sqlmap/wiki/Features |

### 1.7 凭据暴力破解与数据枚举 (enumdb)

| 属性 | 描述 |
|------|------|
| **描述** | 对MySQL进行凭据暴力破解（支持CIDR范围）、数据库/表/列枚举、关键词敏感数据搜索、数据Dump提取为CSV/XLSX、交互式SQL Shell。 |
| **目标数据库** | MySQL |
| **涉及工具** | enumdb |
| **武器化状态** | **已实现** |
| **出处** | https://github.com/m8sec/enumdb/blob/master/README.md |

### 1.8 MySQL驱动高低版本切换

| 属性 | 描述 |
|------|------|
| **描述** | 新增MySQL驱动高低版本切换选项，以兼容不同版本的MySQL数据库。 |
| **目标数据库** | MySQL |
| **涉及工具** | MDUT-Extend |
| **武器化状态** | **已实现** |
| **出处** | https://github.com/DeEpinGh0st/MDUT-Extend-Release/releases/tag/v1.0.0 |

---

## 二、Microsoft SQL Server 攻击手法

> 涉及工具：MDUT、MDUT-Extend、MSDAT、PowerUpSQL、SQLRecon、SharpSQLTools、enumdb、Sylas、sqlmap

### 2.1 发现与枚举类

#### 2.1.1 无认证信息获取

| 属性 | 描述 |
|------|------|
| **描述** | 利用TDS协议和SQL Browser Service，在无需认证的情况下获取远程MSSQL服务器的技术信息，包括数据库版本、实例名等。SQLRecon通过UDP 1434端口查询。 |
| **涉及工具** | MSDAT、SQLRecon |
| **武器化状态** | **已实现** |
| **出处(MSDAT)** | https://github.com/quentinhardy/msdat/blob/master/README.md |
| **出处(SQLRecon)** | https://github.com/skahwah/SQLRecon/blob/main/README.md |

#### 2.1.2 本地SQL Server实例发现

| 属性 | 描述 |
|------|------|
| **描述** | 通过注册表搜索发现本地系统上运行的SQL Server实例。 |
| **涉及工具** | PowerUpSQL |
| **武器化状态** | **已实现** |
| **对应函数** | `Get-SQLInstanceLocal` |
| **出处** | https://github.com/NetSPI/PowerUpSQL/wiki/Discovery-Functions |

#### 2.1.3 域内SQL Server发现 (SPN枚举)

| 属性 | 描述 |
|------|------|
| **描述** | 通过向域控制器查询注册MSSQL服务主体名称(SPN)来发现域内所有SQL Server实例。支持多线程UDP扫描。 |
| **涉及工具** | PowerUpSQL、SQLRecon |
| **武器化状态** | **已实现** |
| **对应函数/模块** | PowerUpSQL: `Get-SQLInstanceDomain`; SQLRecon: `SqlSpns` |
| **出处(PowerUpSQL)** | https://github.com/NetSPI/PowerUpSQL/wiki/Discovery-Functions |
| **出处(SQLRecon)** | https://github.com/skahwah/SQLRecon/blob/main/README.md |

#### 2.1.4 UDP广播/端口扫描发现

| 属性 | 描述 |
|------|------|
| **描述** | 通过向子网广播地址发送UDP请求或UDP端口扫描来发现本地网络上的SQL Server实例，支持多线程。 |
| **涉及工具** | PowerUpSQL |
| **武器化状态** | **已实现** |
| **对应函数** | `Get-SQLInstanceBroadcast`、`Get-SQLInstanceScanUDPThreaded` |
| **出处** | https://github.com/NetSPI/PowerUpSQL/wiki/Discovery-Functions |

#### 2.1.5 数据库配置与权限枚举

| 属性 | 描述 |
|------|------|
| **描述** | 获取数据库版本、数据库列表、用户列表、禁用的用户、存储过程、当前用户角色和权限等配置信息。 |
| **涉及工具** | MSDAT、PowerUpSQL、SQLRecon |
| **武器化状态** | **已实现** |
| **对应模块** | MSDAT: `search`; PowerUpSQL: `Get-SQLServerInfo`、`Invoke-SQLAudit`; SQLRecon: `Whoami`、`Databases`、`Users` |
| **出处** | https://github.com/quentinhardy/msdat、https://github.com/NetSPI/PowerUpSQL、https://github.com/skahwah/SQLRecon |

#### 2.1.6 敏感数据搜索

| 属性 | 描述 |
|------|------|
| **描述** | 在数据库表的列名中搜索敏感数据模式（如password、credential、信用卡号、SSN等），支持显示空列和获取样本数据。 |
| **涉及工具** | MSDAT、PowerUpSQL、SQLRecon、enumdb |
| **武器化状态** | **已实现** |
| **出处** | 多工具支持 |

#### 2.1.7 文件/目录枚举

| 属性 | 描述 |
|------|------|
| **描述** | 利用`xp_subdirs`和`xp_dirtree`存储过程列出指定目录的文件和子目录；利用`xp_fixeddrives`和`xp_availablemedia`枚举磁盘驱动器。 |
| **涉及工具** | MSDAT |
| **武器化状态** | **已实现** |
| **对应模块** | `xpdirectory` |
| **出处** | https://github.com/quentinhardy/msdat/blob/master/README.md |

#### 2.1.8 Schema/表数据转储

| 属性 | 描述 |
|------|------|
| **描述** | 提取数据库的完整Schema信息和所有表数据并保存到文件（CSV/XLSX格式），排除默认数据库。 |
| **涉及工具** | MSDAT、enumdb |
| **武器化状态** | **已实现** |
| **出处** | https://github.com/quentinhardy/msdat、https://github.com/m8sec/enumdb |

#### 2.1.9 链接服务器枚举

| 属性 | 描述 |
|------|------|
| **描述** | 枚举SQL Server上配置的链接服务器(Linked Server)。 |
| **涉及工具** | PowerUpSQL、SQLRecon |
| **武器化状态** | **已实现** |
| **对应函数/模块** | PowerUpSQL: `Get-SQLServerLink`; SQLRecon: `Links`、`CheckRpc` |
| **出处** | https://github.com/NetSPI/PowerUpSQL、https://github.com/skahwah/SQLRecon |

### 2.2 命令执行类

#### 2.2.1 xp_cmdshell 命令执行

| 属性 | 描述 |
|------|------|
| **描述** | 通过xp_cmdshell存储过程在数据库服务器上执行操作系统命令。支持自动启用/禁用xp_cmdshell，支持交互式Shell模式。sqlmap、SharpSQLTools、SQLRecon、Sylas等工具均实现此功能。 |
| **涉及工具** | MDUT、MSDAT、PowerUpSQL、SQLRecon、SharpSQLTools、Sylas、sqlmap |
| **武器化状态** | **已实现** |
| **出处汇总** | MDUT: https://github.com/SafeGroceryStore/MDUT; MSDAT: https://github.com/quentinhardy/msdat; PowerUpSQL: https://github.com/NetSPI/PowerUpSQL; SQLRecon: https://github.com/skahwah/SQLRecon; SharpSQLTools: https://github.com/uknowsec/SharpSQLTools; Sylas: https://github.com/Ryze-T/Sylas; sqlmap: https://github.com/sqlmapproject/sqlmap |

#### 2.2.2 sp_oacreate / OLE Automation 命令执行

| 属性 | 描述 |
|------|------|
| **描述** | 利用`sp_oacreate`创建WScript.Shell COM组件，通过OLE Automation执行系统命令。相比xp_cmdshell更为隐蔽。支持PowerShell反向Shell和文件上传下载。 |
| **涉及工具** | MDUT、MSDAT、PowerUpSQL、SQLRecon、SharpSQLTools、Sylas |
| **武器化状态** | **已实现** |
| **出处汇总** | MDUT: https://github.com/SafeGroceryStore/MDUT; MSDAT: https://github.com/quentinhardy/msdat; PowerUpSQL: https://github.com/NetSPI/PowerUpSQL; SQLRecon: https://github.com/skahwah/SQLRecon; SharpSQLTools: https://github.com/uknowsec/SharpSQLTools; Sylas: https://github.com/Ryze-T/Sylas |

#### 2.2.3 CLR (Common Language Runtime) 程序集执行

| 属性 | 描述 |
|------|------|
| **描述** | 通过MSSQL的CLR集成功能，加载.NET程序集实现代码执行和提权。支持一键激活/恢复CLR组件，可从网络路径加载DLL。SharpSQLTools的CLR模块提供最丰富的功能（文件管理、进程查看、Potato提权、LSASS转储、RDP启用、Shellcode加载等）。 |
| **涉及工具** | MDUT、PowerUpSQL、SQLRecon、SharpSQLTools、Sylas |
| **武器化状态** | **已实现** |
| **出处汇总** | MDUT: https://github.com/SafeGroceryStore/MDUT; PowerUpSQL: https://github.com/NetSPI/PowerUpSQL/wiki/Attacking-SQL-Server-CLR; SQLRecon: https://github.com/skahwah/SQLRecon; SharpSQLTools: https://github.com/uknowsec/SharpSQLTools; Sylas: https://github.com/Ryze-T/Sylas |

#### 2.2.4 SQL Agent Job 命令执行

| 属性 | 描述 |
|------|------|
| **描述** | 利用SQL Server Agent Jobs执行系统命令，支持CMDExec、PowerShell、ActiveX:JScript和ActiveX:VBScript子系统。支持反向Shell功能，可列出和查看Agent Jobs代码。 |
| **涉及工具** | MSDAT、PowerUpSQL、SQLRecon |
| **武器化状态** | **已实现** |
| **出处汇总** | MSDAT: https://github.com/quentinhardy/msdat; PowerUpSQL: https://github.com/NetSPI/PowerUpSQL; SQLRecon: https://github.com/skahwah/SQLRecon |

#### 2.2.5 R/Python 外部脚本命令执行

| 属性 | 描述 |
|------|------|
| **描述** | 通过SQL Server 2016+的外部脚本功能使用R语言或Python执行OS命令。不需要从磁盘读取DLL。 |
| **涉及工具** | PowerUpSQL |
| **武器化状态** | **已实现** |
| **对应函数** | `Invoke-SQLOSCmdR`、`Invoke-SQLOSCmdPython` |
| **目标版本** | MSSQL 2016+ |
| **出处** | https://github.com/NetSPI/PowerUpSQL/wiki/Primary-Attack-Functions |

#### 2.2.6 MSSQL Shellcode 加载

| 属性 | 描述 |
|------|------|
| **描述** | 通过MSSQL加载和执行Shellcode的功能，可直接加载自定义Shellcode到内存执行。SharpSQLTools使用QueueUserAPC注入技术远程加载x64 Shellcode。 |
| **涉及工具** | MDUT-Extend、SharpSQLTools |
| **武器化状态** | **已实现** |
| **出处** | MDUT-Extend: https://github.com/DeEpinGh0st/MDUT-Extend-Release/releases/tag/v1.3.0; SharpSQLTools: https://github.com/uknowsec/SharpSQLTools |

#### 2.2.7 自定义扩展存储过程

| 属性 | 描述 |
|------|------|
| **描述** | 创建自定义DLL扩展存储过程来执行OS命令。生成DLL文件并注册为XP。 |
| **涉及工具** | PowerUpSQL |
| **武器化状态** | **已实现** |
| **对应函数** | `Create-SQLFileXpDll` |
| **出处** | https://github.com/NetSPI/PowerUpSQL/wiki/PowerUpSQL-Cheat-Sheet |

### 2.3 权限提升类

#### 2.3.1 Impersonation 权限提升

| 属性 | 描述 |
|------|------|
| **描述** | 利用不安全的Impersonation配置获取权限提升。SQLRecon支持在所有支持Impersonation的模块中模拟指定用户执行操作。 |
| **涉及工具** | PowerUpSQL、SQLRecon |
| **武器化状态** | **已实现** |
| **出处(PowerUpSQL)** | https://github.com/NetSPI/PowerUpSQL/wiki/Attacking-Insecure-Impersonation-Configurations |
| **出处(SQLRecon)** | https://github.com/skahwah/SQLRecon/wiki/4.-Impersonation-Modules |

#### 2.3.2 Trustworthy 数据库权限提升

| 属性 | 描述 |
|------|------|
| **描述** | 利用标记为TRUSTWORTHY的数据库进行权限提升，从普通数据库用户获取sysadmin权限。 |
| **涉及工具** | MSDAT、PowerUpSQL |
| **武器化状态** | **已实现** |
| **出处(MSDAT)** | https://github.com/quentinhardy/msdat |
| **出处(PowerUpSQL)** | https://github.com/NetSPI/PowerUpSQL/wiki/Attacking-Trustworthy-Databases |

#### 2.3.3 本地管理员到 sysadmin

| 属性 | 描述 |
|------|------|
| **描述** | 利用本地管理员权限模拟SQL Server服务账户，从而获得sysadmin权限。 |
| **涉及工具** | PowerUpSQL |
| **武器化状态** | **已实现** |
| **对应函数** | `Invoke-SQLImpersonateService` |
| **出处** | https://github.com/NetSPI/PowerUpSQL/wiki/Primary-Attack-Functions |

#### 2.3.4 GodPotato 提权 (EfsPotato/BadPotato)

| 属性 | 描述 |
|------|------|
| **描述** | 利用Windows NTLM中继和SeImpersonatePrivilege特权实现从服务账户到SYSTEM的提权。GodPotato相比原版JuicyPotato更稳定且支持更多Windows版本。SharpSQLTools支持EfsPotato和BadPotato两种提权方式。 |
| **涉及工具** | MDUT-Extend、SharpSQLTools |
| **武器化状态** | **已实现** |
| **目标平台** | Windows (需SeImpersonatePrivilege) |
| **出处** | MDUT-Extend: https://github.com/DeEpinGh0st/MDUT-Extend-Release/releases/tag/v1.3.0; SharpSQLTools: https://github.com/uknowsec/SharpSQLTools |

#### 2.3.5 自动权限提升 (综合)

| 属性 | 描述 |
|------|------|
| **描述** | 自动识别配置弱点并尝试获取sysadmin权限。综合多种提权技术（Impersonation、Trustworthy等）。 |
| **涉及工具** | PowerUpSQL |
| **武器化状态** | **已实现** |
| **对应函数** | `Invoke-SQLEscalatePriv` |
| **出处** | https://github.com/NetSPI/PowerUpSQL/wiki/Primary-Attack-Functions |

### 2.4 文件操作类

#### 2.4.1 文件上传/下载

| 属性 | 描述 |
|------|------|
| **描述** | 在MSSQL服务器上进行文件上传和下载。支持OLE Automation、bulkinsert、openrowset等多种方法。SharpSQLTools支持使用OLE Automation上传下载文件。 |
| **涉及工具** | MDUT、MSDAT、SharpSQLTools、Sylas |
| **武器化状态** | **已实现** |
| **出处汇总** | MDUT: https://github.com/SafeGroceryStore/MDUT; MSDAT: https://github.com/quentinhardy/msdat; SharpSQLTools: https://github.com/uknowsec/SharpSQLTools; Sylas: https://github.com/Ryze-T/Sylas |

#### 2.4.2 WebShell 写入 (LOG备份/差异备份)

| 属性 | 描述 |
|------|------|
| **描述** | 通过LOG备份方式或差异备份方式向MSSQL服务器写入WebShell。 |
| **涉及工具** | Sylas |
| **武器化状态** | **已实现** |
| **出处** | https://github.com/Ryze-T/Sylas |

### 2.5 链接服务器与横向移动类

#### 2.5.1 链接服务器爬取 (Linked Server Crawling)

| 属性 | 描述 |
|------|------|
| **描述** | 递归遍历所有可访问的链接服务器路径，枚举SQL Server版本和链接配置权限。支持在链接服务器上执行任意SQL查询（包括xp_cmdshell和xp_dirtree）。可导出Neo4j图形进行可视化。SQLRecon支持链接服务器链攻击（如SQL01 -> SQL02 -> PAYMENTS01）。 |
| **涉及工具** | PowerUpSQL、SQLRecon |
| **武器化状态** | **已实现** |
| **对应函数** | PowerUpSQL: `Get-SQLServerLinkCrawl`; SQLRecon: Linked Chain模块 |
| **出处(PowerUpSQL)** | https://github.com/NetSPI/PowerUpSQL、https://www.netspi.com/blog/technical-blog/network-pentesting/sql-server-link-crawling-powerupsql/ |
| **出处(SQLRecon)** | https://github.com/skahwah/SQLRecon/wiki/6.-Linked-Chain-Modules |

#### 2.5.2 远程SQL请求执行 (Link Hopping)

| 属性 | 描述 |
|------|------|
| **描述** | 通过目标数据库作为跳板，对另一个远程MSSQL服务器执行SQL查询。支持bulkinsert和openrowset方法。 |
| **涉及工具** | MSDAT |
| **武器化状态** | **已实现** |
| **对应模块** | `bulkopen` (`--request-rdb`) |
| **出处** | https://github.com/quentinhardy/msdat |

### 2.6 凭据窃取与密码恢复类

#### 2.6.1 密码哈希提取

| 属性 | 描述 |
|------|------|
| **描述** | 从MSSQL数据库中提取登录密码哈希值。支持所有MSSQL版本（2000/2005/2008/2014/2016/2019）。PowerUpSQL支持通过`-migrate`开关进行本地管理员权限提升后提取。 |
| **涉及工具** | MSDAT、PowerUpSQL |
| **武器化状态** | **已实现** |
| **对应模块/函数** | MSDAT: `passwordstealer`; PowerUpSQL: `Get-SQLServerPasswordHash` |
| **出处** | https://github.com/quentinhardy/msdat、https://github.com/NetSPI/PowerUpSQL |

#### 2.6.2 SMB认证捕获 / UNC路径注入

| 属性 | 描述 |
|------|------|
| **描述** | 通过诱导数据库服务器连接攻击者控制的SMB共享/UNC路径，捕获NetNTLMv2哈希。支持多种触发方式：bulkinsert、openrowset、xp_dirtree、xp_fileexist、xp_getfiledetails。PowerUpSQL的`Invoke-SQLUncPathInjection`可自动完成从SPN查询到哈希捕获的完整流程。 |
| **涉及工具** | MSDAT、PowerUpSQL、SQLRecon |
| **武器化状态** | **已实现** |
| **出处汇总** | MSDAT: https://github.com/quentinhardy/msdat; PowerUpSQL: https://github.com/NetSPI/PowerUpSQL/wiki/Password-Recovery-Functions; SQLRecon: https://github.com/skahwah/SQLRecon |

#### 2.6.3 Windows自动登录密码获取

| 属性 | 描述 |
|------|------|
| **描述** | 通过`xp_regread`从注册表中获取Windows自动登录密码。 |
| **涉及工具** | PowerUpSQL |
| **武器化状态** | **已实现** |
| **对应函数** | `Get-SQLRecoverPwAutoLogon` |
| **出处** | https://github.com/NetSPI/PowerUpSQL/wiki/Password-Recovery-Functions |

#### 2.6.4 Agent Job凭据劫持

| 属性 | 描述 |
|------|------|
| **描述** | 利用Agent Jobs进行域权限提升，劫持SQL Server凭据。 |
| **涉及工具** | PowerUpSQL |
| **武器化状态** | **已实现** |
| **出处** | https://github.com/NetSPI/PowerUpSQL/wiki/Hijacking-SQL-Server-Credentials-using-Agent-Jobs-for-Domain-Privilege-Escalation |

#### 2.6.5 ADSI凭据获取

| 属性 | 描述 |
|------|------|
| **描述** | 从链接的ADSI服务器获取明文凭据。通过上传自定义CLR assembly包含LDAP服务器到SQL Server运行时，利用Agent Jobs将ADSI连接凭据导向本地LDAP服务器获取明文。 |
| **涉及工具** | SQLRecon |
| **武器化状态** | **已实现** |
| **对应模块** | `Adsi` |
| **出处** | https://github.com/skahwah/SQLRecon/blob/main/README.md |
| **技术博客** | https://www.ibm.com/think/x-force/databases-beware-abusing-microsoft-sql-server-with-sqlrecon |

#### 2.6.6 Pass-the-Hash 认证

| 属性 | 描述 |
|------|------|
| **描述** | 使用NT哈希(Pass-the-Hash)通过原始TDS/NTLM进行认证，不需要明文密码或提升权限（不需要SeImpersonate）。 |
| **涉及工具** | SQLRecon |
| **武器化状态** | **已实现** (v4.0新增) |
| **对应参数** | `/a:pth` + `/hash:NT_HASH` |
| **出处** | https://github.com/skahwah/SQLRecon/blob/main/README.md |

### 2.7 持久化类

#### 2.7.1 注册表Run键持久化

| 属性 | 描述 |
|------|------|
| **描述** | 使用`xp_regwrite`过程设置可执行文件在用户登录时自动运行。写入`HKEY_LOCAL_MACHINE\Software\Microsoft\Windows\CurrentVersion\Run`。 |
| **涉及工具** | PowerUpSQL |
| **武器化状态** | **已实现** |
| **对应函数** | `Get-SQLPersistRegRun` |
| **出处** | https://github.com/NetSPI/PowerUpSQL/wiki/Persistence-Functions |

#### 2.7.2 注册表Debugger后门 (IFEO)

| 属性 | 描述 |
|------|------|
| **描述** | 使用`xp_regwrite`为指定可执行文件配置调试器后门（Image File Execution Options），在被调用时运行另一个可执行文件。常用于创建RDP后门（如utilman.exe调用cmd.exe）。 |
| **涉及工具** | PowerUpSQL |
| **武器化状态** | **已实现** |
| **对应函数** | `Get-SQLPersistRegDebugger` |
| **出处** | https://github.com/NetSPI/PowerUpSQL/wiki/Persistence-Functions |

#### 2.7.3 DDL触发器持久化

| 属性 | 描述 |
|------|------|
| **描述** | 使用SQL Server DDL事件触发器创建Windows系统后门。支持通过xp_cmdshell执行任意命令、添加本地OS管理员或添加SQL Server sysadmin。 |
| **涉及工具** | PowerUpSQL |
| **武器化状态** | **已实现** |
| **对应函数** | `Get-SQLPersistTriggerDDL` |
| **出处** | https://github.com/NetSPI/PowerUpSQL/wiki/Persistence-Functions |

#### 2.7.4 启动存储过程持久化

| 属性 | 描述 |
|------|------|
| **描述** | 利用SQL Server启动存储过程实现持久化。 |
| **涉及工具** | PowerUpSQL |
| **武器化状态** | **实验性** |
| **出处** | https://blog.netspi.com/sql-server-persistence-part-1-startup-stored-procedures/ |

### 2.8 SCCM攻击模块 (SQLRecon)

| 属性 | 描述 |
|------|------|
| **描述** | SQLRecon提供15+个专门针对SCCM/ECM数据库的攻击模块，包括：用户枚举、站点枚举、客户端登录信息、加密凭据获取、任务序列获取与解密、凭据解密、管理员权限提升/移除、脚本数据获取、CI数据获取等。 |
| **涉及工具** | SQLRecon |
| **武器化状态** | **已实现** |
| **目标数据库** | MSSQL (SCCM/ECM数据库) |
| **出处** | https://github.com/skahwah/SQLRecon/blob/main/README.md |
| **Wiki** | https://github.com/skahwah/SQLRecon/wiki/7.-SCCM-Modules |

### 2.9 Azure SQL 支持 (SQLRecon)

| 属性 | 描述 |
|------|------|
| **描述** | SQLRecon v4.0+支持Azure SQL Database的认证和操作，包括EntraID认证和Azure本地认证。 |
| **涉及工具** | SQLRecon |
| **武器化状态** | **已实现** |
| **对应参数** | `/a:entraid`、`/a:azurelocal` |
| **出处** | https://github.com/skahwah/SQLRecon/blob/main/README.md |

### 2.10 凭据暴力破解 (enumdb)

| 属性 | 描述 |
|------|------|
| **描述** | 对MSSQL进行凭据暴力破解（支持CIDR范围）、数据库/表/列枚举、关键词敏感数据搜索、数据Dump提取为CSV/XLSX、交互式SQL Shell。 |
| **涉及工具** | enumdb |
| **武器化状态** | **已实现** |
| **出处** | https://github.com/m8sec/enumdb |

### 2.11 端口扫描

| 属性 | 描述 |
|------|------|
| **描述** | 利用openrowset通过数据库服务器进行端口扫描，探测内网其他主机的端口开放状态。支持单端口、多端口和端口范围扫描。 |
| **涉及工具** | MSDAT |
| **武器化状态** | **已实现** |
| **对应模块** | `bulkopen` (`--scan-ports`) |
| **出处** | https://github.com/quentinhardy/msdat |

### 2.12 SQL注入检测与利用 (sqlmap)

| 属性 | 描述 |
|------|------|
| **描述** | sqlmap支持对MSSQL的完整SQL注入检测与利用，包括5种注入技术、数据库指纹识别、用户/密码哈希/权限枚举、全库数据Dump、文件读写、xp_cmdshell命令执行、OS Shell、Meterpreter会话、SMB反射攻击、Windows注册表访问等。 |
| **涉及工具** | sqlmap |
| **武器化状态** | **已实现** |
| **出处** | https://github.com/sqlmapproject/sqlmap/wiki/Features |

### 2.13 配置操作模块 (SQLRecon)

| 属性 | 描述 |
|------|------|
| **描述** | SQLRecon支持启用/禁用多种MSSQL功能组件：RPC、CLR、OLE Automation、xp_cmdshell。 |
| **涉及工具** | SQLRecon |
| **武器化状态** | **已实现** |
| **对应模块** | `EnableRpc`/`DisableRpc`、`EnableClr`/`DisableClr`、`EnableOle`/`DisableOle`、`EnableXp`/`DisableXp` |
| **出处** | https://github.com/skahwah/SQLRecon |

---

## 三、Oracle 攻击手法

> 涉及工具：MDUT、MDUT-Extend、ODAT、Sylas、sqlmap

### 3.1 信息收集/枚举类

#### 3.1.1 SID枚举 (SID Enumeration)

| 属性 | 描述 |
|------|------|
| **描述** | 通过多种方式枚举Oracle数据库的有效SID，包括字典攻击、暴力破解和TNS Listener ALIAS枚举。利用cx_Oracle连接时返回的不同错误码区分有效SID和无效SID。 |
| **涉及工具** | ODAT |
| **武器化状态** | **已实现** |
| **对应模块** | `sidguesser` |
| **出处** | https://github.com/quentinhardy/odat/wiki/sidguesser |

#### 3.1.2 Service Name枚举 (Service Name Enumeration)

| 属性 | 描述 |
|------|------|
| **描述** | 枚举Oracle数据库的有效Service Name，类似于SID枚举但针对Service Name。 |
| **涉及工具** | ODAT |
| **武器化状态** | **已实现** (v5.0+) |
| **对应模块** | `snguesser` |
| **出处** | https://github.com/quentinhardy/odat/wiki/Home |

#### 3.1.3 TNS信息收集 (TNS Listener Reconnaissance)

| 属性 | 描述 |
|------|------|
| **描述** | 无需认证即可与TNS Listener通信，获取数据库别名、版本号和状态信息。 |
| **涉及工具** | ODAT |
| **武器化状态** | **已实现** |
| **对应模块** | `tnscmd` |
| **出处** | https://github.com/quentinhardy/odat/wiki/tnscmd |

#### 3.1.4 数据库信息收集与列名搜索

| 属性 | 描述 |
|------|------|
| **描述** | 获取数据库实例的基本信息（版本、表结构等），搜索可能包含敏感信息（如密码）的列名。 |
| **涉及工具** | ODAT |
| **武器化状态** | **已实现** |
| **对应模块** | `search` |
| **出处** | https://github.com/quentinhardy/odat/wiki/search |

### 3.2 认证攻击类

#### 3.2.1 凭证爆破 (Credential Brute-force)

| 属性 | 描述 |
|------|------|
| **描述** | 通过字典攻击暴力破解Oracle数据库账户的用户名和密码。支持单文件模式、双文件模式、密码喷洒模式、用户名作为密码模式。 |
| **涉及工具** | ODAT |
| **武器化状态** | **已实现** |
| **对应模块** | `passwordguesser` |
| **出处** | https://github.com/quentinhardy/odat/wiki/passwordguesser |

#### 3.2.2 用户名作为密码测试

| 属性 | 描述 |
|------|------|
| **描述** | 利用已认证会话获取所有Oracle用户名列表，然后尝试用用户名本身（大小写）作为密码登录。 |
| **涉及工具** | ODAT |
| **武器化状态** | **已实现** |
| **对应模块** | `userlikepwd` |
| **出处** | https://github.com/quentinhardy/odat/wiki/userlikepwd |

### 3.3 命令执行类

#### 3.3.1 Java存储过程命令执行

| 属性 | 描述 |
|------|------|
| **描述** | 通过Oracle的Java存储过程功能上传并执行Java代码，实现命令执行。MDUT使用JAVA Util导入方式和ShellUtil工具类；ODAT的`java`模块支持`--exec`、`--shell`、`--reverse-shell`等模式。 |
| **涉及工具** | MDUT、ODAT |
| **武器化状态** | **已实现** |
| **出处(MDUT)** | https://github.com/SafeGroceryStore/MDUT |
| **出处(ODAT)** | https://github.com/quentinhardy/odat/wiki/java |

#### 3.3.2 DBMS_SCHEDULER命令执行

| 属性 | 描述 |
|------|------|
| **描述** | 利用Oracle DBMS_SCHEDULER包创建类型为'EXECUTABLE'的调度作业，在数据库服务器上执行操作系统命令。支持`--exec`（无回显）、`--reverse-shell`（反向TCP Shell）、`--make-download`（PowerShell下载）。 |
| **涉及工具** | ODAT、Sylas |
| **武器化状态** | **已实现** |
| **对应模块** | ODAT: `dbmsscheduler`; Sylas: DBMS_SCHEDULER命令执行(无回显) |
| **出处(ODAT)** | https://github.com/quentinhardy/odat/wiki/dbmsscheduler |
| **出处(Sylas)** | https://github.com/Ryze-T/Sylas |

#### 3.3.3 DBMS_XMLQUERY命令执行 (回显)

| 属性 | 描述 |
|------|------|
| **描述** | 通过DBMS_XMLQUERY执行系统命令，可获取回显。 |
| **涉及工具** | Sylas |
| **武器化状态** | **已实现** |
| **出处** | https://github.com/Ryze-T/Sylas |

#### 3.3.4 External Table命令执行

| 属性 | 描述 |
|------|------|
| **描述** | 利用Oracle外部表功能执行操作系统命令或程序。创建指向操作系统文件或程序的外部表，通过SELECT操作触发执行。 |
| **涉及工具** | ODAT |
| **武器化状态** | **已实现** |
| **对应模块** | `externaltable` |
| **出处** | https://github.com/quentinhardy/odat/wiki/externaltable |

#### 3.3.5 ORADBG调试接口命令执行

| 属性 | 描述 |
|------|------|
| **描述** | 利用Oracle调试功能(oradbg接口)执行操作系统命令。 |
| **涉及工具** | ODAT |
| **武器化状态** | **已实现** |
| **对应模块** | `oradbg` |
| **出处** | https://github.com/quentinhardy/odat/wiki/oradbg |

#### 3.3.6 反弹Shell

| 属性 | 描述 |
|------|------|
| **描述** | 通过Oracle Java存储过程实现反弹Shell功能。 |
| **涉及工具** | MDUT |
| **武器化状态** | **已实现** |
| **出处** | https://github.com/SafeGroceryStore/MDUT |

### 3.4 文件操作类

#### 3.4.1 UTL_FILE 文件读/写/删除

| 属性 | 描述 |
|------|------|
| **描述** | 利用Oracle UTL_FILE包进行文件的读取、写入和删除操作。使用UTL_FILE.FOPEN/GET_LINE读取、UTL_FILE.PUT_LINE写入、UTL_FILE.FREMOVE删除。 |
| **涉及工具** | ODAT |
| **武器化状态** | **已实现** |
| **对应模块** | `utlfile` |
| **出处** | https://github.com/quentinhardy/odat/wiki/utlfile |

#### 3.4.2 CTXSYS文件读取

| 属性 | 描述 |
|------|------|
| **描述** | 利用Oracle CTXSYS.DRITHSX.SN功能读取文件。 |
| **涉及工具** | ODAT |
| **武器化状态** | **已实现** |
| **对应模块** | `ctxsys` |
| **出处** | https://github.com/quentinhardy/odat/wiki/ctxsys |

#### 3.4.3 DBMS_XSLPROCESSOR文件上传/下载

| 属性 | 描述 |
|------|------|
| **描述** | 利用Oracle DBMS_XSLPROCESSOR包进行文件上传和下载。 |
| **涉及工具** | ODAT |
| **武器化状态** | **已实现** |
| **对应模块** | `dbmsxslprocessor` |
| **出处** | https://github.com/quentinhardy/odat/wiki/dbmsxslprocessor |

#### 3.4.4 DBMS_ADVISOR文件上传

| 属性 | 描述 |
|------|------|
| **描述** | 利用Oracle DBMS_ADVISOR包上传文件到数据库服务器。 |
| **涉及工具** | ODAT |
| **武器化状态** | **已实现** |
| **对应模块** | `dbmsadvisor` |
| **出处** | https://github.com/quentinhardy/odat/wiki/dbmsadvisor |

#### 3.4.5 DBMS_LOB文件下载

| 属性 | 描述 |
|------|------|
| **描述** | 利用Oracle DBMS_LOB包读取文件（使用BFILE）。 |
| **涉及工具** | ODAT |
| **武器化状态** | **已实现** |
| **对应模块** | `dbmslob` |
| **出处** | https://github.com/quentinhardy/odat/wiki/dbmslob |

#### 3.4.6 文件管理 (MDUT/Sylas)

| 属性 | 描述 |
|------|------|
| **描述** | 支持Oracle服务器的文件上传、下载和管理。MDUT-Extend v1.3.0新增大文件分块上传和下载功能，解决ORA-24345等错误。Sylas支持文件查看、获取数据库目录、文件上传。 |
| **涉及工具** | MDUT、MDUT-Extend、Sylas |
| **武器化状态** | **已实现** |
| **出处(MDUT)** | https://github.com/SafeGroceryStore/MDUT |
| **出处(MDUT-Extend)** | https://github.com/DeEpinGh0st/MDUT-Extend-Release/releases/tag/v1.3.0 |
| **出处(Sylas)** | https://github.com/Ryze-T/Sylas |

### 3.5 权限提升类

#### 3.5.1 CREATE ANY PROCEDURE 提权

| 属性 | 描述 |
|------|------|
| **描述** | 利用CREATE ANY PROCEDURE系统权限提升至DBA或SYS权限。在APEX模式下创建存储过程，利用其高权限执行任意SQL。 |
| **涉及工具** | ODAT |
| **武器化状态** | **已实现** |
| **对应模块** | `privesc --dba-with-create-any-procedure` |
| **出处** | https://github.com/quentinhardy/odat/wiki/privesc |

#### 3.5.2 CREATE/EXECUTE ANY PROCEDURE 提权

| 属性 | 描述 |
|------|------|
| **描述** | 利用CREATE PROCEDURE和EXECUTE ANY PROCEDURE权限以SYS身份执行任意SQL。 |
| **涉及工具** | ODAT |
| **武器化状态** | **已实现** |
| **对应模块** | `privesc --dba-with-execute-any-procedure` |
| **出处** | https://github.com/quentinhardy/odat/wiki/privesc |

#### 3.5.3 CREATE ANY TRIGGER 提权

| 属性 | 描述 |
|------|------|
| **描述** | 利用CREATE ANY TRIGGER和CREATE PROCEDURE权限以SYS身份执行任意SQL。 |
| **涉及工具** | ODAT |
| **武器化状态** | **已实现** |
| **对应模块** | `privesc --dba-with-create-any-trigger` |
| **出处** | https://github.com/quentinhardy/odat/wiki/privesc |

#### 3.5.4 ANALYZE ANY 提权

| 属性 | 描述 |
|------|------|
| **描述** | 利用ANALYZE ANY和CREATE PROCEDURE权限以SYS身份执行任意SQL。 |
| **涉及工具** | ODAT |
| **武器化状态** | **已实现** |
| **对应模块** | `privesc --dba-with-analyze-any` |
| **出处** | https://github.com/quentinhardy/odat/wiki/privesc |

#### 3.5.5 CREATE ANY INDEX 提权

| 属性 | 描述 |
|------|------|
| **描述** | 利用CREATE ANY INDEX和CREATE PROCEDURE权限以SYS身份执行任意SQL。 |
| **涉及工具** | ODAT |
| **武器化状态** | **已实现** |
| **对应模块** | `privesc --dba-with-create-any-index` |
| **出处** | https://github.com/quentinhardy/odat/wiki/privesc |

### 3.6 网络操作类

#### 3.6.1 HTTP请求发送与端口扫描 (UTL_HTTP)

| 属性 | 描述 |
|------|------|
| **描述** | 利用Oracle UTL_HTTP包从数据库服务器发送HTTP请求，可用于端口扫描和内网探测。 |
| **涉及工具** | ODAT |
| **武器化状态** | **已实现** |
| **对应模块** | `utlhttp` |
| **出处** | https://github.com/quentinhardy/odat/wiki/utlhttp |

#### 3.6.2 HTTP请求发送 (HttpUriType)

| 属性 | 描述 |
|------|------|
| **描述** | 利用Oracle HttpUriType对象发送HTTP GET请求和端口扫描。 |
| **涉及工具** | ODAT |
| **武器化状态** | **已实现** |
| **对应模块** | `httpuritype` |
| **出处** | https://github.com/quentinhardy/odat |

#### 3.6.3 TCP端口扫描 (UTL_TCP)

| 属性 | 描述 |
|------|------|
| **描述** | 利用Oracle UTL_TCP包进行TCP端口扫描和原始TCP包发送。 |
| **涉及工具** | ODAT |
| **武器化状态** | **已实现** |
| **对应模块** | `utltcp` |
| **出处** | https://github.com/quentinhardy/odat/wiki/utltcp |

### 3.7 密码哈希提取类

#### 3.7.1 密码哈希提取 (Password Hash Dumping)

| 属性 | 描述 |
|------|------|
| **描述** | 从Oracle数据库中提取用户密码哈希值，支持多种提取方法：直接从SYS.USER$或DBA_USERS表读取、通过DBMS_METADATA.GET_DDL导出、利用DBMS_STATS包获取、利用ORACLE_OCM视图间接获取(CVE-2020-2984)。 |
| **涉及工具** | ODAT |
| **武器化状态** | **已实现** |
| **对应模块** | `passwordstealer` |
| **出处** | https://github.com/quentinhardy/odat/wiki/passwordstealer |

### 3.8 CVE漏洞利用类

#### 3.8.1 TNS投毒 (CVE-2012-1675)

| 属性 | 描述 |
|------|------|
| **描述** | 利用CVE-2012-1675进行TNS Listener投毒攻击，通过向TNS Listener注册恶意数据库实例，将合法用户的连接重定向到攻击者控制的代理服务器，实现中间人攻击和会话劫持。 |
| **CVE编号** | CVE-2012-1675 |
| **涉及工具** | ODAT |
| **武器化状态** | **已实现** |
| **对应模块** | `tnspoison` |
| **出处** | https://github.com/quentinhardy/odat/wiki/tnspoison |

#### 3.8.2 会话密钥嗅探 (CVE-2012-3137)

| 属性 | 描述 |
|------|------|
| **描述** | 利用CVE-2012-3137通过网络嗅探获取Oracle认证会话密钥和盐值，进而通过字典攻击破解密码。需在本地网络接口嗅探TNS认证流量。 |
| **CVE编号** | CVE-2012-3137 |
| **涉及工具** | ODAT |
| **武器化状态** | **已实现** |
| **对应模块** | `stealremotepwds` |
| **目标版本** | Oracle 11g (需root权限) |
| **出处** | https://github.com/quentinhardy/odat/wiki/stealremotepwds |

#### 3.8.3 XQuery漏洞 (CVE-2014-4237)

| 属性 | 描述 |
|------|------|
| **描述** | 利用CVE-2014-4237使仅有SELECT权限的用户可以修改其有SELECT权限的任何表。Oracle XQuery存在漏洞，允许绕过ALTER权限检查。 |
| **CVE编号** | CVE-2014-4237 |
| **涉及工具** | ODAT |
| **武器化状态** | **已实现** |
| **对应模块** | `cve --set-pwd-2014-4237` |
| **出处** | https://github.com/quentinhardy/odat/wiki/cve |

#### 3.8.4 JVM安全绕过 (CVE-2018-3004)

| 属性 | 描述 |
|------|------|
| **描述** | 利用CVE-2018-3004绕过Oracle JVM安全限制，实现任意文件写入。 |
| **CVE编号** | CVE-2018-3004 |
| **涉及工具** | ODAT |
| **武器化状态** | **已实现** |
| **对应模块** | `cve --cve-2018-3004`、`java --create-file-CVE-2018-3004` |
| **出处** | https://github.com/quentinhardy/odat/wiki/cve |

#### 3.8.5 ORACLE_OCM间接密码获取 (CVE-2020-2984)

| 属性 | 描述 |
|------|------|
| **描述** | 利用CVE-2020-2984通过ORACLE_OCM视图间接获取密码哈希。 |
| **CVE编号** | CVE-2020-2984 |
| **涉及工具** | ODAT |
| **武器化状态** | **已实现** |
| **对应模块** | `passwordstealer --get-passwords-ocm` |
| **出处** | https://github.com/quentinhardy/odat/wiki/passwordstealer |

### 3.9 SMB认证捕获

| 属性 | 描述 |
|------|------|
| **描述** | 诱导Oracle数据库服务器向攻击者控制的SMB服务器发起认证请求，捕获SMB/NTLM哈希。需目标为Windows且Oracle服务不能运行在网络服务/系统/本地服务账户下。 |
| **涉及工具** | ODAT |
| **武器化状态** | **已实现** |
| **对应模块** | `smb --capture` |
| **出处** | https://github.com/quentinhardy/odat/wiki/smb |

### 3.10 PL/SQL解包装 (PL/SQL Unwrapping)

| 属性 | 描述 |
|------|------|
| **描述** | 反编译Oracle数据库中被加密的(wrapped) PL/SQL源代码。逆向Oracle的PL/SQL包装算法，还原原始源代码。 |
| **涉及工具** | ODAT |
| **武器化状态** | **已实现** |
| **对应模块** | `unwrapper` |
| **出处** | https://github.com/quentinhardy/odat/wiki/unwrapper |

### 3.11 表数据转储与SQL Shell

| 属性 | 描述 |
|------|------|
| **描述** | 导出数据库表中的数据到本地文件（CSV/XLSX格式）；获取交互式SQL shell直接执行SQL查询。 |
| **涉及工具** | ODAT |
| **武器化状态** | **已实现** |
| **对应模块** | `search --dump`、`search --sql-shell` |
| **出处** | https://github.com/quentinhardy/odat |

### 3.12 SQL注入检测与利用 (sqlmap)

| 属性 | 描述 |
|------|------|
| **描述** | sqlmap支持对Oracle的完整SQL注入检测与利用，包括5种注入技术、数据库指纹识别、用户/密码哈希/权限枚举、全库数据Dump、文件读写等。 |
| **涉及工具** | sqlmap |
| **武器化状态** | **已实现** |
| **出处** | https://github.com/sqlmapproject/sqlmap/wiki/Features |

### 3.13 痕迹清理

| 属性 | 描述 |
|------|------|
| **描述** | 清理Oracle数据库中创建的Java存储过程和函数。 |
| **涉及工具** | MDUT |
| **武器化状态** | **已实现** |
| **出处** | https://github.com/SafeGroceryStore/MDUT |

### 3.14 Windows快速信息收集 (Sylas)

| 属性 | 描述 |
|------|------|
| **描述** | 提供四个快速执行功能：进程信息、用户枚举、补丁信息、系统版本。 |
| **涉及工具** | Sylas |
| **武器化状态** | **已实现** |
| **出处** | https://github.com/Ryze-T/Sylas |

---

## 四、PostgreSQL 攻击手法

> 涉及工具：MDUT、MDUT-Extend、Sylas、sqlmap

### 4.1 UDF提权 (Windows)

| 属性 | 描述 |
|------|------|
| **描述** | 通过上传自定义UDF库到PostgreSQL实现Windows平台提权。 |
| **目标平台** | Windows |
| **涉及工具** | MDUT |
| **武器化状态** | **已实现** |
| **出处** | https://github.com/SafeGroceryStore/MDUT |

### 4.2 COPY FROM PROGRAM 命令执行

| 属性 | 描述 |
|------|------|
| **描述** | 通过PostgreSQL的`COPY ... FROM PROGRAM`语法执行系统命令。PostgreSQL 9.3+支持。Sylas支持通过PostgreSQL执行系统命令。 |
| **目标版本** | PostgreSQL 9.3+ |
| **涉及工具** | MDUT、Sylas、sqlmap |
| **武器化状态** | **已实现** |
| **出处(MDUT)** | https://github.com/SafeGroceryStore/MDUT |
| **出处(Sylas)** | https://github.com/Ryze-T/Sylas |
| **出处(sqlmap)** | https://github.com/sqlmapproject/sqlmap |

### 4.3 文件读取 (CVE方式)

| 属性 | 描述 |
|------|------|
| **描述** | 通过PostgreSQL的CVE漏洞方式实现文件读取功能。MDUT-Extend v1.2.0新增。 |
| **涉及工具** | MDUT-Extend |
| **武器化状态** | **已实现** |
| **出处** | https://github.com/DeEpinGh0st/MDUT-Extend-Release/releases/tag/v1.2.0 |

### 4.4 文件读写 (sqlmap/MDUT/Sylas)

| 属性 | 描述 |
|------|------|
| **描述** | sqlmap支持PostgreSQL的文件读取和写入；MDUT和MDUT-Extend支持PostgreSQL文件管理；Sylas支持使用`pg_read_file()`读取文件、`COPY`命令写入文件。 |
| **涉及工具** | MDUT、MDUT-Extend、Sylas、sqlmap |
| **武器化状态** | **已实现** |
| **出处** | 多工具支持 |

### 4.5 UDF注入命令执行 (sqlmap)

| 属性 | 描述 |
|------|------|
| **描述** | sqlmap通过注入自定义用户定义函数(UDF)执行任意命令并检索标准输出。 |
| **涉及工具** | sqlmap |
| **武器化状态** | **已实现** |
| **出处** | https://github.com/sqlmapproject/sqlmap/wiki/Features |

### 4.6 OS Shell 与 OAST攻击 (sqlmap)

| 属性 | 描述 |
|------|------|
| **描述** | sqlmap支持通过PostgreSQL获取交互式OS Shell、带外TCP连接（Meterpreter/VNC）、DNS外泄攻击、Metasploit Shellcode内存执行。 |
| **涉及工具** | sqlmap |
| **武器化状态** | **已实现** |
| **出处** | https://github.com/sqlmapproject/sqlmap/wiki/Features |

---

## 五、Redis 攻击手法

> 涉及工具：MDUT、MDUT-Extend、RedisEXP、Databasetools、NoSQLMap

### 5.1 未授权访问利用

| 属性 | 描述 |
|------|------|
| **描述** | 连接到无需密码认证的Redis实例，直接执行任意Redis命令，获取完全控制权。 |
| **涉及工具** | RedisEXP、Databasetools |
| **武器化状态** | **已实现** |
| **出处(RedisEXP)** | https://github.com/yuyan-sec/RedisEXP |
| **出处(Databasetools)** | https://github.com/Hel10-Web/Databasetools |

### 5.2 密码爆破

| 属性 | 描述 |
|------|------|
| **描述** | 使用字典文件对Redis进行密码爆破，支持自定义密码字典。 |
| **涉及工具** | RedisEXP、Databasetools |
| **武器化状态** | **已实现** |
| **出处(RedisEXP)** | https://github.com/yuyan-sec/RedisEXP |
| **出处(Databasetools)** | https://github.com/Hel10-Web/Databasetools |

### 5.3 主从复制RCE (Master-Slave Replication RCE)

| 属性 | 描述 |
|------|------|
| **描述** | 利用Redis主从复制机制，将目标Redis实例设置为攻击者控制的恶意主节点的从节点，通过FULLRESYNC同步恶意模块文件(exp.so/exp.dll)到目标，加载模块执行系统命令。**注意：此操作会清空目标Redis数据！** |
| **影响版本** | Redis 4.x、5.x (Redis 6.0+引入ACL限制，Redis 7.0默认禁用模块加载) |
| **涉及工具** | MDUT、RedisEXP、Databasetools |
| **武器化状态** | **已实现** |
| **出处(MDUT)** | https://github.com/SafeGroceryStore/MDUT/blob/main/redis-cus-rogue.py |
| **出处(RedisEXP)** | https://github.com/yuyan-sec/RedisEXP |
| **出处(Databasetools)** | https://github.com/Hel10-Web/Databasetools |

### 5.4 模块加载RCE (Module Load RCE)

| 属性 | 描述 |
|------|------|
| **描述** | 直接在目标Redis上加载恶意的.so或.dll模块文件，注册system.exec命令实现任意命令执行。 |
| **影响版本** | Redis 4.x+ |
| **涉及工具** | RedisEXP |
| **武器化状态** | **已实现** |
| **出处** | https://github.com/yuyan-sec/RedisEXP |

### 5.5 SSH公钥注入 (SSH Key Injection)

| 属性 | 描述 |
|------|------|
| **描述** | 通过Redis的`CONFIG SET dir`和`CONFIG SET dbfilename`命令，将攻击者的SSH公钥写入目标服务器的`/root/.ssh/authorized_keys`文件，实现无密码SSH登录。MDUT-Extend提供无损写入方式。 |
| **目标平台** | Linux |
| **涉及工具** | MDUT、MDUT-Extend、RedisEXP、Databasetools |
| **武器化状态** | **已实现** |
| **出处(MDUT)** | https://github.com/SafeGroceryStore/MDUT |
| **出处(MDUT-Extend)** | https://github.com/DeEpinGh0st/MDUT-Extend-Release/releases/tag/v1.0.0 |
| **出处(RedisEXP)** | https://github.com/yuyan-sec/RedisEXP |
| **出处(Databasetools)** | https://github.com/Hel10-Web/Databasetools |

### 5.6 WebShell写入

| 属性 | 描述 |
|------|------|
| **描述** | 利用Redis的持久化功能，将WebShell内容写入Web服务器可访问的目录，获取Web后门。支持base64编码写入。 |
| **涉及工具** | MDUT、RedisEXP、Databasetools |
| **武器化状态** | **已实现** |
| **出处(MDUT)** | https://github.com/SafeGroceryStore/MDUT |
| **出处(RedisEXP)** | https://github.com/yuyan-sec/RedisEXP |
| **出处(Databasetools)** | https://github.com/Hel10-Web/Databasetools |

### 5.7 计划任务注入 (Crontab Injection)

| 属性 | 描述 |
|------|------|
| **描述** | 通过Redis向Linux系统的计划任务目录(如`/var/spool/cron/`)写入恶意计划任务文件，实现定时命令执行或反弹Shell。 |
| **目标平台** | Linux |
| **涉及工具** | RedisEXP、Databasetools |
| **武器化状态** | **已实现** |
| **出处(RedisEXP)** | https://github.com/yuyan-sec/RedisEXP |
| **出处(Databasetools)** | https://github.com/Hel10-Web/Databasetools |

### 5.8 Lua沙盒绕过 (CVE-2022-0543)

| 属性 | 描述 |
|------|------|
| **描述** | 利用Debian/Ubuntu系统下Redis的Lua脚本沙箱逃逸漏洞执行任意系统命令。该漏洞源于Debian对Lua沙箱的补丁缺陷。 |
| **CVE编号** | CVE-2022-0543 |
| **目标平台** | Debian/Ubuntu发行版 |
| **涉及工具** | MDUT-Extend、RedisEXP、Databasetools |
| **武器化状态** | **已实现** |
| **出处(MDUT-Extend)** | https://github.com/DeEpinGh0st/MDUT-Extend-Release/releases/tag/v1.3.0 |
| **出处(RedisEXP)** | https://github.com/yuyan-sec/RedisEXP |
| **出处(Databasetools)** | https://github.com/Hel10-Web/Databasetools |

### 5.9 Gopher SSRF Payload生成

| 属性 | 描述 |
|------|------|
| **描述** | 生成Gopher协议的Redis命令Payload，用于通过SSRF漏洞攻击内网Redis实例。 |
| **涉及工具** | RedisEXP |
| **武器化状态** | **已实现** |
| **出处** | https://github.com/yuyan-sec/RedisEXP |

### 5.10 DLL劫持 (Windows)

| 属性 | 描述 |
|------|------|
| **描述** | 通过Redis写入恶意的dbghelp.dll文件到redis-server.exe所在目录，利用Windows DLL加载机制实现代码执行。 |
| **目标平台** | Windows |
| **涉及工具** | RedisEXP |
| **武器化状态** | **已实现** |
| **出处** | https://github.com/yuyan-sec/RedisEXP |

### 5.11 反弹Shell

| 属性 | 描述 |
|------|------|
| **描述** | 通过Redis实现反弹Shell功能。MDUT v2.1.0+和MDUT-Extend均支持。 |
| **涉及工具** | MDUT、MDUT-Extend |
| **武器化状态** | **已实现** |
| **出处(MDUT)** | https://github.com/SafeGroceryStore/MDUT |
| **出处(MDUT-Extend)** | https://github.com/DeEpinGh0st/MDUT-Extend-Release |

### 5.12 无损文件读写 (MDUT-Extend)

| 属性 | 描述 |
|------|------|
| **描述** | Redis无损文件读写功能，可在不破坏原有文件内容的情况下读写目标文件。 |
| **涉及工具** | MDUT-Extend |
| **武器化状态** | **已实现** |
| **出处** | https://github.com/DeEpinGh0st/MDUT-Extend-Release/releases/tag/v1.0.0 |

### 5.13 主从复制文件上传

| 属性 | 描述 |
|------|------|
| **描述** | 利用主从复制机制将任意文件上传到目标Redis服务器指定路径，可用于上传WebShell、恶意程序等。 |
| **涉及工具** | RedisEXP |
| **武器化状态** | **已实现** |
| **出处** | https://github.com/yuyan-sec/RedisEXP |

### 5.14 文件存在性检测

| 属性 | 描述 |
|------|------|
| **描述** | 通过Redis判断目标服务器上指定绝对路径的文件是否存在。 |
| **涉及工具** | RedisEXP |
| **武器化状态** | **已实现** |
| **出处** | https://github.com/yuyan-sec/RedisEXP |

### 5.15 CVE-2025-49844 利用

| 属性 | 描述 |
|------|------|
| **描述** | 利用Redis的安全漏洞实现代码执行。 |
| **CVE编号** | CVE-2025-49844 |
| **涉及工具** | MDUT-Extend |
| **武器化状态** | **已实现** |
| **出处** | https://github.com/DeEpinGh0st/MDUT-Extend-Release/releases/tag/v1.3.0 |

### 5.16 Slave-Read-Only控制与痕迹清理

| 属性 | 描述 |
|------|------|
| **描述** | 控制Redis从节点的只读属性；执行`SLAVEOF NO ONE`断开主从复制关系用于攻击后清理。 |
| **涉及工具** | MDUT、RedisEXP |
| **武器化状态** | **已实现** |
| **出处(MDUT)** | https://github.com/SafeGroceryStore/MDUT |
| **出处(RedisEXP)** | https://github.com/yuyan-sec/RedisEXP |

---

## 六、MongoDB 攻击手法

> 涉及工具：MDUT-Extend、NoSQLMap、sqlmap(2026新增)

### 6.1 未授权访问扫描与利用

| 属性 | 描述 |
|------|------|
| **描述** | 扫描目标网络中开放的MongoDB实例，检测是否存在未授权访问。可对整个网段进行批量扫描，发现匿名可访问的MongoDB实例。 |
| **涉及工具** | NoSQLMap、MDUT-Extend |
| **武器化状态** | **已实现** |
| **出处(NoSQLMap)** | https://github.com/codingo/NoSQLMap/blob/master/nsmscan.py |
| **出处(MDUT-Extend)** | https://github.com/DeEpinGh0st/MDUT-Extend-Release/releases/tag/v1.3.0 |

### 6.2 数据库枚举与克隆

| 属性 | 描述 |
|------|------|
| **描述** | 获取MongoDB服务器版本和平台信息、列出所有数据库、枚举集合。将目标MongoDB数据库完整克隆到攻击者控制的MongoDB实例中，实现数据窃取。 |
| **涉及工具** | NoSQLMap、MDUT-Extend |
| **武器化状态** | **已实现** |
| **出处** | https://github.com/codingo/NoSQLMap/blob/master/nsmmongo.py |

### 6.3 NoSQL Web应用注入攻击

| 属性 | 描述 |
|------|------|
| **描述** | 针对使用MongoDB的Web应用进行自动化注入测试，支持布尔盲注、时间盲注、JavaScript注入（$where操作符）、认证绕过。支持GET/POST方法和Burp请求导入。 |
| **涉及工具** | NoSQLMap |
| **武器化状态** | **已实现** |
| **出处** | https://github.com/codingo/NoSQLMap/blob/master/nsmweb.py |

### 6.4 GridFS文件枚举

| 属性 | 描述 |
|------|------|
| **描述** | 通过MongoDB的GridFS功能枚举存储的文件(附件)，可能获取敏感文件。 |
| **涉及工具** | NoSQLMap |
| **武器化状态** | **已实现** |
| **出处** | https://github.com/codingo/NoSQLMap/blob/master/nsmmongo.py |

### 6.5 Meterpreter Shell获取

| 属性 | 描述 |
|------|------|
| **描述** | 通过MongoDB服务获取Meterpreter反向Shell，实现远程控制。 |
| **涉及工具** | NoSQLMap |
| **武器化状态** | **已实现** |
| **条件** | 需要Metasploit Framework配合 |
| **出处** | https://github.com/codingo/NoSQLMap |

### 6.6 命令执行 (MDUT-Extend)

| 属性 | 描述 |
|------|------|
| **描述** | MDUT-Extend v1.3.0新增对MongoDB数据库的利用支持，包括命令执行等操作。 |
| **涉及工具** | MDUT-Extend |
| **武器化状态** | **已实现** |
| **出处** | https://github.com/DeEpinGh0st/MDUT-Extend-Release/releases/tag/v1.3.0 |

### 6.7 NoSQL注入支持 (sqlmap 2026新增)

| 属性 | 描述 |
|------|------|
| **描述** | sqlmap在2026年新增对MongoDB等NoSQL数据库的注入检测和利用支持。 |
| **涉及工具** | sqlmap |
| **武器化状态** | **已实现** |
| **出处** | https://github.com/sqlmapproject/sqlmap/commit/2893fd5c4d8f056f76f73facb6e6e24a25c04c85 |

---

## 七、CouchDB 攻击手法

> 涉及工具：NoSQLMap

### 7.1 未授权访问扫描

| 属性 | 描述 |
|------|------|
| **描述** | 扫描并检测CouchDB实例是否存在未授权访问，无需凭据即可连接并获取数据库信息。 |
| **涉及工具** | NoSQLMap |
| **武器化状态** | **已实现** |
| **对应代码** | nsmcouch.py - couchScan函数 |
| **出处** | https://github.com/codingo/NoSQLMap/blob/master/nsmcouch.py |

### 7.2 数据库/用户/密码哈希枚举

| 属性 | 描述 |
|------|------|
| **描述** | 枚举所有数据库、用户列表和密码哈希(包括SHA1哈希和盐值)，支持字典攻击和暴力破解恢复明文密码。 |
| **涉及工具** | NoSQLMap |
| **武器化状态** | **已实现** |
| **支持哈希版本** | CouchDB < 1.3: `password_sha` (SHA1); CouchDB >= 1.3: `derived_key` (PBKDF2) |
| **出处** | https://github.com/codingo/NoSQLMap/blob/master/nsmcouch.py |

### 7.3 数据库克隆

| 属性 | 描述 |
|------|------|
| **描述** | 利用CouchDB的复制功能将目标数据库完整克隆到攻击者控制的CouchDB实例。 |
| **涉及工具** | NoSQLMap |
| **武器化状态** | **已实现** |
| **出处** | https://github.com/codingo/NoSQLMap/blob/master/nsmcouch.py |

### 7.4 Web应用注入攻击

| 属性 | 描述 |
|------|------|
| **描述** | 通过NoSQLMap的Web应用攻击模块，针对使用CouchDB的Web应用进行注入测试。支持布尔盲注、时间盲注、JavaScript注入。 |
| **涉及工具** | NoSQLMap |
| **武器化状态** | **已实现** |
| **出处** | https://github.com/codingo/NoSQLMap/blob/master/nsmweb.py |

---

## 八、Cassandra 攻击手法

| 属性 | 描述 |
|------|------|
| **描述** | NoSQLMap README中明确指出计划在未来的版本中增加对Cassandra数据库的支持。目前无武器化攻击手法。 |
| **涉及工具** | NoSQLMap (计划中) |
| **武器化状态** | **计划中（尚未实现）** |
| **出处** | https://github.com/codingo/NoSQLMap#readme |

---

## 九、多数据库/通用攻击手法

> 涉及工具：sqlmap（支持40+种数据库后端）

### 9.1 SQL注入检测技术（5种核心技法）

| 技术名称 | 英文名称 | 描述 | 武器化状态 |
|----------|----------|------|------------|
| **基于布尔的盲注** | Boolean-based blind | 替换或追加受影响的参数，使用包含SELECT子语句的SQL语句字符串。通过比较HTTP响应头/正文与原始请求的差异，逐字符推断注入语句的输出。采用二分算法。 | **已实现** |
| **基于时间的盲注** | Time-based blind | 替换或追加受影响的参数，使用使后端DBMS延迟返回一定秒数的查询语句。通过比较HTTP响应时间与原始请求的差异，逐字符推断注入语句的输出。 | **已实现** |
| **基于报错的注入** | Error-based | 替换或追加受影响的参数，使用特定于数据库的错误消息语句。解析HTTP响应搜索包含注入预定义字符链和子查询输出的DBMS错误消息。 | **已实现** |
| **UNION查询注入** | UNION query-based | 追加到受影响的参数，使用以`UNION ALL SELECT`开头的语法有效的SQL语句。适用于Web应用页面直接在for循环中传递SELECT语句输出的场景。 | **已实现** |
| **堆叠查询注入** | Stacked queries | 测试Web应用是否支持堆叠查询，如果支持，则追加分号(`;`)后跟要执行的SQL语句。可用于执行除SELECT之外的其他SQL语句。 | **已实现** |

**出处**: https://github.com/sqlmapproject/sqlmap/wiki/Techniques

### 9.2 数据库指纹识别

| 攻击手法 | 描述 | 武器化状态 |
|----------|------|------------|
| **数据库软件版本识别** | 基于错误消息、banner解析、函数输出比较和特定功能进行广泛的后端数据库软件版本和底层操作系统指纹识别。 | **已实现** |
| **强制指定DBMS** | 用户可以通过`--dbms`选项强制指定后端数据库管理系统名称。 | **已实现** |
| **Web服务器/应用技术识别** | 基本的Web服务器软件和Web应用技术指纹。 | **已实现** |
| **Banner获取** | 支持检索DBMS banner信息。 | **已实现** |
| **会话用户识别** | 支持检索当前会话用户信息和当前数据库信息。 | **已实现** |
| **DBA权限检测** | 支持检查会话用户是否为数据库管理员(DBA)。 | **已实现** |

**出处**: https://github.com/sqlmapproject/sqlmap/wiki/Features

### 9.3 数据提取与枚举

| 攻击手法 | 描述 | 武器化状态 |
|----------|------|------------|
| **用户枚举** | 枚举DBMS所有用户。 | **已实现** |
| **密码哈希提取** | 枚举DBMS用户密码哈希，自动识别格式并支持字典攻击破解。 | **已实现** |
| **权限/角色枚举** | 枚举DBMS用户权限和角色。 | **已实现** |
| **数据库/表/列枚举** | 枚举所有数据库、指定数据库的表、指定表的列。 | **已实现** |
| **Schema枚举** | 枚举完整的数据库schema。 | **已实现** |
| **全库数据Dump** | 支持完整导出数据库表、指定范围的条目或特定列；自动dump所有数据库的schema和条目。 | **已实现** |
| **关键词搜索** | 搜索特定数据库名称、跨所有数据库的特定表、或跨所有表的特定列。 | **已实现** |
| **自定义SQL执行** | 在交互式SQL客户端中运行自定义SQL语句。 | **已实现** |
| **暴力破解表名/列名** | 当无权读取系统表时，暴力破解表名和列名。 | **已实现** |

**出处**: https://github.com/sqlmapproject/sqlmap/wiki/Features

### 9.4 文件系统访问

| 攻击手法 | 描述 | 目标数据库 | 武器化状态 |
|----------|------|------------|------------|
| **文件读取** | 从后端DBMS文件系统读取文件。 | MySQL, PostgreSQL, MSSQL, Oracle等 | **已实现** |
| **文件写入** | 将本地文件写入后端DBMS文件系统。 | MySQL, PostgreSQL, MSSQL, Oracle等 | **已实现** |

**出处**: https://github.com/sqlmapproject/sqlmap/wiki/Features

### 9.5 操作系统命令执行

| 攻击手法 | 描述 | 目标数据库 | 武器化状态 |
|----------|------|------------|------------|
| **UDF注入命令执行** | 注入自定义UDF执行任意命令并检索标准输出。 | MySQL, PostgreSQL | **已实现** |
| **xp_cmdshell命令执行** | 通过MSSQL的xp_cmdshell()存储过程执行命令，自动启用/重建。 | MSSQL | **已实现** |
| **交互式OS Shell** | 提供交互式操作系统shell。 | MySQL, PostgreSQL, MSSQL | **已实现** |

**出处**: https://github.com/sqlmapproject/sqlmap/wiki/Features

### 9.6 OAST攻击与Metasploit集成

| 攻击手法 | 描述 | 目标数据库 | 武器化状态 |
|----------|------|------------|------------|
| **带外TCP连接** | 建立攻击者机器与数据库服务器之间的带外有状态TCP连接（交互式命令提示符、Meterpreter或VNC）。 | MySQL, PostgreSQL, MSSQL | **已实现** |
| **DNS外泄攻击** | 使用DNS域名进行数据外泄。 | 支持的数据库 | **已实现** |
| **Metasploit Shellcode内存执行** | 通过UDF `sys_bineval()` 在数据库内存中执行Metasploit shellcode。 | MySQL, PostgreSQL | **已实现** |
| **独立Payload上传执行** | 通过UDF `sys_exec()` 上传和执行Metasploit独立payload stager。 | MySQL, PostgreSQL, MSSQL | **已实现** |
| **SMB反射攻击** | 通过SMB反射攻击（MS08-068）执行Metasploit shellcode。 | MSSQL (Windows) | **已实现** |
| **sp_replwritetovarbin溢出利用** | 利用MSSQL 2000/2005的sp_replwritetovarbin存储过程堆缓冲区溢出（MS09-004），自动DEP绕过。 | MSSQL 2000/2005 | **已实现** |
| **权限提升** | 通过Metasploit的getsystem命令支持数据库进程用户权限提升。 | Windows | **已实现** |

**出处**: https://github.com/sqlmapproject/sqlmap/wiki/Features

### 9.7 Windows注册表访问

| 攻击手法 | 描述 | 武器化状态 |
|----------|------|------------|
| **注册表读取** | 读取Windows注册表键值。 | **已实现** |
| **注册表写入** | 写入Windows注册表键值数据。 | **已实现** |
| **注册表删除** | 删除Windows注册表键值。 | **已实现** |

**出处**: https://github.com/sqlmapproject/sqlmap/wiki/Features

### 9.8 多目标扫描 (MDUT-Extend)

| 属性 | 描述 |
|------|------|
| **描述** | 新增对多种数据库的批量存活探测功能，可扫描目标网段中存活的数据库服务。 |
| **目标数据库** | MySQL、MSSQL、Oracle、PostgreSQL、Redis、MongoDB |
| **涉及工具** | MDUT-Extend |
| **武器化状态** | **已实现** |
| **出处** | https://github.com/DeEpinGh0st/MDUT-Extend-Release/releases/tag/v1.3.0 |

---

## 十、跨数据库攻击手法对比总表

### 10.1 命令执行类攻击手法对比

| 攻击手法 | MySQL | MSSQL | Oracle | PostgreSQL | Redis | MongoDB |
|----------|-------|-------|--------|-----------|-------|---------|
| UDF提权/命令执行 | MDUT/sqlmap | - | - | sqlmap | - | - |
| xp_cmdshell | - | MDUT/MSDAT/PowerUpSQL/SQLRecon/SharpSQLTools/Sylas/sqlmap | - | - | - | - |
| sp_oacreate/OLE Automation | - | MDUT/MSDAT/PowerUpSQL/SQLRecon/SharpSQLTools/Sylas | - | - | - | - |
| CLR程序集执行 | - | MDUT/PowerUpSQL/SQLRecon/SharpSQLTools/Sylas | - | - | - | - |
| Agent Job命令执行 | - | MSDAT/PowerUpSQL/SQLRecon | - | - | - | - |
| R/Python脚本执行 | - | PowerUpSQL | - | - | - | - |
| Java存储过程 | - | - | MDUT/ODAT | - | - | - |
| DBMS_SCHEDULER | - | - | ODAT/Sylas | - | - | - |
| DBMS_XMLQUERY | - | - | Sylas | - | - | - |
| External Table | - | - | ODAT | - | - | - |
| ORADBG | - | - | ODAT | - | - | - |
| COPY FROM PROGRAM | - | - | - | MDUT/Sylas/sqlmap | - | - |
| 主从复制RCE | - | - | - | - | MDUT/RedisEXP/Databasetools | - |
| 模块加载RCE | - | - | - | - | RedisEXP | - |
| Lua沙盒绕过 | - | - | - | - | MDUT-Extend/RedisEXP/Databasetools | - |
| 命令执行(MongoDB) | - | - | - | - | - | MDUT-Extend |

### 10.2 文件操作类攻击手法对比

| 攻击手法 | MySQL | MSSQL | Oracle | PostgreSQL | Redis |
|----------|-------|-------|--------|-----------|-------|
| 文件读取 | sqlmap | MDUT/MSDAT/SharpSQLTools/Sylas/sqlmap | ODAT/MDUT/Sylas/sqlmap | MDUT-Extend/Sylas/sqlmap | RedisEXP |
| 文件写入/上传 | MDUT/sqlmap | MDUT/MSDAT/SharpSQLTools/Sylas/sqlmap | ODAT/MDUT/Sylas/sqlmap | MDUT/Sylas/sqlmap | RedisEXP/Databasetools |
| 文件下载 | - | MDUT/MSDAT/SharpSQLTools | ODAT/MDUT-Extend | - | - |
| WebShell写入 | - | Sylas | - | Sylas | RedisEXP/Databasetools |
| SSH公钥注入 | - | - | - | - | MDUT/MDUT-Extend/RedisEXP/Databasetools |
| 大文件分块传输 | - | - | MDUT-Extend | - | - |

### 10.3 权限提升与持久化类攻击手法对比

| 攻击手法 | MySQL | MSSQL | Oracle | PostgreSQL | Redis |
|----------|-------|-------|--------|-----------|-------|
| Impersonation提权 | - | PowerUpSQL/SQLRecon | - | - | - |
| Trustworthy提权 | - | MSDAT/PowerUpSQL | - | - | - |
| 本地Admin到sysadmin | - | PowerUpSQL | - | - | - |
| Potato系列提权 | - | MDUT-Extend/SharpSQLTools | - | - | - |
| CREATE ANY PROCEDURE提权 | - | - | ODAT | - | - |
| ANALYZE ANY提权 | - | - | ODAT | - | - |
| 注册表持久化 | - | PowerUpSQL | - | - | - |
| 触发器持久化 | - | PowerUpSQL | - | - | - |
| Shellcode加载 | - | MDUT-Extend/SharpSQLTools | - | - | - |

### 10.4 网络/横向移动类攻击手法对比

| 攻击手法 | MySQL | MSSQL | Oracle | PostgreSQL | Redis |
|----------|-------|-------|--------|-----------|-------|
| SMB认证捕获 | - | MSDAT/PowerUpSQL/SQLRecon | ODAT | - | - |
| 端口扫描 | - | MSDAT | ODAT(UTL_HTTP/UTL_TCP) | - | - |
| HTTP请求发送 | - | - | ODAT(UTL_HTTP/HttpUriType) | - | - |
| 链接服务器爬取 | - | PowerUpSQL/SQLRecon | - | - | - |
| DNS外泄 | sqlmap | sqlmap | sqlmap | sqlmap | - |
| Gopher SSRF | - | - | - | - | RedisEXP |

### 10.5 信息收集/枚举类攻击手法对比

| 攻击手法 | MySQL | MSSQL | Oracle | PostgreSQL | Redis | MongoDB | CouchDB |
|----------|-------|-------|--------|-----------|-------|---------|---------|
| 无认证信息获取 | - | MSDAT/SQLRecon | ODAT(TNS) | - | - | - | - |
| 凭据暴力破解 | enumdb | MSDAT/PowerUpSQL/enumdb | ODAT | - | RedisEXP/Databasetools | - | NoSQLMap |
| 密码哈希提取 | sqlmap | MSDAT/PowerUpSQL | ODAT | sqlmap | - | - | NoSQLMap |
| 数据库枚举 | enumdb/sqlmap | MSDAT/PowerUpSQL/SQLRecon/enumdb | ODAT | sqlmap | - | NoSQLMap/MDUT-Extend | NoSQLMap |
| SPN/域枚举 | - | PowerUpSQL/SQLRecon | - | - | - | - | - |
| 敏感数据搜索 | enumdb/sqlmap | MSDAT/PowerUpSQL/SQLRecon/enumdb | ODAT | sqlmap | - | - | - |
| 数据库克隆 | - | - | - | - | - | NoSQLMap | NoSQLMap |

---

## 十一、汇总统计表

### 11.1 各数据库类型攻击手法数量统计

| 数据库类型 | 攻击手法数量 | 主要来源工具 | 关键CVE |
|-----------|------------|-------------|---------|
| **MySQL** | 8 | MDUT, sqlmap, enumdb | - |
| **Microsoft SQL Server** | 48+ | MDUT, MDUT-Extend, MSDAT, PowerUpSQL, SQLRecon, SharpSQLTools, Sylas, enumdb, sqlmap | MS08-068, MS09-004 |
| **Oracle** | 37+ | MDUT, MDUT-Extend, ODAT, Sylas, sqlmap | CVE-2012-1675, CVE-2012-3137, CVE-2014-4237, CVE-2018-3004, CVE-2020-2984 |
| **PostgreSQL** | 6 | MDUT, MDUT-Extend, Sylas, sqlmap | - |
| **Redis** | 16 | MDUT, MDUT-Extend, RedisEXP, Databasetools | CVE-2022-0543, CVE-2025-49844 |
| **MongoDB** | 7 | MDUT-Extend, NoSQLMap, sqlmap | - |
| **CouchDB** | 4 | NoSQLMap | - |
| **Cassandra** | 0 (计划中) | NoSQLMap | - |
| **多数据库/通用** | 20+ | sqlmap | - |
| **总计** | **142+** | **15+款工具** | **10+个CVE** |

### 11.2 CVE漏洞利用汇总表

| CVE编号 | 目标数据库 | 影响版本 | 利用效果 | 涉及工具 |
|---------|-----------|----------|----------|----------|
| CVE-2012-1675 | Oracle 10g/11g | 默认安装后均受影响 | TNS Listener投毒，中间人攻击 | ODAT |
| CVE-2012-3137 | Oracle 11g | Oracle 11g | 会话密钥嗅探，破解密码 | ODAT |
| CVE-2014-4237 | Oracle 11g/12c | 受影响版本 | SELECT用户可修改任意表 | ODAT |
| CVE-2018-3004 | Oracle 12c/18c/19c | 受影响版本 | JVM安全限制绕过，任意文件写入 | ODAT |
| CVE-2020-2984 | Oracle | 受影响版本 | 通过ORACLE_OCM间接获取密码哈希 | ODAT |
| CVE-2022-0543 | Redis (Debian/Ubuntu) | Debian/Ubuntu上的Redis | Lua沙箱逃逸，命令执行 | MDUT-Extend, RedisEXP, Databasetools |
| CVE-2025-49844 | Redis | 特定Redis版本 | 代码执行 | MDUT-Extend |
| MS08-068 | MSSQL (Windows) | Windows | SMB反射攻击 | sqlmap |
| MS09-004 | MSSQL 2000/2005 | MSSQL 2000/2005 | sp_replwritetovarbin堆缓冲区溢出 | sqlmap |

---

## 十二、工具索引表

### 12.1 工具概览

| 工具名称 | 开发语言 | 主要目标数据库 | 攻击手法数量 | GitHub Stars | 项目地址 |
|----------|----------|--------------|------------|-------------|----------|
| **sqlmap** | Python | 40+种数据库 | 30+ | 37.7k+ | https://github.com/sqlmapproject/sqlmap |
| **ODAT** | Python 3 | Oracle | 36+ | - | https://github.com/quentinhardy/odat |
| **PowerUpSQL** | PowerShell | MSSQL | 40+ | 2.7k+ | https://github.com/NetSPI/PowerUpSQL |
| **SQLRecon** | C# | MSSQL/Azure SQL | 45+ | 813+ | https://github.com/skahwah/SQLRecon |
| **MDUT** | Java | MySQL/MSSQL/Oracle/PostgreSQL/Redis | 20+ | 2.2k+ | https://github.com/SafeGroceryStore/MDUT |
| **MDUT-Extend** | Java | MySQL/MSSQL/Oracle/PostgreSQL/Redis/MongoDB | 25+ | 926+ | https://github.com/DeEpinGh0st/MDUT-Extend-Release |
| **MSDAT** | Python 3 | MSSQL | 23+ | - | https://github.com/quentinhardy/msdat |
| **SharpSQLTools** | C# | MSSQL | 30+ | 965+ | https://github.com/uknowsec/SharpSQLTools |
| **NoSQLMap** | Python | MongoDB/CouchDB | 10+ | 3.3k+ | https://github.com/codingo/NoSQLMap |
| **RedisEXP** | Go | Redis | 13+ | 944+ | https://github.com/yuyan-sec/RedisEXP |
| **Databasetools** | Go | MySQL/MSSQL/Oracle/PostgreSQL/Redis | 10+ | 866+ | https://github.com/Hel10-Web/Databasetools |
| **Sylas** | C# | MSSQL/Oracle/PostgreSQL | 15+ | 545+ | https://github.com/Ryze-T/Sylas |
| **enumdb** | Python 3 | MySQL/MSSQL | 10+ | 222+ | https://github.com/m8sec/enumdb |

### 12.2 各工具支持的攻击手法速查

| 工具 | MySQL | MSSQL | Oracle | PostgreSQL | Redis | MongoDB | CouchDB |
|------|-------|-------|--------|-----------|-------|---------|---------|
| **MDUT** | UDF提权/反弹Shell/NTFS/命令执行/HTTP隧道 | CLR/sp_oacreate/文件管理 | Java执行/反弹Shell/文件管理 | UDF提权/命令执行 | 主从复制/SSH注入/反弹Shell | - | - |
| **MDUT-Extend** | 驱动切换 | GodPotato/Shellcode/综合利用 | 大文件传输 | CVE文件读取 | 无损读写/CVE-2022-0543/CVE-2025-49844/反弹Shell | 命令执行 | - |
| **ODAT** | - | - | 36+种攻击手法(SID枚举/命令执行/文件操作/权限提升/CVE利用等) | - | - | - | - |
| **MSDAT** | - | 23+种攻击手法(xp_cmdshell/OLE/AgentJob/SMB/哈希提取/端口扫描等) | - | - | - | - | - |
| **PowerUpSQL** | - | 40+种攻击手法(发现/提权/命令执行/持久化/链接爬取等) | - | - | - | - | - |
| **SQLRecon** | - | 45+种攻击手法(含SCCM/Azure/PTH/ADSI等高级模块) | - | - | - | - | - |
| **SharpSQLTools** | - | xp_cmdshell/sp_oacreate/CLR提权/LSASS转储/Shellcode加载等 | - | - | - | - | - |
| **Sylas** | - | xp_cmdshell/sp_oacreate/CLR/WebShell写入 | DBMS_XMLQUERY/DBMS_SCHEDULER/文件管理 | COPY命令执行/WebShell写入 | - | - | - |
| **sqlmap** | 完整注入支持/UDF/文件读写/OS Shell | 完整注入支持/xp_cmdshell/文件读写/OS Shell | 完整注入支持/文件读写 | 完整注入支持/COPY/UDF/文件读写 | - | NoSQL注入(2026新增) | - |
| **enumdb** | 暴力破解/枚举/Dump | 暴力破解/枚举/Dump | - | - | - | - | - |
| **NoSQLMap** | - | - | - | - | - | 匿名访问/注入/克隆/GridFS | 未授权访问/枚举/克隆/注入 |
| **RedisEXP** | - | - | - | - | 13种攻击手法(主从复制/模块/SSH/WebShell/Crontab/CVE-2022-0543/Gopher/DLL劫持等) | - | - |
| **Databasetools** | - | - | - | - | 7种攻击手法(交互式Shell/主从复制/Lua绕过/SSH/WebShell/Crontab/爆破) | - | - |

---

## 参考链接汇总

### MySQL相关工具
1. MDUT: https://github.com/SafeGroceryStore/MDUT
2. MDUT中文文档: https://www.yuque.com/u21224612/nezuig
3. MDUT-Extend: https://github.com/DeEpinGh0st/MDUT-Extend-Release

### MSSQL相关工具
4. MSDAT: https://github.com/quentinhardy/msdat
5. PowerUpSQL: https://github.com/NetSPI/PowerUpSQL
6. PowerUpSQL Wiki: https://github.com/NetSPI/PowerUpSQL/wiki
7. SQLRecon: https://github.com/skahwah/SQLRecon
8. SharpSQLTools: https://github.com/uknowsec/SharpSQLTools
9. enumdb: https://github.com/m8sec/enumdb

### Oracle相关工具
10. ODAT: https://github.com/quentinhardy/odat
11. ODAT Wiki: https://github.com/quentinhardy/odat/wiki
12. Sylas: https://github.com/Ryze-T/Sylas

### NoSQL相关工具
13. NoSQLMap: https://github.com/codingo/NoSQLMap
14. RedisEXP: https://github.com/yuyan-sec/RedisEXP
15. Databasetools: https://github.com/Hel10-Web/Databasetools

### 通用SQL注入工具
16. sqlmap: https://github.com/sqlmapproject/sqlmap
17. sqlmap Wiki: https://github.com/sqlmapproject/sqlmap/wiki

### 参考项目
18. redis-rogue-server: https://github.com/n0b0dyCN/redis-rogue-server
19. redis-rce: https://github.com/Ridter/redis-rce
20. RabR (exp.dll/exp.so来源): https://github.com/0671/RabR
21. WarSQLKit: https://github.com/mindspoof/MSSQL-Fileless-Rootkit-WarSQLKit

---

> **免责声明**：本文档仅用于安全研究和防御目的，旨在帮助安全人员了解数据库攻击工具的能力以加强防护。使用这些工具进行未经授权的攻击是违法行为。本文档中的信息基于公开的GitHub仓库和技术博客，仅供合法授权的安全测试和学术研究使用。

> **法律警告**：未经授权访问计算机系统、数据库或网络属于违法行为。使用本文档中描述的技术和工具攻击您不拥有或未经授权测试的系统是非法的，并可能导致刑事起诉。请始终确保您拥有适当的授权，并遵守所有适用的法律和法规。

---

*文档版本: v1.0*
*最后更新: 2025年7月*
*生成方式: 基于5份独立调研报告的综合整理*
