package mocks

import (
	"context"
	"errors"
	"log"
)


var bucketMap =  map[string]bool{}
var db = map[string]File{}
var dataDirectory = map[string]File{}

type File struct {
	FileName string
	FileLocation string
	ContentType string
}
type MockMinioClient struct{
	db map[string][]File
	downloadDirectory map[string]string
}

func makeNewMockClient()*MockMinioClient{
	return &MockMinioClient{
		db: make(map[string][]File),
		downloadDirectory: make(map[string]string),
	}
}

func(mmc *MockMinioClient)createBucket(nameOfBucket string,ctx context.Context)error{
	
	if nameOfBucket==""{
		return errors.New("bucket name should not be empty")
	}
	_,ok:=mmc.db[nameOfBucket]
	if ok{
		return  errors.New("bucket with this name exists")
	}
	mmc.db[nameOfBucket] = []File{}
	return nil
}

func(mmc *MockMinioClient)uploadFile(ctx context.Context,bucketName,fileName,fileLocation,contentType string)error{

	if bucketName==""||fileName==""||fileLocation==""||contentType==""{
		return errors.New("no information about file or bucket should be empty")
	}
	file:=File{
		FileName:fileName ,
		FileLocation: fileLocation,
		ContentType: contentType,
	}
	mmc.db[bucketName] = append(mmc.db[bucketName], file)
	return nil
}
func(mmc *MockMinioClient)downloadFile(ctx context.Context,bucketName,fileName,fileLocation string)error{
	 if bucketName==""||fileName==""||fileLocation==""{
		return errors.New("no information about file or bucket should be empty")
	}
	
	files,ok:=mmc.db[bucketName]	
	if !ok{
		return errors.New("bucket does not exists")
	}
	for _,file:=range files{
		if file.FileName==fileName{
			mmc.downloadDirectory[fileLocation]= file.FileName
			return nil
		}
	}

	
	return errors.New("file you are trying to download does not exists")
}
func(mmc *MockMinioClient)deleteFile(ctx context.Context,bucketName,fileName string)error{
	if bucketName=="" || fileName==""{
		return errors.New("bucketname or file name should not be empty")
	}
	
	for k:= range mmc.db{
		if k==bucketName{
			files:=mmc.db[bucketName]
			for _,f:=range files{
			if f.FileName==fileName{
				index,err:=getIndex(files,fileName)
				if err!=nil{
					log.Fatal(err)
				}
				mmc.db[bucketName] = deleteElement(files,index)
				return nil
			}
		}
		return errors.New("File you are trying to delete does not exsits")
		}
	}
	return errors.New("bucket does not exists")
}

 func getIndex(files[]File,fileName string)(int,error){
	 for i,f:=range files{
		 if f.FileName==fileName{
			 return i,nil
		 }
	 } 
	 return 0,errors.New("element not found")
 }

 func deleteElement(slice []File, index int) []File {
	return append(slice[:index], slice[index+1:]...)
 }
 