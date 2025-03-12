package main

import (
	"context"
	"fmt"
	uploaderv1 "github.com/k1v4/protos/gen/file_uploader"
	"google.golang.org/grpc"
	"log"
	"os"
	"time"
)

func main() {
	fileData, err := os.ReadFile("motivation.png")
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}

	// Теперь fileData — это []byte, который можно передать в gRPC запрос
	//log.Printf("File data: %v", fileData)

	conn, err := grpc.Dial(":50053", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := uploaderv1.NewFileUploaderClient(conn)

	time.Sleep(20 * time.Second)

	response, err := client.UploadFile(context.Background(), &uploaderv1.ImageUploadRequest{
		ImageData: fileData,
		FileName:  "motivation.png",
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(response)

	//resp, err := client.DeleteFile(context.Background(), &uploaderv1.ImageDeleteRequest{
	//	Url: "https://82a3fa46-643f-4a21-8a10-c2889596892b.selstorage.ru/motivation.png",
	//})
	//if err != nil {
	//	panic(err)
	//}
	//
	//fmt.Println(response)
}
