// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/ichunt2019/gcfg.

package gcron_test

import (
	"time"

	"github.com/ichunt2019/gcfg/os/gcron"
	"github.com/ichunt2019/gcfg/os/glog"
)

func Example_cronAddSingleton() {
	gcron.AddSingleton("* * * * * *", func() {
		glog.Println("doing")
		time.Sleep(2 * time.Second)
	})
	select {}
}
