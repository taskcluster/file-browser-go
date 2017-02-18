package browser

import (
	"sync"
)

var registryLock sync.Mutex

var cmdRegistry map[string](func(Command, chan<- *ResultSet)) = make(map[string]func(Command, chan<- *ResultSet))

func RegisterCommand(command string, fun func(Command, chan<- *ResultSet)) {
	registryLock.Lock()
	_, ok := cmdRegistry[command]
	if ok {
		panic("Command already registered")
	}
	cmdRegistry[command] = fun
	registryLock.Unlock()
}

func RunCommand(cmd Command, outChan chan<- *ResultSet) {
	if cmd.Id == "" {
		outChan <- FailedResultSet("", "No id specified")
		return
	}
	fun, registered := cmdRegistry[cmd.Cmd]
	if !registered {
		panic("Command not registered")
	}
	fun(cmd, outChan)
}
