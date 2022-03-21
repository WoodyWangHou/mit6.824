package mr

import (
	"houwang/mr/lib"
	"os"
	"testing"
	"houwang/mr"
	"unicode"
	"strings"
	"strconv"
	"bufio"
	"log"
)

const WorkerTestFilePath string = "/tmp/worker_test_file"
const WorkerTestOutputPath string = "/tmp/worker_output_file"
const WorkerTestData string = "count count count"

func Map(filename string, contents string) []mr.KeyValue {
	// function to detect word separators.
	ff := func(r rune) bool { return !unicode.IsLetter(r) }

	// split contents into an array of words.
	words := strings.FieldsFunc(contents, ff)

	kva := []mr.KeyValue{}
	for _, w := range words {
		kv := mr.KeyValue{w, "1"}
		kva = append(kva, kv)
	}
	return kva
}

func Reduce(key string, values []string) string {
	// return the number of occurrences of this word.
	return strconv.Itoa(len(values))
}

func WorkerTestMain(m *testing.M) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	WorkerSetup()
	code := m.Run()
	WorkerTeardown()
	os.Exit(code)
}

func WorkerSetup() {
	file, _ := os.Create(WorkerTestFilePath)
	file.WriteString(WorkerTestData)
	defer file.Close()

	outputFile, _ := os.Create(WorkerTestOutputPath)
	defer outputFile.Close()
}

func WorkerTeardown() {
	os.Remove(WorkerTestFilePath)
	os.Remove(WorkerTestOutputPath)
}

func TestNonIdleTask(t *testing.T) {
	mrWorker := mr.CreateMrWorker(Map, Reduce)
	task := lib.Task{ TaskState : lib.InProgress }
	err := mrWorker.ExecTask(task)
	if err == nil {
		t.Fatal("Task should be aborted")
	}
}

func TestMapWorker(t *testing.T) {
	mrWorker := mr.CreateMrWorker(Map, Reduce)
	task := lib.Task{
		TaskState: lib.Idle,
		TaskType: lib.Map,
		MapTask: lib.MapTask{
			InputFile: WorkerTestFilePath,
		},
	}
    task.MapTask.OutputFiles = append(task.MapTask.OutputFiles, WorkerTestOutputPath)
	
	mrWorker.ExecTask(task)

    fileHandle, _ := os.Open(WorkerTestOutputPath)
	defer fileHandle.Close()
	fileScanner := bufio.NewScanner(fileHandle)

	for fileScanner.Scan() {
	    line := fileScanner.Text()
		if line != "count 1" {
			t.Fatal("Map job executed wrongly")
		}
	}
}
