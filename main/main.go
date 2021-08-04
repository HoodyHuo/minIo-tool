package main

import (
	"flag"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"log"
	"os"
)

func main() {
	//备份子命令
	backupCmd := flag.NewFlagSet("backup", flag.ExitOnError)
	endpoint := backupCmd.String("e", "", "endpoint like：127.0.0.1:9000")
	accessKeyID := backupCmd.String("i", "", "accessKeyID")
	secretAccessKey := backupCmd.String("p", "", "secretAccessKey")
	baseDir := backupCmd.String("d", "./backup", "savePath like: ./backup or /path/to/save")
	bucketName := backupCmd.String("b", "", "bucketName，if set then backup the bucket only")

	//还原子命令
	restoreCmd := flag.NewFlagSet("restore", flag.ExitOnError)
	restoreEndpoint := restoreCmd.String("e", "", "endpoint like：127.0.0.1:9000")
	restoreAccessKeyID := restoreCmd.String("i", "", "accessKeyID")
	restoreSecretAccessKey := restoreCmd.String("p", "", "secretAccessKey")
	restoreBaseDir := restoreCmd.String("d", "./backup", "savePath like: ./backup or /path/to/save")
	restoreBucketName := restoreCmd.String("b", "", "bucketName，if set then restore the bucket only")

	//删除子命令
	deleteCmd := flag.NewFlagSet("restore", flag.ExitOnError)
	deleteEndpoint := deleteCmd.String("e", "", "endpoint like：127.0.0.1:9000")
	deleteAccessKeyID := deleteCmd.String("i", "", "accessKeyID")
	deleteSecretAccessKey := deleteCmd.String("p", "", "secretAccessKey")
	deleteBucketName := deleteCmd.String("b", "", "bucketName to be delete")

	if len(os.Args) < 2 {
		fmt.Println("expected 'backup' 'restore' 'delete' subcommands")
		printUsage()
		os.Exit(1)
	}
	switch os.Args[1] {
	case "backup":
		err := backupCmd.Parse(os.Args[2:len(os.Args)])
		if err != nil {
			backupCmd.Usage()
			os.Exit(1)
		}
		if *endpoint == "" || *accessKeyID == "" || *secretAccessKey == "" {
			backupCmd.Usage()
			os.Exit(1)
		}
		client := makeClient(*endpoint, *accessKeyID, *secretAccessKey)
		err = backup(client, *baseDir, *bucketName)
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}
		break
	case "restore":
		err := restoreCmd.Parse(os.Args[2:len(os.Args)])
		if err != nil {
			restoreCmd.Usage()
			os.Exit(1)
		}
		if *restoreEndpoint == "" || *restoreAccessKeyID == "" || *restoreSecretAccessKey == "" {
			restoreCmd.Usage()
			os.Exit(1)
		}
		client := makeClient(*restoreEndpoint, *restoreAccessKeyID, *restoreSecretAccessKey)
		restore(client, *restoreBaseDir, *restoreBucketName)
		break
	case "delete":
		err := deleteCmd.Parse(os.Args[2:len(os.Args)])
		if err != nil {
			deleteCmd.Usage()
			os.Exit(1)
		}
		if *deleteBucketName == "" || *deleteEndpoint == "" || *deleteAccessKeyID == "" || *deleteSecretAccessKey == "" {
			deleteCmd.Usage()
			os.Exit(1)
		}
		client := makeClient(*deleteEndpoint, *deleteAccessKeyID, *deleteSecretAccessKey)
		deleteBucket(client, *deleteBucketName)
		break
	default:
		printUsage()
	}
}

func printUsage() {
	fmt.Println("usage like:")
	fmt.Println("minIo-tool backup -e endpoint -i accessKeyID -p secretAccessKey -d pathToSave [-b BucketName]")
	fmt.Println("or")
	fmt.Println("minIo-tool restore -e endpoint -i accessKeyID -p secretAccessKey -d pathToRestore [-b BucketName]")
	fmt.Println("or")
	fmt.Println("minIo-tool delete -e endpoint -i accessKeyID -p secretAccessKey -b bucketName")
}

func makeClient(endpoint string, accessKeyID string, secretAccessKey string) *minio.Client {

	//model:= flag.("backup","minIo-tool backup -e endpoint -i accessKeyID -p secretAccessKey",backup)
	// Initialize minio client object.
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds: credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
	})
	if err != nil {
		log.Fatalln(err)
	}
	return minioClient
}
