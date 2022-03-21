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
	SplitNumber int // The split this task reads in
	InputFile  string
	OutputFiles []string
}

// Reduce Task
type ReduceTask struct {
	// In real MR, reduce needs to be informed of all map tasks' file, here they are smashed into one
	Partition int // partition number this reduce task reads in
	InputFile string
	OutputFile string
}

// file name
const IntermediateFileNamePrefix = "mr-mapped-"

