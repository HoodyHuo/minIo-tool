package main

import (
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/schollz/progressbar/v3"
	"io/fs"
	"log"
	"os"
	"path"
)

func restore(client *minio.Client, dir string, bucketName string) {

	fays := os.DirFS(dir)
	dirEntries, err := fs.ReadDir(fays, ".")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for _, entry := range dirEntries {
		if bucketName != "" && entry.Name() != bucketName {
			continue
		}
		fmt.Printf("isDir:%t name:%s \n", entry.IsDir(), entry.Name())
		if entry.IsDir() == false {
			continue
		}
		restoreBucket(client, entry, path.Join(dir))
	}
}

func restoreBucket(client *minio.Client, bucketDir os.DirEntry, dirPath string) {
	exists, err := client.BucketExists(context.Background(), bucketDir.Name())
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	if exists == false {
		err := client.MakeBucket(context.Background(), bucketDir.Name(), minio.MakeBucketOptions{})
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}
	}
	log.Printf("\rstart restore bucket %s", bucketDir.Name())

	fss := os.DirFS(path.Join(dirPath, bucketDir.Name()))
	var files = []string{}
	var totalSize int64 = 0
	fs.WalkDir(fss, ".", func(subPath string, entry fs.DirEntry, err error) error {
		if entry == nil {
			return nil
		}
		if entry.IsDir() {
			return nil
		}
		info, err := entry.Info()
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}
		totalSize += info.Size()
		files = append(files, subPath)
		return nil
	})

	bar := progressbar.NewOptions64(
		totalSize, //总量
		progressbar.OptionSetDescription("restore "+bucketDir.Name()), //初始描述
		progressbar.OptionSetRenderBlankState(true),
		progressbar.OptionShowBytes(true), //显示字节数据
		progressbar.OptionClearOnFinish(), //进度满了就自动清除
		progressbar.OptionShowCount())     //进度条开始

	for _, filePath := range files {
		uploadInfo, err := client.FPutObject(context.Background(),
			bucketDir.Name(),
			filePath,
			path.Join(dirPath, bucketDir.Name(), filePath),
			minio.PutObjectOptions{})
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}
		bar.Add64(uploadInfo.Size)
	}

	log.Printf("\rbucket %s restore finished, total %d objects (%.2fMb)",
		bucketDir.Name(), len(files), float64(totalSize)/1024/1024)
}
