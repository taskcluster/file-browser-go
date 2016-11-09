/* Check if it works */
const
  child_process = require('child_process'),
  FileBrowser   = require('../lib/FileBrowser'),
  assert        = require('assert');

describe ('Basic', function(){
  // Basic run

  it('can list contents of a directory', async function(done) {
    let shell = child_process.spawn ('./file-browser-go', [] , {
      stdio: ['pipe', 'pipe', 'ignore'],
      detached: true,
      cwd: process.env.HOME
    });
    let fb = new FileBrowser(shell);
    let contents = await fb.ls('/'); 
    console.log(contents);
    await fb.kill();
  });

});
