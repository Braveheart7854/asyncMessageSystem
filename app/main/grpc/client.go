package main

import (
	"context"
	"log"
	"time"

	pb "asyncMessageSystem/app/controller/producer/grpc"
	"google.golang.org/grpc"
)

const (
	address     = "localhost:3334"
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewProducerClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	log.Println(time.Now())
	r, err := c.Notify(ctx, &pb.NoticeRequest{Uid: 12163,Type:1,Data:"test"})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("code: %d", r.GetCode())
	log.Printf("msg: %s", r.GetMsg())
	log.Printf("data: %s", r.GetData())
	//log.Printf("type: %s", reflect.TypeOf(r.GetData()))
	//
	//type data struct {
	//	Uid int64
	//	Type int64
	//	Data string
	//}
	//a := &data{}
	//errs := json.Unmarshal([]byte(r.GetData()),a)
	//
	//fmt.Println(errs)
	//fmt.Println(a.Uid)
	log.Println(time.Now())


}