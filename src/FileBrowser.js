'use strict';
/* API for all file-browser commands */

const
  assert    = require('assert'),
  debug     = require('debug')('browser'),
  Buffer    = require('buffer').Buffer,
  msp       = require('msgpack-lite'),
  through2  = require('through2'),
  fs        = require('fs');

const
  Command         = require('./Command.js'),
  FileOperations  = require('./FileOperations.js'),
  StringDecoder   = require('string_decoder').StringDecoder;

const decoder = new StringDecoder();

class FileBrowser {
  
  constructor(shell) {

    assert(shell);
    assert(shell.stdin);
    assert(shell.stdout);

    this.shell = shell;   

    this.stdin = through2.obj(function (chunk, enc, cb) {
      this.push(msp.encode(chunk));
      cb();
    });

    this.stdout = msp.createDecodeStream();
    this.testOut = msp.createDecodeStream();

    this.stdin.pipe(this.shell.stdin);
    this.stdin.pipe(this.testOut);
    this.shell.stdout.pipe(this.stdout);

    this.stdout.setMaxListeners(0);

    this.testOut.on('data', (data) => {
      debug("Wrote: ", data);
    });
    this.stdout.on('data', debug);
    this.stdout.on('error', debug);
    this.stdin.on('error', debug);

    let self = this;

    ["mv", "cp"].forEach(c => {

      self[c] = (src, dest) => {
        let cmd = Command[c](src, dest);
        self.stdin.write(cmd);
        return self.identifyAndResolve(c, cmd.id);
      }

    });

    ["ls", "rm", "mkdir"].forEach(c => {

      self[c] = (path) => {
        let cmd = Command[c](path);
        self.stdin.write(cmd);
        return self.identifyAndResolve(c, cmd.id);
      }

    });
  }

  
  identifyAndResolve (cmd, id) {
    let self = this;
    return new Promise((resolve, reject) => {
      return self.stdout.on('data', resultSet => {
        let rid = resultSet.id;
        if(id == rid && resultSet.cmd == cmd) {
          return resolve(resultSet);
        }
      });
    });
  }
/*
  async ls(path) {
    let command = Command.ls(path);
    this.stdout.registerCommand(command.id);
    this.stdin.write(command);
    return this.identifyAndResolve("ls", command.id);
  }

  async rm (path) {
    let cmd = Command.rm(path);
    this.stdout.registerCommand(cmd.id);
    this.stdin.write(cmd);
    return this.identifyAndResolve("rm", cmd.id);
  }

  async mv (src, dest) {
    let command = Command.mv(src, dest);
    this.stdout.registerCommand(command.id);
    this.stdin.write(command);
    return this.identifyAndResolve("mv", command.id);
  }

  async mkdir (path) {
    let command = Command.mkdir(path);
    this.stdout.registerCommand(command.id);
    this.stdin.write(command);
    return this.identifyAndResolve("mkdir", path);
  }

  async cp (src, dest) {
    let command = Command.cp(src, dest);
    this.stdout.registerCommand(command.id);
    this.stdin.write(command);
    return this.identifyAndResolve("cp", command.id);
  }
*/
  async putfile (src , dest) {  

    let cmd = [], fail = false;
    let result = {};

    while(cmd.length < 2 && !fail){
      cmd = FileOperations.putfile(src, dest);
      if(cmd === null) return {
        error: "Error opening file for reading"
      };

      for (let i in cmd){

        let c = cmd[i];
        this.stdin.write(c);
        result = await this.identifyAndResolve('putfile', c.id);
        
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
    let block, total = 0;

    this.stdin.write(cmd);
    
    return new Promise(resolve => {
      this.stdout.on('data', res => {
        if(res.id != cmd.id) return;
        if(res.error != ''){
          debug(res.error);
          return resolve(null);
        }
        if (res.fileData.currentPiece == 0){
          total = res.fileData.totalPieces;
          return;
        }
        fs.appendFileSync(dest, res.fileData.data);
        if(res.fileData.currentPiece == total){
          return resolve(dest);
        }
      }).on('error', debug);
    });

  }

  async kill () {
    let prom = new Promise((resolve, reject) => {
      this.shell.on('exit', resolve).on('error', reject);
    });
    // this.stdin.destroy();
    this.stdin.end();
    this.stdout.end();
    this.shell.kill();
    return prom;
  }
  
}

module.exports = FileBrowser;
