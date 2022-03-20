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
	taskState  TaskState
	taskType   TaskType
	mapTask MapTask
	reduceTask ReduceTask
}

func (this *Task) isIdle() bool {
	return this.taskState == Idle
}

// Map Task
type MapTask struct {
	inputFile  File
	outputFiles []File
}

// Reduce Task
type ReduceTask struct {
	// In real MR, reduce needs to be informed of all map tasks' file, here they are smashed into one
	inputFile File
	outputFile File
}

// file name
const IntermediateFileNamePrefix = "mr-mapped-"

