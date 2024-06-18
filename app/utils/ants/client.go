package ants

import (
	"errors"

	"github.com/panjf2000/ants/v2"
)

// AntsInterface 定义一个接口，使 Ants 更易于测试和扩展
type AntsInterface interface {
	Submit(task func()) error
	Release()
	GetStatus() (int, int)
	SubmitTask(task func(params map[string]any) map[string]any, params map[string]any) (map[string]any, error)
}

// Ants 结构体封装了 ants.Pool
type Ants struct {
	pool *ants.Pool
}

// NewAnts 初始化并返回一个 Ants 实例
func NewAnts(poolSize int) (AntsInterface, error) {
	pool, err := ants.NewPool(poolSize)
	if err != nil {
		return nil, err
	}
	return &Ants{pool: pool}, nil
}

// Submit 提交任务到 ants 池
func (a *Ants) Submit(task func()) error {
	return a.pool.Submit(task)
}

// Release 释放 ants 池资源
func (a *Ants) Release() {
	a.pool.Release()
}

// GetStatus 获取 ants.Pool 状态, 返回正在运行的 goroutine 数和池的容量
func (a *Ants) GetStatus() (int, int) {
	running := a.pool.Running()
	capacity := a.pool.Cap()
	return running, capacity
}

// SubmitTask 提交带参数的任务到 ants 池，并等待结果返回，结果为 JSON 格式
func (a *Ants) SubmitTask(task func(params map[string]any) map[string]any, params map[string]any) (map[string]any, error) {
	// 创建一个通道来接收任务结果
	resultChan := make(chan map[string]any, 1)

	// 从池中获取一个 goroutine
	err := a.pool.Submit(func() {
		defer close(resultChan) // 确保通道在任务完成后关闭
		result := task(params)
		resultChan <- result // 将结果发送到通道
	})

	if err != nil {
		panic(err)
	}

	// 等待并返回结果
	result, ok := <-resultChan
	if !ok {
		panic(errors.New("task failed"))
	}

	// // 将结构体转换为 JSON
	// jsonData, err := json.Marshal(result)
	// if err != nil {
	// 	fmt.Println("Error marshalling to JSON:", err)
	// 	panic(errors.New("task failed"))
	// }

	return result, nil
}
