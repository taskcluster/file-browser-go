'use strict';
/* API for all file-browser commands */

const
  through2  = require('through2'),
  assert     = require('assert'),
  fs        = require('fs');

const
  Command         = require('./Command.js'),
  FileOperations   = require('./FileOperations.js'),
  Registry         = require('./Registry.js'),
  StringDecoder   = require('string_decoder').StringDecoder;

const decoder = new StringDecoder();

class FileBrowser {
  
  constructor(shell) {
    assert(shell);
    assert(shell.stdin);
    assert(shell.stdout);

    this.shell = shell;   

    this.stdin = through2({
      objectMode: true
    }, (chunk, encoding, cb) => {
      cb(null, JSON.stringify(chunk));
    });

    this.stdin.pipe(this.shell.stdin);

    this.registry = new Registry(this.shell.stdout);

    let self = this;

    ["mv", "cp"].forEach(c => {

      self[c] = (src, dest) => {
        let cmd = Command[c](src, dest);
        self.registry.registerCommand(cmd.id);
        self.stdin.write(cmd);
        return self.identifyAndResolve(cmd.cmd, cmd.id);
      }

    });

    ["ls", "rm", "mkdir"].forEach(c => {

      self[c] = (path) => {
        let cmd = Command[c](path);
        self.registry.registerCommand(cmd.id);
        self.stdin.write(cmd);
        return self.identifyAndResolve(cmd.cmd, cmd.id);
      }

    });
  }

  
  identifyAndResolve (command, id) {
    return new Promise(resolve => {
      return this.registry.on( command, resultSet => {
        let cmd = resultSet.cmd;
        let rid = resultSet.id;

        if(cmd == command && id == rid ) {
          return resolve(resultSet);
        }
      });
    });
  }
/*
  async ls(path) {
    let command = Command.ls(path);
    this.registry.registerCommand(command.id);
    this.stdin.write(command);
    return this.identifyAndResolve("ls", command.id);
  }

  async rm (path) {
    let cmd = Command.rm(path);
    this.registry.registerCommand(cmd.id);
    this.stdin.write(cmd);
    return this.identifyAndResolve("rm", cmd.id);
  }

  async mv (src, dest) {
    let command = Command.mv(src, dest);
    this.registry.registerCommand(command.id);
    this.stdin.write(command);
    return this.identifyAndResolve("mv", command.id);
  }

  async mkdir (path) {
    let command = Command.mkdir(path);
    this.registry.registerCommand(command.id);
    this.stdin.write(command);
    return this.identifyAndResolve("mkdir", path);
  }

  async cp (src, dest) {
    let command = Command.cp(src, dest);
    this.registry.registerCommand(command.id);
    this.stdin.write(command);
    return this.identifyAndResolve("cp", command.id);
  }
*/
  async putfile (src , dest) {  

    let cmd = [], fail = false;
    let result = {};

    while(cmd.len < 2 && !fail){
      cmd = FileOperations.putfile(src, dest);
      if(cmd === null) return {
        error: "Error opening file for reading"
      };

      for (let i in cmd){

        let c = cmd[i];
        this.registry.registerCommand(c);
        this.stdin.write(c);
        result = await identifyAndResolve(c.cmd, c.id);
        
        if (result.error !== "") {
          fail = true;
          break;
        }
      }
    }

    FileOperations.putfileClean(dest);
    return result;

  }

  async getfile (src , dest) {

    let cmd = Command.getfile(src); 
    let block, totalPieces, i;

    this.registry.registerCommand(cmd.id);
    this.stdin.write(command);
    
    block = await this.identifyAndResolve(cmd.cmd, cmd.id);

    if (block.error !== "" ){
      console.error(block);
      return null;
    }

    totalPieces = initBlock.fileData.totalPieces;
    
    for(i = 0; i < totalPieces; i++) {
      block = await this.identifyAndResolve(cmd.cmd, cmd.id);
      fs.appendFileSync(dest, block.fileData.data);
    }

    return dest;

  }

  async kill () {
    let prom = new Promise((resolve, reject) => {
      this.shell.on('exit', resolve).on('error', reject);
    });
    this.stdin.destroy();
    this.shell.kill();
    return prom;
  }
  
}

module.exports = FileBrowser;
