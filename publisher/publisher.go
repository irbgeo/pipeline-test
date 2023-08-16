package publisher

import (
	"context"
	"fmt"

	"github.com/irbgeo/pipline-test/controller"
)

type publisher struct {
}

func New() *publisher {
	return &publisher{}
}

func (s *publisher) PublishState(ctx context.Context, state controller.State) error {
	fmt.Println(state)
	return nil
}
