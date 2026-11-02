package installer

import (
	"go.uber.org/zap"
	"os/exec"
	"strings"
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
	return e.RunInDir("", app, args...)
}

func (e *Executor) RunInDir(dir string, app string, args ...string) error {
	cmd := exec.Command(app, args...)
	cmd.Dir = dir
	e.logger.Info("executing", zap.String("cmd", cmd.String()), zap.String("dir", dir))
	out, err := cmd.CombinedOutput()
	e.logger.Info("command output")
	for _, line := range strings.Split(string(out), "\n") {
		e.logger.Info(line)
	}
	return err
}
