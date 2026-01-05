package service

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	v1 "github.com/hypercoze/hyper-admin/api/gen/go/system/v1"
	"github.com/hypercoze/hyper-admin/app/system/service/internal/biz"
)

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(NewSystemService)

type SystemService struct {
	v1.UnimplementedSystemServer

	system *biz.SystemUseCase
	log    *log.Helper
}

func NewSystemService(
	system *biz.SystemUseCase,
	logger log.Logger,
) *SystemService {
	return &SystemService{
		system: system,
		log:    log.NewHelper(log.With(logger, "module", "service/system")),
	}
}
