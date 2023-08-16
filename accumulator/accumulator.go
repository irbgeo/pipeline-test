package accumulator

import (
	"context"
	"sync/atomic"

	"github.com/irbgeo/pipline-test/controller"
)

type accumulator struct {
	state int64
}

func New() *accumulator {
	return &accumulator{}
}

func (s *accumulator) UpdateState(ctx context.Context, p controller.ProcessedPackage) error {
	var sum int64
	for _, n := range p {
		sum += int64(n)
	}

	atomic.StoreInt64(&s.state, sum)
	return nil
}

func (s *accumulator) GetState() controller.State {
	return controller.State(atomic.LoadInt64(&s.state))
}
