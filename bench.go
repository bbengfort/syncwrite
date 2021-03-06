package syncwrite

import (
	"sync"
	"time"
)

// Benchmark conducts throughput testing of a log, measuring the number of
// Append operations per second for a fixed number of actions; running one or
// more go routines to conduct the operation on the log.
func Benchmark(log Log, path string, n, threads int) (float64, error) {
	// Open the log at the location
	if err := log.Open(path); err != nil {
		return 0.0, err
	}
	defer log.Close()

	// Create the action
	action := func() error {
		return log.Append([]byte("foo"))
	}

	// Create the parallel constructs
	group := new(sync.WaitGroup)
	errors := make([]error, threads)

	// Run the specified number of Go routines
	start := time.Now()
	for i := 0; i < threads; i++ {
		group.Add(1)
		go func(idx int) {
			errors[idx] = benchmarker(n, action)
			group.Done()
		}(i)
	}

	// Wait for the group to complete
	group.Wait()
	totalLatency := time.Since(start)

	// Check for errors
	for _, err := range errors {
		if err != nil {
			return 0.0, err
		}
	}

	// Compute the throughput
	throughput := float64(n*threads) / totalLatency.Seconds()
	return throughput, nil
}

//===========================================================================
// Benchmarker Go Routines and timer.
//===========================================================================

type action func() error

// run a single go routine that runs the action N times then returns the
// total time it took to run all n actions. If an error occurs in the action,
// then the error is returned.
func benchmarker(n int, f action) (err error) {
	for i := 0; i < n; i++ {
		if err = f(); err != nil {
			return err
		}
	}
	return nil
}
