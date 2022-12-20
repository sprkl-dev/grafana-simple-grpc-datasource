package v2

import (
	"testing"

	v2 "bitbucket.org/innius/grafana-simple-grpc-datasource/pkg/proto/v2"
	v3 "bitbucket.org/innius/grafana-simple-grpc-datasource/pkg/proto/v3"
	"github.com/stretchr/testify/assert"
)

func TestConvert_ToV2Demensions(t *testing.T) {
	v3d := []*v3.Dimension{
		{
			Key:   "Dummy Key 1",
			Value: "Dummy Value 1",
		},
		{
			Key:   "Dummy Key 2",
			Value: "Dummy Value 2",
		},
		{
			Key:   "Dummy Key 3",
			Value: "Dummy Value 3",
		},
	}

	v2d := toV2Dimensions(v3d)

	assert.Len(t, v2d, len(v3d))
	for i := range v2d {
		assert.Equal(t, &v2.Dimension{Key: v3d[i].Key, Value: v3d[i].Value}, v2d[i])
	}
}

func TestConvert_ToV3Config(t *testing.T) {
	v2c := &v2.Config{
		Unit: "Dummy Unit",
	}

	v3c := toV3Config(v2c)
	assert.Equal(t, &v3.Config{Unit: v2c.Unit}, v3c)
}

func TestConvert_ToV3Meta(t *testing.T) {
	v2m := &v2.FrameMeta{
		Type:                   v2.FrameMeta_FrameTypeTimeSeriesMany,
		ExecutedQueryString:    "Dummy Query",
		PreferredVisualization: v2.FrameMeta_VisTypeNodeGraph,
		Notices: []*v2.FrameMeta_Notice{
			{
				Severity: v2.FrameMeta_Notice_NoticeSeverityError,
				Text:     "Dummy Text",
				Link:     "Dummy Link",
				Inspect:  v2.FrameMeta_Notice_InspectTypeMeta,
			},
		},
	}

	v3m := toV3Meta(v2m)
	assert.Equal(t, &v3.FrameMeta{
		Type:                   v3.FrameMeta_FrameType(v2m.Type),
		Notices:                toV3Notices(v2m.Notices),
		PreferredVisualization: v3.FrameMeta_VisType(v2m.PreferredVisualization),
		ExecutedQueryString:    v2m.ExecutedQueryString,
	}, v3m)
}

func TestConvert_ToV3Notices(t *testing.T) {
	v2n := []*v2.FrameMeta_Notice{
		{
			Severity: v2.FrameMeta_Notice_NoticeSeverityError,
			Text:     "Dummy Text 1",
			Link:     "Dummy Link 1",
			Inspect:  v2.FrameMeta_Notice_InspectTypeMeta,
		},
		{
			Severity: v2.FrameMeta_Notice_NoticeSeverityWarning,
			Text:     "Dummy Text 2",
			Link:     "Dummy Link 2",
			Inspect:  v2.FrameMeta_Notice_InspectTypeStats,
		},
	}

	v3n := toV3Notices(v2n)
	assert.Len(t, v3n, len(v2n))

	for i := range v3n {
		assert.Equal(t, &v3.FrameMeta_Notice{
			Severity: v3.FrameMeta_Notice_NoticeSeverity(v2n[i].Severity),
			Text:     v2n[i].Text,
			Link:     v2n[i].Link,
			Inspect:  v3.FrameMeta_Notice_InspectType(v2n[i].Inspect),
		}, v3n[i])
	}
}

func TestConvert_ToV3Labels(t *testing.T) {
	v2l := []*v2.Label{
		{
			Key:   "Dummy Label 1",
			Value: "Dummy Value 1",
		},
		{
			Key:   "Dummy Label 2",
			Value: "Dummy Value 2",
		},
	}

	v3l := toV3Labels(v2l)

	assert.Len(t, v3l, len(v2l))
	for i := range v3l {
		assert.Equal(t, &v3.Label{
			Key:   v2l[i].Key,
			Value: v2l[i].Value,
		}, v3l[i])
	}
}

func TestConvert_ToV3Fields(t *testing.T) {
	v2f := []*v2.Field{
		{
			Name: "Dummy Name 1",
			Labels: []*v2.Label{
				{
					Key:   "Dummy Label 1",
					Value: "Dummy Value 1",
				},
			},
			Config: &v2.Config{
				Unit: "Dummy unit",
			},
			Values: []float64{0, 1, 2, 3, -1},
		},
		{
			Name: "Dummy Name 2",
			Labels: []*v2.Label{
				{
					Key:   "Dummy Label 1",
					Value: "Dummy Value 1",
				},
			},
			Config: &v2.Config{
				Unit: "Dummy unit",
			},
			Values: []float64{99, 1, 2, 3, 4},
		},
	}

	v3f := toV3Fields(v2f)

	assert.Len(t, v3f, len(v2f))
	for i := range v3f {
		assert.Equal(t, &v3.Field{
			Name:   v2f[i].Name,
			Labels: toV3Labels(v2f[i].Labels),
			Config: toV3Config(v2f[i].Config),
			Values: []*v3.MetricValue{
				{
					Value: &v3.MetricValue_DoubleValue{DoubleValue: v2f[i].Values[0]},
				},
				{
					Value: &v3.MetricValue_DoubleValue{DoubleValue: v2f[i].Values[1]},
				},
				{
					Value: &v3.MetricValue_DoubleValue{DoubleValue: v2f[i].Values[2]},
				},
				{
					Value: &v3.MetricValue_DoubleValue{DoubleValue: v2f[i].Values[3]},
				},
				{
					Value: &v3.MetricValue_DoubleValue{DoubleValue: v2f[i].Values[4]},
				},
			},
		}, v3f[i])
	}
}
