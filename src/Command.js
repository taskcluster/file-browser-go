/* Generates JSON commands based on parameters */
'use strict';
let slugid = require('slugid');

let Command = {};

Command.ls = path => {
	return {
		cmd: "ls",
		args: [ path ],
		id: slugid.v4()
	};
}

Command.mkdir = path => {
	return {
		cmd: "mkdir",
		args: [ path ],
		id: slugid.v4()
	};
}

Command.cp = (src, dest) => {
	return {
		cmd: "cp",
		id: slugid.v4(),
		args: [ src, dest ]
	};
}

Command.rm = path => {
	return {
		cmd: "rm",
		id: slugid.v4(),
		args: [ path ]
	};
}

Command.mv = (src, dest) => {
	return {
		cmd: "mv",
		id: slugid.v4(),
		args: [ src, dest ]
	};
}

Command.getfile = path => {
	return {
		cmd: "getfile",
		id: slugid.v4(),
		args: [path]
	};
}

Command.putfile = (path, data) => {
	return {
		cmd: "putfile",
		id: slugid.v4(),
		args: [path],
		data: data
	};
}

module.exports = Command;
