'use strict';
/* API for all file-browser commands */

const
  assert    = require('assert'),
  debug     = require('debug')('browser'),
  Buffer    = require('buffer').Buffer,
  msp       = require('msgpack-lite'),
  through2  = require('through2'),
  Promise   = require('Promise'),
  fs        = require('fs');

const
  Command         = require('./Command.js'),
  FileOperations  = require('./FileOperations.js'),
  EventEmitter    = require('events').EventEmitter,
  StringDecoder   = require('string_decoder').StringDecoder;

const decoder = new StringDecoder();

class FileBrowser extends EventEmitter {
  
  constructor(shell) {

    super();

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

    this.stdout.on('data', (data) => {
      if(!data.id) return;
      debug('Received: ', data);
      this.emit(data.id, data);
    });

    this.stdout.on('error', debug);
    this.stdin.on('error', debug);

    let self = this;

    ["mv", "cp"].forEach(c => {

      self[c] = (src, dest) => {
        let cmd = Command[c](src, dest);
        return self.writeAndResolve(cmd);
      }

    });

    ["ls", "rm", "mkdir"].forEach(c => {

      self[c] = (path) => {
        let cmd = Command[c](path);
        return self.writeAndResolve(cmd);
      }

    });
  }

  
  writeAndResolve (cmd) {
    let self = this;
    return new Promise(resolve => {
      self.stdin.write(cmd);
      return self.once(cmd.id, resolve);
    });
  }

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
        result = await this.writeAndResolve(c);
        
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

    let value = await new Promise(resolve => {
      this.stdin.write(cmd);
      this.on(cmd.id, res => {
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
      });
    });

    this.removeAllListeners(cmd.id);
    debug('Removed listener for getfile id: ', cmd.id);
    return value;

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
