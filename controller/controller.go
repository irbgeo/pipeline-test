package controller

import (
	"context"
	"log/slog"
	"time"
)

type controller struct {
	cancel context.CancelFunc

	source      source
	processor   processor
	accumulator accumulator
	publisher   publisher

	publishTimeout time.Duration
}

type source interface {
	PackageWatcher(ctx context.Context) <-chan Package
}

type processor interface {
	Pipeline(ctx context.Context) (chan<- Package, <-chan ProcessedPackage)
}

type accumulator interface {
	UpdateState(ctx context.Context, p ProcessedPackage) error
	GetState() State
}

type publisher interface {
	PublishState(ctx context.Context, s State) error
}

func New(
	source source,
	processor processor,
	accumulator accumulator,
	publisher publisher,

	publishTimeout time.Duration,
) *controller {
	return &controller{
		source:      source,
		processor:   processor,
		accumulator: accumulator,
		publisher:   publisher,

		publishTimeout: publishTimeout,
	}
}

func (s *controller) Start() {
	var ctx context.Context
	ctx, s.cancel = context.WithCancel(context.Background())

	in, out := s.processor.Pipeline(ctx)

	go s.gettingData(ctx, in)
	go s.processing(ctx, out)
	go s.publishing(ctx)
}

func (s *controller) Stop() {
	s.cancel()
}

func (s *controller) gettingData(ctx context.Context, in chan<- Package) {
	packageCh := s.source.PackageWatcher(ctx)
	for p := range packageCh {
		in <- p
	}
}

func (s *controller) processing(ctx context.Context, out <-chan ProcessedPackage) {
	for r := range out {
		if err := s.accumulator.UpdateState(ctx, r); err != nil {
			slog.Error("failed to update accumulator", "processed package", r, err)
		}
	}
}

func (s *controller) publishing(ctx context.Context) {
	t := time.NewTicker(s.publishTimeout)
	for {
		select {
		case <-t.C:
			if err := s.publisher.PublishState(ctx, s.accumulator.GetState()); err != nil {
				slog.Error("failed to publish state", err)
			}
		case <-ctx.Done():
			t.Stop()
			return
		}
	}
}
