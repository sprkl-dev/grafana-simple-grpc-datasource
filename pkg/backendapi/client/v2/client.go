package v2

import (
	v2 "bitbucket.org/innius/grafana-simple-grpc-datasource/pkg/proto/v2"
	v3 "bitbucket.org/innius/grafana-simple-grpc-datasource/pkg/proto/v3"
	"google.golang.org/grpc"
)

type ClientOptions struct {
	V2Client v2.GrafanaQueryAPIClient
}

func NewClient(conn *grpc.ClientConn, opts ClientOptions) (v3.GrafanaQueryAPIClient, error) {
	if opts.V2Client == nil {
		opts.V2Client = v2.NewGrafanaQueryAPIClient(conn)
	}

	return &adapter{v2Client: opts.V2Client}, nil
}
