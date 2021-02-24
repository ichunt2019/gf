// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/ichunt2019/gcfg.

package gconv_test

import (
	"github.com/ichunt2019/gcfg/test/gtest"
	"github.com/ichunt2019/gcfg/util/gconv"
	"testing"
)

func Test_Unsafe(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s := "I love 小泽玛利亚"
		t.AssertEQ(gconv.UnsafeStrToBytes(s), []byte(s))
	})

	gtest.C(t, func(t *gtest.T) {
		b := []byte("I love 小泽玛利亚")
		t.AssertEQ(gconv.UnsafeBytesToStr(b), string(b))
	})
}
