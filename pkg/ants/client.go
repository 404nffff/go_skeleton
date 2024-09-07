package ants

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"sync"

	"github.com/panjf2000/ants/v2"
)

// AntsInterface 定义了 Ants 结构体应实现的方法集合
// 这个接口使得 Ants 更易于测试和扩展
type AntsInterface interface {
	Submit(task func()) error
	Release()
	GetStatus() (int, int)
	SubmitTask(ctx context.Context, task func(params map[string]any) (map[string]any, error), params map[string]any) (map[string]any, error)
	Push(ctx context.Context, fn any, params ...any) error
	Exec(ctx context.Context) ([]any, error)
}

// Ants 结构体封装了 ants.Pool，提供了更高级的任务管理功能
type Ants struct {
	pool     *ants.Pool   // 底层的 ants 协程池
	taskPool sync.Map     // 用于存储待执行的任务
	mu       sync.RWMutex // 用于保护并发访问
}

// taskFunc 表示一个待执行的任务
type taskFunc struct {
	method any   // 要执行的方法
	params []any // 方法的参数
}

// NewAnts 初始化并返回一个 Ants 实例
// poolSize 参数指定了底层 ants.Pool 的大小
func NewAnts(poolSize int) (AntsInterface, error) {
	pool, err := ants.NewPool(poolSize)
	if err != nil {
		return nil, fmt.Errorf("failed to create ants pool: %w", err)
	}
	return &Ants{
		pool: pool,
	}, nil
}

// Submit 提交一个无参数的任务到 ants 池
func (a *Ants) Submit(task func()) error {
	return a.pool.Submit(task)
}

// Release 释放 ants 池资源
// 在程序结束时应调用此方法以释放资源
func (a *Ants) Release() {
	a.pool.Release()
}

// GetStatus 获取 ants.Pool 的当前状态
// 返回正在运行的 goroutine 数和池的容量
func (a *Ants) GetStatus() (int, int) {
	return a.pool.Running(), a.pool.Cap()
}

// SubmitTask 提交带参数的任务到 ants 池，并等待结果返回
// 支持通过 context 进行取消操作
func (a *Ants) SubmitTask(ctx context.Context, task func(params map[string]any) (map[string]any, error), params map[string]any) (map[string]any, error) {
	resultChan := make(chan struct {
		result map[string]any
		err    error
	}, 1)

	err := a.pool.Submit(func() {
		defer close(resultChan)
		result, err := task(params)
		resultChan <- struct {
			result map[string]any
			err    error
		}{result, err}
	})

	if err != nil {
		return nil, fmt.Errorf("failed to submit task: %w", err)
	}

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case res := <-resultChan:
		return res.result, res.err
	}
}

// Push 方法用于添加任务到任务池
// fn 参数应该是一个函数，params 是该函数的参数
func (a *Ants) Push(ctx context.Context, fn any, params ...any) error {
	if reflect.TypeOf(fn).Kind() != reflect.Func {
		return errors.New("fn must be a function")
	}

	a.mu.Lock()
	defer a.mu.Unlock()

	tasks, _ := a.taskPool.LoadOrStore(ctx, &sync.Map{})
	tasksMap := tasks.(*sync.Map)

	newTask := taskFunc{
		method: fn,
		params: params,
	}

	var taskSlice []taskFunc
	if existingTasks, ok := tasksMap.Load("tasks"); ok {
		taskSlice = existingTasks.([]taskFunc)
	}
	taskSlice = append(taskSlice, newTask)
	tasksMap.Store("tasks", taskSlice)

	return nil
}

// Exec 方法用于执行任务池中的所有任务
// 返回所有任务的结果和可能发生的错误
func (a *Ants) Exec(ctx context.Context) ([]any, error) {
	a.mu.Lock()
	defer func() {
		a.mu.Unlock()

		// 执行完成后，清理任务池
		a.taskPool.Delete(ctx)
	}()

	tasksInterface, ok := a.taskPool.Load(ctx)
	if !ok {
		return nil, nil // 没有任务需要执行
	}

	tasksMap := tasksInterface.(*sync.Map)
	tasksSliceInterface, ok := tasksMap.Load("tasks")
	if !ok {
		return nil, nil // 没有任务需要执行
	}

	tasks := tasksSliceInterface.([]taskFunc)
	taskResults := make([]any, len(tasks))
	errs := make([]error, len(tasks))

	var wg sync.WaitGroup
	wg.Add(len(tasks))

	for i, task := range tasks {
		i, task := i, task // 创建局部变量以在闭包中正确捕获
		err := a.pool.Submit(func() {
			defer wg.Done()

			f := reflect.ValueOf(task.method)
			if f.Kind() != reflect.Func {
				errs[i] = errors.New("task is not a function")
				return
			}

			args := make([]reflect.Value, len(task.params))
			for j, param := range task.params {
				args[j] = reflect.ValueOf(param)
			}

			result := f.Call(args)
			if len(result) > 0 {
				taskResults[i] = result[0].Interface()
			}
			if len(result) > 1 && !result[1].IsNil() {
				errs[i] = result[1].Interface().(error)
			}
		})
		if err != nil {
			errs[i] = fmt.Errorf("failed to submit task: %w", err)
		}
	}

	wg.Wait()

	// 检查是否有错误发生
	for _, err := range errs {
		if err != nil {
			return taskResults, fmt.Errorf("one or more tasks failed: %w", err)
		}
	}

	return taskResults, nil
}
