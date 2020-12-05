# QAsystem

仿知乎问答系统服务端

目前正在完善中。。。

### 1.获取代码

```shell
$ git clone https://github.com/SSunSShine/QAsystem

$ cd QAsystem
```

### 2.下载依赖

```shell
$ go mod tidy
```

### 3.修改配置信息
```shell
$ touch ./conf/configuration.yaml
```

```yaml
# 数据库
db:
  driver: mysql
  addr: root:123456@/qasystem?charset=utf8&parseTime=True&loc=Local

# Redis
redis:
  addr: 127.0.0.1:6379
  password:
  db: 0

# jwt认证密钥
jwtKey: jwt123456

# 端口
address: :8080
```

### 4.初始化并运行

```shell
$ (sudo) go run ./ -init
```
