package bg_process

import (
	"context"
	"crypto/tls"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"log/slog"
	"sync"
	"time"
	"yacloud_revival/internal/token"
	compute "yacloud_revival/third_party/cloudapi/github.com/yandex-cloud/go-genproto/yandex/cloud/compute/v1"
)

//func init() {
//	grpc.EnableTracing = true
//	grpclog.SetLoggerV2(grpclog.NewLoggerV2(os.Stdout, os.Stdout, os.Stdout))
//}

type task struct {
	grpcAddress string
	instanceId  string
	period      time.Duration
	logger      *slog.Logger
	tokenGetter *token.Token
}

func (t *task) CheckInstance() error {

	logger := t.logger.With(
		slog.String("instance_id", t.instanceId),
	)

	iAmToken, err := t.tokenGetter.Get()
	if err != nil {
		return err
	}

	md := metadata.New(map[string]string{
		"authorization": fmt.Sprintf("Bearer %s", iAmToken),
	})

	logger.Info("TOKEN", slog.String("token", iAmToken))

	conn, err := grpc.NewClient(t.grpcAddress, grpc.WithDisableRetry(), grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{InsecureSkipVerify: true})))

	if err != nil {
		logger.Error("cannot create grpc client")
		return err
	}

	defer conn.Close()

	client := compute.NewInstanceServiceClient(conn)
	t.logger.Info("client created, getting instance")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
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
		return err
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
			return err
		}
		logger.Info("instance up is in progress")
	}

	return nil
}

func (t *task) Run(ctx context.Context, wg *sync.WaitGroup) {

	ticker := time.NewTicker(t.period)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			err := t.CheckInstance()
			if err != nil {
				t.logger.Info(
					"error checking instance",
					slog.String("error", err.Error()),
				)
			}
		}
	}
}
