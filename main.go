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

var (
	theHost *snowflake.Host
)

func main() {
	nodeIDStr := os.Getenv(EnvNodeId)
	nodeID, _ := strconv.ParseInt(nodeIDStr, 10, 32)
	var err error = nil
	theHost, err = snowflake.NewHost(nodeID)
	if err != nil {
		panic(err)
	}
	log.Printf("New host as node id: %v", nodeID)
	http.HandleFunc("/GenID", HttpGenID)
	http.ListenAndServe(fmt.Sprintf(":%v", Port), nil)
}

func HttpGenID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	buffer := bytes.NewBuffer(nil)
	id := theHost.Generate()
	buffer.WriteString(id.ToString())
	buffer.WriteRune('\n')
	buffer.WriteString(
		fmt.Sprintf("%v | %v | %v",
			id.Timestamp(), id.NodeID(), id.SeqID()))
	log.Printf("generate new id: %v", id.ToString())
	w.Write(buffer.Bytes())
}
