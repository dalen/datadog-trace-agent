package event

import (
	"github.com/DataDog/datadog-trace-agent/internal/agent"
)

// metricBasedExtractor is an event extractor that decides whether to extract APM events from spans based on
// the value of the event extraction rate metric set on those spans.
type metricBasedExtractor struct{}

// NewMetricBasedExtractor returns an APM event extractor that decides whether to extract APM events from spans based on
// the value of the event extraction rate metric set on those span.
func NewMetricBasedExtractor() Extractor {
	return &metricBasedExtractor{}
}

// Extract decides whether to extract APM events from a span based on the value of the event extraction rate metric set
// on that span. If such a value exists, the extracted event is returned along with this rate and a true value.
// Otherwise, false is returned as the third value and the others are invalid.
//
// NOTE: If priority is UserKeep (manually sampled) any extraction rate bigger than 0 is upscaled to 1 to ensure no
// extraction sampling is done on this event.
func (e *metricBasedExtractor) Extract(s *agent.WeightedSpan, priority agent.SamplingPriority) (*agent.Event, float64, bool) {
	extractionRate, ok := s.GetEventExtractionRate()
	if !ok {
		return nil, 0, false
	}
	if extractionRate > 0 && priority >= agent.PriorityUserKeep {
		// If the trace has been manually sampled, we keep all matching spans
		extractionRate = 1
	}
	return &agent.Event{
		Priority: priority,
		Span:     s.Span,
	}, extractionRate, true
}
