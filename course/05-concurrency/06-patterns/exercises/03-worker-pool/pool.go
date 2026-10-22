package pool

// Result pairs a job with the value work produced for it.
type Result struct {
	Job   int
	Value int
}

// Run starts n workers (at least 1) that take jobs from jobs until it is
// closed and send Result{job, work(job)} on the returned channel. The
// returned channel is closed after every worker has finished.
func Run(n int, jobs <-chan int, work func(int) int) <-chan Result {
	// TODO
	return nil
}
