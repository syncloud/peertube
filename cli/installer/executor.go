package installer

import (
	"go.uber.org/zap"
	"os/exec"
)

type Executor struct {
	logger *zap.Logger
}

func NewExecutor(logger *zap.Logger) *Executor {
	return &Executor{
		logger: logger,
	}
}

func (e *Executor) Run(app string, args ...string) error {
	cmd := exec.Command(app, args...)
	e.logger.Info("executing", zap.String("cmd", cmd.String()))
	out, err := cmd.CombinedOutput()
	e.logger.Info(cmd.String(), zap.ByteString("output", out))
	if err != nil {
		e.logger.Error(cmd.String(), zap.Error(err))
	}
	return err
}
