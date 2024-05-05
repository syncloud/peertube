package installer

import (
	"go.uber.org/zap"
	"os/exec"
)

type Executor struct {
	logger *zap.Logger
}

func (d *Executor) Run(app string, args ...string) error {
	cmd := exec.Command(app, args...)
	d.logger.Info("executing", zap.String("cmd", cmd.String()))
	out, err := cmd.CombinedOutput()
	d.logger.Info(cmd.String(), zap.ByteString("output", out))
	if err != nil {
		d.logger.Error(cmd.String(), zap.String("error", "failed to run command"))
	}
	return err
}
