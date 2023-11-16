package mr

//
// RPC definitions.
//
// remember to capitalize all names.
//

import (
	"os"
	"strconv"
)

//
// example to show how to declare the arguments
// and reply for an RPC.
//

type ExampleArgs struct {
	X int
}

type ExampleReply struct {
	Y int
}

// Add your RPC definitions here.

type GetWorkerNumArgs struct {
}
type GetWorkerNumReply struct {
	id int
}

type ReadLockArgs struct {
	fileName string
}
type ReadLockReply struct {
	tookLock bool
}

type ReadLockReleaseArgs struct {
	fileName string
}
type ReadLockReleaseReply struct {
}

type WriteLockArgs struct {
	fileName string
}
type WriteLockReply struct {
	tookLock bool
}

type WriteLockReleaseArgs struct {
	fileName string
}
type WriteLockReleaseReply struct {
}

type AskTaskArgs struct {
	id int
}
type AskTaskReply struct {
	taskType string
	target   string
}

type RegisterIntermidiateArgs struct {
	fileName string
}
type RegisterIntermidiateReply struct {
}

type GetIntermidiateArgs struct {
}
type GetIntermidiateReply struct {
	intermidiate []string
}

type MapCompletedArgs struct {
	id int
}
type MapCompletedReply struct {
}

type ReduceCompletedArgs struct {
	id int
}
type ReduceCompletedReply struct {
}

// Cook up a unique-ish UNIX-domain socket name
// in /var/tmp, for the coordinator.
// Can't use the current directory since
// Athena AFS doesn't support UNIX-domain sockets.
func coordinatorSock() string {
	s := "/var/tmp/5840-mr-"
	s += strconv.Itoa(os.Getuid())
	return s
}
