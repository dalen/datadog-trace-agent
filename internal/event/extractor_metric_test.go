package event

import (
	"math/rand"
	"testing"

	"github.com/DataDog/datadog-trace-agent/internal/agent"
)

func createTestSpansWithEventRate(eventRate float64) []*agent.WeightedSpan {
	spans := make([]*agent.WeightedSpan, 1000)
	for i, _ := range spans {
		spans[i] = &agent.WeightedSpan{Span: &agent.Span{TraceID: rand.Uint64(), Service: "test", Name: "test"}}
		if eventRate >= 0 {
			spans[i].SetMetric(agent.KeySamplingRateEventExtraction, eventRate)
		}
	}
	return spans
}

func TestMetricBasedExtractor(t *testing.T) {
	tests := []extractorTestCase{
		// Name: <priority>/<extraction rate>
		{"none/missing", createTestSpansWithEventRate(-1), 0, -1},
		{"none/0", createTestSpansWithEventRate(0), 0, 0},
		{"none/0.5", createTestSpansWithEventRate(0.5), 0, 0.5},
		{"none/1", createTestSpansWithEventRate(1), 0, 1},
		{"1/missing", createTestSpansWithEventRate(-1), 1, -1},
		{"1/0", createTestSpansWithEventRate(0), 1, 0},
		{"1/0.5", createTestSpansWithEventRate(0.5), 1, 0.5},
		{"1/1", createTestSpansWithEventRate(1), 1, 1},
		// Priority 2 should have extraction rate of 1 so long as any extraction rate is set and > 0
		{"2/missing", createTestSpansWithEventRate(-1), 2, -1},
		{"2/0", createTestSpansWithEventRate(0), 2, 0},
		{"2/0.5", createTestSpansWithEventRate(0.5), 2, 1},
		{"2/1", createTestSpansWithEventRate(1), 2, 1},
	}

	for _, test := range tests {
		testExtractor(t, NewMetricBasedExtractor(), test)
	}
}
