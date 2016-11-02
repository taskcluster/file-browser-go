/* Generates JSON commands based on parameters */
'use strict';

let Command = {};

Command.ls = path => {
	return {
		cmd: "ls",
		args: [ path ]
	};
}

Command.mkdir = path => {
	return {
		cmd: "mkdir",
		args: [ path ]
	};
}

Command.cp = (oldpath, newpath) => {
	return {
		cmd: "cp",
		args: [ oldpath, newpath ]
	};
}

Command.rm= path => {
	return {
		cmd: "rm",
		args: [ path ]
	};
}

Command.mv = (oldpath, newpath) => {
	return {
		cmd: "mv",
		args: [ oldpath, newpath ]
	};
}

Command.getfile = path => {
	return {
		cmd: "getfile",
		args: [path]
	};
}

Command.putfile = (path, data) => {
	return {
		cmd: "putfile",
		args: [path],
		data: data
	};
}

module.exports = Command;
