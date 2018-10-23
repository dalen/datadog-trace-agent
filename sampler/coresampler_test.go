package sampler

import (
	"testing"
	"time"

	log "github.com/cihub/seelog"
	"github.com/stretchr/testify/assert"
)

func getTestSampler() *Sampler {
	// Disable debug logs in these tests
	log.UseLogger(log.Disabled)

	// No extra fixed sampling, no maximum TPS
	extraRate := 1.0
	maxTPS := 0.0

	return newSampler(extraRate, maxTPS)
}

func TestSamplerAccessRace(t *testing.T) {
	// regression test: even though the sampler is channel protected, it
	// has getters accessing its fields.
	s := newSampler(1, 2)
	go func() {
		for i := 0; i < 10000; i++ {
			s.SetSignatureCoefficients(float64(i), float64(i)/2)
		}
	}()
	for i := 0; i < 5000; i++ {
		s.GetState()
		s.GetAllCountScores()
	}
}

func TestSamplerLoop(t *testing.T) {
	s := getTestSampler()

	exit := make(chan bool)

	go func() {
		s.Run()
		close(exit)
	}()

	s.Stop()

	select {
	case <-exit:
		return
	case <-time.After(time.Second * 1):
		assert.Fail(t, "Sampler took more than 1 second to close")
	}
}
