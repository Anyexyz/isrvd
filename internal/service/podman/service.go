// Package podman 提供 Podman 业务服务层
package podman

import (
	"context"
	"fmt"

	"github.com/rehiy/pango/logman"

	"isrvd/internal/registry"
	pkgpodman "isrvd/pkgs/podman"
)

// Service Podman 业务服务
type Service struct {
	podman *pkgpodman.PodmanService
}

// NewService 创建 Podman 业务服务
func NewService() (*Service, error) {
	svc := registry.PodmanService
	if svc == nil {
		logman.Error("Podman service not initialized")
		return nil, fmt.Errorf("Podman 服务未初始化")
	}
	return &Service{podman: svc}, nil
}

// CheckAvailability 检测 Podman 可用性
func (s *Service) CheckAvailability(ctx context.Context) bool {
	if s.podman == nil {
		return false
	}
	_, err := s.podman.GetInfo(ctx)
	return err == nil
}

// GetPodmanService 返回底层 pkgs/podman.PodmanService
func (s *Service) GetPodmanService() *pkgpodman.PodmanService {
	return s.podman
}
