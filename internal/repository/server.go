package repository

import (
	"encoding/json"
	"io"
	"log"
	"net"
	"strings"
	"sync"

	"github.com/EmeraldLS/kv-db/internal/model"
)

var all_databases = make(map[string]*model.Database, 0)

func Listen(port string) {
	list, err := net.Listen("tcp", "0.0.0.0:"+port)
	if err != nil {
		log.Fatal(err)
	}
	defer list.Close()

	log.Println("Listening on ", list.Addr().String())

	for {
		conn, err := list.Accept()
		if err != nil {
			log.Println(err)
			continue
		}

		go handleConn(conn)

	}

}

func handleConn(conn net.Conn) {

	buf := make([]byte, 1024)
	for {
		n, err := conn.Read(buf)
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Println(err)
		}

		if len(buf) == 0 {
			continue
		}

		req := string(buf[:n])
		reqSlice := strings.Split(req, " ")
		method := strings.ToUpper(reqSlice[0])

		switch method {
		case "CREATE":
			createDB(reqSlice, conn)

		case "INSERT":
			insertInDb(reqSlice, conn)

		case "FIND":
			conn.Write([]byte("you suck"))

		case "CONTENT":
			DbContent(reqSlice, conn)

		default:
			defaultHandler(conn)
		}

	}
}

func createDB(req []string, conn net.Conn) {
	var db *model.Database
	if len(req) == 1 {
		db = model.NewDatabase()
	} else {
		db = model.NewDatabase()
		db.WithName(req[1])
	}

	all_databases[db.GetId()] = db

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

func insertInDb(reqSlice []string, conn net.Conn) {

	type resp_format struct {
		Message string `json:"message"`
	}

	if len(reqSlice) >= 4 {
		db_id := reqSlice[1]
		var mu sync.Mutex

		mu.Lock()
		defer mu.Unlock()

		db := all_databases[db_id]

		if db.GetId() == db_id {
			key := reqSlice[2]
			val := reqSlice[3]
			db.Insert(key, val)

			resp := &resp_format{
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

func DbContent(req []string, conn net.Conn) {
	if len(req) >= 2 {
		id := req[1]
		var mu sync.RWMutex
		mu.RLock()
		defer mu.RUnlock()

		db := all_databases[id]
		if db != nil {
			type resp_format struct {
				ID      string                 `json:"id"`
				Name    string                 `json:"name"`
				Content map[string]interface{} `json:"content"`
			}

			resp := &resp_format{
				ID:      id,
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
			type resp_format struct {
				Message string `json:"message"`
			}

			resp := &resp_format{
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

func defaultHandler(conn net.Conn) {
	type resp_format struct {
		Message string `json:"message"`
	}

	resp := &resp_format{
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
