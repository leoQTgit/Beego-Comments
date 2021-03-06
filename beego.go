// Copyright 2014 beego Author. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package beego

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	// VERSION represent beego web framework version.
	VERSION = "1.6.1"
	// 定义程序运行的两种模式
	// 开发模式
	DEV = "dev"
	// 生产使用模式
	PROD = "prod"
)

//hook function to run
type hookfunc func() error

var (
	hooks = make([]hookfunc, 0) //保存 hookfunc的 hooks切片
)

// AddAPPStartHook is used to register the hookfunc
// The hookfuncs will run in beego.Run()
// such as sessionInit, middlerware start, buildtemplate, admin start
func AddAPPStartHook(hf hookfunc) {
	hooks = append(hooks, hf)
}

// Beego 应用的的入口函数,初始化了hooks函数,判断beego.Run()是否为下列几种情况
// beego.Run() default run on HttpPort
// beego.Run("localhost")
// beego.Run(":8089")
// beego.Run("127.0.0.1:8089")
func Run(params ...string) {
	initBeforeHTTPRun()

	if len(params) > 0 && params[0] != "" {
		strs := strings.Split(params[0], ":")
		if len(strs) > 0 && strs[0] != "" {
			BConfig.Listen.HTTPAddr = strs[0]
		}
		if len(strs) > 1 && strs[1] != "" {
			BConfig.Listen.HTTPPort, _ = strconv.Atoi(strs[1])
		}
	}
	//单例模式， BeeApp只有一个,定义并初始化于 app.go文件中
	BeeApp.Run()
}

func initBeforeHTTPRun() {
	//讲每个函数添加到 hooks中,并且逐个调用
	AddAPPStartHook(registerMime)
	AddAPPStartHook(registerDefaultErrorHandler)
	AddAPPStartHook(registerSession)
	AddAPPStartHook(registerDocs)
	AddAPPStartHook(registerTemplate)
	AddAPPStartHook(registerAdmin)

	for _, hk := range hooks {
		if err := hk(); err != nil {
			panic(err)
		}
	}
}

// TestBeegoInit is for test package init
func TestBeegoInit(ap string) {
	os.Setenv("BEEGO_RUNMODE", "test")
	appConfigPath = filepath.Join(ap, "conf", "app.conf")
	os.Chdir(ap)
	initBeforeHTTPRun()
}
