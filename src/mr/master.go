package mr

import (
	"houwang/lib"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"sync"
)

type TaskStore struct {
	mapTasks []lib.Task
	reduceTasks []lib.Task
	tasksByState map[lib.TaskState][]*lib.Task
	lock sync.Mutex
}

func (this *TaskStore) AddMapTask(task lib.Task) {
	lock.Lock()
	defer lock.Unlock()
	this.mapTasks = append(this.mapTasks, task)
}

func (this *TaskStore) AddReduceTask(task lib.Task) {
	lock.Lock()
	defer lock.Unlock()
	this.reduceTasks = append(this.reduceTasks, task)
}

func (this *TaskStore) PopIdleTask() lib.Task{
	lock.Lock()
	defer lock.Unlock()
	if len(this.tasksByState) == 0 {
		return lib.Task{}
	}
	task := this.tasksByState[lib.Idle][0]
	this.tasksByState[lib.Idle][1:]
	return *task
}

func (this *TaskStore) AreTasksDone() {
	lock.Lock()
	defer lock.Unlock()
	return len(this.tasksByState[lib.Idle]) == 0 && len(this.tasksByState[lib.InProgress]) ==0
}

type Master struct {
	jobConfig lib.JobConfig
	taskStore TaskStore
}

//
// start a thread that listens for RPCs from worker.go
//
func (m *Master) server() {
	rpc.Register(m)
	rpc.HandleHTTP()
	//l, e := net.Listen("tcp", ":1234")
	sockname := masterSock()
	os.Remove(sockname)
	l, e := net.Listen("unix", sockname)
	if e != nil {
		log.Fatal("listen error:", e)
	}
	go http.Serve(l, nil)
}

//
// main/mrmaster.go calls Done() periodically to find out
// if the entire job has finished. This is the heartbeats
//
func (this *Master) Done() bool {
	return this.taskStore.AreTasksDone()
}

/**
* RPC Handles
 */

// Try to assign an idle task to worker,
func AssignTask() AssignTaskResponse {

}

//
// create a Master.
// main/mrmaster.go calls this function.
// nReduce is the number of reduce tasks to use.
//
func MakeMaster(files []string, nReduce int) *Master {
	m := Master{}

	// Your code here.
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	m.server()
	return &m
}
