package biz

import "github.com/google/wire"

// ProviderSet is biz providers.
var ProviderSet = wire.NewSet(NewSystemUseCase)

type SystemUseCase struct {
}

func NewSystemUseCase() *SystemUseCase {
	return &SystemUseCase{}
}
