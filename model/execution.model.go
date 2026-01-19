package model

type ExecuteRequest struct {
	Language string `json:"language"`
	Code     string `json:"code"`
	Input    string `json:"input"`
}

type ExecuteResult struct {
	Status string `json:"status"`
	Stdout string `json:"stdout"`
	Stderr string `json:"stderr"`
}


type JudgeResult struct {
	Status   string // CE, RE, TLE, MLE, success
	Stdout   string
	Stderr   string // includes warnings
	ExitCode int
}