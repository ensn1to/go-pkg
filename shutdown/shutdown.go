package shutdown

import (
	"sync"
)

// ShutdownManager
type ShutdownManager interface {
	GetName() string
	Start(GracefulShutdowner) error
	ShutdownStart() error
	ShutdownFinish() error
}

// ShutdownCallbacker
// OnShutdown will be called when shutdown is requested. The para is
// the name of the shutdownmanger to shudown.
type ShutdownCallbacker interface {
	OnShutdown(string) error
}

// ShutdownFunc kinds of helper type. so you can easily provide anonymous function
// as ShutdownCallbacks
type ShutdownFunc func(string) error

func (shutdownFunc ShutdownFunc) OnShutdown(shutdownManager string) error {
	return shutdownFunc(shutdownManager)
}

type ErrorHandler interface {
	OnError(error)
}

type ErrorHandlerFunc func(error)

func (errorHandlerFunc ErrorHandlerFunc) OnError(err error) {
	errorHandlerFunc(err)
}

// GracefulShutdowner
type GracefulShutdowner interface {
	StartShutdown(ShutdownManager)
	AddShutdownCallback(ShutdownCallbacker)
	ReportError(error)
}

// GracefulShutdown main struct
type GracefulShutdown struct {
	shutdownManagers  []ShutdownManager
	shutdownCallbacks []ShutdownCallbacker
	errorHandler      ErrorHandler
}

func New() *GracefulShutdown {
	return &GracefulShutdown{
		shutdownManagers:  make([]ShutdownManager, 0, 2),
		shutdownCallbacks: make([]ShutdownCallbacker, 0, 10),
	}
}

// Start calls method start() on all added ShutdownManagers.
// The ShutdownManager start to listen shutdown requests
func (gs *GracefulShutdown) Start() error {
	for _, manager := range gs.shutdownManagers {
		if err := manager.Start(gs); err != nil {
			return err
		}
	}

	return nil
}

// AddShutdownManager adds a shutdownManger
func (gs *GracefulShutdown) AddShutdownManager(shutdownManager ShutdownManager) {
	gs.shutdownManagers = append(gs.shutdownManagers, shutdownManager)
}

// AddShutdownCallback adds a shutdownCallback that will be called when shutdown is requested.
func (gs *GracefulShutdown) AddShutdownCallback(shutdownCallback ShutdownCallbacker) {
	gs.shutdownCallbacks = append(gs.shutdownCallbacks, shutdownCallback)
}

// SetErrorHandler sets errorHandlers that will be called
// when error happended in ShutdownManager or in ShutdownCallbacks
func (gs *GracefulShutdown) SetErrorHandler(errorHandler ErrorHandler) {
	gs.errorHandler = errorHandler
}

// StartShutdown called when a shutdownManager starts to shutdown
func (gs *GracefulShutdown) StartShutdown(shutdownManager ShutdownManager) {
	gs.ReportError(shutdownManager.ShutdownStart())

	var wg sync.WaitGroup
	for _, shutdownCallback := range gs.shutdownCallbacks {
		wg.Add(1)
		go func(scb ShutdownCallbacker) {
			defer wg.Done()
			gs.ReportError(scb.OnShutdown(shutdownManager.GetName()))
		}(shutdownCallback)
	}

	wg.Wait()

	gs.ReportError(shutdownManager.ShutdownFinish())
}

// ReportError be used to report errors to ErrorHandler
func (gs *GracefulShutdown) ReportError(err error) {
	if err != nil && gs.errorHandler != nil {
		gs.errorHandler.OnError(err)
	}
}
