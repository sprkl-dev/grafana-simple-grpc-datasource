package v2

import (
	v2 "bitbucket.org/innius/grafana-simple-grpc-datasource/pkg/proto/v2"
	v3 "bitbucket.org/innius/grafana-simple-grpc-datasource/pkg/proto/v3"
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
	"testing"
	"time"
)

type v2Mock struct {
	mock.Mock
}

func (v v2Mock) ListDimensionKeys(ctx context.Context, in *v2.ListDimensionKeysRequest, opts ...grpc.CallOption) (*v2.ListDimensionKeysResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (v v2Mock) ListDimensionValues(ctx context.Context, in *v2.ListDimensionValuesRequest, opts ...grpc.CallOption) (*v2.ListDimensionValuesResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (v v2Mock) ListMetrics(ctx context.Context, in *v2.ListMetricsRequest, opts ...grpc.CallOption) (*v2.ListMetricsResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (v *v2Mock) GetMetricValue(ctx context.Context, in *v2.GetMetricValueRequest, opts ...grpc.CallOption) (*v2.GetMetricValueResponse, error) {
	args := v.Called(ctx, in)
	if v, ok := args.Get(0).(*v2.GetMetricValueResponse); ok {
		return v, args.Error(1)
	}
	return nil, args.Error(1)
}

func (v *v2Mock) GetMetricHistory(ctx context.Context, in *v2.GetMetricHistoryRequest, opts ...grpc.CallOption) (*v2.GetMetricHistoryResponse, error) {
	args := v.Called(ctx, in)
	if v, ok := args.Get(0).(*v2.GetMetricHistoryResponse); ok {
		return v, args.Error(1)
	}
	return nil, args.Error(1)
}

func (v *v2Mock) GetMetricAggregate(ctx context.Context, in *v2.GetMetricAggregateRequest, opts ...grpc.CallOption) (*v2.GetMetricAggregateResponse, error) {
	args := v.Called(ctx, in)
	if v, ok := args.Get(0).(*v2.GetMetricAggregateResponse); ok {
		return v, args.Error(1)
	}
	return nil, args.Error(1)
}

func TestAdapter_GetMetricValue(t *testing.T) {
	ts := timestamppb.New(time.Unix(1000, 0))

	req := &v3.GetMetricValueRequest{
		Dimensions: []*v3.Dimension{
			{
				Key:   "machine",
				Value: "m1",
			},
		},
		Metrics: []string{"foo"},
	}

	l := &v2.Label{
		Key:   "l.key",
		Value: "l.value",
	}
	c := &v2.Config{
		Unit: "c.unit",
	}
	n := &v2.FrameMeta_Notice{
		Severity: v2.FrameMeta_Notice_NoticeSeverityWarning,
		Text:     "n.text",
		Link:     "n.link",
		Inspect:  v2.FrameMeta_Notice_InspectTypeData,
	}
	meta := &v2.FrameMeta{
		Notices:                []*v2.FrameMeta_Notice{n},
		Type:                   v2.FrameMeta_FrameTypeDirectoryListing,
		PreferredVisualization: v2.FrameMeta_VisTypeTable,
		ExecutedQueryString:    "meta.executedquerystring",
	}
	m := &v2Mock{}
	m.On("GetMetricValue", mock.Anything, &v2.GetMetricValueRequest{
		Dimensions: []*v2.Dimension{
			{
				Key:   req.Dimensions[0].Key,
				Value: req.Dimensions[0].Value,
			},
		},
		Metrics: req.Metrics,
	}).Return(&v2.GetMetricValueResponse{
		Frames: []*v2.GetMetricValueResponse_Frame{
			{
				Metric:    req.Metrics[0],
				Timestamp: ts,
				Meta:      meta,
				Fields: []*v2.SingleValueField{
					{
						Name:   "Dummy Name",
						Value:  20,
						Labels: []*v2.Label{l},
						Config: c,
					},
				},
			},
		},
	}, nil)

	sut := &adapter{
		v2Client: m,
	}

	res, err := sut.GetMetricValue(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	m.AssertExpectations(t)
	expected := []*v3.GetMetricValueResponse_Frame{
		{
			Metric: req.Metrics[0],
			Meta:   toV3Meta(meta),
			Fields: []*v3.SingleValueField{
				{
					Name:   "Dummy Name",
					Labels: toV3Labels([]*v2.Label{l}),
					Config: toV3Config(c),
					Value: &v3.MetricValue{
						Value: &v3.MetricValue_DoubleValue{
							DoubleValue: 20,
						},
					},
				},
			},
			Timestamp: ts,
		},
	}
	assert.Equal(t, expected, res.Frames)
}

func TestAdapter_GetMetricHistory(t *testing.T) {
	req := &v3.GetMetricHistoryRequest{
		Dimensions: []*v3.Dimension{
			{
				Key:   "machine",
				Value: "m1",
			},
		},
		Metrics:       []string{"foo"},
		StartDate:     timestamppb.New(time.Unix(1000, 0)),
		EndDate:       timestamppb.New(time.Unix(2000, 0)),
		MaxItems:      30000,
		TimeOrdering:  v3.TimeOrdering_DESCENDING,
		StartingToken: "start-here",
	}

	l := &v2.Label{
		Key:   "l.key",
		Value: "l.value",
	}
	c := &v2.Config{
		Unit: "c.unit",
	}
	n := &v2.FrameMeta_Notice{
		Severity: v2.FrameMeta_Notice_NoticeSeverityWarning,
		Text:     "n.text",
		Link:     "n.link",
		Inspect:  v2.FrameMeta_Notice_InspectTypeData,
	}
	meta := &v2.FrameMeta{
		Notices:                []*v2.FrameMeta_Notice{n},
		Type:                   v2.FrameMeta_FrameTypeDirectoryListing,
		PreferredVisualization: v2.FrameMeta_VisTypeTable,
		ExecutedQueryString:    "meta.executedquerystring",
	}
	m := &v2Mock{}
	m.On("GetMetricHistory", mock.Anything, &v2.GetMetricHistoryRequest{
		Dimensions: []*v2.Dimension{
			{
				Key:   req.Dimensions[0].Key,
				Value: req.Dimensions[0].Value,
			},
		},
		Metrics:       req.Metrics,
		StartDate:     req.StartDate,
		EndDate:       req.EndDate,
		MaxItems:      req.MaxItems,
		TimeOrdering:  v2.TimeOrdering(req.TimeOrdering),
		StartingToken: req.StartingToken,
	}).Return(&v2.GetMetricHistoryResponse{
		Frames: []*v2.Frame{
			{
				Metric:     req.Metrics[0],
				Timestamps: []*timestamppb.Timestamp{req.StartDate},
				Fields: []*v2.Field{
					{
						Name:   "Dummy Name",
						Labels: []*v2.Label{l},
						Config: c,
						Values: []float64{2},
					},
				},
				Meta: meta,
			},
		},
		NextToken: "next-please",
	}, nil)

	sut := &adapter{
		v2Client: m,
	}

	res, err := sut.GetMetricHistory(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	m.AssertExpectations(t)
	expected := []*v3.Frame{
		{
			Metric: req.Metrics[0],
			Fields: []*v3.Field{
				{
					Name:   "Dummy Name",
					Labels: toV3Labels([]*v2.Label{l}),
					Config: toV3Config(c),
					Values: []*v3.MetricValue{{Value: &v3.MetricValue_DoubleValue{DoubleValue: 2}}},
				},
			},
			Timestamps: []*timestamppb.Timestamp{
				req.StartDate,
			},
			Meta: toV3Meta(meta),
		},
	}
	assert.Equal(t, expected, res.Frames)
	assert.Equal(t, "next-please", res.NextToken)
}

func TestAdapter_GetMetricAggregate(t *testing.T) {
	req := &v3.GetMetricAggregateRequest{
		Dimensions: []*v3.Dimension{
			{
				Key:   "machine",
				Value: "m1",
			},
		},
		Metrics:       []string{"foo"},
		AggregateType: v3.AggregateType_COUNT,
		StartDate:     timestamppb.New(time.Unix(1000, 0)),
		EndDate:       timestamppb.New(time.Unix(2000, 0)),
		MaxItems:      30000,
		TimeOrdering:  v3.TimeOrdering_DESCENDING,
		StartingToken: "start-here",
		IntervalMs:    999,
	}

	l := &v2.Label{
		Key:   "l.key",
		Value: "l.value",
	}
	c := &v2.Config{
		Unit: "c.unit",
	}
	n := &v2.FrameMeta_Notice{
		Severity: v2.FrameMeta_Notice_NoticeSeverityWarning,
		Text:     "n.text",
		Link:     "n.link",
		Inspect:  v2.FrameMeta_Notice_InspectTypeData,
	}
	meta := &v2.FrameMeta{
		Notices:                []*v2.FrameMeta_Notice{n},
		Type:                   v2.FrameMeta_FrameTypeDirectoryListing,
		PreferredVisualization: v2.FrameMeta_VisTypeTable,
		ExecutedQueryString:    "meta.executedquerystring",
	}
	m := &v2Mock{}
	m.On("GetMetricAggregate", mock.Anything, &v2.GetMetricAggregateRequest{
		Dimensions: []*v2.Dimension{
			{Key: req.Dimensions[0].Key, Value: req.Dimensions[0].Value},
		},
		Metrics:       req.Metrics,
		AggregateType: v2.AggregateType(req.AggregateType),
		StartDate:     req.StartDate,
		EndDate:       req.EndDate,
		MaxItems:      req.MaxItems,
		TimeOrdering:  v2.TimeOrdering(req.TimeOrdering),
		StartingToken: req.StartingToken,
		IntervalMs:    req.IntervalMs,
	}).Return(&v2.GetMetricAggregateResponse{
		Frames: []*v2.Frame{
			{
				Metric:     req.Metrics[0],
				Timestamps: []*timestamppb.Timestamp{req.StartDate},
				Fields: []*v2.Field{
					{
						Name:   "Dummy Name",
						Labels: []*v2.Label{l},
						Config: c,
						Values: []float64{2},
					},
				},
				Meta: meta,
			},
		},
		NextToken: "next-please",
	}, nil)

	sut := &adapter{
		v2Client: m,
	}

	res, err := sut.GetMetricAggregate(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	m.AssertExpectations(t)
	expected := []*v3.Frame{
		{
			Metric: req.Metrics[0],
			Fields: []*v3.Field{
				{
					Name:   "Dummy Name",
					Labels: toV3Labels([]*v2.Label{l}),
					Config: toV3Config(c),
					Values: []*v3.MetricValue{{Value: &v3.MetricValue_DoubleValue{DoubleValue: 2}}},
				},
			},
			Timestamps: []*timestamppb.Timestamp{req.StartDate},
			Meta:       toV3Meta(meta),
		},
	}
	assert.Equal(t, expected, res.Frames)
	assert.Equal(t, "next-please", res.NextToken)
}
