package connector

import (
	"bitbucket.org/innius/grafana-simple-grpc-datasource/pkg/backendapi/client"
	"bitbucket.org/innius/grafana-simple-grpc-datasource/pkg/framer"
	"bitbucket.org/innius/grafana-simple-grpc-datasource/pkg/models"
	pb "bitbucket.org/innius/grafana-simple-grpc-datasource/pkg/proto/v3"
	"context"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func historyQueryToInput(query models.MetricHistoryQuery) *pb.GetMetricHistoryRequest {
	var dimensions []*pb.Dimension
	for _, d := range query.Dimensions {
		dimensions = append(dimensions, &pb.Dimension{
			Key:   d.Key,
			Value: d.Value,
		})
	}
	metrics := make([]string, len(query.Metrics))
	for i := range query.Metrics {
		metrics[i] = query.Metrics[i].MetricId
	}
	return &pb.GetMetricHistoryRequest{
		Dimensions:    dimensions,
		Metrics:       metrics,
		StartDate:     timestamppb.New(query.TimeRange.From),
		EndDate:       timestamppb.New(query.TimeRange.To),
		StartingToken: query.NextToken,
	}
}

func GetMetricHistory(ctx context.Context, client client.BackendAPIClient, query models.MetricHistoryQuery) (*framer.MetricHistory, error) {
	clientReq := historyQueryToInput(query)

	resp, err := client.GetMetricHistory(ctx, clientReq)

	if err != nil {
		return nil, err
	}
	return &framer.MetricHistory{
		GetMetricHistoryResponse: resp,
		Query:                    query,
	}, nil
}
