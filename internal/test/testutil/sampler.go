package testutil

import (
	"github.com/DataDog/datadog-trace-agent/internal/agent"
	"github.com/DataDog/datadog-trace-agent/internal/sampler"
)

// MockEngine mocks a sampler engine
type MockEngine struct {
	wantSampled bool
	wantRate    float64
}

// NewMockEngine returns a MockEngine for tests
func NewMockEngine(wantSampled bool, wantRate float64) *MockEngine {
	return &MockEngine{wantSampled: wantSampled, wantRate: wantRate}
}

// Sample returns a constant rate
func (e *MockEngine) Sample(_ agent.Trace, _ *agent.Span, _ string) (bool, float64) {
	return e.wantSampled, e.wantRate
}

// Run mocks Engine.Run()
func (e *MockEngine) Run() {
	return
}

// Stop mocks Engine.Stop()
func (e *MockEngine) Stop() {
	return
}

// GetState mocks Engine.GetState()
func (e *MockEngine) GetState() interface{} {
	return nil
}

// GetType mocks Engine.GetType()
func (e *MockEngine) GetType() sampler.EngineType {
	return sampler.NormalScoreEngineType
}
