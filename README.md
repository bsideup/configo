# Configo [![Build Status](https://travis-ci.org/zeroturnaround/configo.svg?branch=master)](https://travis-ci.org/zeroturnaround/configo) [![Goreport](https://goreportcard.com/badge/github.com/zeroturnaround/configo)](https://goreportcard.com/report/github.com/zeroturnaround/configo) [![Join the chat at https://gitter.im/zeroturnaround/configo](https://badges.gitter.im/zeroturnaround/configo.svg)](https://gitter.im/zeroturnaround/configo?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge) [![Approved issues](https://badge.waffle.io/zeroturnaround/configo.svg?label=ready&title=Waffle.io)](http://waffle.io/zeroturnaround/configo)

**Configo** helps running **12factor** (http://12factor.net/config) applications by loading environment variables from different sources.

# Configuration
See wiki for detailed explanation of configuration options, supported sources and more examples:
https://github.com/zeroturnaround/configo/wiki

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
  myAppImage
```
We also have a server to run some background jobs:
```bash
docker run \
  -e "DB_MONGO_URI=mongodb://user:pass@mongo.prod.domain.com/db" \
  -e "DB_REDIS_URI=redis://some.redis.prod.domain.com/0" \
  -e TWITTER_KEY=abcdefg \
  -e SEND_EMAILS=true \
  -e RUN_JOBS=true \
  myAppImage
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
  myAppImage

docker run \
  -e CONFIGO_SOURCE_0='{"type": "http", "format": "yaml", "url": "https://my.server.com/common.yaml"}' \
  -e CONFIGO_SOURCE_100='{"type": "http", "format": "yaml", "url": "https://my.server.com/jobs.yaml"}' \
  myAppImage
```
Since we added **Configo**, it will now load configuration from the sources we specified. In addition, it will merge these settings and configure the environment variables for your application.

# Thanks
* http://projects.spring.io/spring-cloud/ - For inspiration (See Spring Cloud Config).
* https://github.com/kelseyhightower/confd - For some ideas regarding sources.
* https://github.com/spf13/viper - For parses.
* https://github.com/hashicorp/terraform - For FlatMap implementation.
