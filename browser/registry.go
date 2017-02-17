package browser

var cmdRegistry map[string](func (Command, chan interface{}))

func OnePathWrapper (fun func (string, chan interface{}, string)) func (Command, chan interface{}) {
	return func(cmd Command, outChan chan interface{}) {
		if len(cmd.Args) == 0 {
			outChan <- FailedResultSet(cmd.Id, "No path specified")
			return
		}
		fun(cmd.Id, outChan, cmd.Args[0])
	}
}

func TwoPathWrapper (fun func (string, chan interface{}, string, string)) func (Command, chan interface{}) {
	return func(cmd Command, outChan chan interface{}) {
		if len(cmd.Args) < 2 {
			outChan <- FailedResultSet(cmd.Id, "Not enough arguments")
			return
		}
		fun (cmd.Id, outChan, cmd.Args[0], cmd.Args[1])
	}
}

func PutFileWrapper () func (Command, chan interface{}) {
	return func(cmd Command, outChan chan interface{}) {
		if len(cmd.Args) == 0 {
			outChan <- FailedResultSet(cmd.Id, "Not enough arguments")
			return
		}
		PutFile (cmd.Id, outChan, cmd.Args[0], cmd.Data)
	}
}


func InitRegistry() {
	cmdRegistry = make(map[string](func (Command, chan interface{})))
	cmdRegistry["ls"] = OnePathWrapper(List)
	cmdRegistry["rm"] = OnePathWrapper(Remove)
	cmdRegistry["mkdir"] = OnePathWrapper(MakeDirectory)
	cmdRegistry["mv"] = TwoPathWrapper(Move)
	cmdRegistry["cp"] = TwoPathWrapper(Copy)
	cmdRegistry["getfile"] = OnePathWrapper(GetFile)
	cmdRegistry["putfile"] = PutFileWrapper()
}

func RunCommand (cmd string) (func (Command, chan interface{})){
	fun , registered := cmdRegistry[cmd]
	if cmd.Id == "" {
		return func(c Command, outChan chan interface{}){
			outChan <- FailedResultSet("", "No id specified")
		}
	}
	if !registered {
		panic("Command not present")
	}
	return fun
}
