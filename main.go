package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

var (
	endpint        = "oss-cn-beijing.aliyuncs.com"
	ak             = "xx"
	sk             = "xx"
	bucketName     = "clouds-station"
	uploadFilePath = ""
	help           bool
)

func usage() {
	fmt.Fprintf(os.Stderr, `Usage: cloud-station -f <uplaod_file_path>
`)
	flag.PrintDefaults()
	os.Exit(0)
}

func loadParam() error {
	flag.BoolVar(&help, "h", false, "help usage")
	flag.StringVar(&uploadFilePath, "f", "", "upload file path")
	flag.Parse()

	if uploadFilePath == "" {
		usage()
	}

	if help {
		usage()
	}

	return nil
}

func upload(filePath string) error {
	client, err := oss.New(endpint, ak, sk)
	if err != nil {
		return err
	}

	bucket, err := client.Bucket(bucketName)
	if err != nil {
		return err
	}

	// 第一个参数: 上传到oss里面的文件的key(路径)
	// 第二个参数: 需要上传的文件的路径
	err = bucket.PutObjectFromFile(filePath, filePath)
	if err != nil {
		return err
	}

	// 打印下载URL
	// sts, 临时授权token(有效期1天)
	signedURL, err := bucket.SignURL(filePath, oss.HTTPGet, 60*60*24)
	if err != nil {
		return fmt.Errorf("sign file download url error, %s", err)
	}
	fmt.Printf("upload file: %s success\n", uploadFilePath)
	fmt.Printf("下载链接: %s\n", signedURL)
	fmt.Println("\n注意: 文件下载有效期为1天, 中转站保存时间为3天, 请及时下载")
	return nil
}

func main() {
	// 读取用户输入的参数
	loadParam()

	// 执行文件上传
	if err := upload(uploadFilePath); err != nil {
		fmt.Printf("upload file error, %s\n", err)
		os.Exit(1)
	}

}
