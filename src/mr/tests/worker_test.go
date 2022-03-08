package mr

import (
    "testing"
    "regexp"
	"mr/worker"
)

import "../mr"
import "unicode"
import "strings"
import "strconv"

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

//
// The reduce function is called once for each key generated by the
// map tasks, with a list of all the values created for that key by
// any map task.
//
func Reduce(key string, values []string) string {
	// return the number of occurrences of this word.
	return strconv.Itoa(len(values))
}

func createTestFile() string {}
func purgeTestFile() {}

func testWorker(t *testing.T) {
  worker := worker.createMrWorker(Map, Reduce)

}