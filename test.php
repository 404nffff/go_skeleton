<?php

require 'vendor/autoload.php';

use Predis\Client;

// 连接到 Redis
$redis = new Client([
    'scheme' => 'tcp',
    'host'   => '127.0.0.1',
    'port'   => 6381,
]);

// 定义任务
$task = [
    'Type' => 'email:send',
    'Payload' => json_encode([
        'user_id' => 123,
        'template' => 'welcome',
    ]),
];

// 将任务序列化为 JSON
$encodedTask = json_encode($task);

// 生成一个唯一的任务 ID
$taskId = uniqid('task:', true);

// 使用 Redis 事务来确保原子性操作
$redis->multi();

// 将任务添加到 asynq 的默认队列
$redis->hset("asynq:{default}", $taskId, $encodedTask);

// 将任务 ID 添加到处理队列
$redis->lpush("asynq:{default}:pending", $taskId);

// 执行事务
$redis->exec();

echo "Task enqueued with ID: $taskId\n";
