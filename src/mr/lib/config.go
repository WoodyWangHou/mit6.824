package lib

type JobConfig struct {
	InputFiles []string
	SplitFileBytes int
	WorkerHeartBeatIntervalSec int
	ReducePartition int
}