package log

import (
	"context"
	"github.com/sirupsen/logrus"
)

var Logger = logrus.New()

func init() {
	Logger.SetFormatter(&logrus.JSONFormatter{})
	Logger.SetReportCaller(true)
	Logger.SetLevel(logrus.DebugLevel)
}

func WithContext(ctx context.Context) *logrus.Entry {
	if ctx == nil {
		ctx = context.Background()
		return Logger.WithContext(ctx)
	}
	return Logger.WithContext(ctx).WithFields(logrus.Fields{
		"req_id": ctx.Value("X-Request-Id"),
	})
}
