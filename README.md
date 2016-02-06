# Configo [![Build Status](https://travis-ci.org/zeroturnaround/configo.svg?branch=master)](https://travis-ci.org/zeroturnaround/configo) [![Join the chat at https://gitter.im/zeroturnaround/configo](https://badges.gitter.im/zeroturnaround/configo.svg)](https://gitter.im/zeroturnaround/configo?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge) [![Approved issues](https://badge.waffle.io/zeroturnaround/configo.svg?label=ready&title=Waffle.io)](http://waffle.io/zeroturnaround/configo)

**Configo** helps running **12factor** (http://12factor.net/config) applications by loading environment variables from different sources.

# Usage
Imagine having an application that is configurable with environment variables. Let us assume that this is a self-contained (http://12factor.net/processes) **NodeJS** application, and that we have a **Docker** image for it:
```Dockerfile
FROM node

ADD . /app
WORKDIR /app

CMD ["node", "server.js"]
```

## Loading the configuration
Surely you want to deploy this application to dev/qa/production. Some configuration is obviously required. We will use these environment variables for this configuration:
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

Since we have 5 servers in production, we have to configure these environment variables for each server. Would it not be nice to have a single source of configuration and load it for each app? And maybe some shared configuration as well? This is where **Configo** comes in. 

Meet **Configo**!

First, change your Dockerfile ever so slightly:
```diff
FROM node

+RUN curl -L https://github.com/zeroturnaround/configo/releases/download/v0.1.0/configo.linux-amd64 >/usr/local/bin/configo && \
+    chmod +x /usr/local/bin/configo

ADD . /app
WORKDIR /app

-CMD ["node", "server.js"]
+CMD ["configo", "node", "server.js"]
```

> For this example, we will use an URL as a source for **Configo**. Other possible sources can be used - check the configuration section below for more information.

Upload your configuration files to your internal HTTP server:
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

We are now ready to start our applications:
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
Since we added **Configo**, it will now load configuration from the sources we specified. In addition, it will merge these settings and configure the environment variables for your application.

## Mapping configuration
Wait! There is more! You can use Golang templates to manipulate environment variables. Simply set the environment variable with the value prefixed by `CONFIGO:` and it will be executed (the end result will not include the `CONFIGO:` prefix). Here is an example of this:
```bash
docker run \
  -e CONFIGO_SOURCE_0='{"type": "http", "format": "yaml", "url": "https://my.server.com/common.yaml"}' \
  -e DB_REDIS_URI='CONFIGO:{{or .DB_REDIS_URI .REDIS_URI "redis://localhost/0"}}' \
  myAppImage node server.js
```
In this example, we are using the built-in `or` function. This will return the first non-empty argument or the last argument. The entire functionality of Go templates is available for use. More information can be found at https://golang.org/pkg/text/template/. 

# Configuration
Configo can be run without specifying a configuration. Doing this will not affect your application.

To make **Configo** load environment variables, you need to specify at least the environment variable `CONFIGO_SOURCE_N` where N indicates the priority for this source. Your application can have multiple sources and this will use priorities for overriding. Priority #5 overrides #0 and so on. N can be any positive number. You can reserve a space by using high values like 100, 200 etc. This way you do not have to rearrange sources when you decide to insert additional sources in the middle.

Each source is a JSON with at least `type` field.

## Source types
### URL
URL is one of the simplest source types. This will load the file from URL specified and parse it.

| Field    | Description                                                                                                          | Required |
|----------|----------------------------------------------------------------------------------------------------------------------|----------|
| url      | URL to the file                                                                                                      |    Yes   |
| format   | Specifies the format to be used. Values that are allowed include: <ul><li>json</li><li>yaml</li><li>hcl</li><li>toml</li><li>properties</li></ul> |    Yes   |
| insecure | Skip domain and certificate check                                                                                    |    No    |
| tls.cert | PEM-encoded TLS certificate                                                                                          |    No    |
| tls.key  | PEM-encoded TLS key                                                                                                  |    No    |

Example:
```bash
CONFIGO_SOURCE_0='{"type": "http", "format": "json", "url": "http://my.server.com/myAppConfig"}'
```
---

### File
This will load a file from the local file system and parse it.

| Field  | Description                                                                                                          | Required |
|--------|----------------------------------------------------------------------------------------------------------------------|----------|
| path   | Path to the file                                                                                                     |    Yes   |
| format | Specifies the format to be used. Values that are allowed include: <ul><li>json</li><li>yaml</li><li>hcl</li><li>toml</li><li>properties</li></ul> |    Yes   |

Example:
```bash
CONFIGO_SOURCE_0='{"type": "file", "path": "/etc/myApp/test.yml", "format": "yaml"}'
```
---

### Redis
This will use a Redis hashmap as a source.

| Field | Description                          | Required |
|-------|--------------------------------------|----------|
| uri   | Redis URI for connection.            |    Yes   |
| key   | Redis key (should point to hashmap). |    Yes   |

Example:
```bash
CONFIGO_SOURCE_0='{"type": "redis", "uri": "redis://56.42.168.12:6390/0", "key": "myAppConfig"}'
```
---

### Consul
This will use Consul as a source.

| Field   | Description                   | Required |
|---------|-------------------------------|----------|
| address | Where to connect.             |    Yes   |
| prefix  | Will be used in the KV query. |    Yes   |
| scheme  | Scheme to use.                |    No    |

Example:
```bash
CONFIGO_SOURCE_0='{"type": "consul", "address": "consul.prod.corp.com:8500", "prefix": "myAppConfig"}'
```
---

### Etcd
This will use Etcd as a source.

| Field      | Description                      | Required |
|------------|----------------------------------|----------|
| endpoints  | Array of endpoints to connect.   |    Yes   |
| prefix     | Which prefix to use.             |    Yes   |
| keepPrefix | When true, will keep prefix.     |    No    |

Example:
```bash
CONFIGO_SOURCE_0='{"type": "etcd", "endpoints": ["http://etcd.corp.com:4001"], "prefix": "myApp/"}'
```
---

### DynamoDB
This will use DynamoDB as a source.

| Field     | Description                                                            | Required |
|-----------|------------------------------------------------------------------------|----------|
| table     | Table name.                                                            |    Yes   |
| key       | Hash key (range keys are not supported).                               |    Yes   |
| endpoint  | Endpoint to be used for the connection.                                |    No    |
| region    | Which region to connect to (default: us-west-1).                       |    No    |
| accessKey | AWS accessKey. Will use default credentials when not set.              |    No    |
| secretKey | AWS secretKey. Will use default credentials when accessKey is not set. |    No    |

Example:
```bash
CONFIGO_SOURCE_0='{"type": "dynamodb", "table": "configs", "key": "myApp"}'
```
---

# Installation
## Prebuilt binaries
Precompiled binaries are available as GitHub releases at https://github.com/zeroturnaround/configo/releases/latest. 

## Build it yourself (with Docker)
If you have Docker installed, run this command to build binaries for all platforms:

```bash
$ docker run -it --rm -v "$PWD":/go/src/github.com/zeroturnaround/configo \
  -w /go/src/github.com/zeroturnaround/configo golang:1.5 make godep_restore build_all
```

## Build it yourself (without Docker)
You need to have Golang 1.5 and GNU Make installed if you want to build **Configo**.

```bash
$ make godep_restore build_all
```

Binaries for each platform will be available in the `/bin/` folder:
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
* http://projects.spring.io/spring-cloud/ - For inspiration (See Spring Cloud Config).
* https://github.com/kelseyhightower/confd - For some ideas regarding sources.
* https://github.com/spf13/viper - For parses.
* https://github.com/hashicorp/terraform - For FlatMap implementation.
