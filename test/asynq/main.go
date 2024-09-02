package main  // 确保这里是 package main

import (
    "context"
    "encoding/json"
    "fmt"
    "log"
    "time"

    "github.com/hibiken/asynq"
)

const (
    TypeEmailDelivery = "email:deliver"
)

type EmailDeliveryPayload struct {
    UserID int
    Email  string
}

func main() {  // 确保有 main() 函数
    // 创建一个 Redis 连接
    redisConnection := asynq.RedisClientOpt{
        Addr: "localhost:6381",
    }

    // 初始化一个任务客户端
    client := asynq.NewClient(redisConnection)
    defer client.Close()

    // 创建一个任务
    payload, _ := json.Marshal(EmailDeliveryPayload{UserID: 123, Email: "user@example.com"})
    task := asynq.NewTask(TypeEmailDelivery, payload)

    // 将任务加入队列
    info, err := client.Enqueue(task)
    if err != nil {
        log.Fatalf("could not enqueue task: %v", err)
    }
    fmt.Printf("Enqueued task: id=%s queue=%s\n", info.ID, info.Queue)

    // 创建一个任务服务器（消费者）
    srv := asynq.NewServer(
        redisConnection,
        asynq.Config{
            Concurrency: 10,
            Queues: map[string]int{
                "default": 1,
            },
        },
    )

    // 定义任务处理函数
    mux := asynq.NewServeMux()
    mux.HandleFunc(TypeEmailDelivery, HandleEmailDeliveryTask)

    // 启动任务服务器
    if err := srv.Run(mux); err != nil {
        log.Fatalf("could not run server: %v", err)
    }
}

func HandleEmailDeliveryTask(ctx context.Context, t *asynq.Task) error {
    var p EmailDeliveryPayload
    if err := json.Unmarshal(t.Payload(), &p); err != nil {
        return fmt.Errorf("json.Unmarshal failed: %v", err)
    }
    log.Printf("Sending Email to User: user_id=%d, email=%s", p.UserID, p.Email)
    // 在这里实现发送邮件的逻辑
    time.Sleep(5 * time.Second) // 模拟耗时操作
    return nil
}
