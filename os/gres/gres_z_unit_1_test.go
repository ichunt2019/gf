// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/ichunt2019/gcfg.

package gres_test

import (
	"github.com/ichunt2019/gcfg/os/gtime"
	"strings"
	"testing"

	"github.com/ichunt2019/gcfg/os/gfile"

	"github.com/ichunt2019/gcfg/debug/gdebug"
	"github.com/ichunt2019/gcfg/os/gres"
	"github.com/ichunt2019/gcfg/test/gtest"
)

func Test_PackToGoFile(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		srcPath := gdebug.TestDataPath("files")
		goFilePath := gdebug.TestDataPath("testdata.go")
		pkgName := "testdata"
		err := gres.PackToGoFile(srcPath, goFilePath, pkgName)
		t.Assert(err, nil)
	})
}

func Test_Pack(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		srcPath := gdebug.TestDataPath("files")
		data, err := gres.Pack(srcPath)
		t.Assert(err, nil)

		r := gres.New()

		err = r.Add(string(data))
		t.Assert(err, nil)
		t.Assert(r.Contains("files/"), true)
	})
}

func Test_PackToFile(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		srcPath := gdebug.TestDataPath("files")
		dstPath := gfile.TempDir(gtime.TimestampNanoStr())
		err := gres.PackToFile(srcPath, dstPath)
		t.Assert(err, nil)

		defer gfile.Remove(dstPath)

		r := gres.New()
		err = r.Load(dstPath)
		t.Assert(err, nil)
		t.Assert(r.Contains("files"), true)
	})
}

func Test_PackMulti(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		srcPath := gdebug.TestDataPath("files")
		goFilePath := gdebug.TestDataPath("data/data.go")
		pkgName := "data"
		array, err := gfile.ScanDir(srcPath, "*", false)
		t.Assert(err, nil)
		err = gres.PackToGoFile(strings.Join(array, ","), goFilePath, pkgName)
		t.Assert(err, nil)
	})
}

func Test_PackWithPrefix1(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		srcPath := gdebug.TestDataPath("files")
		goFilePath := gfile.TempDir("testdata.go")
		pkgName := "testdata"
		err := gres.PackToGoFile(srcPath, goFilePath, pkgName, "www/gf-site/test")
		t.Assert(err, nil)
	})
}

func Test_PackWithPrefix2(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		srcPath := gdebug.TestDataPath("files")
		goFilePath := gfile.TempDir("testdata.go")
		pkgName := "testdata"
		err := gres.PackToGoFile(srcPath, goFilePath, pkgName, "/var/www/gf-site/test")
		t.Assert(err, nil)
	})
}
