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
	"bytes"
	"errors"
	"net/http"
)

const (
	SSLAnalytics = "https://ssl.google-analytics.com/collect"
)

type config struct {
	tid string
	cid string
}

func postAnalytics(c config) error {
	var payload bytes.Buffer
	payload.WriteString("v=1")
	payload.WriteString("&tid=" + c.tid)
	payload.WriteString("&cid=" + c.cid)
	payload.WriteString("&t=event")
	payload.WriteString("&ds=web")

	req, err := http.NewRequest("POST", SSLAnalytics, &payload)
	if err != nil {
		return err
	}

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return errors.New("Data was not uploaded error: " + resp.Status)
	}

	return nil
}
