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

func MasterTestMain(m *testing.M) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	MasterSetup()
	code := m.Run()
	MasterTeardown()
	os.Exit(code)
}

func MasterSetup() {
	file, _ := os.Create(WorkerTestFilePath)
	file.WriteString(WorkerTestData)
	defer file.Close()

	outputFile, _ := os.Create(WorkerTestOutputPath)
	defer outputFile.Close()
}

func MasterTeardown() {
	os.Remove(WorkerTestFilePath)
	os.Remove(WorkerTestOutputPath)
}

func TestTaskStore(t *testing.T) {
	taskStore := mr.TaskStore{}
}
