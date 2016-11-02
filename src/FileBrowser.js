/* API for all file-browser commands */

let through2 = require('through2');

let Command = require('./Command.js');
let FileOperations = require('./FileOperations.js');

class FileBrowser {
  
  constructor(shell) {
    this.shell = shell;   

    this.stdin = through2({
      objectMode: true
    }, (chunk, encoding, cb) => {
      cb(null, JSON.stringify(chunk));
    });

    this.stdout = through2((chunk, encoding, cb) => {
      cb(null, JSON.parse(chunk));
    });

    this.stdin.pipe(this.shell.stdin);
    this.shell.stdout.pipe(this.stdout);
  }

// Sync commands
  
  identifyAndResolve (command, path) {
    return new Promise((resolve, reject) => {
      this.stdout.on('data', resultSet => {
        let cmd = resultSet.cmd;
        let resPath = resultSet.path;

        if(cmd == command && resPath == path) {
          return resolve(resultSet);
        }
      });
    });
  }

  async ls(path) {
    let command = Command.ls(path);

    this.stdin.write(command);
    return identifyAndResolve("ls", path);
  }

  async rm (path) {
    this.stdin.write(Command.rm(path));
    return identifyAndResolve("rm", path);
  }

  async mv (oldpath, newpath) {
    this.stdin.write(Command.mv(oldpath, newpath));
    return identifyAndResolve("mv", newpath);
  }

  async mkdir (path) {
    this.stdin.write(Command.mkdir(path));
    return identifyAndResolve("mkdir", path);
  }

  // Async commands
  async cp (oldpath, newpath) {
    this.stdin.write(Command.cp(oldpath, newpath));
    return identifyAndResolve("cp", newpath);
  }

  // getfile and putfile
  
}
