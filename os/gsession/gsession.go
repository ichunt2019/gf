// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/ichunt2019/gcfg.

// Package gsession implements manager and storage features for sessions.
package gsession

import (
	"errors"
	"github.com/ichunt2019/gcfg/util/guid"
)

var (
	ErrorDisabled = errors.New("this feature is disabled in this storage")
)

// NewSessionId creates and returns a new and unique session id string,
// which is in 36 bytes.
func NewSessionId() string {
	return guid.S()
}
