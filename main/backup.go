package main

import (
	"context"
	"github.com/minio/minio-go/v7"
	"github.com/schollz/progressbar/v3"
	"io"
	"log"
	"os"
	"path"
	"strings"
	"sync"
)

/**
备份
*/
func backup(minioClient *minio.Client, dir string, bucketName string) error {

	if bucketName == "" {
		buckets, err := minioClient.ListBuckets(context.Background())
		if err != nil {
			log.Println(err)
			return err
		}
		for _, bucket := range buckets {
			backupBucket(minioClient, dir, bucket.Name)
		}
	} else {
		backupBucket(minioClient, dir, bucketName)
	}
	return nil
}

func backupBucket(minioClient *minio.Client, dir string, bucketName string) {
	log.Printf("\rstart backup bucket %s", bucketName)
	//桶路径
	bucketDir := path.Join(dir, bucketName)
	err := os.MkdirAll(bucketDir, 0777)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	//桶内文件
	objectsInfo := minioClient.ListObjects(context.Background(), bucketName, minio.ListObjectsOptions{WithMetadata: true, Recursive: true})
	objectNames := []minio.ObjectInfo{}
	var totalSize int64 = 0
	for object := range objectsInfo {
		objectNames = append(objectNames, object)
		totalSize += object.Size
	}
	bar := progressbar.NewOptions64(
		totalSize, //总量
		progressbar.OptionSetDescription("backup "+bucketName), //初始描述
		progressbar.OptionSetRenderBlankState(true),
		progressbar.OptionShowBytes(true), //显示字节数据
		progressbar.OptionClearOnFinish(), //进度满了就自动清除
		progressbar.OptionShowCount())     //进度条开始

	//创建chan
	objectsChan := make(chan minio.ObjectInfo, 4)
	wg := sync.WaitGroup{}
	wg.Add(len(objectNames))
	go dispatchBackupObjects(objectsChan, bucketName, bucketDir, bar, minioClient, &wg)

	for _, obj := range objectNames {
		objectsChan <- obj
	}
	wg.Wait()
	log.Printf("\rbucket %s backup finished, total %d objects (%.2fMb)", bucketName, len(objectNames), float64(totalSize)/1024/1024)
}

func dispatchBackupObjects(objectsChan chan minio.ObjectInfo, bucketName string, bucketDir string, bar *progressbar.ProgressBar, minioClient *minio.Client, wg *sync.WaitGroup) {
	for info := range objectsChan {
		backupObject(info, bucketName, bucketDir, bar, minioClient, wg)
	}
}

func backupObject(object minio.ObjectInfo, bucketName string, bucketDir string, bar *progressbar.ProgressBar, minioClient *minio.Client, wg *sync.WaitGroup) {
	getObject, err := minioClient.GetObject(context.Background(), bucketName, object.Key, minio.GetObjectOptions{})
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	//构建路径
	fullPath := path.Join(bucketDir, object.Key)
	p := fullPath[0:strings.LastIndex(fullPath, "/")]
	err2 := os.MkdirAll(p, 0777)
	if err2 != nil {
		log.Println(err2)
		os.Exit(1)
	}
	//创建文件
	localFile, err3 := os.Create(fullPath)
	if err3 != nil {
		log.Println(err3)
		os.Exit(1)
	}
	//写文件
	if _, err4 := io.Copy(localFile, getObject); err != nil {
		log.Println(err4)
		os.Exit(1)
	}
	//更新进度条、任务组
	err5 := bar.Add64(object.Size)
	if err5 != nil {
		log.Println(err5)
		os.Exit(1)
	}
	wg.Done()
}
