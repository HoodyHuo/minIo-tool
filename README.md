# MinIO-tool
> base on golang
## TODO
- 添加并发处理
- 优化异常处理机制

## 功能
- 备份  
- 还原
- 删除

## 备份桶
````sh
minIo-tool backup -e endpoint -i accessKeyID -p secretAccessKey -d pathToSave [-b BucketName]
````
1.桶会按文件目录放置于 `pathToSave`目录下  
2.如果指定桶名称则只会备份指定桶
  
eg.
````shell
minIo-tool backup -e 200.200.201.10:19000 -i minioadmin -p minioadmin -d ./backup
````
## 还原桶
````sh
minIo-tool restore -e endpoint -i accessKeyID -p secretAccessKey -d pathToRestore [-b BucketName]
````
1.以`pathToRestore`目录下的第一层文件夹作为桶名称，创建桶并上传目录下文件  
2.如果指定桶名称则只会备份`pathToRestore`目录下`BucketName`名称的文件夹


eg.
````shell
minIo-tool restore -e 200.200.201.10:19000 -i minioadmin -p minioadmin -d ./backup
````
## 删除桶
````sh
minIo-tool delete -e endpoint -i accessKeyID -p secretAccessKey -b bucketName
````
eg.
````shell
minIo-tool delete -e 200.200.201.10:19000 -i minioadmin -p minioadmin -b th-bucket
````