# WordGate Hugo同步工具

这个工具用于将Hugo站点的产品数据和应用配置同步到WordGate API服务器。

## 功能特点

- 自动从Hugo配置文件中读取WordGate配置信息
- 从Hugo内容文件中提取产品信息
- 将产品数据同步到WordGate API服务器
- 将应用配置同步到WordGate API服务器
- 支持干运行模式，可以预览要同步的数据而不实际发送

## 安装

### 从源代码安装

1. 克隆仓库：

   ```bash
   git clone https://github.com/allnationconnect/wordgate.git
   cd wordgate
   ```

2. 安装命令行工具：

   ```bash
   go install ./cmd/wordgate
   ```

   这将在您的`$GOPATH/bin`目录下安装`wordgate`命令。确保此目录已添加到您的系统PATH中。

### 使用预编译二进制文件

1. 从[releases页面](https://github.com/allnationconnect/wordgate_sdk/releases)下载适合您系统的最新版本。

2. 解压缩下载的文件：

   ```bash
   tar -xzf wordgate_[版本]_[系统]_[架构].tar.gz
   ```

3. 将解压出的`wordgate`二进制文件移动到您的PATH目录中：

   ```bash
   sudo mv wordgate /usr/local/bin/
   ```

## 配置

在配置文件中添加以下WordGate配置：

```yaml
wordgate:
  base_url: "http://127.0.0.1:8080/api"  # WordGate API服务器地址
  appCode: "site_owner"                 # 应用代码
  app_secret: "your-app-secret"         # 应用密钥
  enable_payment: true                   # 是否启用支付功能

  # 产品配置
  products:
    files:
      - "content/courses/*.md"    # 使用glob模式匹配内容文件
    items:
      - code: "COURSE001"           # 产品代码
        name: "示例课程1"          # 产品名称
        price: 9900              # 价格(分)  
  
  # 应用基本信息
  app:
    name: "我的应用"               # 应用名称
    description: "应用描述"        # 应用描述
    currency: "CNY"              # 结算货币

  # 会员系统配置
  membership:
    tiers:
      - code: "FREE"             # 会员等级代码
        name: "免费会员"          # 会员等级名称
        level: 0                 # 等级值(用于排序)
        is_default: true         # 是否默认等级

  # 应用配置
  config:
    # 邮件配置
    smtp:
      host: "smtp.example.com"
      port: 587
      username: "noreply@example.com"
      password: "your-password"
      from_email: "noreply@example.com"
      from_name: "您的网站名称"
      reply_to_email: "support@example.com"
      
    # 短信配置
    sms:
      provider: "aliyun"
      api_key: "your-api-key"
      api_secret: "your-api-secret"
      sign_name: "您的网站名称"
      
    # 安全配置
    security:
      session_expire: 86400  # 会话过期时间（秒）
      code_expire: 600       # 验证码过期时间（秒）
      
    # 支付配置
    payment:
      antom:
        enabled: true
        client_id: "your-client-id"
        antom_public_key: "antom-public-key"
        your_public_key: "your-public-key"
        your_private_key: "your-private-key"
        is_sandbox: true
        domain: "example.com"
```

### 关于产品代码

当配置`code: "slug"`时，工具会自动使用文件名(不含扩展名)作为产品代码。例如：

- 文件 `content/courses/introduction-to-python.md` 将生成产品代码 `introduction-to-python`
- 文件 `content/workshops/summer-camp.md` 将生成产品代码 `summer-camp`

这是一个特殊处理，利用了Hugo的命名规则，方便您快速引用产品。

## 关于异常处理

该工具对前置元数据(Front Matter)处理进行了特殊优化：

1. **多种分隔符支持**：同时支持YAML(`---`)和TOML(`+++`)格式的前置元数据
2. **容错处理**：即使缺少分隔符或格式不规范，也能尝试提取有用信息
3. **自动补全**：对于缺失必要字段的内容文件：
   - 如果`code`配置为"slug"但找不到对应字段，将自动使用文件名
   - 缺失名称时，将使用代码作为名称
   - 缺失价格时，默认为0
   - 缺失有效期时，默认为365天

这些改进使得工具能够更加健壮地处理各种不规范的内容文件，减少同步失败的情况。

## 命令行选项

### 基本使用

```bash
wordgate
```

此命令将使用当前目录中的配置文件（按优先顺序查找：`config.yaml`、`config.yml`、`hugo.yaml`、`hugo.yml`）。

### 指定配置文件

```bash
wordgate -config=path/to/config.yaml
```

指定配置文件的路径。

### 干运行模式

```bash
wordgate -dry-run
```

只显示要同步的数据，不发送实际请求到服务器。这对于验证配置和内容非常有用。

### 生成示例配置

```bash
wordgate -print-demo
```

打印一个完整的示例配置文件和Markdown内容文件示例。您可以将输出重定向到文件：

```bash
wordgate -print-demo > config.yaml
```

## 示例用法

### 首次使用

1. 生成示例配置文件：

   ```bash
   wordgate -print-demo > config.yaml
   ```

2. 编辑配置文件，填入您的应用信息和API凭据。

3. 运行干运行模式，验证配置：

   ```bash
   wordgate -dry-run
   ```

4. 实际执行同步：

   ```bash
   wordgate
   ```

## 帮助信息

```bash
./wordgate -help
```

## 产品数据格式

从Hugo内容文件中提取的产品数据将以以下格式同步到API服务器：

```json
{
  "products": [
    {
      "code": "product-code",
      "name": "产品名称",
      "price": 9900
    },
    ...
  ]
}
```

## 常见问题

### 找不到产品数据？

- 确保您的内容文件包含正确的前置元数据（Front Matter）
- 检查`content_types`和`sections`配置是否正确
- 使用`-dry-run`参数查看工具找到的产品

### API请求失败？

- 检查`base_url`、`app_code`和`app_secret`是否正确
- 确保API服务器正常运行并可访问
- 查看服务器日志以获取更多信息

## 贡献

欢迎通过Pull Request贡献代码或提交Issue报告问题。 