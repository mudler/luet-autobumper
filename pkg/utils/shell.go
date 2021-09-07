package utils

import (
	"os"
	"os/exec"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

func RunSH(cmd string, opts ...func(cmd *exec.Cmd)) (string, error) {
	log.Debug("Executing", cmd)
	c := exec.Command("sh", "-c", cmd)
	c.Env = os.Environ()

	for _, o := range opts {
		o(c)
	}

	out, err := c.CombinedOutput()
	if err != nil {
		return "", errors.Wrapf(err, "failed running shell command: %s", string(out))
	}
	return string(out), err
}
