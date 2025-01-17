# Captcha  #

[![Carbon Release](https://img.shields.io/github/release/golang-module/base64Captcha.svg)](https://github.com/golang-module/base64Captcha/releases)
[![Go Test](https://github.com/golang-module/base64Captcha/actions/workflows/test.yml/badge.svg)](https://github.com/golang-module/base64Captcha/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/golang-module/base64Captcha)](https://goreportcard.com/report/github.com/golang-module/base64Captcha)
[![Carbon Doc](https://img.shields.io/badge/go.dev-reference-brightgreen?logo=go&logoColor=white&style=flat)](https://pkg.go.dev/github.com/golang-module/base64Captcha)
[![License](https://img.shields.io/github/license/golang-module/base64Captcha)](https://github.com/golang-module/base64Captcha/blob/master/LICENSE)

#### 项目简介

一个轻量级、语义化、对开发者友好的 `golang` 验证码库，支持任何 Unicode 字符，可以轻松定制以支持语音、数学、中文、韩语、日语等

> Fork 于 [mojotv/base64Captcha](https://github.com/mojotv/base64Captcha)，由于原仓库不再更新和维护，所以拉取了一个分支，在此对原作者表示感谢。
#### 仓库地址

[github.com/golang-module/base64Captcha](https://github.com/golang-module/base64Captcha "github.com/golang-module/base64Captcha")

#### 安装使用

##### go version >= 1.16

```go
go get -u github.com/golang-module/base64Captcha
```

#### 用法示例

```go
import "github.com/golang-module/base64Captcha/store

// 使用内存存储(默认)
base64Store := store.DefaultMemoryStore

// 使用 sync map 存储
base64Store := store.DefaultSyncMapStore
```

##### 生成纯数字验证码
```go
// 使用默认配置
base64Driver := driver.DefaultDriverDigit

// 使用自定义配置
base64Driver := driver.NewDriverDigit(driver.DriverDigit{
    Width:    240, // 宽度
    Height:   80,  // 高度
    Length:   6,   // 长度
    MaxSkew:  0.7, // 随机弧度
    DotCount: 80,  // 点数量
})
```

##### 生成、验证
```go
captcha := base64Captcha.NewCaptcha(base64Driver, base64Store)

// 生成验证码
id, src, answer, err = captcha.Generate()

// 验证验证码
captcha.Verify(id, answer, true)
```
