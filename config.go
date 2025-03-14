package wordgate_sdk

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// WordgateConfig 表示完整的Wordgate配置
// 包含应用信息、API认证信息、产品配置、会员配置等
type WordgateConfig struct {
	// BaseURL API服务的基础URL
	BaseURL string `yaml:"base_url" json:"base_url"`
	// AppCode 应用代码，用于API认证
	AppCode string `yaml:"app_code" json:"app_code"`
	// AppSecret 应用密钥，用于API认证
	AppSecret string `yaml:"app_secret" json:"app_secret"`
	// EnablePayment 是否启用支付功能
	EnablePayment bool `yaml:"enable_payment" json:"enable_payment"`
	// Products 产品相关配置
	Products ProductConfig `yaml:"products" json:"products"`
	// App 应用基本信息
	App AppInfo `yaml:"app" json:"app"`
	// Config 应用配置
	Config AppConfig `yaml:"config" json:"config"`
	// Membership 会员系统配置
	Membership MembershipConfig `yaml:"membership" json:"membership"`
}

// AppInfo 应用基本信息
type AppInfo struct {
	// Name 应用名称
	Name string `yaml:"name" json:"name"`
	// Description 应用描述
	Description string `yaml:"description" json:"description"`
	// Currency 结算货币代码(如CNY、USD等)
	Currency string `yaml:"currency" json:"currency"`
}

// ProductConfig 产品配置
type ProductConfig struct {
	// Files 文件匹配模式列表，用于从文件中提取产品信息
	// 路径相对于配置文件所在目录
	Files []string `yaml:"files" json:"files"`
	// Items 直接在配置文件中定义的产品列表
	Items []Product `yaml:"items" json:"items"`
}

// Product 产品定义
type Product struct {
	// Code 产品代码，唯一标识一个产品
	Code string `yaml:"code" json:"code"`
	// Name 产品名称
	Name string `yaml:"name" json:"name"`
	// Price 产品价格(单位:分)
	Price int `yaml:"price" json:"price"`
}

// AppConfig 应用配置
type AppConfig struct {
	// SMTP 邮件配置
	SMTP SMTPConfig `yaml:"smtp" json:"smtp"`
	// SMS 短信配置
	SMS SMSConfig `yaml:"sms" json:"sms"`
	// Security 安全配置
	Security SecurityConfig `yaml:"security" json:"security"`
	// Payment 支付配置
	Payment PaymentConfig `yaml:"payment" json:"payment"`
	// Site 网站配置
	Site SiteConfig `yaml:"site" json:"site"`
}

// SiteConfig 网站配置
type SiteConfig struct {
	// BaseURL 网站基础URL
	BaseURL string `yaml:"base_url" json:"base_url"`
	// PayPagePath 支付页面路径
	PayPagePath string `yaml:"pay_page_path" json:"pay_page_path"`
	// PayResultPagePath 支付结果页面路径
	PayResultPagePath string `yaml:"pay_result_page_path" json:"pay_result_page_path"`
}

// GeneratePaymentURL 生成支付页面URL
func (c *SiteConfig) GeneratePaymentURL(orderNo string) string {
	// 如果路径已经是完整URL，则直接使用
	path := c.PayPagePath
	if path == "" {
		path = "/pay"
	}

	if strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://") {
		return fmt.Sprintf("%s?order_no=%s", path, orderNo)
	}

	// 规范化baseURL和path的连接
	base := c.BaseURL
	if !strings.HasSuffix(base, "/") {
		base += "/"
	}

	if strings.HasPrefix(path, "/") {
		path = path[1:]
	}

	url := base + path
	return fmt.Sprintf("%s?order_no=%s", url, orderNo)
}

// GeneratePayResultURL 生成支付结果页面URL
func (c *SiteConfig) GeneratePayResultURL(orderNo string, queryParams map[string]string) string {
	path := c.PayResultPagePath
	if path == "" {
		path = "/pay-result"
	}

	baseURL := ""
	if strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://") {
		baseURL = path
	} else {
		// 规范化baseURL和path的连接
		base := c.BaseURL
		if !strings.HasSuffix(base, "/") {
			base += "/"
		}

		if strings.HasPrefix(path, "/") {
			path = path[1:]
		}

		baseURL = base + path
	}

	// 构建带查询参数的URL
	result := fmt.Sprintf("%s?order_no=%s", baseURL, orderNo)

	// 添加其他查询参数
	if len(queryParams) > 0 {
		u, err := url.Parse(result)
		if err == nil {
			q := u.Query()
			for k, v := range queryParams {
				q.Set(k, v)
			}
			u.RawQuery = q.Encode()
			result = u.String()
		}
	}

	return result
}

// SMTPConfig 邮件配置
type SMTPConfig struct {
	// Host SMTP服务器地址
	Host string `yaml:"host" json:"host"`
	// Port SMTP服务器端口
	Port int `yaml:"port" json:"port"`
	// Username SMTP用户名
	Username string `yaml:"username" json:"username"`
	// Password SMTP密码
	Password string `yaml:"password" json:"password"`
	// FromEmail 发件人邮箱
	FromEmail string `yaml:"from_email" json:"from_email"`
	// FromName 发件人名称
	FromName string `yaml:"from_name" json:"from_name"`
	// ReplyToEmail 回复邮箱
	ReplyToEmail string `yaml:"reply_to_email" json:"reply_to_email"`
}

// SMSConfig 短信配置
type SMSConfig struct {
	// Provider 短信服务提供商
	Provider string `yaml:"provider" json:"provider"`
	// APIKey API密钥
	APIKey string `yaml:"api_key" json:"api_key"`
	// APISecret API密钥对应的密钥
	APISecret string `yaml:"api_secret" json:"api_secret"`
	// SignName 短信签名
	SignName string `yaml:"sign_name" json:"sign_name"`
}

// SecurityConfig 安全配置
type SecurityConfig struct {
	// SessionExpire 会话过期时间(秒)
	SessionExpire int `yaml:"session_expire" json:"session_expire"`
	// CodeExpire 验证码过期时间(秒)
	CodeExpire int `yaml:"code_expire" json:"code_expire"`
}

// PaymentConfig 支付配置
type PaymentConfig struct {
	// GatewayMode 网关模式配置
	GatewayMode GatewayModeConfig `yaml:"gateway_mode" json:"gateway_mode"`
	// Antom Antom支付配置
	Antom AntomConfig `yaml:"antom" json:"antom"`
}

// GatewayModeConfig 网关模式配置
type GatewayModeConfig struct {
	// Enabled 是否启用网关模式
	Enabled bool `yaml:"enabled" json:"enabled"`
	// NotifyURL 通知URL
	NotifyURL string `yaml:"notify_url" json:"notify_url"`
	// RedirectURL 重定向URL
	RedirectURL string `yaml:"redirect_url" json:"redirect_url"`
}

// AntomConfig Antom支付配置
type AntomConfig struct {
	// Enabled 是否启用Antom支付
	Enabled bool `yaml:"enabled" json:"enabled"`
	// ClientID Antom客户端ID
	ClientID string `yaml:"client_id" json:"client_id"`
	// AntomPublicKey Antom公钥
	AntomPublicKey string `yaml:"antom_public_key" json:"antom_public_key"`
	// YourPublicKey 您的公钥
	YourPublicKey string `yaml:"your_public_key" json:"your_public_key"`
	// YourPrivateKey 您的私钥
	YourPrivateKey string `yaml:"your_private_key" json:"your_private_key"`
	// IsSandbox 是否使用沙箱环境
	IsSandbox bool `yaml:"is_sandbox" json:"is_sandbox"`
	// Domain 域名
	Domain string `yaml:"domain" json:"domain"`
}

// MembershipConfig 会员系统配置
type MembershipConfig struct {
	// Tiers 会员等级列表
	Tiers []MembershipTier `yaml:"tiers" json:"tiers"`
}

// MembershipTier 会员等级
type MembershipTier struct {
	// Code 会员等级代码
	Code string `yaml:"code" json:"code"`
	// Name 会员等级名称
	Name string `yaml:"name" json:"name"`
	// Level 等级值，用于排序
	Level int `yaml:"level" json:"level"`
	// IsDefault 是否默认等级
	IsDefault bool `yaml:"is_default" json:"is_default"`
	// Prices 会员价格配置列表
	Prices []MembershipPrice `yaml:"prices" json:"prices"`
}

// MembershipPrice 会员价格
type MembershipPrice struct {
	// PeriodType 周期类型，如 month、year 等
	PeriodType string `yaml:"period_type" json:"period_type"`
	// Price 价格(单位:分)
	Price int64 `yaml:"price" json:"price"`
	// OriginalPrice 原价(单位:分)，用于显示折扣信息
	OriginalPrice int64 `yaml:"original_price" json:"original_price"`
}

// LoadConfig 从文件加载Wordgate配置
//
// filePath 参数指定配置文件的路径
// 返回加载的配置和可能的错误
func LoadConfig(filePath string) (*WordgateConfig, error) {
	// 读取配置文件内容
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %w", err)
	}

	// 根据文件扩展名选择解析方法
	ext := filepath.Ext(filePath)
	switch strings.ToLower(ext) {
	case ".yaml", ".yml":
		return loadFromYAML(data)
	case ".json":
		return loadFromJSON(data)
	default:
		return nil, fmt.Errorf("不支持的配置文件格式: %s", ext)
	}
}

// loadFromYAML 从YAML数据加载配置
//
// data 参数包含YAML格式的配置数据
// 返回解析的配置和可能的错误
func loadFromYAML(data []byte) (*WordgateConfig, error) {
	// 创建顶级配置结构
	var topLevelConfig struct {
		Wordgate WordgateConfig `yaml:"wordgate"`
	}

	// 解析YAML数据
	err := yaml.Unmarshal(data, &topLevelConfig)
	if err != nil {
		return nil, fmt.Errorf("解析YAML配置失败: %w", err)
	}

	// 检查是否包含wordgate配置
	if isEmpty(topLevelConfig.Wordgate) {
		return nil, fmt.Errorf("配置文件中缺少wordgate配置或配置不完整")
	}

	return &topLevelConfig.Wordgate, nil
}

// loadFromJSON 从JSON数据加载配置
//
// data 参数包含JSON格式的配置数据
// 返回解析的配置和可能的错误
func loadFromJSON(data []byte) (*WordgateConfig, error) {
	// 创建顶级配置结构
	var topLevelConfig struct {
		Wordgate WordgateConfig `json:"wordgate"`
	}

	// 解析JSON数据
	err := json.Unmarshal(data, &topLevelConfig)
	if err != nil {
		return nil, fmt.Errorf("解析JSON配置失败: %w", err)
	}

	// 检查是否包含wordgate配置
	if isEmpty(topLevelConfig.Wordgate) {
		return nil, fmt.Errorf("配置文件中缺少wordgate配置或配置不完整")
	}

	return &topLevelConfig.Wordgate, nil
}

// isEmpty 判断WordgateConfig是否为空
//
// 当所有关键字段都为空时，认为配置为空
func isEmpty(config WordgateConfig) bool {
	return config.BaseURL == "" && config.AppCode == "" && config.AppSecret == ""
}

// ValidateConfig 验证配置是否有效
//
// 检查必要的配置字段是否存在，返回可能的错误
func ValidateConfig(config *WordgateConfig) error {
	if config.BaseURL == "" {
		return fmt.Errorf("缺少必要配置: BaseURL")
	}
	if config.AppCode == "" {
		return fmt.Errorf("缺少必要配置: AppCode")
	}
	if config.AppSecret == "" {
		return fmt.Errorf("缺少必要配置: AppSecret")
	}
	return nil
}
