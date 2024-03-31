package grpcserv

import (
	"context"
	"fmt"
	"time"

	grpcretry "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/retry"
	protoc "github.com/marckuusha/protoc/gen/go/hwserv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	api protoc.HWServClient
}

func New(ctx context.Context, addr string, retriesCount int, timeout time.Duration) (*Client, error) {

	const op = "grpc.Client.New"

	retryOpts := []grpcretry.CallOption{
		grpcretry.WithCodes(codes.NotFound, codes.Aborted, codes.DeadlineExceeded),
		grpcretry.WithMax(uint(retriesCount)),
		grpcretry.WithPerRetryTimeout(timeout),
	}

	// logsOpts := []grpclog.Option{
	// 	grpclog.WithLogOnEvents(grpclog.PayloadReceived, grpclog.PayloadSent),
	// }

	cc, err := grpc.DialContext(ctx, addr, grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(
			//	grpclog.UnaryClientInterceptor(InterseptorLogget(log), logsOpts...),
			grpcretry.UnaryClientInterceptor(retryOpts...),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("%s : %s", op, err)
	}

	return &Client{
		api: protoc.NewHWServClient(cc),
	}, nil
}

func (c *Client) SendMsg(ctx context.Context, name string) (string, error) {
	resp, err := c.api.SendMsg(ctx, &protoc.SendMsgReq{
		UserName: name,
	})
	if err != nil {
		return "", fmt.Errorf("cannot send msg grpc: %w", err)
	}

	return resp.Msg, nil
}
