package mr

import (
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"sync"
)

type Coordinator struct {
	// Your definitions here.

	locks map[string]*rwLock
	mutex sync.Mutex
	currentId int

	mapTaskQueue []*MapTask
	mapTaskCompleted int
	mapTaskN int

	reduceTaskQueue []*ReduceTask
	ReduceTaskCompleted int
	ReduceTaskN int

	intermidiate []string

	runningTasks map[string]string
}

type MapTask struct {
	fileName string
}

type ReduceTask struct {
	reduceN int
}


type rwLock struct {
	mutex   sync.Mutex
	readers int
	rLock   bool
	wLock   bool
}

func NewCoordinator(files []string, reduceN int) *Coordinator {
	var mTaskQueue []*MapTask = make([]*MapTask, 0, 255)
	for _, file := range files {
		mTaskQueue = append(taskQueue, &MapTask{
			fileName: file,
		})
	}

	var rTaskQueue []*ReduceTask = make([]*ReduceTask, 0, int(reduceN * 1.2))
	for i := 0; i < reduceN; i++ {
		taskQueue = append(taskQueue, &ReduceTask{
			reduceN: i,
		})
	}

	return &Coordinator{
		runningTasks: make(map[string]string),
		currentId: 0,
		locks: make(map[string]*rwLock),
	}
}

func (m *Coordinator) getLock(fileName string) *rwLock {
	m.mutex.Lock()
	lock, exists := m.locks[fileName]
	if !exists {
		lock = &rwLock{}
		m.locks[fileName] = lock
	}
	m.mutex.Unlock()
	return lock
}

func (c *Coordinator) dequeue() *Task{
	var t = 
}

// Your code here -- RPC handlers for the worker to call.

// an example RPC handler.
//
// the RPC argument and reply types are defined in rpc.go.
func (c *Coordinator) Example(args *ExampleArgs, reply *ExampleReply) error {
	reply.Y = args.X + 1
	return nil
}

func (c *Coordinator) GetWorkerNum(args *GetWorkerNumArgs, reply *GetWorkerNumReply) error{
	c.mutex.Lock()
	reply.id = c.currentId
	c.currentId++
	c.mutex.Unlock()
	return nil
}

// R/W locks
func (c *Coordinator) ReadLock(args *ReadLockArgs, reply *ReadLockReply) error {
	lock := c.getLock(args.fileName)
	if lock.wLock {
		reply.tookLock = false
		return nil
	}

	lock.mutex.Lock()
	reply.tookLock = true
	lock.rLock = true
	lock.readers++
	lock.mutex.Unlock()
	return nil
}

func (c *Coordinator) ReleaseReadLock(args *ReadLockReleaseArgs, reply *ReadLockReleaseReply) error {
	lock := c.getLock(args.fileName)
	if lock.readers <= 0 {
		return nil
	}
	lock.mutex.Lock()
	lock.readers--
	if lock.readers == 0 {
		lock.rLock = false
	}
	lock.mutex.Unlock()
	return nil
}

func (c *Coordinator) WriteLock(args *WriteLockArgs, reply *WriteLockReply) error {
	lock := c.getLock(args.fileName)
	if lock.rLock || lock.wLock {
		reply.tookLock = false
		return nil
	}
	lock.mutex.Lock()
	lock.wLock = true
	reply.tookLock = true
	lock.mutex.Unlock()
	return nil
}

func (c *Coordinator) ReleaseWriteLock(args *WriteLockReleaseArgs, reply *WriteLockReleaseReply) error {
	lock := c.getLock(args.fileName)
	lock.mutex.Lock()
	lock.wLock = false
	lock.mutex.Unlock()
	return nil
}

// Task Allocation

func (c *Coordinator) AskTask(args *AskTaskArgs, reply *AskTaskReply) error {

	if c.mapTaskItr < len(c.mapTasks) {
		c.mutex.Lock()
		reply.taskType = "map"
		reply.target = c.mapTasks[c.mapTaskItr]
		c.mapTaskItr++
		c.mutex.Unlock()
		return nil
	}

	if c.reduceTaskItr < len(c.reduceTasks){
		c.mutex.Lock()
		reply.taskType = "reduce"
		reply.target = string(c.reduceTasks[c.reduceTaskItr])
		c.reduceTaskItr++
		c.mutex.Unlock()
		return nil
	}

	return nil
}

func (c *Coordinator) RegisterIntermidiate(args *RegisterIntermidiateArgs, reply *ReadLockReleaseReply) error{
	c.mutex.Lock()
	c.reduceTasks = append(c.reduceTasks, args.fileName)
	c.mutex.Unlock()
	return nil
}

// start a thread that listens for RPCs from worker.go
func (c *Coordinator) server() {
	rpc.Register(c)
	rpc.HandleHTTP()
	//l, e := net.Listen("tcp", ":1234")
	sockname := coordinatorSock()
	os.Remove(sockname)
	l, e := net.Listen("unix", sockname)
	if e != nil {
		log.Fatal("listen error:", e)
	}
	go http.Serve(l, nil)
}

// main/mrcoordinator.go calls Done() periodically to find out
// if the entire job has finished.
func (c *Coordinator) Done() bool {
	ret := false

	// Your code here.

	return ret
}

// create a Coordinator.
// main/mrcoordinator.go calls this function.
// nReduce is the number of reduce tasks to use.
func MakeCoordinator(files []string, nReduce int) *Coordinator {
	c := Coordinator{}

	// Your code here.

	c.server()
	return &c
}
