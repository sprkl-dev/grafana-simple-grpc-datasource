package models

import (
	"encoding/json"
	"strings"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
)

type MetricAggregateQuery struct {
	MetricBaseQuery
	AggregateType string `json:"aggregateType"`
}

func UnmarshalToMetricAggregateQuery(dq *backend.DataQuery) (*MetricAggregateQuery, error) {
	query := &MetricAggregateQuery{}
	if err := json.Unmarshal(dq.JSON, query); err != nil {
		return nil, err
	}

	// add on the DataQuery params
	query.TimeRange = dq.TimeRange
	query.Interval = dq.Interval
	query.MaxDataPoints = dq.MaxDataPoints
	query.QueryType = dq.QueryType

	return query, nil
}

func (q MetricAggregateQuery) FormatDisplayName() string {
	if q.MetricBaseQuery.DisplayName == "" {
		return ""
	}
	ctx := newContext(q.MetricBaseQuery)
	ctx["aggregate"] = strings.ToLower(q.AggregateType)

	s, err := parseDisplayNameExpr(ctx, q.MetricBaseQuery.DisplayName)
	if err != nil {
		return ""
	}
	return s
}
