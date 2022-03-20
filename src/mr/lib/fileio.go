package lib

import (
	"bufio"
	"log"
	"os"
	"errors"
	"bytes"
	"fmt"
	"sync"
)

type FileIO struct{
	file *os.File
	scanner *bufio.Scanner
	writer *bufio.Writer
	writerLock sync.Mutex
}

func CreateFileIO(fileLocation string) FileIO {
	fileIO := FileIO{}
    file, err := os.OpenFile(fileLocation, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		log.Println("Failed to open the file: ", err.Error())
		return fileIO
	}
	fileIO.file = file
	fileIO.scanner = bufio.NewScanner(file)
	fileIO.writer = bufio.NewWriter(file)
	return fileIO
}

// This function return current line and advance the scanner to next token, returns error
// if it is currently the last token
func (this *FileIO) ReadLine() (string, error) {
	if this.scanner == nil {
		return "", errors.New("Failed to read a file, FileIO does not have a reader")
	}
	hasMore := this.scanner.Scan()
	text := this.scanner.Text()
	if !hasMore {
		return "", this.scanner.Err()
	}

	return text, nil
}

func (this *FileIO) ReadAll() (string, error) {
	if this.scanner == nil {
		return "", errors.New("Failed to read a file, FileIO does not have a reader")
	}
	var outputBuf bytes.Buffer
	for this.scanner.Scan() {
		outputBuf.WriteString(this.scanner.Text())
	}

	if err := this.scanner.Err(); err != nil {
		fmt.Println("shouldn't see an error scanning a string: ", err.Error())
		return "", err
	}
    // reset scanner
	this.scanner = bufio.NewScanner(this.file)
	return outputBuf.String(), nil
}

func (this *FileIO) AppendString(content string) {
	if this.writer == nil {
		log.Println("FileIO does not have a writer, ignore appendString")
		return
	}
    mu := this.writerLock.Lock()
	defer mu.Unlock()

	writtenLen, err := this.writer.WriteString(content)
	if writtenLen != len(content) || err != nil {
		log.Println("Failed to write to the file. Written bytes: ", writtenLen, "; err: ", err.Error())
	}
}

func (this *FileIO) Close() {
	this.writer.Flush()
	this.file.Close()
}
