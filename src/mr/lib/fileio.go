package lib

import (
	"bufio"
	"log"
	"os"
)

type FileIO struct{}

// The functions returns nil if there is any issue opening / writing file
func (this *FileIO) ReadFileBuffered(fileLocation string) (*bufio.Reader, *os.File) {
	file, err := os.Open(fileLocation)
	if err != nil {
		log.Println("Failed to open file at: ", fileLocation)
		return nil, nil
	}
	return bufio.NewReader(file), file
}

func (this *FileIO) WriteFileBuffered(fileLocation string) (*bufio.Writer, *os.File) {
	file, err := os.OpenFile(fileLocation, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		if os.IsNotExist(err) {
			createdFile, err := os.Create(fileLocation)
			if err != nil {
				log.Println("Cannot create file: ", fileLocation, ". With err: ", err.Error())
				return nil, nil
			}
			return bufio.NewWriter(createdFile), file
		}

		log.Println("Cannot create or write to file: ", fileLocation)
		return nil, nil
	}

	return bufio.NewWriter(file), file
}
