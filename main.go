package main

import (
	"context"
	"log"
)




func main()  {
	client,err:=NewMinioClient()
	if err!=nil{
		log.Fatal("Error occured while connecting to minio : ",err)
	}
	err=client.uploadFile(context.Background(),"cloud-fs","go-info","./data/go_information.txt","text/plain")
	if err!=nil{
		log.Fatal("Error occured while uploading file")
	}
}