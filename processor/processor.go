package processor

import (
	"context"
	"slices"

	"github.com/irbgeo/pipline-test/controller"
)

type processor struct {
	workers chan struct{}
}

func New(workersNumber int) *processor {
	return &processor{workers: make(chan struct{}, workersNumber)}
}

func (s *processor) Pipeline(ctx context.Context) (chan<- controller.Package, <-chan controller.ProcessedPackage) {
	in, out := make(chan controller.Package), make(chan controller.ProcessedPackage)

	go func() {
		for {
			p := <-in

			s.workers <- struct{}{}
			go func() {
				defer func() {
					<-s.workers
				}()

				out <- find3Largest(p)
			}()
		}
	}()

	return in, out
}

func find3Largest(data []int) []int {
	if len(data) < 3 {
		return nil
	}

	slices.Sort(data)
	return data[len(data)-3:]
}
