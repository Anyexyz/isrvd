package server

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"isrvd/internal/helper"
	svcPodman "isrvd/internal/service/podman"
	pkgpodman "isrvd/pkgs/podman"
)

// podmanInfo 获取 Podman 概览信息
func (app *App) podmanInfo(c *gin.Context) {
	result, err := app.podmanSvc.Info(c.Request.Context())
	if err != nil {
		helper.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	helper.RespondSuccess(c, "Podman info retrieved", result)
}

// podmanListContainers 列出容器
func (app *App) podmanListContainers(c *gin.Context) {
	all := c.DefaultQuery("all", "false") == "true"
	result, err := app.podmanSvc.ListContainers(c.Request.Context(), all)
	if err != nil {
		helper.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	helper.RespondSuccess(c, "Containers listed successfully", result)
}

// podmanCreateContainer 创建容器
func (app *App) podmanCreateContainer(c *gin.Context) {
	var req pkgpodman.ContainerCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}
	result, err := app.podmanSvc.CreateContainer(c.Request.Context(), req)
	if err != nil {
		helper.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	helper.RespondSuccess(c, "容器创建成功", result)
}

// podmanContainerAction 容器操作
func (app *App) podmanContainerAction(c *gin.Context) {
	req := pkgpodman.ContainerActionRequest{
		ID: c.Param("id"),
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := app.podmanSvc.ContainerAction(c.Request.Context(), req); err != nil {
		helper.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	helper.RespondSuccess(c, "Container "+req.Action+" successfully", nil)
}

// podmanContainerLogs 获取容器日志
func (app *App) podmanContainerLogs(c *gin.Context) {
	req := pkgpodman.ContainerLogsRequest{
		ID:     c.Param("id"),
		Tail:   c.DefaultQuery("tail", "100"),
		Follow: c.DefaultQuery("follow", "false") == "true",
	}
	if req.ID == "" {
		helper.RespondError(c, http.StatusBadRequest, "container id is required")
		return
	}
	result, err := app.podmanSvc.ContainerLogs(c.Request.Context(), req)
	if err != nil {
		helper.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	helper.RespondSuccess(c, "Container logs retrieved", result)
}

// ─── 镜像 ───

// podmanListImages 列出镜像
func (app *App) podmanListImages(c *gin.Context) {
	all := c.DefaultQuery("all", "false") == "true"
	result, err := app.podmanSvc.ListImages(c.Request.Context(), all)
	if err != nil {
		helper.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	helper.RespondSuccess(c, "Images listed successfully", result)
}

// podmanImageAction 镜像操作
func (app *App) podmanImageAction(c *gin.Context) {
	req := svcPodman.ImageActionRequest{
		ID: c.Param("id"),
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := app.podmanSvc.ImageAction(c.Request.Context(), req); err != nil {
		helper.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	helper.RespondSuccess(c, "Image "+req.Action+" successfully", nil)
}

// podmanPullImage 拉取镜像
func (app *App) podmanPullImage(c *gin.Context) {
	var req pkgpodman.ImagePullRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}
	result, err := app.podmanSvc.PullImage(c.Request.Context(), req)
	if err != nil {
		helper.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	helper.RespondSuccess(c, "镜像拉取成功", result)
}

// podmanTagImage 为镜像添加标签
func (app *App) podmanTagImage(c *gin.Context) {
	var req pkgpodman.ImageTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := app.podmanSvc.TagImage(c.Request.Context(), req); err != nil {
		helper.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	helper.RespondSuccess(c, "镜像打标签成功", nil)
}

// podmanBuildImage 构建镜像
func (app *App) podmanBuildImage(c *gin.Context) {
	var req pkgpodman.ImageBuildRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}
	result, err := app.podmanSvc.BuildImage(c.Request.Context(), req)
	if err != nil {
		helper.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	helper.RespondSuccess(c, "镜像构建成功", result)
}

// ─── 网络 ───

// podmanListNetworks 列出网络
func (app *App) podmanListNetworks(c *gin.Context) {
	result, err := app.podmanSvc.ListNetworks(c.Request.Context())
	if err != nil {
		helper.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	helper.RespondSuccess(c, "Networks listed successfully", result)
}

// podmanNetworkAction 网络操作
func (app *App) podmanNetworkAction(c *gin.Context) {
	req := svcPodman.NetworkActionRequest{
		ID: c.Param("id"),
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := app.podmanSvc.NetworkAction(c.Request.Context(), req); err != nil {
		helper.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	helper.RespondSuccess(c, "Network "+req.Action+" successfully", nil)
}

// podmanCreateNetwork 创建网络
func (app *App) podmanCreateNetwork(c *gin.Context) {
	var req pkgpodman.NetworkCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}
	result, err := app.podmanSvc.CreateNetwork(c.Request.Context(), req)
	if err != nil {
		helper.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	helper.RespondSuccess(c, "网络创建成功", result)
}

// ─── 卷 ───

// podmanListVolumes 列出卷
func (app *App) podmanListVolumes(c *gin.Context) {
	result, err := app.podmanSvc.ListVolumes(c.Request.Context())
	if err != nil {
		helper.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	helper.RespondSuccess(c, "Volumes listed successfully", result)
}

// podmanVolumeAction 卷操作
func (app *App) podmanVolumeAction(c *gin.Context) {
	req := svcPodman.VolumeActionRequest{
		Name: c.Param("name"),
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := app.podmanSvc.VolumeAction(c.Request.Context(), req); err != nil {
		helper.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	helper.RespondSuccess(c, "Volume "+req.Action+" successfully", nil)
}

// podmanCreateVolume 创建卷
func (app *App) podmanCreateVolume(c *gin.Context) {
	var req pkgpodman.VolumeCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}
	result, err := app.podmanSvc.CreateVolume(c.Request.Context(), req)
	if err != nil {
		helper.RespondError(c, http.StatusInternalServerError, err.Error())
		return
	}
	helper.RespondSuccess(c, "卷创建成功", result)
}
