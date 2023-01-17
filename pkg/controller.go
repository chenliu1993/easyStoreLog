package pkg

import (
	"os"
	"os/signal"
	"syscall"
)

type Controller struct {
	sigs chan os.Signal
}

func New() *Controller {
	c := make(chan os.Signal)

	//nolint:go-staticcheck
	signal.Notify(c, syscall.SIGHUP, syscall.SIGTERM, syscall.SIGINT)
	return &Controller{
		sigs: c,
	}
}

func (ctrl *Controller) Start() {

}
