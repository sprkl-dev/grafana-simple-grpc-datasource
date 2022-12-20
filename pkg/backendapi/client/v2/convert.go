package v2

import (
	v2 "bitbucket.org/innius/grafana-simple-grpc-datasource/pkg/proto/v2"
	v3 "bitbucket.org/innius/grafana-simple-grpc-datasource/pkg/proto/v3"
)

// ##################################
// Convert from V3 --> V2 functions
// ##################################

func toV2Dimensions(dims []*v3.Dimension) []*v2.Dimension {
	if dims == nil {
		return nil
	}

	d := make([]*v2.Dimension, len(dims))
	for i := range dims {
		v := dims[i]
		d[i] = &v2.Dimension{
			Key:   v.Key,
			Value: v.Value,
		}
	}
	return d
}

// ##################################
// Convert from V2 --> V3 functions
// ##################################

func toV3Fields(fields []*v2.Field) []*v3.Field {
	if fields == nil {
		return nil
	}

	f := make([]*v3.Field, len(fields))
	for i, v := range fields {
		f[i] = &v3.Field{
			Name:   v.Name,
			Config: toV3Config(v.Config),
			Labels: toV3Labels(v.Labels),
			Values: make([]*v3.MetricValue, len(v.Values)),
		}

		if v.Values == nil {
			continue
		}

		for j, value := range v.Values {
			f[i].Values[j] = &v3.MetricValue{
				Value: &v3.MetricValue_DoubleValue{
					DoubleValue: value,
				},
			}
		}
	}

	return f
}

func toV3Config(cfg *v2.Config) *v3.Config {
	if cfg == nil {
		return nil
	}

	return &v3.Config{
		Unit: cfg.Unit,
	}
}

func toV3Labels(labels []*v2.Label) []*v3.Label {
	if labels == nil {
		return nil
	}

	l := make([]*v3.Label, len(labels))
	for i, v := range labels {
		l[i] = &v3.Label{
			Key:   v.Key,
			Value: v.Value,
		}
	}

	return l
}

func toV3SingleValueFields(fields []*v2.SingleValueField) []*v3.SingleValueField {
	if fields == nil {
		return nil
	}

	f := make([]*v3.SingleValueField, len(fields))
	for i, v := range fields {
		f[i] = &v3.SingleValueField{
			Name:   v.Name,
			Labels: toV3Labels(v.Labels),
			Config: toV3Config(v.Config),
			Value: &v3.MetricValue{
				Value: &v3.MetricValue_DoubleValue{
					DoubleValue: v.Value,
				},
			},
		}
	}

	return f
}

func toV3Notices(v2notices []*v2.FrameMeta_Notice) []*v3.FrameMeta_Notice {
	if v2notices == nil {
		return nil
	}

	v3notices := make([]*v3.FrameMeta_Notice, len(v2notices))
	for i, notice := range v2notices {
		v3notices[i] = &v3.FrameMeta_Notice{
			Severity: v3.FrameMeta_Notice_NoticeSeverity(notice.Severity),
			Text:     notice.Text,
			Link:     notice.Link,
			Inspect:  v3.FrameMeta_Notice_InspectType(notice.Inspect),
		}
	}

	return v3notices
}

func toV3Meta(v2meta *v2.FrameMeta) *v3.FrameMeta {
	if v2meta == nil {
		return nil
	}

	return &v3.FrameMeta{
		Type:                   v3.FrameMeta_FrameType(v2meta.Type),
		PreferredVisualization: v3.FrameMeta_VisType(v2meta.PreferredVisualization),
		ExecutedQueryString:    v2meta.ExecutedQueryString,
		Notices:                toV3Notices(v2meta.Notices),
	}
}
