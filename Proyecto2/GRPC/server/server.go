package main

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	pb "server/proto"

	_ "github.com/go-sql-driver/mysql"
	"google.golang.org/grpc"
)

var ctx = context.Background()
var db *sql.DB

type server struct {
	pb.UnimplementedGetInfoServer
}

const (
	port = ":3001"
)

type Data struct {
	Name  string
	Album string
	Year  string
	Rank  string
}

func (s *server) ReturnInfo(ctx context.Context, in *pb.RequestId) (*pb.ReplyInfo, error) {

	fmt.Println("Recibí de cliente: ", in.GetName())
	data := Data{
		Name:  in.GetName(),
		Album: in.GetAlbum(),
		Year:  in.GetYear(),
		Rank:  in.GetRank(),
	}
	fmt.Println(data)

	return &pb.ReplyInfo{Info: "Se recibió la información"}, nil
}

func sendDataProducer(json_data Data) {
	data, err := json.Marshal(json_data)
	if err != nil {
		log.Fatalln(err)
	}

	urlProd := os.Getenv("PRODURL")

	resp, err := http.Post(urlProd, "application/json", bytes.NewBuffer(data))
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()
}

func main() {
	listen, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalln(err)
	}
	s := grpc.NewServer()
	pb.RegisterGetInfoServer(s, &server{})

	if err := s.Serve(listen); err != nil {
		log.Fatalln(err)
	}
}
