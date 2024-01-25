package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"sync"

	"github.com/EmeraldLS/kv-db/server/model"
	"github.com/EmeraldLS/kv-db/server/utils"
)

var databaseCompass = make(map[string]*model.Database, 0)

func CreateDatabase(req []string, conn net.Conn) {
	var db *model.Database
	if len(req) == 1 {
		db = model.NewDatabase()
	} else {
		db = model.NewDatabase()
		db.WithName(req[1])
	}
	mu := sync.Mutex{}
	mu.Lock()
	defer mu.Unlock()

	databaseCompass[db.GetId()] = db

	response := model.DatabaseContentResponse{
		Id:   db.GetId(),
		Name: db.GetName(),
	}

	data, err := json.Marshal(&response)
	if err != nil {
		log.Println("Unable to marshal into json: ", err)
	}

	_, err = conn.Write(data)
	if err != nil {
		log.Println("Unable to write data to the connection: ", err)
	}
}

func InserOneDatabase(reqSlice []string, conn net.Conn) {
	if len(reqSlice) >= 4 {
		db_id := reqSlice[1]
		var mu sync.Mutex

		mu.Lock()
		defer mu.Unlock()

		db := databaseCompass[db_id]

		if db.GetId() == db_id {
			key := reqSlice[2]
			val := strings.Join(reqSlice[3:], " ")
			db.Insert(key, val)

			resp := &model.DefaultResponseFormat{
				Message: "document inserted successfully",
			}
			utils.SendResponse(resp, conn)
		}
	}
}

func FindOneDatabase(req []string, conn net.Conn) {
	if len(req) >= 2 {
		id := req[1]
		var mu sync.RWMutex
		mu.RLock()
		defer mu.RUnlock()

		db, ok := databaseCompass[id]
		if ok {

			resp := &model.DatabaseContentResponse{
				Id:      id,
				Name:    db.GetName(),
				Content: db.GetContent(),
			}

			data, err := json.Marshal(resp)
			if err != nil {
				log.Println("Unable to marshal into json:", err)
			}

			_, err = conn.Write(data)
			if err != nil {
				log.Println("unable to send data: ", err)
			}

		} else {

			resp := &model.DefaultResponseFormat{
				Message: "no database with provided id",
			}
			utils.SendResponse(resp, conn)

		}
	} else {
		resp := &model.DefaultResponseFormat{
			Message: "database id not provided",
		}
		utils.SendResponse(resp, conn)
	}
}

func Defaulthandler(conn net.Conn) {

	resp := &model.DefaultResponseFormat{
		Message: "invalid operation",
	}
	utils.SendResponse(resp, conn)
}

func SaveDatabase(req []string, conn net.Conn) {
	if len(req) >= 2 {
		id := req[1]
		mu := sync.Mutex{}

		mu.Lock()
		defer mu.Unlock()

		db := databaseCompass[id]
		if db != nil {

			resp := &model.DatabaseContentResponse{
				Id:      id,
				Name:    db.GetName(),
				Content: db.GetContent(),
			}

			respByte, err := json.MarshalIndent(resp, "", " ")
			if err != nil {
				log.Fatal(err)
			}

			err = os.MkdirAll("./databases", os.ModePerm)
			if err != nil {
				log.Fatal(err)
			}

			file, err := os.OpenFile(fmt.Sprintf("./databases/db-%s.json", id), os.O_CREATE|os.O_RDWR, os.ModePerm)
			if err != nil {
				log.Fatal(err)
			}

			_, err = file.Write(respByte)
			if err != nil {
				log.Fatal(err)
			}

			respMsg := model.NewDefaultResponseFormat()
			respMsg.WithMessage("database saved successfully")
			utils.SendResponse(respMsg, conn)

		}
	} else {
		respMsg := model.NewDefaultResponseFormat()
		utils.SendResponse(respMsg, conn)
	}
}

func DeleteDatabase(req []string, conn net.Conn) {
	if len(req) == 2 {
		id := req[1]
		mu := sync.Mutex{}

		mu.Lock()
		defer mu.Unlock()

		_, ok := databaseCompass[id]
		if ok {
			delete(databaseCompass, id)
			respMsg := model.NewDefaultResponseFormat()
			respMsg.WithMessage("database delete successfully")

			utils.SendResponse(respMsg, conn)

		} else {
			respMsg := model.NewDefaultResponseFormat()
			respMsg.WithMessage("database not found")

			utils.SendResponse(respMsg, conn)
		}
	} else {
		respMsg := model.NewDefaultResponseFormat()
		utils.SendResponse(respMsg, conn)
	}
}
