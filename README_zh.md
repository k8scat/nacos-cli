# nacos-cli

一个用于管理 Nacos 服务器配置、服务和命名空间的命令行工具。

## 安装

```bash
go install github.com/k8scat/nacos-cli@latest
```

或者从源代码构建：

```bash
git clone https://github.com/k8scat/nacos-cli.git
cd nacos-cli
go build
```

## 使用方法

### 配置

您可以使用命令行参数提供配置信息，或在 `$HOME/.nacos-cli.yaml` 创建配置文件。

配置文件示例：

```yaml
server: http://localhost:8848
username: nacos
password: nacos
namespace: public
```

全局参数：

```
--config string       配置文件路径（默认为 $HOME/.nacos-cli.yaml）
--namespace string    Nacos 命名空间 ID
--password string     Nacos 服务器密码
--server string       Nacos 服务器地址（默认为 "http://localhost:8848"）
--username string     Nacos 服务器用户名
```

### 示例

#### 配置管理

获取配置：

```bash
# 获取配置并打印到控制台
nacos-cli config get --data-id=example --group=DEFAULT_GROUP

# 获取配置并保存到文件
nacos-cli config get --data-id=example --group=DEFAULT_GROUP --output=config.yaml
```

发布配置：

```bash
# 通过命令行发布配置
nacos-cli config publish --data-id=example --group=DEFAULT_GROUP --content="hello world"

# 从文件发布配置
nacos-cli config publish --data-id=example --group=DEFAULT_GROUP --file=config.yaml

# 从标准输入发布配置
cat config.yaml | nacos-cli config publish --data-id=example --group=DEFAULT_GROUP
```

删除配置：

```bash
nacos-cli config delete --data-id=example --group=DEFAULT_GROUP
```

监听配置变更：

```bash
nacos-cli config listen --data-id=example --group=DEFAULT_GROUP --content="current content"
```

获取配置历史记录：

```bash
nacos-cli config history --data-id=example --group=DEFAULT_GROUP --pretty
```

获取配置历史详情：

```bash
nacos-cli config history-detail --data-id=example --group=DEFAULT_GROUP --nid=123 --pretty
```

获取先前配置版本：

```bash
nacos-cli config previous --data-id=example --group=DEFAULT_GROUP --pretty
```

#### 服务管理

获取服务信息：

```bash
# 获取服务信息
nacos-cli service get --name=example-service

# 获取服务信息并美化输出
nacos-cli service get --name=example-service --pretty
```

注册实例：

```bash
# 注册简单实例
nacos-cli service register --name=example-service --ip=192.168.1.1 --port=8080

# 注册带元数据的实例
nacos-cli service register --name=example-service --ip=192.168.1.1 --port=8080 --metadata='{"version":"1.0.0"}'

# 注册带集群的实例
nacos-cli service register --name=example-service --ip=192.168.1.1 --port=8080 --cluster=DEFAULT
```

注销实例：

```bash
nacos-cli service deregister --name=example-service --ip=192.168.1.1 --port=8080
```

修改实例：

```bash
nacos-cli service modify-instance --name=example-service --ip=192.168.1.1 --port=8080 --weight=2.0 --metadata='{"version":"2.0.0"}'
```

列出实例：

```bash
nacos-cli service list-instances --name=example-service --pretty
```

获取实例详情：

```bash
nacos-cli service get-instance --name=example-service --ip=192.168.1.1 --port=8080 --pretty
```

发送实例心跳：

```bash
nacos-cli service beat --name=example-service --ip=192.168.1.1 --port=8080
```

创建服务：

```bash
nacos-cli service create --name=example-service --group=DEFAULT_GROUP --protect=0.7 --metadata='{"description":"example service"}'
```

删除服务：

```bash
nacos-cli service delete --name=example-service --group=DEFAULT_GROUP
```

更新服务：

```bash
nacos-cli service update --name=example-service --group=DEFAULT_GROUP --protect=0.8 --metadata='{"description":"updated example service"}'
```

列出服务：

```bash
nacos-cli service list --page=1 --size=20 --pretty
```

获取系统开关：

```bash
nacos-cli service get-switches --pretty
```

更新系统开关：

```bash
nacos-cli service update-switch --entry=distro.enableAll --value=true
```

获取系统指标：

```bash
nacos-cli service metrics --pretty
```

获取服务器列表：

```bash
nacos-cli service servers --pretty
```

获取集群领导者：

```bash
nacos-cli service leader --pretty
```

更新实例健康状态：

```bash
nacos-cli service update-health --name=example-service --ip=192.168.1.1 --port=8080 --healthy=true
```

#### 命名空间管理

列出命名空间：

```bash
# 列出所有命名空间
nacos-cli namespace list

# 列出所有命名空间并美化输出
nacos-cli namespace list --pretty
```

创建命名空间：

```bash
nacos-cli namespace create --id=dev --name="开发环境" --desc="用于开发"
```

修改命名空间：

```bash
nacos-cli namespace modify --id=dev --name="开发" --desc="开发环境"
```

删除命名空间：

```bash
nacos-cli namespace delete --id=dev
```

## 许可证

[MIT](https://github.com/k8scat/nacos-cli/blob/main/LICENSE)
