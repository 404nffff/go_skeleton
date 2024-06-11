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


