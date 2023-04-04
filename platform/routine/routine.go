package routine

import (
	"context"
	"sync"
	"time"

	"sso/platform/logger"
	"sso/platform/wait"

	"go.uber.org/zap"
)

// Routine contains the information of a routine.
type Routine struct {
	// Name is the name of the routine.
	Name string
	// Operation is the operation to be executed.
	Operation func(ctx context.Context, log logger.Logger)
	// NoWait indicates to not add the routine to the wait group.
	NoWait bool
	// Timeout is the timeout of the routine.
	//
	// The default value is 1 minute.
	Timeout time.Duration
}
type ContextKey any

// ExecuteRoutine executes a routine and recovers from panic.
// It also adds the routine to the wait group if NoWait is false.
// The routine is cancelled if it exceeds the timeout specified.
func ExecuteRoutine(ctx context.Context, routine Routine, log logger.Logger) {
	if !routine.NoWait {
		wait.RoutineWaitGroup.Add(1)
	}

	go func(wtGroup *sync.WaitGroup) {
		log := log.Named("routine")
		if routine.Timeout == 0 {
			routine.Timeout = 5 * time.Minute
		}
		ctxWithTimeout, cancel := context.WithTimeout(context.Background(), routine.Timeout)

		newCtx := context.WithValue(context.WithValue(
			context.WithValue(
				context.WithValue(
					ctxWithTimeout, ContextKey("request-start-time"), time.Now(),
				), ContextKey("x-ws-request-id"), ctx.Value("x-ws-request-id"),
			), ContextKey("x-user"), ctx.Value("x-user"),
		), ContextKey("x-request-id"), ctx.Value("x-request-id"))

		defer func() {
			cancel()
			if !routine.NoWait {
				wtGroup.Done()
			}
		}()

		c := make(chan int)

		go func(ctx context.Context, log logger.Logger) {
			defer func() {
				if r := recover(); r != nil {
					log.Error(ctx, "routine panicked",
						zap.String("routine-name", routine.Name),
						zap.Any("panic", r))
					c <- 1
				}
			}()

			routine.Operation(ctx, log)
			c <- 1
		}(newCtx, log)

		select {
		case <-newCtx.Done():
			log.Error(newCtx, "routine cancelled")
		case <-c:
		}
	}(wait.RoutineWaitGroup)
}
