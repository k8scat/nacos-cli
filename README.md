# nacos-cli

A command-line tool for managing Nacos server configurations, services, and namespaces.

## Installation

```bash
go install github.com/k8scat/nacos-cli@latest
```

Or build from source:

```bash
git clone https://github.com/k8scat/nacos-cli.git
cd nacos-cli
go build
```

## Usage

### Configuration

You can use command-line flags to provide configuration information or create a config file at `$HOME/.nacos-cli.yaml`.

Config file example:

```yaml
server: http://localhost:8848
username: nacos
password: nacos
namespace: public
```

Global flags:

```
--config string       Config file (default is $HOME/.nacos-cli.yaml)
--namespace string    Nacos namespace ID
--password string     Nacos server password
--server string       Nacos server address (default "http://localhost:8848")
--username string     Nacos server username
```

### Examples

#### Configuration Management

Get configuration:

```bash
# Get configuration and print to console
nacos-cli config get --data-id=example --group=DEFAULT_GROUP

# Get configuration and save to file
nacos-cli config get --data-id=example --group=DEFAULT_GROUP --output=config.yaml
```

Publish configuration:

```bash
# Publish configuration from command line
nacos-cli config publish --data-id=example --group=DEFAULT_GROUP --content="hello world"

# Publish configuration from file
nacos-cli config publish --data-id=example --group=DEFAULT_GROUP --file=config.yaml

# Publish configuration from stdin
cat config.yaml | nacos-cli config publish --data-id=example --group=DEFAULT_GROUP
```

Delete configuration:

```bash
nacos-cli config delete --data-id=example --group=DEFAULT_GROUP
```

Listen for configuration changes:

```bash
nacos-cli config listen --data-id=example --group=DEFAULT_GROUP --content="current content"
```

Get configuration history:

```bash
nacos-cli config history --data-id=example --group=DEFAULT_GROUP --pretty
```

Get configuration history detail:

```bash
nacos-cli config history-detail --data-id=example --group=DEFAULT_GROUP --nid=123 --pretty
```

Get previous configuration version:

```bash
nacos-cli config previous --data-id=example --group=DEFAULT_GROUP --pretty
```

#### Service Management

Get service information:

```bash
# Get service information
nacos-cli service get --name=example-service

# Get service information with pretty output
nacos-cli service get --name=example-service --pretty
```

Register instance:

```bash
# Register a simple instance
nacos-cli service register --name=example-service --ip=192.168.1.1 --port=8080

# Register with metadata
nacos-cli service register --name=example-service --ip=192.168.1.1 --port=8080 --metadata='{"version":"1.0.0"}'

# Register with cluster
nacos-cli service register --name=example-service --ip=192.168.1.1 --port=8080 --cluster=DEFAULT
```

Deregister instance:

```bash
nacos-cli service deregister --name=example-service --ip=192.168.1.1 --port=8080
```

Modify instance:

```bash
nacos-cli service modify-instance --name=example-service --ip=192.168.1.1 --port=8080 --weight=2.0 --metadata='{"version":"2.0.0"}'
```

List instances:

```bash
nacos-cli service list-instances --name=example-service --pretty
```

Get instance details:

```bash
nacos-cli service get-instance --name=example-service --ip=192.168.1.1 --port=8080 --pretty
```

Send instance heartbeat:

```bash
nacos-cli service beat --name=example-service --ip=192.168.1.1 --port=8080
```

Create service:

```bash
nacos-cli service create --name=example-service --group=DEFAULT_GROUP --protect=0.7 --metadata='{"description":"example service"}'
```

Delete service:

```bash
nacos-cli service delete --name=example-service --group=DEFAULT_GROUP
```

Update service:

```bash
nacos-cli service update --name=example-service --group=DEFAULT_GROUP --protect=0.8 --metadata='{"description":"updated example service"}'
```

List services:

```bash
nacos-cli service list --page=1 --size=20 --pretty
```

Get system switches:

```bash
nacos-cli service get-switches --pretty
```

Update system switch:

```bash
nacos-cli service update-switch --entry=distro.enableAll --value=true
```

Get system metrics:

```bash
nacos-cli service metrics --pretty
```

Get servers list:

```bash
nacos-cli service servers --pretty
```

Get cluster leader:

```bash
nacos-cli service leader --pretty
```

Update instance health:

```bash
nacos-cli service update-health --name=example-service --ip=192.168.1.1 --port=8080 --healthy=true
```

#### Namespace Management

List namespaces:

```bash
# List all namespaces
nacos-cli namespace list

# List all namespaces with pretty output
nacos-cli namespace list --pretty
```

Create namespace:

```bash
nacos-cli namespace create --id=dev --name="Development Environment" --desc="Used for development"
```

Modify namespace:

```bash
nacos-cli namespace modify --id=dev --name="Development" --desc="Development environment"
```

Delete namespace:

```bash
nacos-cli namespace delete --id=dev
```

## License

[MIT](https://github.com/k8scat/nacos-cli/blob/main/LICENSE)
