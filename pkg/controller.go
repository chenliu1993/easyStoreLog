package pkg

import (
	"compress/gzip"
	"context"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var bucket_name = "worker1-flowlog"
var numCalcsCreated int32
var bufferSize = 1024 * 1024

var bytesPool = &sync.Pool{
	// Not thread safe, therefore atomic to make sure of it
	New: func() interface{} {
		atomic.AddInt32(&numCalcsCreated, 1)
		buf := make([]byte, bufferSize)
		return &buf
	},
}

func (ctrl *Controller) GenAWSS3Client() {
	// Load the Shared AWS Configuration (~/.aws/config)
	cfg, err := config.LoadDefaultConfig(ctrl.ctx)
	if err != nil {
		log.Printf("[ERROR]: s3 client config error: %v", err)
		return
	}
	log.Printf("Get config from env")
	// Create an Amazon S3 service client
	ctrl.s3client = s3.NewFromConfig(cfg)

}

type Controller struct {
	ctx    context.Context
	cancel context.CancelFunc

	logs Record

	s3client *s3.Client
}

func NewController() (*Controller, error) {
	log.Println("New a controller")
	ctx, cancel := context.WithCancel(context.Background())

	ctrl := &Controller{
		ctx:    ctx,
		cancel: cancel,
		logs:   NewRecord([]string{}),
	}
	// Load the Shared AWS Configuration (~/.aws/config)
	cfg, err := config.LoadDefaultConfig(ctrl.ctx)
	if err != nil {
		log.Printf("[ERROR]: s3 client config error: %v", err)
		return nil, err
	}
	ctrl.s3client = s3.NewFromConfig(cfg)
	return ctrl, nil
}

func (ctrl *Controller) GetContext() context.Context {
	return ctrl.ctx
}

func (ctrl *Controller) GetContextCancel() context.CancelFunc {
	return ctrl.cancel
}

func (ctrl *Controller) Start(storedPath string) {
	log.Println("Start Controller")
	go ctrl.collectLogsFromS3(ctrl.ctx, storedPath)
	time.Sleep(5 * time.Second)
	go ctrl.storeLogsIntoDB(ctrl.ctx)
}

func (ctrl *Controller) Stop(stopCh <-chan os.Signal) {
	s := <-stopCh
	log.Printf("Got termination signal: %v, stopping controller", s)
	ctrl.cancel()
	log.Printf("Pool has been generated %d times, Now Stopped", numCalcsCreated)
}

func (ctrl *Controller) collectLogsFromS3(ctx context.Context, storedPath string) {
	// Max to 1000 objects
	// TODO: How to read the next 1000 objects?
	listoutput, err := ctrl.s3client.ListObjectsV2(ctrl.ctx, &s3.ListObjectsV2Input{
		Bucket: aws.String(bucket_name),
	})
	if err != nil {
		log.Printf("[ERROR]: List s3 objects failed: %v", err)
		return
	}

	for _, content := range listoutput.Contents {
		// key is not overlapped
		key := aws.ToString(content.Key)
		log.Printf("Deal with %s now", key)
		targetDir := filepath.Join(storedPath, filepath.Dir(key))
		log.Printf("create parent folder: %s", targetDir)
		if err := os.MkdirAll(targetDir, 0660); err != nil {
			log.Printf("[ERROR]: Create the prefix path failed: %v", err)
			return
		}

		paths := strings.Split(key, "/")

		if paths[len(paths)-1] == "" {
			log.Printf("%s is a path skip", key)
			continue
		}

		// Download bucket
		getoutput, err := ctrl.s3client.GetObject(ctrl.ctx, &s3.GetObjectInput{
			Bucket: &bucket_name,
			Key:    &key,
		})
		if err != nil {
			log.Printf("[ERROR]: get s3 objects %s failed: %v", key, err)
			return
		}
		log.Printf("get s3 object succeeded: %s", key)

		// targetDir must have storedPath in advance
		file, err := os.Create(filepath.Join(targetDir, strings.TrimSuffix(paths[len(paths)-1], ".gz")))
		if err != nil {
			log.Printf("[ERROR]: create path %s failed: %v", key, err)
			return
		}

		reader, err := gzip.NewReader(getoutput.Body)
		defer getoutput.Body.Close()
		defer reader.Close()
		if err != nil {
			log.Printf("[ERROR]: uncompress %s failed: %v", key, err)
			return
		}

		buffer, _ := bytesPool.Get().(*[]byte)
		// bytesBuffer := buffer.Bytes()
		cleanBuffer := func(buf *[]byte, readBytesNum int) {
			log.Printf("reset buffer with %d", len(*buf))
			cleanByteSlice(buffer, readBytesNum)
			bytesPool.Put(buf)
		}
		n, err := reader.Read(*buffer)
		if err != nil {
			if err != io.EOF {
				log.Printf("[ERROR]: read %s failed: %v", key, err)
				cleanBuffer(buffer, n)
				return
			}
		}

		_, err = file.Write((*buffer)[:n])
		if err != nil {
			log.Printf("[ERROR]: write %s locally failed: %v", key, err)
			cleanBuffer(buffer, n)
			return
		}

		cleanBuffer(buffer, n)
	}

}

func (ctrl *Controller) storeLogsIntoDB(ctx context.Context) {
}
