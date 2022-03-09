package lib

import (
	"bufio"
	"log"
	"os"
	"errors"
)

type FileIO struct{
	file *os.File
	reader *bufio.Scanner
	writer *bufio.Writer
}

func CreateFileIO(fileLocation string) FileIO {
	fileIO := FileIO{}
    file, err := os.OpenFile(fileLocation, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		log.Println("Failed to open the file: ", err.Error())
		return fileIO
	}
	fileIO.file = file
	fileIO.reader = bufio.NewScanner(file)
	fileIO.writer = bufio.NewWriter(file)
	return fileIO
}

// This function return current line and advance the scanner to next token, returns error
// if it is currently the last token
func (this *FileIO) ReadLine() (string, error) {
	if this.reader == nil {
		return "", errors.New("Failed to read a file, FileIO does not have a reader")
	}
	hasMore := this.reader.Scan()
	text := this.reader.Text()
	if !hasMore {
		return "", this.reader.Err()
	}

	return text, nil
}

func (this *FileIO) AppendString(content string) {
	if this.writer == nil {
		log.Println("FileIO does not have a writer, ignore appendString")
		return
	}

	writtenLen, err := this.writer.WriteString(content)
	if writtenLen != len(content) || err != nil {
		log.Println("Failed to write to the file. Written bytes: ", writtenLen, "; err: ", err.Error())
	}
}

func (this *FileIO) Close() {
	this.writer.Flush()
	this.file.Close()
}
