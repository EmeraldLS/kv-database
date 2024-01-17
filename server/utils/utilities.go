package utils

import (
	"encoding/json"
	"log"
	"net"

	"github.com/EmeraldLS/kv-db/server/model"
)

func SendResponse(respMsg *model.DefaultResponseFormat, conn net.Conn) {
	respMsgByte, err := json.Marshal(respMsg)
	if err != nil {
		log.Fatal(err)
	}
	_, err = conn.Write(respMsgByte)
	if err != nil {
		log.Fatal(err)
	}
}
