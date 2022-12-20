package v2

import (
	"context"

	v2 "bitbucket.org/innius/grafana-simple-grpc-datasource/pkg/proto/v2"
	v3 "bitbucket.org/innius/grafana-simple-grpc-datasource/pkg/proto/v3"
	"google.golang.org/grpc"
)

type adapter struct {
	v2Client v2.GrafanaQueryAPIClient
}

func (b *adapter) ListDimensionKeys(ctx context.Context, in *v3.ListDimensionKeysRequest, opts ...grpc.CallOption) (*v3.ListDimensionKeysResponse, error) {
	in2 := &v2.ListDimensionKeysRequest{
		Filter:             in.Filter,
		SelectedDimensions: toV2Dimensions(in.SelectedDimensions),
	}

	res2, err := b.v2Client.ListDimensionKeys(ctx, in2, opts...)
	if err != nil {
		return nil, err
	}

	r := make([]*v3.ListDimensionKeysResponse_Result, len(res2.Results))
	if res2.Results == nil {
		return &v3.ListDimensionKeysResponse{}, nil
	}

	for i, v := range res2.Results {
		r[i] = &v3.ListDimensionKeysResponse_Result{
			Key:         v.Key,
			Description: v.Description,
		}
	}

	return &v3.ListDimensionKeysResponse{
		Results: r,
	}, nil
}

func (b *adapter) ListDimensionValues(ctx context.Context, in *v3.ListDimensionValuesRequest, opts ...grpc.CallOption) (*v3.ListDimensionValuesResponse, error) {
	inv2 := &v2.ListDimensionValuesRequest{
		DimensionKey: in.DimensionKey,
		Filter:       in.Filter,
	}
	res, err := b.v2Client.ListDimensionValues(ctx, inv2, opts...)
	if err != nil {
		return nil, err
	}

	if res.Results == nil {
		return &v3.ListDimensionValuesResponse{}, nil
	}

	r := make([]*v3.ListDimensionValuesResponse_Result, len(res.Results))
	for i, v := range res.Results {
		r[i] = &v3.ListDimensionValuesResponse_Result{
			Value:       v.Value,
			Description: v.Description,
		}
	}

	return &v3.ListDimensionValuesResponse{
		Results: r,
	}, nil
}

func (b *adapter) ListMetrics(ctx context.Context, in *v3.ListMetricsRequest, opts ...grpc.CallOption) (*v3.ListMetricsResponse, error) {
	inv2 := &v2.ListMetricsRequest{
		Dimensions: toV2Dimensions(in.Dimensions),
		Filter:     in.Filter,
	}
	res, err := b.v2Client.ListMetrics(ctx, inv2, opts...)
	if err != nil {
		return nil, err
	}

	if res.Metrics == nil {
		return &v3.ListMetricsResponse{}, nil
	}

	r := make([]*v3.ListMetricsResponse_Metric, len(res.Metrics))
	for i, v := range res.Metrics {
		r[i] = &v3.ListMetricsResponse_Metric{
			Name:        v.Name,
			Description: v.Description,
		}
	}
	return &v3.ListMetricsResponse{
		Metrics: r,
	}, nil
}

func (b *adapter) GetMetricAggregate(ctx context.Context, in *v3.GetMetricAggregateRequest, opts ...grpc.CallOption) (*v3.GetMetricAggregateResponse, error) {
	in2 := &v2.GetMetricAggregateRequest{
		Dimensions:    toV2Dimensions(in.Dimensions),
		Metrics:       in.Metrics,
		AggregateType: v2.AggregateType(in.AggregateType),
		StartDate:     in.StartDate,
		EndDate:       in.EndDate,
		MaxItems:      in.MaxItems,
		TimeOrdering:  v2.TimeOrdering(in.TimeOrdering),
		StartingToken: in.StartingToken,
		IntervalMs:    in.IntervalMs,
	}

	res, err := b.v2Client.GetMetricAggregate(ctx, in2, opts...)
	if err != nil {
		return nil, err
	}

	if res.Frames == nil {
		return &v3.GetMetricAggregateResponse{
			NextToken: res.NextToken,
		}, nil
	}

	frames := make([]*v3.Frame, len(res.Frames))
	for i, frame := range res.Frames {
		frames[i] = &v3.Frame{
			Metric:     frame.Metric,
			Timestamps: frame.Timestamps,
			Meta:       toV3Meta(frame.Meta),
			Fields:     toV3Fields(frame.Fields),
		}
	}

	return &v3.GetMetricAggregateResponse{
		Frames:    frames,
		NextToken: res.NextToken,
	}, nil
}

func (b *adapter) GetMetricValue(ctx context.Context, in *v3.GetMetricValueRequest, opts ...grpc.CallOption) (*v3.GetMetricValueResponse, error) {
	in2 := &v2.GetMetricValueRequest{
		Dimensions: toV2Dimensions(in.Dimensions),
		Metrics:    in.Metrics,
	}

	res, err := b.v2Client.GetMetricValue(ctx, in2, opts...)
	if err != nil {
		return nil, err
	}

	if res.Frames == nil {
		return &v3.GetMetricValueResponse{}, nil
	}

	frames := make([]*v3.GetMetricValueResponse_Frame, len(res.Frames))
	for i, frame := range res.Frames {
		frames[i] = &v3.GetMetricValueResponse_Frame{
			Metric:    frame.Metric,
			Timestamp: frame.Timestamp,
			Meta:      toV3Meta(frame.Meta),
			Fields:    toV3SingleValueFields(frame.Fields),
		}
	}

	return &v3.GetMetricValueResponse{
		Frames: frames,
	}, nil
}

func (b *adapter) GetMetricHistory(ctx context.Context, in *v3.GetMetricHistoryRequest, opts ...grpc.CallOption) (*v3.GetMetricHistoryResponse, error) {
	in2 := &v2.GetMetricHistoryRequest{
		Dimensions:    toV2Dimensions(in.Dimensions),
		Metrics:       in.Metrics,
		StartDate:     in.StartDate,
		EndDate:       in.EndDate,
		MaxItems:      in.MaxItems,
		TimeOrdering:  v2.TimeOrdering(in.TimeOrdering),
		StartingToken: in.StartingToken,
	}
	res, err := b.v2Client.GetMetricHistory(ctx, in2, opts...)
	if err != nil {
		return nil, err
	}

	if res.Frames == nil {
		return &v3.GetMetricHistoryResponse{
			NextToken: res.NextToken,
		}, nil
	}

	frames := make([]*v3.Frame, len(res.Frames))
	for i, frame := range res.Frames {
		frames[i] = &v3.Frame{
			Metric:     frame.Metric,
			Timestamps: frame.Timestamps,
			Meta:       toV3Meta(frame.Meta),
			Fields:     toV3Fields(frame.Fields),
		}
	}

	return &v3.GetMetricHistoryResponse{
		Frames:    frames,
		NextToken: res.NextToken,
	}, nil
}
