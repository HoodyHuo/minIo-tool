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
			err := backupBucket(minioClient, dir, bucket.Name)
			if err != nil {
				return err
			}
		}
	} else {
		err := backupBucket(minioClient, dir, bucketName)
		if err != nil {
			return err
		}
	}
	return nil
}

func backupBucket(minioClient *minio.Client, dir string, bucketName string) error {
	//桶路径
	log.Printf("\rstart backup bucket %s", bucketName)
	bucketDir := path.Join(dir, bucketName)
	err := os.MkdirAll(bucketDir, 0777)
	if err != nil {
		return err
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

	for _, obj := range objectNames {
		err4, done := funcName(minioClient, bucketName, obj, bucketDir, bar)
		if done {
			return err4
		}
	}

	log.Printf("\rbucket %s backup finished, total %d objects (%.2fMb)", bucketName, len(objectNames), float64(totalSize)/1024/1024)
	return nil
}

func funcName(minioClient *minio.Client, bucketName string, obj minio.ObjectInfo, bucketDir string, bar *progressbar.ProgressBar) (error, bool) {
	objectInfo, err := minioClient.GetObject(context.Background(), bucketName, obj.Key, minio.GetObjectOptions{})
	if err != nil {
		return err, true
	}
	fullPath := path.Join(bucketDir, obj.Key)
	p := fullPath[0:strings.LastIndex(fullPath, "/")]
	err2 := os.MkdirAll(p, 0777)
	if err2 != nil {
		return err2, true
	}

	localFile, err3 := os.Create(fullPath)
	if err3 != nil {
		return err3, true
	}
	if _, err4 := io.Copy(localFile, objectInfo); err != nil {
		return err4, true
	}
	err5 := bar.Add64(obj.Size)
	if err5 != nil {
		return err5, true
	}
	return nil, false
}
