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
	err=client.downloadFile(context.Background(),"cloud-fs","go-info","./downloads/downloaded_go_information.txt",)
	if err!=nil{
		log.Fatal("Error occured while uploading file")
	}
}