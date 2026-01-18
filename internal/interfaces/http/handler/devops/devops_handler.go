package devops

import (
	"FLOWGO/internal/application/dto"
	"FLOWGO/internal/application/service/devops"
	"FLOWGO/internal/interfaces/http/handler"
	"net/http"

	"github.com/gin-gonic/gin"
)

type DevOpsHandler struct {
	handler.BaseHandler
	devopsService *devops.DevOpsService
}

func NewDevOpsHandler(devopsService *devops.DevOpsService) *DevOpsHandler {
	return &DevOpsHandler{
		devopsService: devopsService,
	}
}

// ConfigRepo 配置仓库
func (h *DevOpsHandler) ConfigRepo(c *gin.Context) {
	var req dto.ConfigRepoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.HandleBadRequest(c, "Invalid request parameters")
		return
	}

	// Remove ProjectID binding as it is global
	// req.ProjectID = ...

	resp, err := h.devopsService.ConfigRepo(c.Request.Context(), req)
	if err != nil {
		h.HandleInternalError(c, err.Error())
		return
	}

	h.HandleSuccess(c, resp)
}

// GetSummary 获取 DevOps 概览
func (h *DevOpsHandler) GetSummary(c *gin.Context) {
	// No projectID param needed
	resp, err := h.devopsService.GetSummary(c.Request.Context())
	if err != nil {
		h.HandleInternalError(c, err.Error())
		return
	}

	h.HandleSuccess(c, resp)
}

// HandleGitHubWebhook 处理 GitHub Webhook
func (h *DevOpsHandler) HandleGitHubWebhook(c *gin.Context) {
	// 简化：这里应该解析 GitHub 特定的 Payload 并转换为通用 Payload
	// 暂时只做一个基本的结构体绑定演示
	var ghPayload struct {
		Ref        string `json:"ref"`
		HeadCommit struct {
			ID      string `json:"id"`
			Message string `json:"message"`
			Author  struct {
				Name string `json:"name"`
			} `json:"author"`
		} `json:"head_commit"`
		Repository struct {
			HTMLURL string `json:"html_url"`
		} `json:"repository"`
	}

	if err := c.ShouldBindJSON(&ghPayload); err != nil {
		h.HandleBadRequest(c, "Invalid GitHub payload")
		return
	}

	payload := dto.WebhookPayload{
		RepoURL:   ghPayload.Repository.HTMLURL, // 注意：这里需要与配置的 RepoURL 匹配
		Ref:       ghPayload.Ref,
		CommitSHA: ghPayload.HeadCommit.ID,
		CommitMsg: ghPayload.HeadCommit.Message,
		Author:    ghPayload.HeadCommit.Author.Name,
		Status:    "success", // Push 事件默认算 update, 如果是 WorkflowRun 事件则有 status
	}

	if err := h.devopsService.HandleWebhook(c.Request.Context(), payload); err != nil {
		h.HandleInternalError(c, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

// TriggerDeployment 触发部署
func (h *DevOpsHandler) TriggerDeployment(c *gin.Context) {
	if err := h.devopsService.TriggerDeployment(c.Request.Context(), "backend"); err != nil {
		h.HandleInternalError(c, err.Error())
		return
	}
	h.HandleSuccess(c, nil)
}
