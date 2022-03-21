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
	mapTasksByState    map[lib.TaskState][]lib.Task
	reduceTasksByState map[lib.TaskState][]lib.Task
	mapLock            sync.Mutex
	reduceLock         sync.Mutex
}

func (this *TaskStore) AddMapTask(task lib.Task) {
	mapLock.Lock()
	defer mapLock.Unlock()
	mapTasksByState[task.TaskState] = append(mapTasksByState[task.TaskState], task)
}

func (this *TaskStore) AddReduceTask(task lib.Task) {
	reduceLock.Lock()
	defer reduceLock.Unlock()
	reduceTasksByState[task.TaskState] = append(reduceTasksByState[task.TaskState], task)
}

func (this *TaskStore) PopIdleMapTask() lib.Task {
	mapLock.Lock()
	defer mapLock.Unlock()
	if len(this.mapTasksByState) == 0 {
		return nil
	}
	task := this.mapTasksByState[lib.Idle][0]
	this.mapTasksByState[lib.Idle][1:]
	return task
}

func (this *TaskStore) PopIdleReduceTask() lib.Task {
	reduceLock.Lock()
	defer reduceLock.Unlock()
	if len(this.reduceTasksByState) == 0 {
		return nil
	}
	task := this.reduceTasksByState[lib.Idle][0]
	this.reduceTasksByState[lib.Idle][1:]
	return task
}

func (this *TaskStore) HasMapTask() bool {
	mapLock.Lock()
	defer mapLock.Unlock()
	return len(this.mapTasksByState[lib.Idle]) != 0 
		|| len(this.mapTasksByState[lib.InProgress]) != 0
}

func (this *TaskStore) HasReduceTask() bool {
	reduceLock.Lock()
	defer reduceLock.Unlock()
	return len(this.reduceTasksByState[lib.Idle]) != 0 
		|| len(this.reduceTasksByState[lib.InProgress]) != 0 
}

func (this *TaskStore) AreTasksDone() bool {
	return !this.HasMapTask() && !this.HasReduceTask()
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
func (this *Master) AssignTask() AssignTaskResponse {
	response := AssignTaskResponse{}
	if this.taskstore.HasMapTask() {
		task := this.taskStore.PopIdleMapTask()
		if task != nil {
			response.Task = task
		}
		return response;
	}

	if this.hasReduceTask() {
		task := this.taskStore.PopIdleReduceTask()
		if task != nil {
			response.Task = task
		}
		return response;
	}

	return response
}

// init the job
func (this *Master) initialize(files []string, nReduce int) {
	this.jobConfig = lib.JobConfig{}
	// create Map Tasks

	// create Reduce Tasks
}

//
// create a Master.
// main/mrmaster.go calls this function.
// nReduce is the number of reduce tasks to use.
//
func MakeMaster(files []string, nReduce int) *Master {
	m := Master{}

	// Your code here.
	m.initialize(files, nReduce)

	m.server()
	return &m
}
