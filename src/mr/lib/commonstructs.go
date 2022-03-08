package lib

type File struct {
	fileName string
	startKey string
	size     int64
}

// Define the TaskState of a worker
type TaskState int64

const (
	Idle TaskState = iota
	InProgress
	Completed
)

// Task
type TaskType int64

const (
	Map TaskType = iota
	Reduce
)

type Task struct {
	taskState      TaskState
	taskType   TaskType
	inputFile  File
	outputFile File
}

func (this *Task) isIdle() bool {
	return this.taskState == Idle
}
