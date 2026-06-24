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

func (r *Request) Match(data map[string]interface{}, matcher *operators.Matcher) (bool, []string) {
	itemStr, ok := r.getMatchPart(matcher.Part, data)

	switch matcher.GetType() {
	case operators.DSLMatcher:
		return matcher.Result(matcher.MatchDSL(data)), []string{}
	case operators.SizeMatcher:
		if !ok {
			return false, []string{}
		}
		return matcher.Result(matcher.MatchSize(len(itemStr))), []string{}
	case operators.WordsMatcher:
		if !ok {
			return false, []string{}
		}
		return matcher.ResultWithMatchedSnippet(matcher.MatchWords(itemStr, data))
	case operators.RegexMatcher:
		if !ok {
			return false, []string{}
		}
		return matcher.ResultWithMatchedSnippet(matcher.MatchRegex(itemStr))
	case operators.BinaryMatcher:
		if !ok {
			return false, []string{}
		}
		return matcher.ResultWithMatchedSnippet(matcher.MatchBinary(itemStr))
	default:
		return false, []string{}
	}
}

func (r *Request) Extract(data map[string]interface{}, extractor *operators.Extractor) map[string]struct{} {
	itemStr, ok := r.getMatchPart(extractor.Part, data)

	switch extractor.GetType() {
	case operators.KValExtractor:
		return extractor.ExtractKval(data)
	case operators.DSLExtractor:
		return extractor.ExtractDSL(data)
	case operators.RegexExtractor:
		if !ok {
			return nil
		}
		return extractor.ExtractRegex(itemStr)
	default:
		return nil
	}
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
		ExtractedResults: wrapped.OperatorsResult.OutputExtracts,
		Metadata:         wrapped.OperatorsResult.PayloadValues,
		Timestamp:        time.Now(),
		IP:               common.ToString(wrapped.InternalEvent["ip"]),
	}
}
