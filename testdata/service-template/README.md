# Local Service-Template Docker Tests

This fixture starts minimal Redis, PostgreSQL, MySQL, and SSH services on
localhost-only high ports and runs zombie service templates after a successful
login.

Ports:

- Redis: `127.0.0.1:16379`, password `zombie_redis_pass`
- PostgreSQL: `127.0.0.1:15432`, `zombie:zombie_pg_pass`
- MySQL: `127.0.0.1:13306`, `zombie:zombie_mysql_pass`
- SSH: `127.0.0.1:10022`, `zombie:zombie_ssh_pass` (built from `./ssh`)

Run from the repository root:

```bash
go test -tags docker ./integration -run TestServiceTemplateDocker -count=1
```

By default the `Existing` scope loads real templates from
`../proton/templates/services`. Override this checkout layout with:

```bash
ZOMBIE_SERVICE_TEMPLATES_DIR=/path/to/proton/templates/services go test -tags docker ./integration -run TestServiceTemplateDocker/Existing -count=1
```

Run a focused scope:

```bash
go test -tags docker ./integration -run TestServiceTemplateDocker/Smoke -count=1
go test -tags docker ./integration -run TestServiceTemplateDocker/Payload -count=1
go test -tags docker ./integration -run TestServiceTemplateDocker/PostExploit -count=1
go test -tags docker ./integration -run TestServiceTemplateDocker/Existing -count=1
go test -tags docker ./integration -run TestServiceTemplateDocker/Exploit -count=1
```

The Docker integration test cleans the fixture with `docker compose down -v`
after the run. Manual cleanup:

```bash
docker compose -f ./testdata/service-template/compose.yaml down -v
```

The Redis fixture enables protected config changes so Redis file-write
templates can be reproduced locally. Keep it bound to localhost-only ports.

Template layout:

- `templates/smoke`: baseline post-auth execution checks.
- `templates/payload`: command-line `--payload` behavior checks.
- `templates/post-exploit`: Docker-only post-exploit capability checks.
- `templates/exploit`: Docker-only high-risk exploit reproductions.

Exploit templates use Nuclei-style `variables`, and CLI overrides use
Nuclei-compatible `-V/--var key=value`. Request-level payloads can be supplied
or overridden with repeated `--payload key=value` flags:

```bash
go run . -i mysql://root:zombie_mysql_root@127.0.0.1:13306 --no-honeypot --no-unauth --service-template ./testdata/service-template/templates/exploit/mysql-outfile-local.yaml -V outfile_path=/tmp/zombie_mysql_custom.txt -V outfile_marker=zombie-custom-ok
go run . -i redis://:zombie_redis_pass@127.0.0.1:16379 --no-honeypot --no-unauth --service-template ./testdata/service-template/templates/payload/redis-payload-cli.yaml --payload payload_key=zombie:payload:cli-a --payload payload_key=zombie:payload:cli-b
```

Service templates also support request-level `payloads` plus `attack`
(`sniper`, `pitchfork`, `clusterbomb`) and both `{{name}}` and `§name§`
placeholders. If `-V` and `--payload` use the same key, `--payload` wins for
payload iteration.

The existing-template smoke currently covers these existing templates from
`../proton/templates/services`:

- `redis/redis-info-gather.yaml`
- `redis/redis-config-check.yaml`
- `redis/redis-sensitive-keys.yaml`
- `redis/redis-mdut-rogue-prereq-check.yaml`
- `postgresql/postgresql-info-gather.yaml`
- `postgresql/postgresql-mdut-capability-check.yaml`
- `mysql/mysql-info-gather.yaml`
- `mysql/mysql-credential-columns.yaml`
- `mysql/mysql-user-enum.yaml`
- `mysql/mysql-udf-check.yaml`
- `mysql/mysql-mdut-capability-check.yaml`
- `ssh/ssh-info-gather.yaml`
- `ssh/ssh-env-files.yaml`
- `ssh/ssh-credential-files.yaml`
- `ssh/ssh-docker-escape.yaml`

The local post-exploit smoke currently covers these capability paths:

- Redis: sensitive key/value read, token key discovery, RDB-backed file write.
- MySQL: application credential table read, `/etc/passwd` read through
  `LOAD_FILE`, `SELECT ... INTO OUTFILE` write.
- PostgreSQL: `/etc/passwd` read through `pg_read_file`, `COPY FROM PROGRAM`
  command execution and file write.
- SSH: command execution, remote file write/read, `.env` secret read, credential
  file path discovery.
