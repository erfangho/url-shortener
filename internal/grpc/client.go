package grpc

import (
	"context"

	pb "github.com/erfangho/url-shortener/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type AnalyticsClient struct {
	client pb.AnalyticsServiceClient
	conn   *grpc.ClientConn
}

func NewAnalyticsClient(addr string) (*AnalyticsClient, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		return nil, err
	}

	return &AnalyticsClient{
		client: pb.NewAnalyticsServiceClient(conn),
		conn:   conn,
	}, nil
}

func (c *AnalyticsClient) RecordClick(ctx context.Context, event *pb.ClickEvent) error {
	_, err := c.client.RecordClick(ctx, event)

	return err
}

func (c *AnalyticsClient) Close() {
	c.conn.Close()
}
