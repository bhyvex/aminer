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
	"log"
	"net/http"
	"time"

	"github.com/minio/cli"
	"gopkg.in/mgo.v2/bson"
)

// SSLAnalytics
const (
	SSLAnalytics = "https://ssl.google-analytics.com/collect"
)

func updateGoogleAnalytics(c *configV1, referer, path string) error {
	var payload bytes.Buffer
	payload.WriteString("v=1")
	// Tracking id UA-XXXXXXXX-1
	payload.WriteString("&tid=" + c.TID)
	// User unique id UUID
	payload.WriteString("&cid=" + c.CID)
	// Type of hit
	payload.WriteString("&t=pageview")
	// Data source
	payload.WriteString("&ds=web")
	// data referrer
	payload.WriteString("&dr=" + mustURLEncodeName(referer))
	// document path
	payload.WriteString("&dp=" + mustURLEncodeName(path))
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

func runAnalyticsCmd(c *cli.Context) {
	conf, err := loadConfigV1()
	if err != nil {
		log.Fatal(err)
	}
	s := connectToMongo(c)
	defer s.Close()

	var result LogMessage
	iter := db.Find(bson.M{"http.request.method": "GET"}).Iter()
	for iter.Next(&result) {
		if time.Since(result.StartTime) < time.Duration(24*time.Hour) {
			err = updateGoogleAnalytics(conf, result.HTTP.Request.Referer(), result.HTTP.Request.RequestURI)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}
