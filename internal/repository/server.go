package repository

import (
	"io"
	"log"
	"net"
	"strings"

	"github.com/EmeraldLS/kv-db/internal/handler"
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
			handler.CreateDB(reqSlice, conn)

		case "INSERT":
			handler.InserOneInDatabase(reqSlice, conn)

		case "FIND":
			conn.Write([]byte("you suck"))

		case "FINDONE":
			handler.FindOne(reqSlice, conn)

		case "SAVE":
			handler.Save(reqSlice, conn)

		default:
			handler.Defaulthandler(conn)
		}

	}
}
