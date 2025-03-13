package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/allnationconnect/wordgate_sdk/sdk"
)

// 示例配置文件内容
const demoConfig = `# Wordgate配置示例
wordgate:
  base_url: "https://api.wordgate.example.com" # API服务地址
  appCode: "my-app-code"                       # 应用代码
  app_secret: "my-app-secret"                  # 应用密钥
  enable_payment: true                         # 是否启用支付功能

  # 应用基本信息
  app:
    name: "我的应用"                     # 应用名称
    description: "这是我的应用描述"      # 应用描述
    currency: "CNY"                    # 结算货币(CNY,USD等)

  # 产品配置
  products:
    # 从文件中提取产品信息(相对于配置文件所在目录)
    files:
      - "content/courses/*.md"         # 匹配内容文件的模式
    
    # 直接在配置中定义产品
    items:
      - code: "COURSE001"                # 产品代码
        name: "示例课程1"               # 产品名称
        price: 9900                    # 价格(分)
      - code: "COURSE002"
        name: "示例课程2"
        price: 19900

  # 会员系统配置
  membership:
    tiers:
      - code: "FREE"                   # 会员等级代码
        name: "免费会员"                # 会员等级名称
        level: 0                       # 等级值(用于排序)
        is_default: true               # 是否默认等级
      - code: "PRO"
        name: "专业会员"
        level: 1
        is_default: false
        prices:                        # 会员价格配置
          - period_type: "month"       # 周期类型(month,year等)
            price: 9900                # 价格(分)
            original_price: 12900      # 原价(分)
          - period_type: "year"
            price: 99900
            original_price: 129900

  # 应用配置
  config:
	  # 网站配置
    site:
      base_url: "https://example.com"   # 网站基础URL
      pay_page_path: "/pay"             # 支付页面路径
      pay_result_page_path: "/pay-result" # 支付结果页面路径
    # 邮件配置
    smtp:
      host: "smtp.example.com"         # SMTP服务器地址
      port: 587                        # SMTP服务器端口
      username: "noreply@example.com"  # SMTP用户名
      password: "smtp-password"        # SMTP密码
      from_email: "noreply@example.com" # 发件人邮箱
      from_name: "我的应用"              # 发件人名称
      reply_to_email: "support@example.com" # 回复邮箱
    
    # 短信配置
    sms:
      provider: "aliyun"               # 短信服务提供商
      api_key: "sms-api-key"           # API密钥
      api_secret: "sms-api-secret"     # API密钥对应的密钥
      sign_name: "我的应用"              # 短信签名
    
    # 安全配置
    security:
      session_expire: 86400            # 会话过期时间(秒)
      code_expire: 300                 # 验证码过期时间(秒)
    
    # 支付配置
    payment:
      antom:
        enabled: true                  # 是否启用
        client_id: "antom-client-id"   # 客户端ID
        antom_public_key: "antom-public-key" # Antom公钥
        your_public_key: "your-public-key"   # 您的公钥
        your_private_key: "your-private-key" # 您的私钥
        is_sandbox: true               # 是否使用沙箱环境
        domain: "example.com"          # 域名
`

// 示例Markdown内容文件
const demoMarkdown = `---
title: "示例课程"
date: 2023-01-01
description: "这是一个示例课程的描述"
product:
  code: "COURSE003"
  name: "Markdown中定义的课程"
  price: 29900
---

# 示例课程

这是示例课程的内容。
`

func main() {
	// 命令行参数
	syncConfigPath := flag.String("sync-config", "", "指定要同步的配置文件路径")
	dryRun := flag.Bool("dry-run", false, "只输出要同步的数据，不发送实际请求到服务器 (需要和-sync-config一起使用)")
	printDemo := flag.Bool("print-demo", false, "打印样板配置文件")
	help := flag.Bool("help", false, "显示帮助信息")

	// 解析命令行参数
	flag.Parse()

	// 如果用户选择打印帮助信息，或者没有提供任何参数，显示帮助
	if *help || (flag.NFlag() == 0 && len(flag.Args()) == 0) {
		printHelp()
		return
	}

	// 如果用户选择打印样板配置文件，则打印并退出
	if *printDemo {
		fmt.Println(demoConfig)
		fmt.Println("\n# 示例Markdown内容文件 (保存为content/courses/course.md)")
		fmt.Println(demoMarkdown)
		return
	}

	// 验证参数组合
	if *dryRun && *syncConfigPath == "" {
		fmt.Println("错误: -dry-run 参数必须与 -sync-config 参数一起使用")
		fmt.Println("\n使用 -help 查看帮助信息")
		os.Exit(1)
	}

	// 处理同步逻辑
	if *syncConfigPath != "" {
		// 获取配置文件的绝对路径和所在目录
		absConfigPath, err := filepath.Abs(*syncConfigPath)
		if err != nil {
			log.Fatalf("无法获取配置文件的绝对路径: %v", err)
		}
		configDir := filepath.Dir(absConfigPath)
		fmt.Printf("配置文件目录: %s\n", configDir)

		// 使用SDK加载配置
		config, err := sdk.LoadConfig(absConfigPath)
		if err != nil {
			log.Fatalf("加载配置失败: %v", err)
		}

		// 创建SDK客户端
		client := sdk.NewClient(config, configDir)

		// 根据模式执行操作
		if *dryRun {
			// 干运行模式，只显示数据不发送请求
			fmt.Println("\n== 干运行模式: 不发送API请求 ==")

			result, err := client.DryRun()
			if err != nil {
				log.Fatalf("干运行失败: %v", err)
			}

			// 打印应用配置摘要
			fmt.Println("\n应用配置摘要:")
			fmt.Printf("  名称: %s\n", result.AppConfig.Name)
			fmt.Printf("  描述: %s\n", result.AppConfig.Description)
			fmt.Printf("  货币: %s\n", result.AppConfig.Currency)

			// 打印会员等级信息
			fmt.Println("\n会员等级配置:")
			if len(result.Memberships) == 0 {
				fmt.Println("❌ 没有找到会员等级配置，请检查配置文件")
			} else {
				fmt.Printf("找到 %d 个会员等级配置\n", len(result.Memberships))
				for i, tier := range result.Memberships {
					fmt.Printf("%d. %s (%s) 等级:%d 默认:%v 价格:%d\n",
						i+1, tier.Name, tier.Code, tier.Level, tier.IsDefault,
						len(tier.Prices))
				}
			}

			// 打印产品信息
			fmt.Println("\n产品信息:")
			if len(result.Products) == 0 {
				fmt.Println("没有找到任何产品")
			} else {
				fmt.Printf("共找到 %d 个产品\n", len(result.Products))
				for i, product := range result.Products {
					fmt.Printf("  %d. %s (%d) - %s\n",
						i+1, product.Name, product.Price, product.Code)
				}
			}
		} else {
			// 实际同步模式
			fmt.Println("\n== 开始同步流程 ==")

			// 使用SDK执行完整同步
			result, err := client.SyncAll()
			if err != nil {
				log.Fatalf("同步失败: %v", err)
			}

			// 打印同步结果
			printSyncResult(result)
		}
	} else {
		// 如果没有指定同步配置，也没有其他有效的参数组合，显示帮助
		printHelp()
	}
}

// printSyncResult 打印同步结果
func printSyncResult(result *sdk.SyncAllResponse) {
	fmt.Println("\n== 同步结果 ==")
	if !result.Success {
		fmt.Printf("❌ 同步失败: %s\n", result.ErrorMessage)
		return
	}
	fmt.Println("✅ 同步成功!")

	// 1. 应用配置同步结果
	fmt.Println("\n1. 应用配置:")
	fmt.Printf("  名称: %s\n", result.AppConfig.Name)
	fmt.Printf("  状态: 成功\n")
	if result.AppConfig.Message != "" {
		fmt.Printf("  消息: %s\n", result.AppConfig.Message)
	}

	// 2. 会员等级同步结果
	fmt.Println("\n2. 会员等级:")
	fmt.Printf("  状态: %s\n", boolToSuccessString(result.Memberships.Success))
	fmt.Printf("  总数: %d\n", result.Memberships.Total)
	fmt.Printf("  创建: %d\n", result.Memberships.Created)
	fmt.Printf("  更新: %d\n", result.Memberships.Updated)
	fmt.Printf("  未变: %d\n", result.Memberships.Unchanged)
	fmt.Printf("  失败: %d\n", result.Memberships.Failed)

	if len(result.Memberships.Errors) > 0 {
		fmt.Println("  错误:")
		for _, err := range result.Memberships.Errors {
			fmt.Printf("    - [%s] %s: %s\n", err.Code, err.TierCode, err.Message)
		}
	}

	// 3. 产品同步结果
	fmt.Println("\n3. 产品:")
	fmt.Printf("  状态: %s\n", boolToSuccessString(result.Products.Success))
	fmt.Printf("  总数: %d\n", result.Products.Total)
	fmt.Printf("  创建: %d\n", result.Products.Created)
	fmt.Printf("  更新: %d\n", result.Products.Updated)
	fmt.Printf("  未变: %d\n", result.Products.Unchanged)
	fmt.Printf("  失败: %d\n", result.Products.Failed)

	if len(result.Products.Errors) > 0 {
		fmt.Println("  错误:")
		for _, err := range result.Products.Errors {
			fmt.Printf("    - [%s] %s: %s\n", err.Code, err.ProductCode, err.Message)
		}
	}

	fmt.Println("\n== 同步流程完成 ==")
	fmt.Println("✓ 应用配置同步成功")
	fmt.Println("✓ 会员等级同步成功")
	fmt.Println("✓ 产品同步成功")
}

// boolToSuccessString 将布尔值转换为成功/失败字符串
func boolToSuccessString(value bool) string {
	if value {
		return "成功"
	}
	return "失败"
}

// printHelp 打印帮助信息
func printHelp() {
	fmt.Println("Wordgate - 产品和应用配置同步工具")
	fmt.Println("\n用法:")
	fmt.Printf("  %s [参数]\n\n", os.Args[0])
	fmt.Println("参数:")
	fmt.Println("  -sync-config string  指定要同步的配置文件路径")
	fmt.Println("  -dry-run            只显示要同步的数据，不发送实际请求到服务器 (需要与-sync-config一起使用)")
	fmt.Println("  -print-demo         打印示例配置文件，可重定向到文件")
	fmt.Println("  -help               显示此帮助信息")
	fmt.Println("\n示例:")
	fmt.Println("  基本使用:             wordgate -sync-config=config.yaml")
	fmt.Println("  干运行模式:           wordgate -sync-config=config.yaml -dry-run")
	fmt.Println("  生成示例配置:         wordgate -print-demo > config.yaml")
	fmt.Println("\n注意:")
	fmt.Println("  -sync-config 参数指定要同步的配置文件路径")
	fmt.Println("  -dry-run 参数必须与 -sync-config 参数一起使用")
}
