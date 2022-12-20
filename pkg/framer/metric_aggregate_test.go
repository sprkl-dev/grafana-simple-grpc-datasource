package framer

import (
	"bitbucket.org/innius/grafana-simple-grpc-datasource/pkg/models"
	pb "bitbucket.org/innius/grafana-simple-grpc-datasource/pkg/proto/v3"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/timestamppb"
	"testing"
	"time"
)

func TestMetricAggregate_Frames(t *testing.T) {
	ts := time.Date(2022, 01, 19, 16, 03, 10, 00, time.Local)

	sut := MetricAggregate{
		GetMetricAggregateResponse: &pb.GetMetricAggregateResponse{
			Frames: []*pb.Frame{
				{
					Metric: "foo",
					Timestamps: []*timestamppb.Timestamp{
						timestamppb.New(ts),
					},
					Fields: []*pb.Field{
						{
							Name:   "field_1",
							Labels: []*pb.Label{{Key: "zone", Value: "a"}},
							Config: nil,
							Values: []*pb.MetricValue{
								{
									Value: &pb.MetricValue_DoubleValue{
										DoubleValue: 10,
									},
								},
							},
						},
					},
				},
				{
					Metric: "bar",
					Timestamps: []*timestamppb.Timestamp{
						timestamppb.New(ts),
					},
					Fields: []*pb.Field{
						{
							Name:   "",
							Labels: nil,
							Config: nil,
							Values: []*pb.MetricValue{
								{
									Value: &pb.MetricValue_DoubleValue{
										DoubleValue: 20,
									},
								},
							},
						},
					},
				},
			},
			NextToken: "next-please",
		},
		Query: models.MetricBaseQuery{
			Dimensions: []models.Dimension{
				{
					Key:   "machine",
					Value: "m1",
				},
			},
			DisplayName: `{{machine}}-{{metric}}-{{zone}}-{{aggregate}}-{{field}}`,
		},
		AggregateType: pb.AggregateType_AVERAGE,
	}

	res, err := sut.Frames()

	assert.NoError(t, err)
	t.Run("the result should contain two frames", func(t *testing.T) {
		assert.Len(t, res, 2)
	})
	t.Run("the data frame should have a name", func(t *testing.T) {
		assert.Equal(t, "foo", res[0].Name)
	})

	t.Run("the format name expression should be applied", func(t *testing.T) {
		assert.Equal(t, "m1-foo-a-avg-field_1", res[0].Fields[1].Config.DisplayNameFromDS)
	})

	t.Run("the data frame should have a NextToken", func(t *testing.T) {
		assert.Equal(t, "next-please", res[0].Meta.Custom.(models.Metadata).NextToken)
	})
}
