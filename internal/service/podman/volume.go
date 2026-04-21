package podman

import (
	"context"

	pkgpodman "isrvd/pkgs/podman"
)

// VolumeActionRequest 卷操作请求
type VolumeActionRequest struct {
	Name     string `json:"name" binding:"required"`
	Action   string `json:"action" binding:"required"`
}

// ListVolumes 列出卷
func (s *Service) ListVolumes(ctx context.Context) (any, error) {
	return s.podman.ListVolumes(ctx)
}

// VolumeAction 卷操作
func (s *Service) VolumeAction(ctx context.Context, req VolumeActionRequest) error {
	return s.podman.RemoveVolume(ctx, req.Name)
}

// CreateVolume 创建卷
func (s *Service) CreateVolume(ctx context.Context, req pkgpodman.VolumeCreateRequest) (string, error) {
	return s.podman.CreateVolume(ctx, req)
}
