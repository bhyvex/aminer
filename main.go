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
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/minio/cli"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// LogMessage is a serializable json log message
type LogMessage struct {
	StartTime     time.Time
	Duration      time.Duration
	StatusMessage string // human readable http status message
	ContentLength string // human readable content length

	// HTTP detailed message
	HTTP struct {
		ResponseHeaders http.Header
		Request         *http.Request
	}
}

var session *mgo.Session

func connectToMongo(c *cli.Context) {
	session, err := mgo.Dial(c.GlobalString("server"))
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)
}

func runPopulateCmd(c *cli.Context) {
	f, err := os.Open(c.Args()[0])
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	connectToMongo(c)
	// make this configurable
	m := session.DB("test").C("downloads")
	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)

	var message LogMessage
	for scanner.Scan() {
		json.Unmarshal([]byte(scanner.Text()), &message)
		err = m.Insert(&message)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func runFindCmd(c *cli.Context) {
	connectToMongo(c)
	// make this configurable
	m := session.DB("test").C("downloads")
	result := LogMessage{}
	iter := m.Find(bson.M{"http.request.method": "GET"}).Iter()
	for iter.Next(&result) {
		if result.HTTP.Request.Method != "GET" {
			continue
		}
		if strings.Contains(result.HTTP.Request.RemoteAddr, "50.204.118.154") {
			continue
		}
		if strings.Contains(result.HTTP.Request.RemoteAddr, "10.134.253.170") {
			continue
		}
		fmt.Print(result.HTTP.Request.Method)
		fmt.Print("    ")
		fmt.Print(result.HTTP.Request.RemoteAddr)
		fmt.Print("    ")
		fmt.Print(result.HTTP.Request.RequestURI)
		fmt.Println("    ")
	}
}

var commands = []cli.Command{
	findCmd,
	populateCmd,
}

var findCmd = cli.Command{
	Name:        "find",
	Description: "find all documents for a map",
	Action:      runFindCmd,
}

var populateCmd = cli.Command{
	Name:        "populate",
	Description: "populate your mongo instance with new data",
	Action:      runPopulateCmd,
}

var flags = []cli.Flag{
	cli.StringFlag{
		Name:  "server",
		Value: "localhost",
		Usage: "IP/HOSTNAME of your mongodb instance",
	},
}

func main() {
	app := cli.NewApp()
	app.Usage = "A miner for your minio access logs"
	app.Version = "0.0.1"
	app.Commands = commands
	app.Flags = flags
	app.Author = "Minio.io"

	app.RunAndExitOnError()
}
