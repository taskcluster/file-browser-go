'use strict';

const
  EventEmitter    = require('events').EventEmitter,
  StringDecoder   = require('string_decoder').StringDecoder,
  _               = require('lodash'),
  FileOperations  = require('./FileOperations');

const decoder = new StringDecoder();

class Registry extends EventEmitter {

  constructor(outputStream){
    super();

    this.outputStream = outputStream;
    this.register = [];

    this.curBuff = "";
    this.pairCount = 0;
    this.curBuffIndex = 0;

    this.outputStream.on('data', data => {
      this.curBuff += decoder.write(data);

      while(this.curBuffIndex < this.curBuff.length){
        if(this.curBuff[this.curBuffIndex] === '{') 
          this.pairCount++;
        else if(this.curBuff[this.curBuffIndex] === '}') 
          this.pairCount--;

        this.curBuffIndex++;

        if(this.pairCount === 0 && this.curBuff[this.curBuffIndex - 1] === '}'){
          var i = 0;
          var residue = "";
          var tempBuff = "";

          while(i < this.curBuffIndex){
            tempBuff += this.curBuff[i]; 
            i++;
          }
          while(i < this.curBuff.length){
            residue += this.curBuff[i];
            i++;
          }

          this.curBuff = residue;
          this.curBuffIndex = 0;

          this.processString(tempBuff);
        }
      }

    });

    this.outputStream.on('end', () => {
      if(this.curBuff === "") return;
      this.processString(this.curBuff);
    });

  }

  processString(chunk){
    chunk = chunk.trim();
    if (chunk === "") {
      return;
    }
    var json = {};
    try{
      json = JSON.parse(chunk); 
    }catch(err){
      console.log("Chunk: ", chunk);
      console.error(err);
    }  
    if(!json.id || !json.cmd) {
      this.emit('error', json);
      return;
    }

    if(_.indexOf(this.register, json.id) === -1){
      this.emit('error', json);
      return;
    }
    _.remove(this.register, (n) => { n == json.id });
    // Not the last chunk returned by getfile
    this.emit(json.cmd, json);
    return;
  }

  registerCommand(id){
    this.register.push(id);
  }

}

module.exports = Registry;
