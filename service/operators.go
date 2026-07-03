package service

import (
	"time"

	"github.com/chainreactors/neutron/common"
	"github.com/chainreactors/neutron/operators"
	"github.com/chainreactors/neutron/protocols"
)

func (r *Request) getMatchPart(part string, data protocols.InternalEvent) (string, bool) {
	switch part {
	case "", "body", "all", "data":
		part = "response"
	}
	item, ok := data[part]
	if !ok {
		return "", false
	}
	return common.ToString(item), true
}

func (r *Request) Match(data map[string]interface{}, matcher *operators.Matcher) (bool, []operators.MatchHit) {
	return protocols.MakeDefaultMatchFunc(data, matcher, func(part string) (string, bool) {
		return r.getMatchPart(part, data)
	})
}

func (r *Request) Extract(data map[string]interface{}, extractor *operators.Extractor) map[string]struct{} {
	return protocols.MakeDefaultExtractFunc(data, extractor, func(part string) (string, bool) {
		return r.getMatchPart(part, data)
	})
}

func (r *Request) responseToDSLMap(response, serviceName, host string) protocols.InternalEvent {
	return protocols.InternalEvent{
		"response": response,
		"service":  serviceName,
		"host":     host,
		"type":     "service",
	}
}

func (r *Request) MakeResultEvent(wrapped *protocols.InternalWrappedEvent) []*protocols.ResultEvent {
	return protocols.MakeDefaultResultEvent(r, wrapped)
}

func (r *Request) MakeResultEventItem(wrapped *protocols.InternalWrappedEvent) *protocols.ResultEvent {
	return &protocols.ResultEvent{
		TemplateID:       common.ToString(wrapped.InternalEvent["template-id"]),
		Type:             common.ToString(wrapped.InternalEvent["type"]),
		Host:             common.ToString(wrapped.InternalEvent["host"]),
		Matched:          common.ToString(wrapped.InternalEvent["matched"]),
		ExtractedResults: wrapped.OperatorsResult.OutputExtracts(),
		Metadata:         wrapped.OperatorsResult.PayloadValues,
		Timestamp:        time.Now(),
		IP:               common.ToString(wrapped.InternalEvent["ip"]),
	}
}
