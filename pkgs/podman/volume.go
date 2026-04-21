package podman

import (
	"context"

	"github.com/docker/docker/api/types/volume"
	"github.com/rehiy/pango/logman"
)

// VolumeInfo 卷信息
type VolumeInfo struct {
	Name       string            `json:"name"`
	Driver     string            `json:"driver"`
	MountPoint string            `json:"mountPoint"`
	Labels     map[string]string `json:"labels,omitempty"`
}

// ListVolumes 获取卷列表
func (s *PodmanService) ListVolumes(ctx context.Context) ([]*VolumeInfo, error) {
	volumes, err := s.client.VolumeList(ctx, volume.ListOptions{})
	if err != nil {
		logman.Error("List volumes failed", "error", err)
		return nil, err
	}

	var result []*VolumeInfo
	for _, vol := range volumes.Volumes {
		result = append(result, &VolumeInfo{
			Name:       vol.Name,
			Driver:     vol.Driver,
			MountPoint: vol.Mountpoint,
			Labels:     vol.Labels,
		})
	}

	return result, nil
}

// VolumeCreateRequest 创建卷请求
type VolumeCreateRequest struct {
	Name   string            `json:"name" binding:"required"`
	Driver string            `json:"driver"`
	Labels map[string]string `json:"labels"`
}

// CreateVolume 创建卷
func (s *PodmanService) CreateVolume(ctx context.Context, req VolumeCreateRequest) (string, error) {
	driver := req.Driver
	if driver == "" {
		driver = "local"
	}

	options := volume.CreateOptions{
		Name:   req.Name,
		Driver: driver,
		Labels: req.Labels,
	}

	resp, err := s.client.VolumeCreate(ctx, options)
	if err != nil {
		logman.Error("Create volume failed", "name", req.Name, "error", err)
		return "", err
	}

	logman.Info("Volume created", "name", req.Name)
	return resp.Name, nil
}

// RemoveVolume 删除卷
func (s *PodmanService) RemoveVolume(ctx context.Context, name string) error {
	err := s.client.VolumeRemove(ctx, name, true)
	if err != nil {
		logman.Error("Remove volume failed", "name", name, "error", err)
		return err
	}

	logman.Info("Volume removed", "name", name)
	return nil
}
