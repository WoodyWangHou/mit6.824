package mr

import (
	"houwang/mr/lib"
	"fmt"
	"hash/fnv"
	"log"
	"net/rpc"
)

//
// Map functions return a slice of KeyValue.
//
type KeyValue struct {
	Key   string
	Value string
}

type MrWorker struct {
	// MrWorker does not track its own state
	task    lib.Task
	mapf    func(string, string) []KeyValue
	reducef func(string, []string) string
}

func (this *MrWorker) runMap() {
  // implement map
  inputFile := this.task.inputFile
  outputFile := this.task.outputFile
  fileio := lib.CreateFileIO(inputFile.fileName)
  content, err := fileio.ReadAll()
  if err != nil {
	  log.Fatal("Failed to read the input file: ", inputFile.fileName)
  }

  // execute mapf
  intermediate := this.mapf(inputFile.fileName, content)

  // Write intermediate to storage
  outputFileIO = lib.CreateFileIO(outputFile.fileName)
  for keyval : intermediate {
	outStr := fmt.Sprintf("%v %v", keyval.Key, keyval.Value)
	outputFileIO.AppendString(outStr)
  }
}

func (this *MrWorker) runReduce() {
  // implement reduce
}

func (this *MrWorker) setTask(task lib.Task) {
	if !task.isIdle() {
		this.task = task
	}
}

func (this *MrWorker) exec() {
	if !this.task.isIdle() {
		log.Println("Task is already running: ", this.task)
		return
	}

	switch this.task.taskType {
	case lib.Map:
		this.runMap()
	case lib.Reduce:
		this.runReduce()
	}
	return
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

	// Your worker implementation here.

	// uncomment to send the Example RPC to the master.
	mrWorker := createMrWorker(mapf, reducef)

}

//
// example function to show how to make an RPC call to the master.
//
// the RPC argument and reply types are defined in rpc.go.
//
// func CallExample() {

// 	// declare an argument structure.
// 	args := ExampleArgs{}

// 	// fill in the argument(s).
// 	args.X = 99

// 	// declare a reply structure.
// 	reply := ExampleReply{}

// 	// send the RPC request, wait for the reply.
// 	call("Master.Example", &args, &reply)

// 	// reply.Y should be 100.
// 	fmt.Printf("reply.Y %v\n", reply.Y)
// }

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
