package shutdown

import (
	"os"
	"os/signal"
	"syscall"
)

const ShutdownManagerName = "PosixSignalManager"

// PosixSignalManager is a kind of ShutdownManager
type PosixSignalManager struct {
	signals []os.Signal
}

// NewPosixSignalManager initializes the PosixSignalManager
// with signals which to shutdown.
// Default signal is SIGINT and SIGTERM
func NewPosixSignalManager(sig ...os.Signal) *PosixSignalManager {
	if len(sig) == 0 {
		sig = make([]os.Signal, 2)
		sig[0] = syscall.SIGINT
		sig[1] = syscall.SIGTERM
	}

	return &PosixSignalManager{
		signals: sig,
	}
}

func (posixSignalManager *PosixSignalManager) GetName() string {
	return ShutdownManagerName
}

// Start starts to listening the posix signals
func (posixSignalManager *PosixSignalManager) Start(gs GracefulShutdowner) error {
	go func() {
		// buffer channel but not sync channel
		c := make(chan os.Signal, 1)
		signal.Notify(c, posixSignalManager.signals...)

		// block until a signal is received
		<-c

		gs.StartShutdown(posixSignalManager)
	}()

	return nil
}

func (posixSignalManager *PosixSignalManager) ShutdownStart() error {
	return nil
}

func (posixSignalManager *PosixSignalManager) ShutdownFinish() error {
	os.Exit(0)

	return nil
}
