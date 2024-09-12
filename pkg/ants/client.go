package ants

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/panjf2000/ants/v2"
)

// AntsInterface 是一个定义使用 ants 管理和执行任务的方法的接口。
type AntsInterface interface {
	// Submit 提交一个任务给 ants 执行。
	// 它接受一个任务函数作为参数，并在提交失败时返回一个错误。
	Submit(task func()) error

	// Release 释放 ants 使用的所有资源。
	Release()

	// GetStatus 返回 ants 的当前状态，包括正在运行和等待的 goroutine 数量。
	// 它返回两个整数，分别表示正在运行和等待的 goroutine 数量。
	GetStatus() (int, int)

	// SubmitTask 提交一个带有额外参数的任务给 ants 执行。
	// 它接受一个任务函数、一个参数映射，并返回一个结果映射和一个错误，如果提交失败。
	SubmitTask(ctx context.Context, task func(params map[string]any) (map[string]any, error), params map[string]any) (map[string]any, error)

	// Push 将一个带有额外参数的任务推送到任务队列中，由 ants 执行。
	// 它接受一个上下文、一个函数和可变参数，并返回一个字符串和一个错误。
	Push(ctx context.Context, fn any, params ...any) (string, error)

	// Exec 使用 ants 执行任务队列中的任务。
	// 它接受一个上下文，并返回一个结果映射和一个错误。
	Exec(ctx context.Context) (map[string]any, map[string]error)

	// SetTimeout 设置 ants 的超时时间。
	// 它接受一个时间持续时间作为参数。
	SetTimeout(timeout time.Duration)
}

// 协程池
type Ants struct {
	pool        *ants.Pool    // 协程池
	taskPool    sync.Map      // 任务池
	mu          sync.RWMutex  // 读写锁
	taskTimeOut time.Duration // 任务超时时间
}

// 任务函数
type taskFunc struct {
	method any    // 函数
	params []any  // 参数
	id     string // 任务ID
}

var taskFuncPool = sync.Pool{
	New: func() interface{} {
		return new(taskFunc)
	},
}

func NewAnts(poolSize int) (AntsInterface, error) {
	pool, err := ants.NewPool(poolSize)
	if err != nil {
		return nil, fmt.Errorf("failed to create ants pool: %w", err)
	}
	return &Ants{pool: pool, taskTimeOut: 10 * time.Second}, nil
}

func (a *Ants) Submit(task func()) error {
	return a.pool.Submit(task)
}

func (a *Ants) Release() {
	a.pool.Release()
	a.taskPool = sync.Map{}
}

func (a *Ants) GetStatus() (int, int) {
	return a.pool.Running(), a.pool.Cap()
}

// SubmitTask 提交一个带有额外参数的任务给 ants 执行。
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

// Push 将一个带有额外参数的任务推送到任务队列中，由 ants 执行。
func (a *Ants) Push(ctx context.Context, fn any, params ...any) (string, error) {
	if reflect.TypeOf(fn).Kind() != reflect.Func {
		return "", errors.New("fn must be a function")
	}

	select {
	case <-ctx.Done():
		return "", ctx.Err()
	default:
	}

	a.mu.Lock()
	defer a.mu.Unlock()

	tasks, _ := a.taskPool.LoadOrStore(ctx, &sync.Map{})
	tasksMap := tasks.(*sync.Map)
	id := uuid.New().String()

	newTask := taskFuncPool.Get().(*taskFunc)
	newTask.method = fn
	newTask.params = params
	newTask.id = id

	var taskSlice []taskFunc
	if existingTasks, ok := tasksMap.Load("tasks"); ok {
		taskSlice = existingTasks.([]taskFunc)
	}
	taskSlice = append(taskSlice, *newTask)
	tasksMap.Store("tasks", taskSlice)

	return id, nil
}

// Exec 使用 ants 执行任务队列中的任务。
func (a *Ants) Exec(ctx context.Context) (map[string]any, map[string]error) {
	a.mu.Lock()
	defer func() {
		a.mu.Unlock()
		a.taskPool.Delete(ctx)

		//设置超时时间
		a.taskTimeOut = 10 * time.Second
	}()

	tasksInterface, ok := a.taskPool.Load(ctx)
	if !ok {
		return nil, nil
	}

	tasksMap := tasksInterface.(*sync.Map)
	tasksSliceInterface, ok := tasksMap.Load("tasks")
	if !ok {
		return nil, nil
	}

	tasks := tasksSliceInterface.([]taskFunc)
	taskResults := sync.Map{}
	errs := make(map[string]error, len(tasks)+1)

	var wg sync.WaitGroup
	wg.Add(len(tasks))

	for _, task := range tasks {
		task := task

		errs[task.id] = nil

		err := a.pool.Submit(func() {
			defer wg.Done()
			defer taskFuncPool.Put(&task)

			done := make(chan struct{})
			go func() {
				defer close(done)

				f := reflect.ValueOf(task.method)
				if f.Kind() != reflect.Func {
					errs[task.id] = errors.New("task is not a function")
					return
				}

				args := make([]reflect.Value, len(task.params)+1)
				args[0] = reflect.ValueOf(ctx)
				for j, param := range task.params {
					args[j+1] = reflect.ValueOf(param)
				}

				result := f.Call(args)
				if len(result) > 0 {
					taskResults.Store(task.id, result[0].Interface())
				}
				if len(result) > 1 && !result[1].IsNil() {
					errs[task.id] = result[1].Interface().(error)
				}
			}()

			select {
			case <-done:
			case <-ctx.Done():
				errs[task.id] = ctx.Err()
			}
		})
		if err != nil {
			errs[task.id] = fmt.Errorf("failed to submit task: %w", err)
		}
	}

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(a.taskTimeOut):

		for errKey, _ := range errs {
			errs[errKey] = errors.New("task timeout")
		}

		return nil, errs
	case <-ctx.Done():

		for errKey, _ := range errs {
			errs[errKey] = ctx.Err()
		}

		return nil, errs
	}

	results := make(map[string]any)
	taskResults.Range(func(key, value interface{}) bool {
		results[key.(string)] = value
		return true
	})

	return results, errs
}

// SetTimeout 设置 ants 的超时时间。
func (a *Ants) SetTimeout(timeout time.Duration) {
	a.taskTimeOut = timeout
}
