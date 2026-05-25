// Package signal provides OS signal interception and graceful shutdown
// coordination for patchwork-deploy.
//
// # Overview
//
// Handler wraps the standard library signal.Notify pattern and cancels
// a context.CancelFunc when an interrupt or termination signal arrives.
//
// ShutdownCoordinator builds on Handler by adding a configurable drain
// timeout: after the signal fires it waits for in-progress work to complete
// before declaring the shutdown forced.
//
// # Usage
//
//	ctx, cancel := context.WithCancel(context.Background())
//	sc := signal.NewShutdownCoordinator(os.Stdout, 30*time.Second, cancel)
//	finished := sc.Start(ctx)
//	// ... start workers using ctx ...
//	<-finished
package signal
