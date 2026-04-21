package podman

import (
	"context"
	"io"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/rehiy/pango/logman"
)

// ImageInfo 镜像信息
type ImageInfo struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Tags        []string `json:"tags"`
	Created     int64    `json:"created"`
	Size        int64    `json:"size"`
	VirtualSize int64    `json:"virtualSize"`
}

// ListImages 获取镜像列表
func (s *PodmanService) ListImages(ctx context.Context, all bool) ([]*ImageInfo, error) {
	images, err := s.client.ImageList(ctx, types.ImageListOptions{All: all})
	if err != nil {
		logman.Error("List images failed", "error", err)
		return nil, err
	}

	var result []*ImageInfo
	for _, img := range images {
		var tags []string
		if len(img.RepoTags) > 0 {
			tags = img.RepoTags
		} else {
			tags = []string{"<none>:<none>"}
		}

		result = append(result, &ImageInfo{
			ID:          img.ID[:12],
			Name:        getImageName(tags),
			Tags:        tags,
			Created:     img.Created,
			Size:        img.Size,
			VirtualSize: img.VirtualSize,
		})
	}

	return result, nil
}

// RemoveImage 删除镜像
func (s *PodmanService) RemoveImage(ctx context.Context, id string) error {
	_, err := s.client.ImageRemove(ctx, id, types.ImageRemoveOptions{Force: true, PruneChildren: true})
	if err != nil {
		logman.Error("Remove image failed", "id", id, "error", err)
		return err
	}

	logman.Info("Image removed", "id", id)
	return nil
}

// ImagePullRequest 拉取镜像请求
type ImagePullRequest struct {
	Image      string `json:"image" binding:"required"`
	Registry   string `json:"registry"`
	Username   string `json:"username"`
	Password   string `json:"password"`
}

// PullImage 拉取镜像
func (s *PodmanService) PullImage(ctx context.Context, req ImagePullRequest) error {
	reader, err := s.client.ImagePull(ctx, req.Image, types.ImagePullOptions{})
	if err != nil {
		logman.Error("Pull image failed", "image", req.Image, "error", err)
		return err
	}
	defer reader.Close()

	_, err = io.ReadAll(reader)
	if err != nil {
		logman.Error("Read pull image response failed", "error", err)
		return err
	}

	logman.Info("Image pulled", "image", req.Image)
	return nil
}

// ImageTagRequest 镜像标签请求
type ImageTagRequest struct {
	ID        string `json:"id" binding:"required"`
	NewTag    string `json:"newTag" binding:"required"`
}

// TagImage 为镜像添加标签
func (s *PodmanService) TagImage(ctx context.Context, req ImageTagRequest) error {
	err := s.client.ImageTag(ctx, req.ID, req.NewTag)
	if err != nil {
		logman.Error("Tag image failed", "id", req.ID, "tag", req.NewTag, "error", err)
		return err
	}

	logman.Info("Image tagged", "id", req.ID, "tag", req.NewTag)
	return nil
}

// ImageBuildRequest 构建镜像请求
type ImageBuildRequest struct {
	Context    []byte `json:"context" binding:"required"`
	Dockerfile string `json:"dockerfile" binding:"required"`
	Tags       []string `json:"tags" binding:"required"`
}

// BuildImage 构建镜像
func (s *PodmanService) BuildImage(ctx context.Context, req ImageBuildRequest) error {
	options := types.ImageBuildOptions{
		Dockerfile: req.Dockerfile,
		Tags:       req.Tags,
		Remove:     true,
	}

	reader, err := s.client.ImageBuild(ctx, strings.NewReader(string(req.Context)), options)
	if err != nil {
		logman.Error("Build image failed", "error", err)
		return err
	}
	defer reader.Body.Close()

	_, err = io.ReadAll(reader.Body)
	if err != nil {
		logman.Error("Read build image response failed", "error", err)
		return err
	}

	logman.Info("Image built", "tags", req.Tags)
	return nil
}

// 辅助函数：获取镜像名称
func getImageName(tags []string) string {
	if len(tags) > 0 && tags[0] != "<none>:<none>" {
		return tags[0]
	}
	return "<none>:<none>"
}
