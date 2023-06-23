package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/xiaobao520123/onlyID/snowflake"
)

const (
	EnvNodeId = "NODE_ID"
	Port      = 80
)

const (
	ApiGenId = "/GenID"
)

var (
	theHost *snowflake.Host
)

func main() {
	// Read host's node id from env variables
	nodeIDStr := os.Getenv(EnvNodeId)
	nodeID, _ := strconv.ParseInt(nodeIDStr, 10, 32)
	var err error = nil

	// Create a new host
	theHost, err = snowflake.NewHost(nodeID)
	if err != nil {
		panic(err)
	}
	log.Printf("New host as node id: %v", nodeID)

	// Start HTTP Server
	http.HandleFunc(ApiGenId, HttpGenID)
	http.ListenAndServe(fmt.Sprintf(":%v", Port), nil)
}

func HttpGenID(w http.ResponseWriter, r *http.Request) {
	// Works only in method GET
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	buffer := bytes.NewBuffer(nil)

	// Generate an ID
	id := theHost.Generate()

	// Output ID
	buffer.WriteString(id.ToString())
	buffer.WriteRune('\n')

	// Output the parts of ID
	buffer.WriteString(
		fmt.Sprintf("%v | %v | %v",
			id.Timestamp(), id.NodeID(), id.SeqID()))

	log.Printf("generate new id: %v", id.ToString())

	// Write HTTP response
	w.Write(buffer.Bytes())
}
