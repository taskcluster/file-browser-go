'use strict';
/* API for all file-browser commands */

const
  assert    = require('assert'),
  debug     = require('debug')('browser'),
  Buffer    = require('buffer').Buffer,
  msp       = require('msgpack-lite'),
  through2  = require('through2'),
  Promise   = require('Promise'),
  lock      = require('lock')(),
  _         = require('lodash'),
  slugid    = require('slugid'),
  fs        = require('fs');

const
  Command         = require('./Command.js');
  // FileOperations  = require('./FileOperations.js');
  // StringDecoder   = require('string_decoder').StringDecoder;

const CHUNKSIZE = 4072;

const lck = key => new Promise(resolve => lock(key, resolve));

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

    // Table that maps from id to resolver
    this._cb = {}

    this.stdout = msp.createDecodeStream();
    // this.testOut = msp.createDecodeStream();

    this.stdin.pipe(this.shell.stdin);
    // this.stdin.pipe(this.testOut);
    this.shell.stdout.pipe(this.stdout);

    // this.stdout.setMaxListeners(0);

    // this.testOut.on('data', (data) => {
    //   debug("Wrote: ", data);
    // });

    this.stdout.on('data', (data) => {
      if(!data.id) return;
      debug('Received: ', data);
      if (!this._cb[data.id]) {
        debug("message for unknown id: ", data.id);
        return;
      }

      this._cb[data.id](data);
    });

    this.stdout.on('error', debug);
    this.stdin.on('error', debug);

    let self = this;

    ["mv", "cp"].forEach(c => {

      self[c] = async (src, dest) => {
        let cmd = Command[c](src, dest);
        let result = await self.writeAndResolve(cmd);
        debug('Received: ', data);
        delete this._cb[cmd.id];
        return result;
      }

    });

    ["ls", "rm", "mkdir"].forEach(c => {

      self[c] = async (path) => {
        let cmd = Command[c](path);
        let result = await self.writeAndResolve(cmd);
        debug('Received: ', data);
        delete this._cb[cmd.id];
        return result;
      }

    });
  }

  
  writeAndResolve (cmd) {
    let self = this;
    return new Promise(resolve => {
      self.stdin.write(cmd);
      this._cb[cmd.id] = resolve;
    });
  }

  putfile (srcStream , dest) {  

    let self = this;

    return new Promise((resolve, reject) => {

      // Lock to guarantee chunks are written in the correct order
      let lk_id = slugid.v4();
      // debug('Lock id:', lk_id);

      srcStream.on('error', err => {
        debug(err);
        reject(err);
      });

      srcStream.on('data', async data => {
        let unlock = await lck(lk_id);
        let fail = false;
        try{

          if (typeof data === 'string') {
            data = Buffer.from(data);
          }

          let chunks = _.chunk(data.toJSON().data, CHUNKSIZE);

          for (let i in chunks) {
            let ch = Buffer.from(chunks[i]);
            let cmd = Command.putfile(dest, ch);
            // debug(cmd);
            let result = await self.writeAndResolve(cmd);
            // debug(result);
            if (result.error != '') {
              fail = true;
              break;
            }
          }

        }finally{
          unlock();
          if (fail) {
            reject('Operation failed');
          }
        }
      });

      srcStream.on('end', resolve);
    });

  }

  getfile (src , outStream) {

    let cmd = Command.getfile(src); 
    let self = this;

    return new Promise(resolve => {
      self.stdin.write(cmd);
      self._cb[cmd.id] = (data) => {
        debug('Getfile :', data);
        if (data.error != '') {
          return resolve(false);
        }
        if (data.fileData.currentPiece == 0) {
          return;
        }
        outStream.write(data.fileData.data);
        if (data.fileData.totalPieces == data.fileData.currentPiece) {
          delete self._cb[cmd.id];
          return resolve(true);
        }
      }
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
