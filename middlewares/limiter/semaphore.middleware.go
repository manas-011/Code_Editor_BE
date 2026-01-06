package limiter

var ExecSemaphore = make(chan struct{}, 5) // max 5 concurrent executions
