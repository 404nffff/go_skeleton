package mongo

import (
	"context"
	"fmt"
	"tool/app/global/variable"

	"go.mongodb.org/mongo-driver/event"
	"go.uber.org/zap"
)

// CustomLogger implements mongo.Logger interface for custom logging.
type CustomLogger struct{}

func (cl CustomLogger) Log(ctx context.Context, msg string, args ...interface{}) {
	variable.Logs.Info(fmt.Sprintf(msg, args...))
}

// NewMonitor returns a CommandMonitor for logging MongoDB commands.
func NewMonitor() *event.CommandMonitor {
	return &event.CommandMonitor{
		Started: func(ctx context.Context, evt *event.CommandStartedEvent) {
			variable.Logs.Info("MongoDB Command Started",
				zap.String("Database", evt.DatabaseName),
				zap.String("Command", evt.CommandName),
				zap.String("CommandDetails", evt.Command.String()),
				zap.Int64("RequestID", evt.RequestID),
				zap.String("ConnectionID", evt.ConnectionID),
			)
		},
		Succeeded: func(ctx context.Context, evt *event.CommandSucceededEvent) {
			variable.Logs.Info("MongoDB Command Succeeded",
				zap.String("Command", evt.CommandName),
				zap.Int64("RequestID", evt.RequestID),
				zap.Duration("Duration", evt.Duration),
				//zap.String("Reply", evt.Reply.String()),
			)
		},
		Failed: func(ctx context.Context, evt *event.CommandFailedEvent) {
			variable.Logs.Error("MongoDB Command Failed",
				zap.String("Command", evt.CommandName),
				zap.Int64("RequestID", evt.RequestID),
				zap.Duration("Duration", evt.Duration),
				zap.String("Failure", evt.Failure),
			)
		},
	}
}
