// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/ichunt2019/gcfg.

// Package gbuild manages the build-in variables from "gf build".
package gbuild

import (
	"github.com/ichunt2019/gcfg"
	"github.com/ichunt2019/gcfg/container/gvar"
	"github.com/ichunt2019/gcfg/encoding/gbase64"
	"github.com/ichunt2019/gcfg/internal/intlog"
	"github.com/ichunt2019/gcfg/internal/json"
	"github.com/ichunt2019/gcfg/util/gconv"
	"runtime"
)

var (
	builtInVarStr = ""                       // Raw variable base64 string.
	builtInVarMap = map[string]interface{}{} // Binary custom variable map decoded.
)

func init() {
	if builtInVarStr != "" {
		err := json.Unmarshal(gbase64.MustDecodeString(builtInVarStr), &builtInVarMap)
		if err != nil {
			intlog.Error(err)
		}
		builtInVarMap["gfVersion"] = gf.VERSION
		builtInVarMap["goVersion"] = runtime.Version()
		intlog.Printf("build variables: %+v", builtInVarMap)
	} else {
		intlog.Print("no build variables")
	}
}

// Info returns the basic built information of the binary as map.
// Note that it should be used with gf-cli tool "gf build",
// which injects necessary information into the binary.
func Info() map[string]string {
	return map[string]string{
		"gf":   GetString("gfVersion"),
		"go":   GetString("goVersion"),
		"git":  GetString("builtGit"),
		"time": GetString("builtTime"),
	}
}

// Get retrieves and returns the build-in binary variable with given name.
func Get(name string, def ...interface{}) interface{} {
	if v, ok := builtInVarMap[name]; ok {
		return v
	}
	if len(def) > 0 {
		return def[0]
	}
	return nil
}

// GetVar retrieves and returns the build-in binary variable of given name as gvar.Var.
func GetVar(name string, def ...interface{}) *gvar.Var {
	return gvar.New(Get(name, def...))
}

// GetString retrieves and returns the build-in binary variable of given name as string.
func GetString(name string, def ...interface{}) string {
	return gconv.String(Get(name, def...))
}

// Map returns the custom build-in variable map.
func Map() map[string]interface{} {
	return builtInVarMap
}
