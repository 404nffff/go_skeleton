## go 框架 骨架封装

### 参考项目
https://github.com/qifengzhang007/GinSkeleton

### 使用组件
1. gin web框架
2. gorm 数据库orm
3. viper 配置文件
4. logrus 日志分割
5. zap 日志组件
6. ants 协程池
7. go.mongodb.org/mongo-driver/mongo mongodb驱动
8. github.com/go-redis/redis/v8 redis驱动
9. github.com/gin-contrib/sessions session管理
10. github.com/bradfitz/gomemcache/memcache  memcache驱动


### 启动命令
```shell
go run cmd/api/main.go [start| start debug| stop | restart]
```

start: 启动后台服务

start debug: 启动后台服务，日志输出到控制台

stop: 停止后台服务

restart: 重启后台服务


### 热更新
1. 使用air工具进行热更新
2. 安装air
```shell
go install github.com/air-verse/air@latest
```


### 静态资源打包 embed
1. 使用go1.16版本以上

2. 打包静态资源
```go
//go:embed admin/layouts/*.tmpl
//go:embed admin/*.tmpl
//go:embed admin/user/*.tmpl
```



### docker部署
1. 构建镜像 
项目根目录下执行
```shell
docker build -t deploy/build/dockerfile .
```
2. 运行容器
```shell
docker run -itd -e APP_DEBUG=false -v "config":/app/config -v "logs":/app/logs -p 8080:8080 xxxx
```
config 目录下存放配置文件
logs 目录下存放日志文件