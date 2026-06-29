# MDUT service-template analysis

Source reviewed: `../MDUT`, mainly:

- `MDAT-DEV/src/main/java/Util/MysqlSqlUtil.java`
- `MDAT-DEV/src/main/java/Util/MssqlSqlUtil.java`
- `MDAT-DEV/src/main/java/Util/PostgreSqlUtil.java`
- `MDAT-DEV/src/main/java/Util/OracleSqlUtil.java`
- `MDAT-DEV/src/main/java/Dao/RedisDao.java`

## Current service template coverage

Existing templates under `../proton/templates/services` already cover more than version checks:

- Redis: info/config checks, sensitive key discovery, crontab/webshell/authorized_keys write templates.
- PostgreSQL: info gathering plus superuser and unsafe language checks.
- MySQL: info gathering, credential-column discovery, user enumeration, UDF surface checks, file-read checks.
- MSSQL: info gathering plus `xp_cmdshell` availability checks.
- Oracle: info gathering plus privilege and role checks.
- SSH: info gathering, env/credential file discovery, Docker escape signals.
- FTP/LDAP/Memcached/MongoDB: basic discovery and sensitive-data checks.

## Template-friendly MDUT capabilities

These MDUT features are read-only or mostly read-only and fit the current `db`/`kv` service-template model.

| Service | MDUT capability | Template status | Notes |
| --- | --- | --- | --- |
| MySQL | Version/OS/arch via `CONCAT_WS`, plugin dir, `secure_file_priv`, FILE privilege, existing UDF function count | Added `mysql/mysql-mdut-capability-check.yaml` | Determines whether UDF/file-write paths are plausible without writing a library. |
| PostgreSQL | `server_version`, current user, superuser, `pg_execute_server_program`, server-file roles, PL/Python languages | Added `postgresql/postgresql-mdut-capability-check.yaml` | Identifies whether `COPY FROM PROGRAM` or unsafe language routes are available without creating objects. |
| MSSQL | Sysadmin, `xp_cmdshell`, OLE Automation, CLR, current DB trustworthy flag, MDUT CLR proc presence | Added `mssql/mssql-mdut-capability-check.yaml` | Does not enable options or import assemblies. |
| Oracle | Version, DBA status, scheduler/Java/procedure privileges, DBA/Java roles, existing `SHELLRUN`/`FILERUN` functions | Added `oracle/oracle-mdut-capability-check.yaml` | Does not create Java source, grant Java permissions, or create scheduler jobs. |
| Redis | Version/arch, dir/dbfilename, replica read-only flag, loaded modules | Added `redis/redis-mdut-rogue-prereq-check.yaml` | Captures rogue-module prerequisites without `SLAVEOF`, `MODULE LOAD`, or payload delivery. |

## High-risk or parameterized MDUT capabilities

These should not be default smoke templates. They need explicit opt-in, user-supplied parameters, and cleanup support.

- MySQL: `SELECT ... INTO DUMPFILE`, `CREATE FUNCTION`, `sys_eval`, `backshell`, NTFS ADS directory creation.
- PostgreSQL: large-object UDF injection/export, `CREATE FUNCTION ... LANGUAGE C`, `COPY FROM PROGRAM`, command-result tables.
- MSSQL: enabling `xp_cmdshell`/OLE/CLR, command execution, SQL Agent jobs, file upload/download/delete, CLR assembly import.
- Oracle: Java source import, Java permission grants, `SHELLRUN`/`FILERUN`, DBMS_SCHEDULER executable jobs, reverse shell helpers.
- Redis: write crontab/webshell/authorized_keys, rogue master replication, `MODULE LOAD`, `system.exec`, reverse shell.

## Local verification scope

The current Docker fixture in `testdata/service-template` can verify Redis, PostgreSQL, MySQL, and SSH templates. It now includes smoke coverage for these MDUT-inspired templates:

- `redis/redis-mdut-rogue-prereq-check.yaml`
- `postgresql/postgresql-mdut-capability-check.yaml`
- `mysql/mysql-mdut-capability-check.yaml`

MSSQL and Oracle templates are load-checked by `go test ./...` through `service/load_test.go`, but they do not yet have local Docker smoke coverage. Add separate fixtures only if the heavier images and license/runtime requirements are acceptable.

## Exploit reproduction status

`go test -tags docker ./integration -run TestServiceTemplateDocker/Exploit -count=1` performs an opt-in local exploit reproduction against Docker-only targets:

| Service | Template | Reproduced behavior | Verification |
| --- | --- | --- | --- |
| Redis | `templates/exploit/redis-write-webshell-local.yaml` | `CONFIG SET dir/dbfilename` plus `SAVE` writes `/var/www/html/shell.php` | Container file contains `<?php system($_GET[0]);?>` |
| MySQL | `templates/exploit/mysql-outfile-local.yaml` | `SELECT ... INTO OUTFILE` writes `/tmp/zombie_mysql_outfile.txt` | `LOAD_FILE` through zombie and container `cat` both return `zombie-mysql-outfile-ok` |
| PostgreSQL | `templates/exploit/postgres-copy-program-local.yaml` | `COPY FROM PROGRAM` executes `id` as the PostgreSQL OS user | Zombie extracts `uid=...`, and `/tmp/zombie_pg_cmd.txt` exists in the container |

Redis 7 blocks protected config changes by default, so the local Redis fixture explicitly starts with `--enable-protected-configs yes`. Keep this fixture bound to localhost.

`go test -tags docker ./integration -run TestServiceTemplateDocker/PostExploit -count=1` performs broader post-auth capability validation against the same Docker-only targets:

| Service | Template | Capability verified |
| --- | --- | --- |
| Redis | `templates/post-exploit/redis-post-exploit-local.yaml` | Sensitive value read, token-key discovery, local RDB file write |
| MySQL | `templates/post-exploit/mysql-post-exploit-local.yaml` | App credential query, `/etc/passwd` read with `LOAD_FILE`, `INTO OUTFILE` write |
| PostgreSQL | `templates/post-exploit/postgres-post-exploit-local.yaml` | `/etc/passwd` read with `pg_read_file`, `COPY FROM PROGRAM` command execution |
| SSH | `templates/post-exploit/ssh-post-exploit-local.yaml` | Command execution, file write/read, env secret and credential-path discovery |

## Implementation gaps before high-risk templates

To safely support MDUT's exploitation actions as templates, zombie should add:

- More exploit templates using the newly supported Nuclei-style `variables`, request-level `payloads`, `attack`, `-V/--var key=value`, and repeated `--payload key=value` CLI overrides for command, attacker host/port, public key, target path, payload path, and encoding.
- Per-template risk labels or explicit allow flags, for example `--service-template-risk critical`.
- Cleanup/finalizer blocks that always run after state-changing Redis/MSSQL/MySQL/PostgreSQL/Oracle actions.
- A dry-run or capability-only mode for CI and broad scanning.
