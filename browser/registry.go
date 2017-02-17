package browser

var cmdRegistry map[string](func(Command, chan *ResultSet))

func OnePathWrapper(fun func(string, chan *ResultSet, string)) func(Command, chan *ResultSet) {
	return func(cmd Command, outChan chan *ResultSet) {
		if len(cmd.Args) == 0 {
			outChan <- FailedResultSet(cmd.Id, "No path specified")
			return
		}
		fun(cmd.Id, outChan, cmd.Args[0])
	}
}

func TwoPathWrapper(fun func(string, chan *ResultSet, string, string)) func(Command, chan *ResultSet) {
	return func(cmd Command, outChan chan *ResultSet) {
		if len(cmd.Args) < 2 {
			outChan <- FailedResultSet(cmd.Id, "Not enough arguments")
			return
		}
		fun(cmd.Id, outChan, cmd.Args[0], cmd.Args[1])
	}
}

func PutFileWrapper() func(Command, chan *ResultSet) {
	return func(cmd Command, outChan chan *ResultSet) {
		if len(cmd.Args) == 0 {
			outChan <- FailedResultSet(cmd.Id, "Not enough arguments")
			return
		}
		PutFile(cmd.Id, outChan, cmd.Args[0], cmd.Data)
	}
}

func InitRegistry() {
	cmdRegistry = make(map[string](func(Command, chan *ResultSet)))
	cmdRegistry["ls"] = OnePathWrapper(List)
	cmdRegistry["rm"] = OnePathWrapper(Remove)
	cmdRegistry["mkdir"] = OnePathWrapper(MakeDirectory)
	cmdRegistry["mv"] = TwoPathWrapper(Move)
	cmdRegistry["cp"] = TwoPathWrapper(Copy)
	cmdRegistry["getfile"] = OnePathWrapper(GetFile)
	cmdRegistry["putfile"] = PutFileWrapper()
}

func RunCommand(cmd Command) func(Command, chan *ResultSet) {
	if cmd.Id == "" {
		return func(c Command, outChan chan *ResultSet) {
			outChan <- FailedResultSet("", "No id specified")
		}
	}
	fun, registered := cmdRegistry[cmd.Cmd]
	if !registered {
		panic("Command not present")
	}
	return fun
}
