package source

import (
	"context"
	"math/rand"
	"time"

	"github.com/irbgeo/pipline-test/controller"
)

type source struct {
	timeout time.Duration
}

func New(timeout time.Duration) *source {
	return &source{timeout: timeout}
}

func (s *source) PackageWatcher(ctx context.Context) <-chan controller.Package {
	ch := make(chan controller.Package)

	go func() {
		t := time.NewTicker(s.timeout)
		for {
			select {
			case <-t.C:
				ch <- generateData()
			case <-ctx.Done():
				t.Stop()
				close(ch)
				return
			}
		}
	}()

	return ch
}

func generateData() []int {
	data := make([]int, 10)
	for i := range data {
		data[i] = rand.Intn(100)
	}

	return data
}
