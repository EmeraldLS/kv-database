package repository

import (
	"io"
	"log"
	"net"
	"strings"

	"github.com/EmeraldLS/kv-db/server/handler"
)

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
	log.Println("New Connection from ", conn.RemoteAddr().String())
	var buf = make([]byte, 1024)
	for {
		n, err := conn.Read(buf)
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Println(err)
		}

		req := string(buf[:n])
		reqSlice := strings.Split(req, " ")
		method := strings.ToUpper(reqSlice[0])

		switch method {
		case "CREATE":
			handler.CreateDatabase(reqSlice, conn)

		case "INSERT":
			handler.InserOneDatabase(reqSlice, conn)

		case "FIND":
			conn.Write([]byte(""))

		case "FINDONE":
			handler.FindOneDatabase(reqSlice, conn)

		case "SAVE":
			handler.SaveDatabase(reqSlice, conn)

		case "DELETE":
			handler.DeleteDatabase(reqSlice, conn)

		default:
			handler.Defaulthandler(conn)
		}

	}
}
