// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/ichunt2019/gcfg.

package utils

import (
	"github.com/ichunt2019/gcfg/internal/command"
)

const (
	debugKey                 = "gf.debug"             // Debug key for checking if in debug mode.
	StackFilterKeyForGoFrame = "/github.com/ichunt2019/gcfg/" // Stack filtering key for all GoFrame module paths.
)

var (
	// isDebugEnabled marks whether GoFrame debug mode is enabled.
	isDebugEnabled = false
)

func init() {
	// Debugging configured.
	value := command.GetOptWithEnv(debugKey)
	if value == "" || value == "0" || value == "false" {
		isDebugEnabled = false
	} else {
		isDebugEnabled = true
	}
}

// IsDebugEnabled checks and returns whether debug mode is enabled.
// The debug mode is enabled when command argument "gf.debug" or environment "GF_DEBUG" is passed.
func IsDebugEnabled() bool {
	return isDebugEnabled
}
