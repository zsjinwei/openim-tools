// Copyright © 2024 OpenIM open source community. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package memamq

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"time"
)

var (
	ErrStop = errors.New("push failed: queue is stopped")
	ErrFull = errors.New("push failed: queue is full")
)

const (
	pushWait = time.Second * 3
)

// AsyncQueue is the interface responsible for asynchronous processing of functions.
//type AsyncQueue interface {
//	Initialize(processFunc func(), workerCount int, bufferSize int)
//	Push(task func()) error
//}

// MemoryQueue is an implementation of the AsyncQueue interface using a channel to process functions.
type MemoryQueue struct {
	taskChan  chan func()
	wg        sync.WaitGroup
	isStopped atomic.Bool
	count     atomic.Int64
	//stopMutex sync.Mutex // Mutex to protect access to isStopped
}

func NewMemoryQueue(workerCount int, bufferSize int) *MemoryQueue {
	if workerCount < 1 || bufferSize < 1 {
		panic("workerCount and bufferSize must be greater than 0")
	}
	mq := &MemoryQueue{}                   // Create a new instance of MemoryQueue
	mq.initialize(workerCount, bufferSize) // Initialize it with specified parameters
	return mq
}

// Initialize sets up the worker nodes and the buffer size of the channel,
// starting internal goroutines to handle tasks from the channel.
func (mq *MemoryQueue) initialize(workerCount int, bufferSize int) {
	mq.taskChan = make(chan func(), bufferSize) // Initialize the channel with the provided buffer size.
	// Start multiple goroutines based on the specified workerCount.
	for i := 0; i < workerCount; i++ {
		mq.wg.Add(1)
		go func() {
			defer mq.wg.Done()
			for task := range mq.taskChan {
				task() // Execute the function
			}
		}()
	}
}

// Push submits a function to the queue.
// Returns an error if the queue is stopped or if the queue is full.
func (mq *MemoryQueue) Push(task func()) error {
	mq.count.Add(1)
	defer mq.count.Add(-1)
	if mq.isStopped.Load() {
		return ErrStop
	}
	timer := time.NewTimer(pushWait)
	defer timer.Stop()
	select {
	case mq.taskChan <- task:
		return nil
	case <-timer.C: // Timeout to prevent deadlock/blocking
		return ErrFull
	}
}

func (mq *MemoryQueue) PushCtx(ctx context.Context, task func()) error {
	mq.count.Add(1)
	defer mq.count.Add(-1)
	if mq.isStopped.Load() {
		return ErrStop
	}
	select {
	case mq.taskChan <- task:
		return nil
	case <-ctx.Done():
		return context.Cause(ctx)
	}
}

func (mq *MemoryQueue) BatchPushCtx(ctx context.Context, tasks ...func()) (int, error) {
	mq.count.Add(1)
	defer mq.count.Add(-1)
	if mq.isStopped.Load() {
		return 0, ErrStop
	}
	for i := range tasks {
		select {
		case <-ctx.Done():
			return i, context.Cause(ctx)
		case mq.taskChan <- tasks[i]:
		}
	}
	return len(tasks), nil
}

func (mq *MemoryQueue) NotWaitPush(task func()) error {
	mq.count.Add(1)
	defer mq.count.Add(-1)
	if mq.isStopped.Load() {
		return ErrStop
	}
	select {
	case mq.taskChan <- task:
		return nil
	default:
		return ErrFull
	}
}

// Stop is used to terminate the internal goroutines and close the channel.
func (mq *MemoryQueue) Stop() {
	if !mq.isStopped.CompareAndSwap(false, true) {
		return
	}
	mq.waitSafeClose()
	close(mq.taskChan)
	mq.wg.Wait()
}

func (mq *MemoryQueue) waitSafeClose() {
	if mq.count.Load() == 0 {
		return
	}
	ticker := time.NewTicker(time.Second / 10)
	defer ticker.Stop()
	for range ticker.C {
		if mq.count.Load() == 0 {
			return
		}
	}
}
