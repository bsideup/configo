# Configo [![Build Status](https://travis-ci.org/zeroturnaround/configo.svg?branch=master)](https://travis-ci.org/zeroturnaround/configo) [![Join the chat at https://gitter.im/zeroturnaround/configo](https://badges.gitter.im/zeroturnaround/configo.svg)](https://gitter.im/zeroturnaround/configo?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)

configo helps run 12factor (http://12factor.net/config) apps by loading environment variables from different sources.

# Usage
Imagine you have an app and it's configurable with environment variables. Let's assume it's self-contained (http://12factor.net/processes) NodeJS app, and we have a Docker image for it:
```Dockerfile
FROM node
ADD . /app
WORKDIR /app

CMD ["node", "server.js"]
```

## Loading configuration
Now you want to deploy it to dev/qa/production. Obviously some configuration is required. We will use environment variables for it:
```bash
docker run \
  -e "DB_MONGO_URI=mongodb://user:pass@mongo.prod.domain.com/db" \
  -e "DB_REDIS_URI=redis://some.redis.prod.domain.com/0" \
  -e "GOOGLE_ANALYTICS_KEY=UA-XXXXX-Y" \
  -e TWITTER_KEY=abcdefg \
  -e SEND_EMAILS=true \
  myAppImage node server.js
```
We also have a server to run some background jobs:
```bash
docker run \
  -e "DB_MONGO_URI=mongodb://user:pass@mongo.prod.domain.com/db" \
  -e "DB_REDIS_URI=redis://some.redis.prod.domain.com/0" \
  -e TWITTER_KEY=abcdefg \
  -e SEND_EMAILS=true \
  -e RUN_JOBS=true \
  myAppImage node server.js
```

Since we have 5 servers on production we have to configure this environment variables on each server. Would be nice to have one single source of configuration and load it for each app, right? And maybe some shared configuration as well? 

Meet Configo!
First, slightly change your Dockerfile:
```Dockerfile
FROM node

# Download Configo binary
RUN curl -L https://github.com/zeroturnaround/configo/releases/download/v0.1.0/configo.linux-amd64 >/usr/local/bin/configo && \
    chmod +x /usr/local/bin/configo

ADD . /app
WORKDIR /app

# Add it before your command
CMD ["configo", "node", "server.js"]
```

> For this example we will use url as a source for configo, but there are more, check "Configuration" section.

Now place your config files on your internal http server:
```bash
$ curl -sSL https://my.server.com/common.yaml
db:
  mongo:
    uri: mongodb://user:pass@mongo.prod.domain.com/db
  redis:
    uri: redis://some.redis.prod.domain.com/0
twitter:
  key: abcdefg
send_emails: true
```
```bash
$ curl -sSL https://my.server.com/server.yaml
google.analytics.key: UA-XXXXX-Y
```
```bash
$ curl -sSL https://my.server.com/jobs.yaml
run_jobs: true
```

We're ready to start our apps:
```bash
docker run \
  -e CONFIGO_SOURCE_0='{"type": "http", "format": "yaml", "url": "https://my.server.com/common.yaml"}' \
  -e CONFIGO_SOURCE_100='{"type": "http", "format": "yaml", "url": "https://my.server.com/server.yaml"}' \
  myAppImage node server.js

docker run \
  -e CONFIGO_SOURCE_0='{"type": "http", "format": "yaml", "url": "https://my.server.com/common.yaml"}' \
  -e CONFIGO_SOURCE_100='{"type": "http", "format": "yaml", "url": "https://my.server.com/jobs.yaml"}' \
  myAppImage node server.js
```
Once we added Configo, it will load configuration from the sources we configured, will merge them and configure environment variables for your application.

## Mapping configuration
But there is more! You can use Golang templates for environment variables manipulation. Just set environment variable with value prefixed with `CONFIGO:`, and it will be executed (result will not include `CONFIGO:` prefix). Consider following example:
```bash
docker run \
  -e CONFIGO_SOURCE_0='{"type": "http", "format": "yaml", "url": "https://my.server.com/common.yaml"}' \
  -e DB_REDIS_URI='CONFIGO:{{or .DB_REDIS_URI .REDIS_URI "redis://localhost/0"}}' \
  myAppImage node server.js
```
In this example we're using built-in `or` function. It will return first non-empty argument or the last argument. All functionality of Go templates is available. Check documentation: https://golang.org/pkg/text/template/

# Configuration
Minimal configuration for Configo is a command to run. It will mean "Pass–through" mode where Configo will not affect execution of your app (except mappings, but their value should be prefixed with Configo-specific string).

To make it load environment variables you should specify at least one environment variable `CONFIGO_SOURCE_N`, where N - priority of this source. Your app could have multiple sources, it will use priorities for overrides, where priority #5 overrides #0. N could be any positive number. You can "reserve" a space by using some high numbers like 100, 200, 1000, so you don't have to rearrange sources if you decided yo insert one more in the middle.

Each source is a JSON with at least `type` field.

## Source types
### URL
One of the simplest source types. Will load file from specified URL and parse it.

| Field  | Description                                                                                                          | Required |
|--------|----------------------------------------------------------------------------------------------------------------------|----------|
| url    | url to the file                                                                                                      |    yes   |
| format | which format to use. Allowed values: <ul><li>json</li><li>yaml</li><li>hcl</li><li>toml</li><li>properties</li></ul> |    yes   |

Example:
```bash
CONFIGO_SOURCE_0='{"type": "http", "format": "json", "url": "http://my.server.com/myAppConfig"}'
```
---

### File
Will load file from the local system and parse it.

| Field  | Description                                                                                                          | Required |
|--------|----------------------------------------------------------------------------------------------------------------------|----------|
| path   | path to the file                                                                                                     |    yes   |
| format | which format to use. Allowed values: <ul><li>json</li><li>yaml</li><li>hcl</li><li>toml</li><li>properties</li></ul> |    yes   |

Example:
```bash
CONFIGO_SOURCE_0='{"type": "file", "path": "/etc/myApp/test.yml", "format": "yaml"}'
```
---

### Redis
Will use Redis's hashmap as a source.

| Field | Description                         | Required |
|-------|-------------------------------------|----------|
| uri   | Redis URI for connection            |    yes   |
| key   | Redis key (should point to hashmap) |    yes   |

Example:
```bash
CONFIGO_SOURCE_0='{"type": "redis", "uri": "redis://56.42.168.12:6390/0", "key": "myAppConfig"}'
```
---

### Consul
Will use Consul as a source.

| Field   | Description              | Required |
|---------|--------------------------|----------|
| address | where to connect         |    yes   |
| prefix  | will be used in KV query |    yes   |
| scheme  | scheme to use            |    no    |

Example:
```bash
CONFIGO_SOURCE_0='{"type": "consul", "address": "consul.prod.corp.com:8500", "prefix": "myAppConfig"}'
```
---

### Etcd
Will use Etcd as a source.

| Field      | Description                   | Required |
|------------|-------------------------------|----------|
| endpoints  | array of endpoints to connect |    yes   |
| prefix     | which prefix to use           |    yes   |
| keepPrefix | if true, will keep prefix     |    no    |

Example:
```bash
CONFIGO_SOURCE_0='{"type": "etcd", "endpoints": ["http://etcd.corp.com:4001"], "prefix": "myApp/"}'
```
---

### DynamoDB
Will use DynamoDB as a source.

| Field     | Description                                                          | Required |
|-----------|----------------------------------------------------------------------|----------|
| table     | table name                                                           |    yes   |
| key       | Hash key (Range keys are not supported yet)                          |    yes   |
| endpoint  | which endpoint to use for connection                                 |    no    |
| region    | which region to connect (default is "us-west-1")                     |    no    |
| accessKey | AWS accessKey. Will use default credentials if is not set.           |    no    |
| secretKey | AWS secretKey. Will use default credentials if accessKey is not set. |    no    |

Example:
```bash
CONFIGO_SOURCE_0='{"type": "dynamodb", "table": "configs", "key": "myApp"}'
```
---

# Installation
## Prebuilt binaries
Precompiled binaries are available as GitHub releases:
https://github.com/zeroturnaround/configo/releases/latest

## Build it yourself (with Docker)
If you have Docker installed then run this command to build binaries for all platforms:

```bash
$ docker run -it --rm -v "$PWD":/go/src/github.com/zeroturnaround/configo \
  -w /go/src/github.com/zeroturnaround/configo golang:1.5 make godep_restore build_all
```

## Build it yourself (without Docker)
You should have Golang 1.5 and GNU Make installed if you want to build Configo.

```bash
$ make godep_restore build_all
```

Binaries for each platform will be available under `/bin/` folder:
```
$ tree bin/
bin/
├── configo.darwin-386
├── configo.darwin-amd64
├── configo.linux-386
├── configo.linux-amd64
├── configo.windows-386
└── configo.windows-amd64

0 directories, 6 files
```

# Thanks
* http://projects.spring.io/spring-cloud/ - for inspiration (See Spring Cloud Config)
* https://github.com/kelseyhightower/confd - for some ideas about sources
* https://github.com/spf13/viper - for parses
* https://github.com/hashicorp/terraform - for FlatMap implementation