package browser

import (
	"sync"
)

var registryLock sync.Mutex

var cmdRegistry = make(map[string]func(Command, chan<- *ResultSet))

func RegisterCommand(command string, fun func(Command, chan<- *ResultSet)) {
	registryLock.Lock()
	defer registryLock.Unlock()
	_, ok := cmdRegistry[command]
	if ok {
		panic("Command already registered")
	}
	cmdRegistry[command] = fun
}

func RunCommand(cmd Command, outChan chan<- *ResultSet) {
	if cmd.Id == "" {
		outChan <- FailedResultSet("", "No id specified")
		return
	}
	// Locking only critical section
	registryLock.Lock()
	fun, registered := cmdRegistry[cmd.Cmd]
	registryLock.Unlock()

	if !registered {
		panic("Command not registered")
	}
	fun(cmd, outChan)
}
