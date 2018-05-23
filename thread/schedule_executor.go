package thread

import (
	"errors"
	"time"
)

var (
	CancellationException = errors.New("the computation was cancelled")
	ExecutionException = errors.New("the computation threw an exception")
	InterruptedException = errors.New("the current thread was interrupted while waiting")
	TimeoutException = errors.New("the wait timed out")
	RejectedExecutionException = errors.New("the task cannot be scheduled for execution")
)

type Future interface {
	/*
	mayInterruptIfRunning - true if the thread executing this task should be interrupted; otherwise, in-progress tasks are allowed to complete
	false if the task could not be cancelled, typically because it has already completed normally; true otherwise
	 */
	cancel(mayInterruptIfRunning bool)bool

	/*
	Returns true if this task was cancelled before it completed normally.
	true if this task was cancelled before it completed
	 */
	isCancelled()bool

	/*
	Returns true if this task completed. Completion may be due to normal termination, an exception, or cancellation -- in all of these cases, this method will return true.
	true if this task completed
	 */
	isDone()

	/*
		Waits if necessary for the computation to complete, and then retrieves its result.
		the computed result
	 */
	get()(interface{},error)

	/*
		Waits if necessary for at most the given time for the computation to complete, and then retrieves its result, if available.
	 */
	getWithTimeOut(timeout time.Duration)(interface{},error)

}

type ScheduleFuture interface {


	Future
	/*
	Returns the remaining delay associated with this object, in the given time unit.
	the remaining delay; zero or negative values indicate that the delay has already elapsed
	 */
	getDelay()time.Duration
}

type CallableBlock struct {
	Tasks func([]interface{})([]interface{},error)
	Params []interface{}
}

type ScheduledExecutorService interface {
	schedule(delay time.Duration,block *CallableBlock)ScheduleFuture

	scheduleAtFixedRate(initialDelay,period time.Duration,block *CallableBlock)ScheduleFuture

	scheduleWithFixedDelay(initalDelay,delay time.Duration,block *CallableBlock)ScheduleFuture
}