package bg_process

import (
	"context"
	"crypto/tls"
	"fmt"
	"log/slog"
	"os"
	"sync"
	"time"
	compute "yacloud_revival/third_party/cloudapi/github.com/yandex-cloud/go-genproto/yandex/cloud/compute/v1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
)

type task struct {
	grpcAddress string
	instanceId  string
	period      time.Duration
	logger      *slog.Logger
}

// func init() {
// 	grpc.EnableTracing = true
// 	grpclog.SetLoggerV2(grpclog.NewLoggerV2(os.Stdout, os.Stdout, os.Stdout))
// }

func (t *task) Run(wg *sync.WaitGroup) {
	md := metadata.New(map[string]string{
		"authorization": fmt.Sprintf("Bearer %s", os.Getenv("IAM_TOKEN")),
	})

	for {
		time.Sleep(t.period)

		logger := t.logger.With(
			slog.String("instance_id", t.instanceId),
		)

		conn, err := grpc.NewClient(t.grpcAddress, grpc.WithDisableRetry(), grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{InsecureSkipVerify: true})))
		if err != nil {
			logger.Error("cannot create grpc client")
			continue
		}

		client := compute.NewInstanceServiceClient(conn)
		t.logger.Info("client created, getting instance")

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
		defer cancel()

		ctx = metadata.NewOutgoingContext(ctx, md)

		instance, err := client.Get(ctx, &compute.GetInstanceRequest{
			InstanceId: t.instanceId,
			View:       compute.InstanceView_BASIC,
		})
		if err != nil {
			logger.Error(
				"cannot get instance",
				slog.String("error", err.Error()))
			continue
		}
		logger.Info(
			"got instance",
			slog.String("status", compute.Instance_Status_name[int32(instance.Status)]),
		)
		if instance.Status == compute.Instance_STOPPED {
			logger.Info("instance was stopped, upping instance")
			_, err := client.Start(ctx, &compute.StartInstanceRequest{InstanceId: t.instanceId})
			if err != nil {
				logger.Error(
					"cannot start instance",
					slog.String("error", err.Error()),
				)
				continue
			}
			logger.Info("instance up is in progress")
		}

		conn.Close()
	}

	wg.Done()
}
