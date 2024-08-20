package main

import (
	"context"
	"errors"
	"fmt"


	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)


type MinioClient struct {
	client *minio.Client
}


func NewMinioClient() (*MinioClient,error){
	endpoint := "play.min.io"
    accessKeyID := "Q3AM3UQ867SPQQA43P2F"
    secretAccessKey := "zuf+tfteSlswRu7BJ86wekitnifILbZam1KYY3TG"
    useSSL := true
	minioClient,err:=minio.New(endpoint,&minio.Options{
		Creds: credentials.NewStaticV4(accessKeyID,secretAccessKey,""),
		Secure:useSSL,
	})
	if err!=nil{
		return nil,err
	}
	return &MinioClient{
		client: minioClient,
	},nil
}


func (mc *MinioClient)createBucket(nameOfBucket string,ctx context.Context)error{

	result,err:=mc.client.BucketExists(ctx,nameOfBucket)
	if err!=nil{
		return err
	}
	if !result{
		return errors.New("bucket with this name already exists")
	}else{ 
	err:=mc.client.MakeBucket(ctx,nameOfBucket,minio.MakeBucketOptions{})
	if err!=nil{
		return err
	}
	}
	fmt.Printf("Bucket created with name : %s",nameOfBucket)
	return nil
}

func (mc *MinioClient)uploadFile(ctx context.Context,bucketName,fileName,fileLocation,contentType string)error{
	result,err:=mc.client.FPutObject(ctx,bucketName,fileName,fileLocation,
	minio.PutObjectOptions{ContentType: contentType})
	if err!=nil{
		return err
	}
	fmt.Println("file uploaded : ",result.Key)
	return nil
}


func (mc *MinioClient)downloadFile(ctx context.Context,bucketName,fileName,fileLocation string)error{
    err:=mc.client.FGetObject(ctx,bucketName,fileName,fileLocation,minio.GetObjectOptions{})
    if err!=nil{
        return err
    }
    fmt.Printf("%s downloaded in to %s",fileName,fileLocation)
    return nil
}
