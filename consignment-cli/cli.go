package main

import (
	"context"
	"encoding/json"
	"errors"
	"google.golang.org/grpc"
	"io/ioutil"
	"log"
	"os"
	pb "shippy/consignment-service/proto/consignment"
)

const (
	ADDR              = "localhost:50051"
	DEFUALT_INFO_FILE = "consignment.json"
)

func main() {
	conn, err := grpc.Dial(ADDR, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("connect error:%v", err)
	}
	defer conn.Close()

	client := pb.NewShippingServiceClient(conn)
	infoFile := DEFUALT_INFO_FILE

	if len(os.Args) >= 1 {
		if os.Args[1] == "get" {
			resp,err := client.GetConsignments(context.Background(),&pb.GetRequest{})
			if err != nil{
				log.Fatalf("failed to list consignments:%v",err)
			}
			for _,v := range resp.GetConsignments() {
				log.Printf("%+v",v)
			}
			return
		}
		infoFile = os.Args[1]
	}

	consignment, err := parseFile(infoFile)
	if err != nil {
		log.Fatalf("parse info file error:%v", err)
	}

	resp, err := client.CreateConsignment(context.Background(), consignment)
	if err != nil {
		log.Fatalf("create consignment error:%v", err)
	}
	log.Printf("created: %t", resp.Created)
}

func parseFile(fileName string) (*pb.Consignment, error) {
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	var consignment pb.Consignment
	err = json.Unmarshal(data, &consignment)
	if err != nil {
		return nil, errors.New("json file unmarshal failed")
	}
	return &consignment, nil
}
