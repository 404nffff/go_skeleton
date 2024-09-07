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

type AntsInterface interface {
	Submit(task func()) error
	Release()
	GetStatus() (int, int)
	SubmitTask(ctx context.Context, task func(params map[string]any) (map[string]any, error), params map[string]any) (map[string]any, error)
	Push(ctx context.Context, fn any, params ...any) (string, error)
	Exec(ctx context.Context, timeout time.Duration) (map[string]any, error)
}

type Ants struct {
	pool     *ants.Pool
	taskPool sync.Map
	mu       sync.RWMutex
}

type taskFunc struct {
	method any
	params []any
	id     string
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
	return &Ants{pool: pool}, nil
}

func (a *Ants) Submit(task func()) error {
	return a.pool.Submit(task)
}

func (a *Ants) Release() {
	a.pool.Release()
}

func (a *Ants) GetStatus() (int, int) {
	return a.pool.Running(), a.pool.Cap()
}

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

func (a *Ants) Exec(ctx context.Context, timeout time.Duration) (map[string]any, error) {
	a.mu.Lock()
	defer func() {
		a.mu.Unlock()
		a.taskPool.Delete(ctx)
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
	errs := make([]error, len(tasks))

	var wg sync.WaitGroup
	wg.Add(len(tasks))

	for i, task := range tasks {
		i, task := i, task
		err := a.pool.Submit(func() {
			defer wg.Done()
			defer taskFuncPool.Put(&task)

			done := make(chan struct{})
			go func() {
				defer close(done)

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
					taskResults.Store(task.id, result[0].Interface())
				}
				if len(result) > 1 && !result[1].IsNil() {
					errs[i] = result[1].Interface().(error)
				}
			}()

			select {
			case <-done:
			case <-ctx.Done():
				errs[i] = ctx.Err()
			}
		})
		if err != nil {
			errs[i] = fmt.Errorf("failed to submit task: %w", err)
		}
	}

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(timeout):
		return nil, fmt.Errorf("execution timed out after %v", timeout)
	case <-ctx.Done():
		return nil, ctx.Err()
	}

	var finalErr error
	for _, err := range errs {
		if err != nil {
			finalErr = errors.Join(finalErr, err)
		}
	}

	results := make(map[string]any)
	taskResults.Range(func(key, value interface{}) bool {
		results[key.(string)] = value
		return true
	})

	return results, finalErr
}
