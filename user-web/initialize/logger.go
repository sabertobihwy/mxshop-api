package initialize

import "go.uber.org/zap"

func InitLogger() {
	// set a global logger to use zap.S() -> secure access by goroutines with Mutex lock
	logger, _ := zap.NewDevelopment()
	zap.ReplaceGlobals(logger)
}
