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

func TestReadBuffered(t *testing.T) {
	fileio := lib.FileIO{}
	reader, file := fileio.ReadFileBuffered(TestFilePath)
	if reader == nil {
		t.Fatal("Failed to open the file")
	}
	testData, _ := reader.ReadString(byte('\n'))
	if testData != TestData {
		t.Fatal("the reader content: ", testData, " - does not matches the test data: ", TestData)
	}
	file.Close()
}

func TestWriteBuffered(t *testing.T) {
	fileio := lib.FileIO{}
	writer, file := fileio.WriteFileBuffered(TestFilePath)
	if writer == nil {
		t.Fatal("Failed to open the file")
	}
	additional := " additional" 
	strSize := len(additional)
	len, err := writer.WriteString(additional)
	if err != nil {
		t.Fatal("Failed to write to file: ", err.Error())
	}
	if len != strSize {
		t.Fatal("The test data is not fully written")
	}
	writer.Flush()
	file.Close()
    
	newFile, _ := os.Open(TestFilePath)
	newData, err := ioutil.ReadAll(newFile)
	if err != nil {
		t.Fatal("failed to open the test file")
	}

	if string(newData) != "test data additional" {
		t.Fatal("Failed to write to file, file content: ", string(newData))
	}
}
