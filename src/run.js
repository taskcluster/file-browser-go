const
  child_process = require('child_process'),
  FileBrowser   = require('./FileBrowser');

try{
  let shell = child_process.spawn('../file-browser-go', [] , {
    stdio: ['pipe','pipe', 'ignore'],
    detached: true
  });
}catch(e){
  console.error(e);
}
