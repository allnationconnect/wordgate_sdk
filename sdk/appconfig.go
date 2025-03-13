package sdk

import (
	"fmt"
)

// AppConfigSyncResponse 应用配置同步响应
type AppConfigSyncResponse struct {
	// Success 表示同步是否成功
	Success bool `json:"success"`
	// Name 应用名称
	Name string `json:"name"`
	// Message 操作消息
	Message string `json:"message"`
}

// SyncAppConfig 同步应用配置
//
// 将配置中定义的应用配置同步到服务器
// 返回同步结果和可能的错误
func (c *Client) SyncAppConfig() (*AppConfigSyncResponse, error) {
	// 创建包含App基本信息和配置信息的请求体
	requestData := struct {
		// 应用基本信息
		Name        string    `json:"name"`
		Description string    `json:"description"`
		Currency    string    `json:"currency"`
		Config      AppConfig `json:"config"` // 将配置作为一个嵌套字段
	}{
		// 应用基本信息
		Name:        c.Config.App.Name,
		Description: c.Config.App.Description,
		Currency:    c.Config.App.Currency,
		// 应用配置 - 直接使用嵌套结构
		Config: c.Config.Config,
	}

	// 创建响应结果
	var response AppConfigSyncResponse

	// 发送请求并解析响应
	if err := c.apiRequestJSON("POST", "/app/config", requestData, &response); err != nil {
		return nil, fmt.Errorf("同步应用配置失败: %w", err)
	}

	// 设置成功标志
	response.Success = true

	return &response, nil
}
