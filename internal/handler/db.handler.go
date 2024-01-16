package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"sync"

	"github.com/EmeraldLS/kv-db/internal/model"
)

var databaseCompass = make(map[string]*model.Database, 0)

func CreateDB(req []string, conn net.Conn) {
	var db *model.Database
	if len(req) == 1 {
		db = model.NewDatabase()
	} else {
		db = model.NewDatabase()
		db.WithName(req[1])
	}

	databaseCompass[db.GetId()] = db

	type db_detail struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}

	response := db_detail{
		ID:   db.GetId(),
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

func InserOneInDatabase(reqSlice []string, conn net.Conn) {
	if len(reqSlice) >= 4 {
		db_id := reqSlice[1]
		var mu sync.Mutex

		mu.Lock()
		defer mu.Unlock()

		db := databaseCompass[db_id]

		if db.GetId() == db_id {
			key := reqSlice[2]
			val := strings.Join(reqSlice[3:], " ")

			fmt.Println("Value: ", val)

			var buf bytes.Buffer
			_, err := buf.Write([]byte("xx"))
			if err != nil {
				log.Fatal(err)
			}

			json.NewEncoder(&buf)
			db.Insert(key, val)

			resp := &model.DefaultResponseFormat{
				Message: "document inserted successfully",
			}
			data, err := json.Marshal(resp)
			if err != nil {
				log.Println("unable to marshal into json: ", err)
			}

			_, err = conn.Write(data)
			if err != nil {
				log.Println("unable to write data: ", err)
			}
		}
	}
}

func FindOne(req []string, conn net.Conn) {
	if len(req) >= 2 {
		id := req[1]
		var mu sync.RWMutex
		mu.RLock()
		defer mu.RUnlock()

		db := databaseCompass[id]
		fmt.Println(db)
		if db != nil {

			dbContent := db.GetContent()

			var buf bytes.Buffer

			for _, v := range dbContent {
				n, err := buf.Read([]byte(v.(string)))
				if err == io.EOF {
					break
				}
				if err != nil {
					log.Fatal(err)
				}

				log.Printf("Read %d into the buffer", n)
			}

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
			data, err := json.Marshal(resp)
			if err != nil {
				log.Println("Unable to marshal into json:", err)
			}

			_, err = conn.Write(data)
			if err != nil {
				log.Println("unable to send data: ", err)
			}

		}
	}
}

func Defaulthandler(conn net.Conn) {

	resp := &model.DefaultResponseFormat{
		Message: "invalid operation",
	}

	data, err := json.Marshal(resp)
	if err != nil {
		log.Println("Unable to marshal into json:", err)
	}

	_, err = conn.Write(data)
	if err != nil {
		log.Println("unable to send data: ", err)
	}
}

func FailAndSendErrResponseToConn(conn net.Conn, msg string, err error) {
	if err != nil {
		var resp = model.NewDefaultResponseFormat(msg)
		respByte, _ := json.Marshal(resp)

		_, err = conn.Write(respByte)
		if err != nil {
			log.Fatal(err)
		}
		return
	}
}

func Save(req []string, conn net.Conn) {
	if len(req) == 2 {
		id := req[1]
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

			file, err := os.OpenFile(fmt.Sprintf("db-%s.json", id), os.O_APPEND|os.O_CREATE|os.O_RDWR, os.ModePerm)
			if err != nil {
				log.Fatal(err)
			}

			_, err = file.Write(respByte)
			if err != nil {
				log.Fatal(err)
			}

			respMsg := model.NewDefaultResponseFormat("db content saved successfully")
			respMsgByte, err := json.Marshal(respMsg)
			if err != nil {
				log.Fatal(err)
			}
			_, err = conn.Write(respMsgByte)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}
