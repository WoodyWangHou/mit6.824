package mr

import (
	"houwang/mr/lib"
	"io/ioutil"
	"os"
	"testing"
)

const TestFilePath string = "/tmp/test_file"
const TestData string = "test data"

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}

func setup() {
	file, _ := os.Create(TestFilePath)
	file.WriteString(TestData)
	file.Close()
}

func teardown() {
	os.Remove(TestFilePath)
}

func reset() {
	teardown()
	setup()
}

func TestReadLine(t *testing.T) {
	fileio := lib.CreateFileIO(TestFilePath)
	defer fileio.Close()
	line, _ := fileio.ReadLine();
	if line != TestData {
		t.Fatal("the reader content: ", line, " - does not matches the test data: ", TestData)
	}
}

func TestWriteLine(t *testing.T) {
	fileio := lib.CreateFileIO(TestFilePath)
	additional := " additional" 
	fileio.AppendString(additional)
	fileio.Close()
    
	newFile, _ := os.Open(TestFilePath)
	newData, err := ioutil.ReadAll(newFile)
	if err != nil {
		t.Fatal("failed to open the test file")
	}

	if string(newData) != "test data additional" {
		t.Fatal("Failed to write to file, file content: ", string(newData))
	}
}

func TestReadAll(t *testing.T) {
	reset()
	fileio := lib.CreateFileIO(TestFilePath)
	additional := " additional" 
	fileio.AppendString(additional)
	fileio.Close()

    newFileio := lib.CreateFileIO(TestFilePath)
	line, _ := newFileio.ReadAll();
	if line != TestData + additional {
		t.Fatal("the reader content: ", line, " - does not matches the test data: ", TestData)
	}
}
