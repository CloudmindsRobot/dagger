# Quick Start

[TOC]

## 一、部署说明

### Docker-Compose

- 进入./deploy/docker-compose 目录编辑 docker-compose.yaml 文件

```
$ cat ./deploy/docker-compose/docker-compose.yaml

version: "3.8"
services:
  ui:
    build:
      context: .
      dockerfile: ./docker-compose/Dockerfile.ui
      target: ui
    image: quay.io/cloudminds/dagger-ui:0.1.0-beta
    container_name: dagger-ui
    depends_on:
      - backend
    ports:
      - "8080:8080"
    networks:
      - dagger
  backend:
    build:
      context: .
      dockerfile: ./docker-compose/Dockerfile.backend
      target: backend
    image: quay.io/cloudminds/dagger-backend:0.1.0-beta
    container_name: dagger-backend
    env:
      #loki服务地址
      LOKI_SERVER: http://1.1.1.1:3100
    ports:
      - "8000:8000"
    networks:
      - dagger
    command: ["sh", "-c", "./dagger"]
    volumes:
      - 'static_data:/usr/src/app/static:rw'
      - 'sqlite_data:/usr/src/app/db:rw'
volumes:
  sqlite_data:
    driver: local
  static_data:
    driver: local

networks:
  dagger:
    driver: bridge
```

- 启动服务

```
$ docker-compose up -d
```

### Helm

更新中...

## 二、使用说明

### 主界面

初始化部署完成后第一次登陆需要注册一个账号用于登陆系统，按照指示注册即可。

<img src="images/quickstart/login.jpg" width="35%" height="35%">

登陆 Dagger 后主界面如下

<img src="images/quickstart/dashboard.jpg" width="50%" height="50%">

### 查询日志

如果 dagger 和 loki 对接成功后，点击条件过滤栏会得到 loki 当前的所有日志流的 label，点击某一个 label 后会自动加载出对应的值

<img src="images/quickstart/labels.jpg" width="50%" height="50%">

### 过滤日志

对日志进行过滤可以直接在条件过滤栏中填入要过滤的`字符串`或者`正则表达式`,如果有多个管道匹配项则重复填入即可

<img src="images/quickstart/filter.jpg" width="50%" height="50%">

### 快速查询

对于需要经常查看的日志 label 组合，点击`保存条件`就可以存储在个人的查询记录里面，如果需要快速定位日志，点击`快速查询`会弹出过去保存的 label 组合，选择自己需要的即可

<img src="images/quickstart/history.jpg" width="50%" height="50%">

### 查询时间

dagger 默认查询时间为最近 5 分钟，需要自定义时间，点击`时间选择`会弹出时间选择器。提供`精细时间`和`最近时间`两种方式选择方式

<img src="images/quickstart/times.jpg" width="40%" height="40%">

### 日志输出限制

Dagger 默认对日志查询限制 2000 行限制，当前支持 500、1000、2000、5000 和 10000 行的日志输出。更大的日志输出行，我们建议结合日志过滤功能使用。

<img src="images/quickstart/live.jpg" width="50%" height="50%">

### 过滤日志级别

当前 dagger 默认支持五种日志级别的过滤，`INFO`,`DEBUG`,`WARN`,`ERROR`和`UNKNOWN`，单机打它们任何一个按钮即可切换到日志级别的视图。

<img src="images/quickstart/level.jpg" width="30%" height="30%">

他们的默认过滤规则如下：

- INFO： [I]、[info]、【info】、info、level=info
- DEBUG： [D]、[deug]、【debug】、debug、level=debug
- WARN： [W]、[warn]、[warning]、【warn】、【warning】、warn、warning、level=warn、level=warning
- ERROR： [E]、[error]、【error】、error、level=error
- UNKNOWN: 未匹配到的日志

> 以上过滤均不区分大小写

### 日志实时推送

当选择好日志的过滤规则后，点击界面右下角的绿色播放按钮即可打开日志实时推送功能，默认最新日志会染成黄色

<img src="images/quickstart/live.jpg" width="50%" height="50%">

### 日志下载

dagger 提供两种方式保留当前查询的日志

- 直接下载日志

点击又下角`+`号会弹出下载按钮，直接下载即可

<img src="images/quickstart/download.jpg" width="10%" height="10%">

- 将日志保存在 dagger

点击右上角`保存`即可将查询的日志的快照保存在 dagger 当中

<img src="images/quickstart/snapshot_saver.jpg" width="50%" height="50%">

### 查询历史

点击界面左上角的侧边栏管理器，可以进入`查询历史`界面，这里面保存用户的历次查询记录，你可以在这个界面点击查询快速切到你想查看的日志流当中。

<img src="images/quickstart/query_history.jpg" width="50%" height="50%">

### 日志快照

点击界面左上角的侧边栏管理器，可以进入`日志快照`界面，这里面保存了用户的日志快照。这里面你可以进行查看快照、快照下载、删除快照等管理工作。

<img src="images/quickstart/snapshot.jpg" width="50%" height="50%">
