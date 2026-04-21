package podman

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/network"
	"github.com/rehiy/pango/logman"
)

// NetworkInfo 网络信息
type NetworkInfo struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Driver      string            `json:"driver"`
	Scope       string            `json:"scope"`
	Subnet      string            `json:"subnet"`
	Gateway     string            `json:"gateway"`
	Containers  int               `json:"containers"`
	Labels      map[string]string `json:"labels,omitempty"`
}

// ListNetworks 获取网络列表
func (s *PodmanService) ListNetworks(ctx context.Context) ([]*NetworkInfo, error) {
	networks, err := s.client.NetworkList(ctx, types.NetworkListOptions{})
	if err != nil {
		logman.Error("List networks failed", "error", err)
		return nil, err
	}

	var result []*NetworkInfo
	for _, net := range networks {
		subnet := ""
		gateway := ""
		if len(net.IPAM.Config) > 0 {
			subnet = net.IPAM.Config[0].Subnet
			gateway = net.IPAM.Config[0].Gateway
		}

		result = append(result, &NetworkInfo{
			ID:         net.ID[:12],
			Name:       net.Name,
			Driver:     net.Driver,
			Scope:      net.Scope,
			Subnet:     subnet,
			Gateway:    gateway,
			Containers: len(net.Containers),
			Labels:     net.Labels,
		})
	}

	return result, nil
}

// NetworkCreateRequest 创建网络请求
type NetworkCreateRequest struct {
	Name     string            `json:"name" binding:"required"`
	Driver   string            `json:"driver" binding:"required"`
	Subnet   string            `json:"subnet"`
	Gateway  string            `json:"gateway"`
	Labels   map[string]string `json:"labels"`
}

// CreateNetwork 创建网络
func (s *PodmanService) CreateNetwork(ctx context.Context, req NetworkCreateRequest) (string, error) {
	ipamConfig := &network.IPAMConfig{}
	if req.Subnet != "" {
		ipamConfig.Subnet = req.Subnet
	}
	if req.Gateway != "" {
		ipamConfig.Gateway = req.Gateway
	}

	ipam := &network.IPAM{}
	if req.Subnet != "" || req.Gateway != "" {
		ipam.Config = []network.IPAMConfig{*ipamConfig}
	}

	driver := req.Driver
	if driver == "" {
		driver = "bridge"
	}

	options := types.NetworkCreate{
		Driver: driver,
		IPAM:   ipam,
		Labels: req.Labels,
	}

	resp, err := s.client.NetworkCreate(ctx, req.Name, options)
	if err != nil {
		logman.Error("Create network failed", "name", req.Name, "error", err)
		return "", err
	}

	logman.Info("Network created", "name", req.Name, "id", resp.ID)
	return resp.ID, nil
}

// RemoveNetwork 删除网络
func (s *PodmanService) RemoveNetwork(ctx context.Context, id string) error {
	err := s.client.NetworkRemove(ctx, id)
	if err != nil {
		logman.Error("Remove network failed", "id", id, "error", err)
		return err
	}

	logman.Info("Network removed", "id", id)
	return nil
}
