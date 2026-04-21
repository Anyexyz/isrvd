package podman

import (
	"context"

	pkgpodman "isrvd/pkgs/podman"
)

// ImageActionRequest 镜像操作请求
type ImageActionRequest struct {
	ID     string `json:"id" binding:"required"`
	Action string `json:"action" binding:"required"`
}

// ImagePushRequest 推送镜像请求
type ImagePushRequest struct {
	Image      string `json:"image" binding:"required"`
	Registry   string `json:"registry"`
	Username   string `json:"username"`
	Password   string `json:"password"`
}

// ImagePullFromRegistryRequest 从仓库拉取镜像请求
type ImagePullFromRegistryRequest struct {
	Registry   string `json:"registry" binding:"required"`
	Image      string `json:"image" binding:"required"`
	Username   string `json:"username"`
	Password   string `json:"password"`
}

// ListImages 列出镜像
func (s *Service) ListImages(ctx context.Context, all bool) (any, error) {
	return s.podman.ListImages(ctx, all)
}

// ImageAction 镜像操作
func (s *Service) ImageAction(ctx context.Context, req ImageActionRequest) error {
	return s.podman.RemoveImage(ctx, req.ID)
}

// PullImage 拉取镜像
func (s *Service) PullImage(ctx context.Context, req pkgpodman.ImagePullRequest) (string, error) {
	err := s.podman.PullImage(ctx, req)
	if err != nil {
		return "", err
	}
	return req.Image, nil
}

// TagImage 为镜像添加标签
func (s *Service) TagImage(ctx context.Context, req pkgpodman.ImageTagRequest) error {
	return s.podman.TagImage(ctx, req)
}

// BuildImage 构建镜像
func (s *Service) BuildImage(ctx context.Context, req pkgpodman.ImageBuildRequest) (string, error) {
	err := s.podman.BuildImage(ctx, req)
	if err != nil {
		return "", err
	}
	return req.Tags[0], nil
}
