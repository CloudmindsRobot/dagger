<img align="center" width="200" height="160" src="docs/logo.png">

[![Build Status](https://github.com/CloudmindsRobot/dagger/workflows/build/badge.svg)](https://github.com/CloudmindsRobot/dagger/actions)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![Go Report Card](https://goreportcard.com/badge/github.com/CloudmindsRobot/dagger)](https://goreportcard.com/report/github.com/CloudmindsRobot/dagger)

## [README of English](https://github.com/CloudmindsRobot/dagger/blob/main/README-EN.md)

# Dagger - A Log management system based on loki

Dagger 是一个基于 Loki 的日志查询和管理系统，它是由达闼科技（ CloudMinds ）云团队的`大禹基础设施平台`派生出来的一个项目。Dagger 运行在 Loki 前端，具备日志查询、搜索，保存和下载等特性，适用于云原生场景下的容器日志管理场景。

前端基于 vue.js 开发，使用 vuetify ui 框架进行设计，采用 axios 异步框架进行前后端交互。

后端基于 go gin 框架开发，gorm 作为数据框架，sqlite 作为数据存储端，采用 go websocket 桥接的方式进行实时数据处理。

本地快速部署请参见[Quick Start](#jump)

# ScreenShots

<img src="docs/screenshot.gif">

# Features

- 日志查询，简单的标签查询，无需复杂的 LogQL 语法
- 支持查询标签保存和快速查询，同时支持查询历史
- 支持日志查询结果快照保存，防止 loki 日志轮转后找不到记录
- 支持日志实时播放和下载功能

# RoadMap

- **日志告警**

  - 支持界面添加日志 rules 编辑、保存和推送
  - 兼容 AlertManager API，接受 Loki Ruler 的告警事件推送
  - 日志事件聚合，分析和告警
  - 支持邮件、阿里云 SMS、阿里云 Voice 告警

- **度量**
  - 支持 LogQL 的方式在前端添加自定义日志度量

# Usage

完整演示视频（[youtube](https://youtu.be/1qc8_nZA_dM)、[bilibili](https://www.bilibili.com/video/BV1Jr4y1w7qz/)）

# Deployment

- <span id = "jump">[Quick Start](docs/quick_start.md)</span>
  - [日志客户端 Fluentd 配置](docs/fluentd_config.md)
  - 日志客户端 promtail 配置
  - [Loki 建议](docs/Loki_best_practice.md)

* 分布式部署

# 1.0 upgrade to 2.0

- [1.0 upgrade to 2.0](https://github.com/CloudmindsRobot/dagger/releases/tag/2.0.1-alpha)

# Release Notes

- [版本发布历史](https://github.com/CloudmindsRobot/dagger/releases)

# Q&A

# Support

扫描二维码关注公众号回复【入群】加入微信讨论组

<img align="left" width="200" height="200" src="docs/qrcode.jpg">
