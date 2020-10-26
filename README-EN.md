<img align="center" width="200" height="160" src="docs/logo.png">

[![Build Status](https://github.com/CloudmindsRobot/dagger/workflows/build/badge.svg)](https://github.com/CloudmindsRobot/dagger/actions)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

## [README of Chinese/中文文档](https://github.com/CloudmindsRobot/dagger/blob/main/README.md)
# Dagger - A Log management system based on loki

Dagger is a Cloud Native log query and management system based on Grafana Loki. 

It's developed by Cloudminds Inc Cloud team, evolved from Cloudminds Devops platform "DaYu", and it has been used in production environment for more than half a year in Cloudminds.

Dagger runs on the top of Grafana Loki as the user interface, which provides multiple containers log query(mutiple labels supported), search(regex supported) and download(mutiple containers) functions.

Dagger is made for Kubernetes and docker mutiple containers log query requirement. It's a Cloud Native log management production. The main advantage of Dagger is that it's really light weighted, easy to use and dev-friendly.

Dagger frontend is developed with vue.js, designed by vuetify ui framework. Axios framework is used for async data transfer between frontend and backend.
Dagger backend is developed with go gin framework, as gorm is used for data framework, sqlite for metadata database and go websocket for realtime data process.

A quick start for local deployment: [Quick Start](#jump)

# ScreenShots

<img src="docs/screenshot.gif">

# Features

- **Manage mutiple Grafana Loki Instances**

  - Mutiple Loki instances could be managed in one Dagger, and Log query request could be dealt accrossing mutiple Loki instances.

- **Log Alarm**

  - Support adding log metric with LogQL in Dagger web interface
  - Support Loki-Ruler. Ruler could be edit and pushed in Dagger web interface
  - Support Log Alarm event subscription, aggregation and muti ways notification.

- **Easy to deploy**
  - One click to deploy Dagger, Loki and mutiple log client

# Usage

Full demo vedio:（[youtube](https://youtu.be/1qc8_nZA_dM)、[bilibili](https://www.bilibili.com/video/BV1Jr4y1w7qz/)）

# Deployment

- <span id = "jump">[Quick Start](docs/quick_start.md)</span>
  - [Log collection client Fluentd config](docs/fluentd_config.md)
  - Log collection client Promtail config
  - [Loki suggestions](docs/Loki_best_practice.md)

* Distributed Deployment

# Release Notes

# Q&A

# Support

You can scan this Wechat qrcode blew to follow our Wechat account and reply "Dagger" to join our discussion group

<img align="left" width="200" height="200" src="docs/qrcode.jpg">
