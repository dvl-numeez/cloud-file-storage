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
	err=client.deleteFile(context.Background(),"cloud-fs","go-info")
	if err!=nil{
		log.Fatal("error occured file deleting a file : ",err)
	}
	
}