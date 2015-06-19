/*
 * Aminer (C) 2014, 2015 Minio, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"os/user"
	"path/filepath"
	"runtime"
	"sync"

	"github.com/minio/mc/pkg/quick"
)

// configV1
type configV1 struct {
	Version string
	Tid     string
	Cid     string
}

// cached variables should *NEVER* be accessed directly from outside this file.
var cache sync.Pool

func getConfigPath() (string, error) {
	u, err := user.Current()
	if err != nil {
		return "", err
	}
	// For windows the path is slightly different
	switch runtime.GOOS {
	case "windows":
		return filepath.Join(u.HomeDir, "miner\\miner.json"), nil
	default:
		return filepath.Join(u.HomeDir, ".miner/miner.json"), nil
	}
}

func loadConfigV1() (*configV1, error) {
	configFile, err := getConfigPath()
	if err != nil {
		return nil, err
	}
	// Cached in private global variable.
	if v := cache.Get(); v != nil { // Use previously cached config.
		return v.(quick.Config).Data().(*configV1), nil
	}
	conf := newConfigV1()
	qconf, err := quick.New(conf)
	if err != nil {
		return nil, err
	}
	err = qconf.Load(configFile)
	if err != nil {
		return nil, err
	}
	cache.Put(qconf)
	return qconf.Data().(*configV1), nil
}

func newConfigV1() *configV1 {
	conf := new(configV1)
	conf.Version = "0.0.1"
	return conf
}
