package v1

import (
	v2client "bitbucket.org/innius/grafana-simple-grpc-datasource/pkg/backendapi/client/v2"
	v1 "bitbucket.org/innius/grafana-simple-grpc-datasource/pkg/proto/v1"
	v3 "bitbucket.org/innius/grafana-simple-grpc-datasource/pkg/proto/v3"
	"google.golang.org/grpc"
)

func NewClient(conn *grpc.ClientConn) (v3.GrafanaQueryAPIClient, error) {
	// INFO: using adapter from v1 to v2 and adapter from v2 to v3 (v1 -> v2 -> v3)
	v2c := &adapter{v1Client: v1.NewGrafanaQueryAPIClient(conn)}
	return v2client.NewClient(conn, v2client.ClientOptions{V2Client: v2c})
}
