package lib

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
	TaskState  TaskState
	TaskType   TaskType
	MapTask MapTask
	ReduceTask ReduceTask
}

func (this *Task) IsIdle() bool {
	return this.TaskState == Idle
}

// Map Task
type MapTask struct {
	InputFile  string
	OutputFiles []string
}

// Reduce Task
type ReduceTask struct {
	// In real MR, reduce needs to be informed of all map tasks' file, here they are smashed into one
	InputFile string
	OutputFile string
}

// file name
const IntermediateFileNamePrefix = "mr-mapped-"

