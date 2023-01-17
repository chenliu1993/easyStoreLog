package pkg

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func (ctrl *Controller) GenAWSS3Client() {
	// Load the Shared AWS Configuration (~/.aws/config)
	cfg, err := config.LoadDefaultConfig(ctrl.ctx)
	if err != nil {
		log.Printf("s3 client config error: %v", err)
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

func New() (*Controller, error) {
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
		log.Printf("s3 client config error: %v", err)
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

func (ctrl *Controller) Start() {
	log.Println("Start Controller")
	go ctrl.collectLogsFromS3(ctrl.ctx)
	go ctrl.storeLogsIntoDB(ctrl.ctx)
}

func (ctrl *Controller) Stop(stopCh <-chan os.Signal) {
	s := <-stopCh
	log.Printf("Got termination signal: %v, stopping controller", s)
	ctrl.cancel()
	log.Println("Stopped")
}

func (ctrl *Controller) collectLogsFromS3(ctx context.Context) {
	// Max to 1000 objects
	// TODO: How to read the next 1000 objects?
	output, err := ctrl.s3client.ListObjectsV2(ctrl.ctx, &s3.ListObjectsV2Input{
		Bucket: aws.String("worker1-flowlog"),
	})
	if err != nil {
		log.Printf("List s3 objects failed: %v", err)
		return
	}

}

func (ctrl *Controller) storeLogsIntoDB(ctx context.Context) {

}
