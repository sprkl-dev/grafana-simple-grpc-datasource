package framer

import (
	"bitbucket.org/innius/grafana-simple-grpc-datasource/pkg/models"
	pb "bitbucket.org/innius/grafana-simple-grpc-datasource/pkg/proto/v3"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/timestamppb"
	"testing"
	"time"
)

func TestGetHistoryResponseFrameConversion(t *testing.T) {
	frame := &pb.Frame{
		Metric:     "my-metric",
		Timestamps: []*timestamppb.Timestamp{timestamppb.New(time.Unix(1000, 0)), timestamppb.New(time.Unix(2000, 0)), timestamppb.New(time.Unix(3000, 0))},
		Fields: []*pb.Field{
			{
				Name:   "v1",
				Labels: nil,
				Values: []*pb.MetricValue{
					{
						Value: &pb.MetricValue_DoubleValue{DoubleValue: 10},
					},
					{
						Value: &pb.MetricValue_DoubleValue{DoubleValue: 20},
					},
					{
						Value: &pb.MetricValue_DoubleValue{DoubleValue: 30},
					},
				},
			},
		},
	}

	sut := MetricHistory{
		GetMetricHistoryResponse: &pb.GetMetricHistoryResponse{
			Frames:    []*pb.Frame{frame},
			NextToken: "",
		},
		Query: models.MetricHistoryQuery{},
	}

	res, err := sut.Frames()
	assert.NoError(t, err)

	assert.Len(t, res, 1)

}

func TestMetricHistory_Frames(t *testing.T) {
	ts := time.Date(2022, 01, 19, 16, 03, 10, 00, time.Local)

	sut := MetricHistory{
		GetMetricHistoryResponse: &pb.GetMetricHistoryResponse{
			Frames: []*pb.Frame{
				{
					Metric: "foo",

					Fields: []*pb.Field{
						{
							Name: "field_1",
							Labels: []*pb.Label{
								{
									Key:   "zone",
									Value: "a",
								},
							},
							Config: &pb.Config{
								Unit: "℃",
							},
							Values: []*pb.MetricValue{
								{
									Value: &pb.MetricValue_DoubleValue{DoubleValue: 10},
								},
							},
						},
					},
					Timestamps: []*timestamppb.Timestamp{timestamppb.New(ts)},
				},
				{
					Metric: "bar",
					Fields: []*pb.Field{
						{
							Name:   "",
							Labels: nil,
							Config: nil,
							Values: []*pb.MetricValue{
								{
									Value: &pb.MetricValue_DoubleValue{DoubleValue: 20},
								},
							},
						},
					},
					Timestamps: []*timestamppb.Timestamp{timestamppb.New(ts)},
				},
			},
			NextToken: "next-please",
		},
		Query: models.MetricHistoryQuery{
			MetricBaseQuery: models.MetricBaseQuery{
				Dimensions: []models.Dimension{
					{
						Key:   "machine",
						Value: "m1",
					},
				},
				DisplayName: `{{machine}}-{{metric}}-{{zone}}-{{field}}`,
			},
		},
	}

	res, err := sut.Frames()

	assert.NoError(t, err)

	t.Run("the result should contain two frames", func(t *testing.T) {
		assert.Len(t, res, 2)
	})
	t.Run("the data frame should have a name", func(t *testing.T) {
		assert.Equal(t, "foo", res[0].Name)
	})
	t.Run("the data field", func(t *testing.T) {
		dataField := res[0].Fields[1]
		t.Run("the format name expression should be applied", func(t *testing.T) {
			assert.Equal(t, "m1-foo-a-field_1", dataField.Config.DisplayNameFromDS)
		})
		t.Run("the data field should have the expected unit", func(t *testing.T) {
			assert.Equal(t, "℃", dataField.Config.Unit)

		})
		t.Run("should have labels", func(t *testing.T) {
			assert.Equal(t, data.Labels{"zone": "a"}, dataField.Labels)
		})
	})

	t.Run("the data frame should have a NextToken", func(t *testing.T) {
		assert.Equal(t, "next-please", res[0].Meta.Custom.(models.Metadata).NextToken)
	})

}
