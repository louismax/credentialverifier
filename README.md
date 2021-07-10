# credentialverifier

[简体中文](README.md)

> 证件号验证器，现已支持的证件号码：
> * 中国大陆居民身份证号码校验
> * 中国港澳居民居住证号码校验
> * 中国台湾居民居住证号码校验

[![gocover.io](https://gocover.io/_badge/github.com/louismax/credentialverifier)](https://gocover.io/github.com/louismax/credentialverifier)
[![Go Report Card](https://goreportcard.com/badge/github.com/louismax/credentialverifier)](https://goreportcard.com/report/github.com/louismax/credentialverifier)
[![GoDoc](https://godoc.org/github.com/louismax/credentialverifier?status.svg)](https://godoc.org/github.com/louismax/credentialverifier)
[![GitHub release](https://img.shields.io/github/tag/louismax/credentialverifier.svg)](https://github.com/louismax/credentialverifier/releases)
[![GitHub license](https://img.shields.io/github/license/louismax/credentialverifier.svg)](https://github.com/louismax/credentialverifier/blob/master/LICENSE)

## 功能

* 中国身份证号验证
* 获取身份证号信息
* 升级 15 位身份证号为 18 位
* 伪造符合校验的身份证号

## 环境要求

* Go >= 1.11

## 安装

``` shell script
$ go get -u github.com/louismax/credentialverifier
```

## 使用

``` go
package main

import (
	"github.com/louismax/credentialverifier"
	"gopkg.in/ffmt.v1"
)

func main() {

	// 验证身份证号合法性
	ffmt.P(credentialverifier.IsValid("610122198310134420", false)) // 非严格模式验证大陆居民身份证18位
	ffmt.P(credentialverifier.IsValid("610122198310134420", true))  // 严格模式验证大陆居民身份证18位
	ffmt.P(credentialverifier.IsValid("15000019970319128X", false)) // 大陆居民身份证末位是X18位
	ffmt.P(credentialverifier.IsValid("610104620927690", false))    // 大陆居民身份证15位
	ffmt.P(credentialverifier.IsValid("810000199408230021", false)) // 港澳居民居住证18位
	ffmt.P(credentialverifier.IsValid("830000199201300022", false)) // 台湾居民居住证18位

	// 获取身份证号信息
	ffmt.P(credentialverifier.GetInfo("610122198310134420", false)) // 非严格模式获取身份证号信息
	ffmt.P(credentialverifier.GetInfo("610122198310134420", true))  // 严格模式获取身份证号信息
	// []interface {}[
	// 	github.com/louismax/credentialverifier.IdInfo{          // 身份证号信息
	// 		AddressCode: int(440308)                           // 地址码
	// 		Abandoned:   int(0)                                // 地址码是否废弃：1为废弃的，0为正在使用的
	// 		Address:     string("广东省深圳市盐田区")             // 地址
	// 		AddressTree: []string[                             // 省市区三级列表
	//			string("广东省")                                    // 省
	//			string("深圳市")                                    // 市
	//			string("盐田区")                                    // 区
	//		]
	// 		Birthday:      <1999-01-10 00:00:00 +0000 UTC>     // 出生日期
	// 		Constellation: string("摩羯座")                     // 星座
	// 		ChineseZodiac: string("卯兔")                       // 生肖
	// 		Sex:           int(1)                              // 性别：1为男性，0为女性
	// 		Length:        int(18)                             // 号码长度
	// 		CheckBit:      string("2")                         // 校验码
	// 	}
	// 	<nil>                                              // 错误信息
	// ]

	// 生成可通过校验的假身份证号
	ffmt.P(credentialverifier.FakeId()) // 随机生成
	ffmt.P(credentialverifier.FakeRequireId(true, "台湾省", "200001", 1)) // 生成出生于2000年1月台湾省的男性居民居住证

	// 15位号码升级为18位
	ffmt.P(credentialverifier.UpgradeId("610104620927690"))
	// []interface {}[
	// 	string("610104196209276908") // 升级后号码
	// 	<nil>                        // 错误信息
	// ]
}
```

## 参考资料

* [中华人民共和国公民身份号码](https://zh.wikipedia.org/wiki/中华人民共和国公民身份号码)
* [中华人民共和国民政部：行政区划代码](http://www.mca.gov.cn/article/sj/xzqh/)
* [中华人民共和国行政区划代码历史数据集](https://github.com/jxlwqq/address-code-of-china)
* [国务院办公厅关于印发《港澳台居民居住证申领发放办法》的通知](http://www.gov.cn/zhengce/content/2018-08/19/content_5314865.htm)
* [港澳台居民居住证](https://zh.wikipedia.org/wiki/港澳台居民居住证)

## 协议

MIT 许可证（MIT）。有关更多信息，请参见[协议文件](LICENSE)。
