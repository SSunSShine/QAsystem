# QAsystem

仿知乎问答系统服务端

目前正在完善中。。。

## 功能

+ [x] 用户可以注册、登录
+ [x] 用户可以发布问题
+ [x] 用户可以看到别人发布的问题

+ [x] 用户可以回答问题，可以修改回答和删除回答
+ [x] 用户可以修改和删除问题

+ [x] 用户可以看到热门问题
+ [x] 用户可以对回答点赞和踩
+ [x] 用户可以在个人中心看到自己发布的问题，回答和赞的内容

+ [x] 用户可以对回答进行评论
+ [x] 用户可以修改和删除评论

注：hot = (浏览×0.1 + 回答×0.5 + 评论×0.2 点赞×0.2)×1000
   每5分钟更新后台的数据。

## Docker 快速部署

```shell
$ docker-compose up
```

## 常规启动

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
$ vim ./conf/configuration.yaml
```

```yaml
# 数据库
db:
  driver: mysql
  addr: mysql8019:root1234@/qasystem?charset=utf8&parseTime=True&loc=Local

# Redis
redis:
  addr: redis-db:6379
  password:
  db: 0

# jwt认证密钥
jwtKey: 设置jwt密钥

# 端口
address: :8080
```

### 4.初始化并运行

```shell
$ (sudo) go run ./ -init
```

## 技术栈

| 技术              | 链接                                             |  描述             |
| ----------------- | ------------------------------------------------|:---------------: |
| gin               | https://github.com/gin-gonic/gin                |  web框架          |
| gorm              | https://github.com/jinzhu/gorm                  |  数据库orm框架     |
| viper             | https://github.com/spf13/viper                  |  配置管理工具      |
| go-redis          | https://github.com/go-redis/redis/v8            |  redis工具        |
| validator         | https://github.com/go-playground/validator/v10  |  数据校验工具      |
| jwt-go            | https://github.com/dgrijalva/jwt-go             |  jwt认证工具        |
| bcrypt            | https://golang.org/x/crypto/bcrypt              |  密码加密工具      |
| httpexpect        | https://github.com/gavv/httpexpect/v2           |  API测试工具       |
| docker            | https://docs.docker.com/                        |  容器部署          |
