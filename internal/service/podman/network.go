package podman

import (
	"context"

	pkgpodman "isrvd/pkgs/podman"
)

// NetworkActionRequest 网络操作请求
type NetworkActionRequest struct {
	ID     string `json:"id" binding:"required"`
	Action string `json:"action" binding:"required"`
}

// ListNetworks 列出网络
func (s *Service) ListNetworks(ctx context.Context) (any, error) {
	return s.podman.ListNetworks(ctx)
}

// NetworkAction 网络操作
func (s *Service) NetworkAction(ctx context.Context, req NetworkActionRequest) error {
	return s.podman.RemoveNetwork(ctx, req.ID)
}

// CreateNetwork 创建网络
func (s *Service) CreateNetwork(ctx context.Context, req pkgpodman.NetworkCreateRequest) (string, error) {
	return s.podman.CreateNetwork(ctx, req)
}
