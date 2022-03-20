package mr

import (
	"houwang/mr/lib"
	"fmt"
	"hash/fnv"
	"log"
	"net/rpc"
	"sort"
	"time"
	"math/rand"
)

//
// Map functions return a slice of KeyValue.
//
type KeyValue struct {
	Key   string
	Value string
}

// for sorting by key.
type ByKey []KeyValue

// for sorting by key.
func (a ByKey) Len() int           { return len(a) }
func (a ByKey) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByKey) Less(i, j int) bool { return a[i].Key < a[j].Key }


type MrWorker struct {
	// MrWorker does not track its own state
	task    lib.Task
	mapf    func(string, string) []KeyValue
	reducef func(string, []string) string
}

func getPartitionNumber(key string, reduceN int) int {
	hash := ihash(key)
	return hash % reduceN
}

func (this *MrWorker) runMap() {
  // implement map
  inputFile := this.task.mapTask.inputFile
  inputFileio := lib.CreateFileIO(inputFile.fileName)
  content, err := inputFileio.ReadAll()
  if err != nil {
	  log.Fatal("Failed to read the input file: ", inputFile.fileName)
  }

  // execute mapf
  intermediateKV := this.mapf(inputFile.fileName, content)

  // Write intermediateKV to storage
  outputFiles := this.task.mapTask.outputFiles
  reduceN := len(outputFiles)
  for keyval : intermediateKV {
	outputFileName := outputFiles[getPartitionNumber(keyval.Key, reduceN)]
	outputFileIO = lib.CreateFileIO(outputFileName)
	outStr := fmt.Sprintf("%v %v", keyval.Key, keyval.Value)
	outputFileIO.AppendString(outStr)
  }
}

func (this *MrWorker) runReduce() {
  // implement reduce
  inputFile := this.task.reduceTask.inputFile
  fileio := lib.CreateFileIO(inputFile.fileName)

  var intermediate []KeyValue
  // read key-vals into mem
  for line, err := fileio.ReadLine(); err != nil {
	  keyVal := KeyValue{}
	  fmt.Sscanf(line, "%v %v", &keyVal)
	  append(intermediate, keyVal)
  }

  // sort
  sort.Sort(ByKey(intermediate))

  // run reducef
  outputFile := this.task.reduceTask.outputFile
  outputFileIO := fileio.CreateFileIO(outputFile.fileName)
  i := 0
	for i < len(intermediate) {
		j := i + 1
		for j < len(intermediate) && intermediate[j].Key == intermediate[i].Key {
			j++
		}
		values := []string{}
		for k := i; k < j; k++ {
			values = append(values, intermediate[k].Value)
		}
		output := reducef(intermediate[i].Key, values)

		// this is the correct format for each line of Reduce output.
		outputFileIO.AppendString(fmt.Sprintf("%v %v\n", intermediate[i].Key, output))

		i = j
	}
  
}

func (this *MrWorker) ExecTask(task lib.Task) {
	if !this.task.isIdle() {
		log.Println("Task is already running: ", this.task)
		return
	}
	this.task = task
	this.task.taskState = lib.TaskState.InProgress

	switch this.task.taskType {
	case lib.Map:
		this.runMap()
	case lib.Reduce:
		this.runReduce()
	}
	this.task.taskState = lib.TaskState.Completed
	return
}

func (this *MrWorker) GetTaskState() lib.TaskState {
	return this.task.taskState
}

func createMrWorker(mapf func(string, string) []KeyValue,
	reducef func(string, []string) string) *MrWorker {
	mrWorker := new(MrWorker)
	mrWorker.mapf = mapf
	mrWorker.reducef = reducef
	return mrWorker
}

//
// use ihash(key) % NReduce to choose the reduce
// task number for each KeyValue emitted by Map.
//
func ihash(key string) int {
	h := fnv.New32a()
	h.Write([]byte(key))
	return int(h.Sum32() & 0x7fffffff)
}

//
// main/mrworker.go calls this function.
//
func Worker(mapf func(string, string) []KeyValue,
	reducef func(string, []string) string) {
	mrWorker := createMrWorker(mapf, reducef)
	
	// Worker join the job, query master for a task
	for {
		var assignTaskResp AssignTaskResponse
		succeed := call("Master.AssignTask", nil, &assignTaskResp)
		
		if succeed {
			mrWorker.ExecTask(assignTaskResp.task)
			// report finished
			var taskStateRequest UpdateTaskStateRequest
			taskStateRequest.taskState = mrWorker.GetTaskState()
			call("Master.UpdateTaskState", nil, &taskStateRequest)
		}

		rand.Seed(time.Now().UnixNano())
		sleepTime := rand.Intn(4)
		time.Sleep(time.Duration(sleepTime)*time.Second)
	}
    
	

//
// send an RPC request to the master, wait for the response.
// usually returns true.
// returns false if something goes wrong.
//
func call(rpcname string, args interface{}, reply interface{}) bool {
	// c, err := rpc.DialHTTP("tcp", "127.0.0.1"+":1234")
	sockname := masterSock()
	c, err := rpc.DialHTTP("unix", sockname)
	if err != nil {
		log.Fatal("dialing:", err)
	}
	defer c.Close()

	err = c.Call(rpcname, args, reply)
	if err == nil {
		return true
	}

	fmt.Println(err)
	return false
}
