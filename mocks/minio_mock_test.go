package mocks

import (
	"context"
	"errors"
	"reflect"

	"testing"
)




func TestCreateBucket(t *testing.T) {
	client:=makeNewMockClient()
	t.Run("creating a bucket",func(t *testing.T) {
		nameOfBuckets:=[]string{"cloud-fs","filesystem","fs","cf","cloud-bucket"}
		for _,bucket:=range nameOfBuckets{
			err:=client.createBucket(bucket,context.TODO())
		if err!=nil{
			t.Error(err)
		}
		_,ok:=client.db[bucket]
		if !ok{
			t.Errorf("Expected value should be %v actual value was %v",true,ok)
		}
		}
		
	})
	t.Run("creating a bucket with the name that exists",func(t *testing.T) {
		nameOfBucket:="cloud-fs"
		err:=client.createBucket(nameOfBucket,context.TODO())
		expected:=errors.New("bucket with this name exists")
		if err.Error()!=expected.Error(){
			t.Errorf("Expected error : %s Actual error : %s",expected.Error(),err.Error())
		}
	})
	t.Run("creating a bucket with no name",func(t *testing.T){
		expected:=errors.New("bucket name should not be empty")
		nameOfBucket:=""
		err:=client.createBucket(nameOfBucket,context.TODO())
		if err.Error()!=expected.Error(){
			t.Errorf("Expected error : %s Actual error : %s",expected.Error(),err.Error())
		}

	})

	
}

func TestUploadFile(t *testing.T) {
	client:=makeNewMockClient()
	t.Run("uploading a file",func(t *testing.T){
		nameOfBucket:="cloud-fs"
		expected:=File{
			FileName: "go_info",
			FileLocation:"./data/go.txt" ,
			ContentType: "text/plain",
		}
		err:=client.uploadFile(context.TODO(),nameOfBucket,expected.FileName,expected.FileLocation,expected.ContentType)
		if err!=nil{
			t.Error(err)
		}
		files:=client.db[nameOfBucket]
		for _,file:= range files{
			if file.FileName!=expected.FileName{
				t.Errorf("Expected : %s Got :%s ",expected.FileName,file.FileName)
			}
		}
	})
	t.Run("uploading a file with not enough information about it",func(t *testing.T){
		cases:=[]File{{"","",""},{"go_info","",""},{"","/data/text",""},{"","","text/plain"}}
		expected:=errors.New("no information about file or bucket should be empty")
		for _,c:=range cases{
			err:=client.uploadFile(context.TODO(),"cloud-fs",c.FileName,c.FileLocation,c.ContentType)
			if err.Error()!=expected.Error(){
				t.Errorf("Expected error : %s Actual error : %s",expected.Error(),err.Error())
			}
		}
	})
	t.Run("uploading  a file to bucket with no name",func(t *testing.T) {
		err:=client.uploadFile(context.Background(),"","go_info","./data/go.txt","text/plain")
		if err==nil{
			t.Error("Expecting an error but did not get")
		}
	})
}

func TestDownloadFile(t *testing.T) {
	type File struct{
		FileName string
		FileLocation string
		OutputDir string
	}
	client:=makeNewMockClient()
	err:=client.uploadFile(context.TODO(),"cloud-fs","go_info","data/go.txt","text/plain")
		if err!=nil{
			t.Error(err)
		}
	t.Run("downloading from the bucket which does not exists",func(t *testing.T) {
		err:=client.downloadFile(context.TODO(),"","go_info","data/doc.txt")
		if err==nil{
			t.Error("Expected an error but did not get it")
		}
	})
	t.Run("not enough file information",func(t *testing.T) {
		files:=[]File{{"","",""},{"go_info","","data/go.txt"},{"go","fs",""},}
		for _,f:=range files{
			err:=client.downloadFile(context.TODO(),f.FileName,f.FileLocation,f.OutputDir)
		if err==nil{
			t.Error("Expected an error but did not get it")
		}	
		}
	})
	t.Run("Downloading a file",func(t *testing.T){
		fileLocation:="data/go.txt"
		err:=client.downloadFile(context.TODO(),"cloud-fs","go_info",fileLocation)
		if err!=nil{
			t.Error(err)
		}
		expected:=db["cloud_fs"]
		got:=dataDirectory[fileLocation]
		if expected.FileName!=got.FileName{
			t.Errorf("Expected : %s Got : %s",expected.FileName,got.FileName)
		}


	})
	t.Run("Downloading a file which does not exists in the bucket",func(t *testing.T){
		fileLocation:="data/go.txt"
		err:=client.downloadFile(context.TODO(),"cloud-fs","info",fileLocation)	
		expected:=errors.New("file you are trying to download does not exists")
		if err.Error()!=expected.Error(){
			t.Errorf("Expected  error %s Got : %s",expected.Error(),err.Error())
		}
	})
}


func TestDeleteFile(t *testing.T) {
	client:=makeNewMockClient()
	data:=[]struct{
		bucketName string
		file []File
	}{
		{"cloud-fs",[]File{{"go_info","data/go.txt","text/plain"}}},
		{"amazon-fs",[]File{{"java_info","data/java.txt","text/plain"}}},
		{"apple-fs",[]File{{"swift_info","data/swift.txt","text/plain"}}},
		{"meta-fs",[]File{{"typescript_info","data/typescript.txt","text/plain"}}},
	}
	for _,d:= range data{
		client.db[d.bucketName] = d.file
	}

	t.Run("Deleting file with invalid or wrong bucket name",func(t *testing.T){
		err:=client.deleteFile(context.TODO(),"abc","go.txt")
		expected:=errors.New("bucket does not exists")
		if err.Error()!=expected.Error(){
			t.Errorf("Expected error : %s Actual error : %s",expected.Error(),err.Error())
		}
	})
	t.Run("Deleting file that does not exists in the bucket",func(t *testing.T){
		err:=client.deleteFile(context.TODO(),"cloud-fs","java_info")
		expected:=errors.New("File you are trying to delete does not exsits")
		if err.Error()!=expected.Error(){
			t.Errorf("Expected error : %s Actual error : %s",expected.Error(),err.Error())
		}
	})
	t.Run("Deleting a file",func(t *testing.T){
		err:=client.deleteFile(context.TODO(),"cloud-fs","go_info")
		if err!=nil{
			t.Error(err)
		}
		got:=client.db["cloud-fs"]
		expected:=[]File{}
		if !reflect.DeepEqual(got,expected){
			t.Errorf("Expected %v Actual %v ",expected,got)
		}
		
	})
	t.Run("Not providing enough information",func(t *testing.T){
		cases:=[]struct{
			bucketName string
			fileName string
		}{
			{"",""},
			{"cloud-fs",""},
			{"","swift_info"},
		}
		for _,c:=range cases{
			err:=client.deleteFile(context.TODO(),c.bucketName,c.fileName)
			expected:=errors.New("bucketname or file name should not be empty")
			if err.Error()!=expected.Error(){
				t.Errorf("Expected error : %s Actual error : %s",expected.Error(),err.Error())
			}
		}
	})
	
}