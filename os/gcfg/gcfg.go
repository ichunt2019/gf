// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/ichunt2019/gf.

// Package gcfg provides reading, caching and managing for configuration.
package gcfg

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/ichunt2019/gf/internal/intlog"
	"github.com/ichunt2019/gf/os/gcmd"
	"github.com/ichunt2019/gf/text/gstr"

	"github.com/ichunt2019/gf/os/gres"

	"github.com/ichunt2019/gf/container/garray"
	"github.com/ichunt2019/gf/container/gmap"
	"github.com/ichunt2019/gf/encoding/gjson"
	"github.com/ichunt2019/gf/os/gfile"
	"github.com/ichunt2019/gf/os/gfsnotify"
	"github.com/ichunt2019/gf/os/glog"
	"github.com/ichunt2019/gf/os/gspath"
)

const (
	DefaultConfigFile = "config.toml" // The default configuration file name.
	cmdEnvKey         = "gf.gcfg"     // Configuration key for command argument or environment.
)

// Configuration struct.
type Config struct {
	defaultName   string           // Default configuration file name.
	searchPaths   *garray.StrArray // Searching path array.
	jsonMap       *gmap.StrAnyMap  // The pared JSON objects for configuration files.
	violenceCheck bool             // Whether do violence check in value index searching. It affects the performance when set true(false in default).
}

var (
	supportedFileTypes = []string{"toml", "yaml", "json", "ini", "xml"}
	resourceTryFiles   = []string{"", "/", "config/", "config", "/config", "/config/"}
)

// New returns a new configuration management object.
// The parameter <file> specifies the default configuration file name for reading.
func New(file ...string) *Config {
	name := DefaultConfigFile
	if len(file) > 0 {
		name = file[0]
	} else {
		// Custom default configuration file name from command line or environment.
		if customFile := gcmd.GetOptWithEnv(fmt.Sprintf("%s.file", cmdEnvKey)).String(); customFile != "" {
			name = customFile
		}
	}
	c := &Config{
		defaultName: name,
		searchPaths: garray.NewStrArray(true),
		jsonMap:     gmap.NewStrAnyMap(true),
	}
	// Customized dir path from env/cmd.
	if customPath := gcmd.GetOptWithEnv(fmt.Sprintf("%s.path", cmdEnvKey)).String(); customPath != "" {
		if gfile.Exists(customPath) {
			_ = c.SetPath(customPath)
		} else {
			if errorPrint() {
				glog.Errorf("Configuration directory path does not exist: %s", customPath)
			}
		}
	} else {
		// Dir path of main package.
		if mainPath := gfile.MainPkgPath(); mainPath != "" && gfile.Exists(mainPath) {
			_ = c.AddPath(mainPath)
		}
		// Dir path of binary.
		if selfPath := gfile.SelfDir(); selfPath != "" && gfile.Exists(selfPath) {
			_ = c.AddPath(selfPath)
		}
		// Dir path of working dir.
		_ = c.AddPath(gfile.Pwd())
	}
	return c
}

func (c *Config) getSearchPaths() []string {
	var (
		searchPaths = c.searchPaths.Slice()
		mainPkgPath = gfile.MainPkgPath()
	)
	if mainPkgPath != "" {
		if !gstr.InArray(searchPaths, mainPkgPath) {
			searchPaths = append([]string{mainPkgPath}, searchPaths...)
		}
	}
	return searchPaths
}

// filePath returns the absolute configuration file path for the given filename by <file>.
func (c *Config) filePath(file ...string) (path string) {
	name := c.defaultName
	if len(file) > 0 {
		name = file[0]
	}
	path = c.FilePath(name)
	if path == "" {
		var (
			searchPaths = c.getSearchPaths()
			buffer      = bytes.NewBuffer(nil)
		)

		if len(searchPaths) > 0 {
			buffer.WriteString(fmt.Sprintf("[gcfg] cannot find config file \"%s\" in following paths:", name))
			index := 1
			for _, v := range searchPaths {
				v = gstr.TrimRight(v, `\/`)
				buffer.WriteString(fmt.Sprintf("\n%d. %s", index, v))
				index++
				buffer.WriteString(fmt.Sprintf("\n%d. %s", index, v+gfile.Separator+"config"))
				index++
			}
		} else {
			buffer.WriteString(fmt.Sprintf("[gcfg] cannot find config file \"%s\" with no path set/add", name))
		}
		if errorPrint() {
			glog.Error(buffer.String())
		}
	}
	return path
}

// SetPath sets the configuration directory path for file search.
// The parameter <path> can be absolute or relative path,
// but absolute path is strongly recommended.
func (c *Config) SetPath(path string) error {
	var (
		isDir    = false
		realPath = ""
	)
	if file := gres.Get(path); file != nil {
		realPath = path
		isDir = file.FileInfo().IsDir()
	} else {
		// Absolute path.
		realPath = gfile.RealPath(path)
		if realPath == "" {
			// Relative path.
			c.searchPaths.RLockFunc(func(array []string) {
				for _, v := range array {
					if path, _ := gspath.Search(v, path); path != "" {
						realPath = path
						break
					}
				}
			})
		}
		if realPath != "" {
			isDir = gfile.IsDir(realPath)
		}
	}
	// Path not exist.
	if realPath == "" {
		buffer := bytes.NewBuffer(nil)
		if c.searchPaths.Len() > 0 {
			buffer.WriteString(fmt.Sprintf("[gcfg] SetPath failed: cannot find directory \"%s\" in following paths:", path))
			c.searchPaths.RLockFunc(func(array []string) {
				for k, v := range array {
					buffer.WriteString(fmt.Sprintf("\n%d. %s", k+1, v))
				}
			})
		} else {
			buffer.WriteString(fmt.Sprintf(`[gcfg] SetPath failed: path "%s" does not exist`, path))
		}
		err := errors.New(buffer.String())
		if errorPrint() {
			glog.Error(err)
		}
		return err
	}
	// Should be a directory.
	if !isDir {
		err := fmt.Errorf(`[gcfg] SetPath failed: path "%s" should be directory type`, path)
		if errorPrint() {
			glog.Error(err)
		}
		return err
	}
	// Repeated path check.
	if c.searchPaths.Search(realPath) != -1 {
		return nil
	}
	c.jsonMap.Clear()
	c.searchPaths.Clear()
	c.searchPaths.Append(realPath)
	intlog.Print("SetPath:", realPath)
	return nil
}

// SetViolenceCheck sets whether to perform hierarchical conflict checking.
// This feature needs to be enabled when there is a level symbol in the key name.
// It is off in default.
//
// Note that, turning on this feature is quite expensive, and it is not recommended
// to allow separators in the key names. It is best to avoid this on the application side.
func (c *Config) SetViolenceCheck(check bool) {
	c.violenceCheck = check
	c.Clear()
}

// AddPath adds a absolute or relative path to the search paths.
func (c *Config) AddPath(path string) error {
	var (
		isDir    = false
		realPath = ""
	)
	// It firstly checks the resource manager,
	// and then checks the filesystem for the path.
	if file := gres.Get(path); file != nil {
		realPath = path
		isDir = file.FileInfo().IsDir()
	} else {
		// Absolute path.
		realPath = gfile.RealPath(path)
		if realPath == "" {
			// Relative path.
			c.searchPaths.RLockFunc(func(array []string) {
				for _, v := range array {
					if path, _ := gspath.Search(v, path); path != "" {
						realPath = path
						break
					}
				}
			})
		}
		if realPath != "" {
			isDir = gfile.IsDir(realPath)
		}
	}
	if realPath == "" {
		buffer := bytes.NewBuffer(nil)
		if c.searchPaths.Len() > 0 {
			buffer.WriteString(fmt.Sprintf("[gcfg] AddPath failed: cannot find directory \"%s\" in following paths:", path))
			c.searchPaths.RLockFunc(func(array []string) {
				for k, v := range array {
					buffer.WriteString(fmt.Sprintf("\n%d. %s", k+1, v))
				}
			})
		} else {
			buffer.WriteString(fmt.Sprintf(`[gcfg] AddPath failed: path "%s" does not exist`, path))
		}
		err := errors.New(buffer.String())
		if errorPrint() {
			glog.Error(err)
		}
		return err
	}
	if !isDir {
		err := fmt.Errorf(`[gcfg] AddPath failed: path "%s" should be directory type`, path)
		if errorPrint() {
			glog.Error(err)
		}
		return err
	}
	// Repeated path check.
	if c.searchPaths.Search(realPath) != -1 {
		return nil
	}
	c.searchPaths.Append(realPath)
	intlog.Print("AddPath:", realPath)
	return nil
}

// GetFilePath returns the absolute path of the specified configuration file.
// If <file> is not passed, it returns the configuration file path of the default name.
// If the specified configuration file does not exist,
// an empty string is returned.
func (c *Config) FilePath(file ...string) (path string) {
	name := c.defaultName
	if len(file) > 0 {
		name = file[0]
	}
	searchPaths := c.getSearchPaths()
	// Searching resource manager.
	if !gres.IsEmpty() {
		for _, v := range resourceTryFiles {
			if file := gres.Get(v + name); file != nil {
				path = file.Name()
				return
			}
		}
		for _, prefix := range searchPaths {
			for _, v := range resourceTryFiles {
				if file := gres.Get(prefix + v + name); file != nil {
					path = file.Name()
					return
				}
			}
		}
	}
	// Searching the file system.
	for _, prefix := range searchPaths {
		prefix = gstr.TrimRight(prefix, `\/`)
		if path, _ = gspath.Search(prefix, name); path != "" {
			return
		}
		if path, _ = gspath.Search(prefix+gfile.Separator+"config", name); path != "" {
			return
		}
	}
	return
}

// SetFileName sets the default configuration file name.
func (c *Config) SetFileName(name string) *Config {
	c.defaultName = name
	return c
}

// GetFileName returns the default configuration file name.
func (c *Config) GetFileName() string {
	return c.defaultName
}

// Available checks and returns whether configuration of given <file> is available.
func (c *Config) Available(file ...string) bool {
	var name string
	if len(file) > 0 && file[0] != "" {
		name = file[0]
	} else {
		name = c.defaultName
	}
	if c.FilePath(name) != "" {
		return true
	}
	if GetContent(name) != "" {
		return true
	}
	return false
}

// getJson returns a *gjson.Json object for the specified <file> content.
// It would print error if file reading fails. It return nil if any error occurs.
func (c *Config) getJson(file ...string) *gjson.Json {
	var name string
	if len(file) > 0 && file[0] != "" {
		name = file[0]
	} else {
		name = c.defaultName
	}
	r := c.jsonMap.GetOrSetFuncLock(name, func() interface{} {
		var (
			content  = ""
			filePath = ""
		)
		// The configured content can be any kind of data type different from its file type.
		isFromConfigContent := true
		if content = GetContent(name); content == "" {
			isFromConfigContent = false
			filePath = c.filePath(name)
			if filePath == "" {
				return nil
			}
			if file := gres.Get(filePath); file != nil {
				content = string(file.Content())
			} else {
				content = gfile.GetContents(filePath)
			}
		}
		// Note that the underlying configuration json object operations are concurrent safe.
		var (
			j   *gjson.Json
			err error
		)
		dataType := gfile.ExtName(name)
		if gjson.IsValidDataType(dataType) && !isFromConfigContent {
			j, err = gjson.LoadContentType(dataType, content, true)
		} else {
			j, err = gjson.LoadContent(content, true)
		}
		if err == nil {
			j.SetViolenceCheck(c.violenceCheck)
			// Add monitor for this configuration file,
			// any changes of this file will refresh its cache in Config object.
			if filePath != "" && !gres.Contains(filePath) {
				_, err = gfsnotify.Add(filePath, func(event *gfsnotify.Event) {
					c.jsonMap.Remove(name)
				})
				if err != nil && errorPrint() {
					glog.Error(err)
				}
			}
			return j
		} else {
			if errorPrint() {
				if filePath != "" {
					glog.Criticalf(`[gcfg] Load config file "%s" failed: %s`, filePath, err.Error())
				} else {
					glog.Criticalf(`[gcfg] Load configuration failed: %s`, err.Error())
				}
			}
		}
		return nil
	})
	if r != nil {
		return r.(*gjson.Json)
	}
	return nil
}
