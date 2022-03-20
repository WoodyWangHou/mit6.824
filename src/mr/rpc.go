package mr

//
// RPC definitions.
//
// remember to capitalize all names.
//

import "os"
import "strconv"
import "houwang/mr/lib"

//
// example to show how to declare the arguments
// and reply for an RPC.
//

type AssignTaskResponse struct {
	task lib.Task
}

type UpdateTaskStateRequest struct {
	taskState lib.TaskState
}

// Add your RPC definitions here.


// Cook up a unique-ish UNIX-domain socket name
// in /var/tmp, for the master.
// Can't use the current directory since
// Athena AFS doesn't support UNIX-domain sockets.
func masterSock() string {
	s := "/var/tmp/824-mr-"
	s += strconv.Itoa(os.Getuid())
	return s
}
