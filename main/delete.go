package main

import (
	"bufio"
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/schollz/progressbar/v3"
	"log"
	"os"
)

func deleteBucket(minioClient *minio.Client, bucketName string) {
	makeSure(bucketName)
	if checkExists(minioClient, bucketName) == false {
		log.Printf("%s not exists in MinIO\n", bucketName)
		os.Exit(1)
	}
	deleteObjects(minioClient, bucketName)
	deleteBucketEmpty(minioClient, bucketName)

}

func makeSure(name string) {
	fmt.Printf("Are you sure delete bucket %s ? [Yes/no]\n", name)
	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		line := input.Text()
		if line == "Yes" {
			return
		} else {
			os.Exit(1)
		}
	}
}

func deleteBucketEmpty(client *minio.Client, name string) {
	err := client.RemoveBucket(context.Background(), name)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	log.Printf("\rbucket %s deleted", name)
}

func deleteObjects(minioClient *minio.Client, bucketName string) {
	objectsInfo := minioClient.ListObjects(context.Background(), bucketName, minio.ListObjectsOptions{WithMetadata: true, Recursive: true})
	var objectNames = []minio.ObjectInfo{}
	var totalSize int64 = 0
	for object := range objectsInfo {
		objectNames = append(objectNames, object)
		totalSize += object.Size
	}

	bar := progressbar.NewOptions(
		len(objectNames), //总量
		progressbar.OptionSetDescription("delete Objects in "+bucketName), //初始描述
		progressbar.OptionSetRenderBlankState(true),
		progressbar.OptionShowBytes(true), //显示字节数据
		progressbar.OptionClearOnFinish(), //进度满了就自动清除
		progressbar.OptionShowCount())     //进度条开始

	for _, objectName := range objectNames {
		err := minioClient.RemoveObject(context.Background(), bucketName, objectName.Key, minio.RemoveObjectOptions{})
		if err != nil {
			log.Printf("Object %s remove fail, %s", objectName.Key, err)
			os.Exit(1)
		}
		bar.Add(1)
	}

	log.Printf("\rbucket %s Objects delete finished, total %d objects (%.2fMb)", bucketName, len(objectNames), float64(totalSize)/1024/1024)
}

func checkExists(minioClient *minio.Client, bucketName string) bool {
	exists, err := minioClient.BucketExists(context.Background(), bucketName)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	return exists
}
