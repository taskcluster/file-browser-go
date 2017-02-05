/* Check if it works */
const
  child_process = require('child_process'),
  FileBrowser   = require('../lib/FileBrowser'),
  assert        = require('assert');

let browser, shell;

describe ('Basic', function(){

	before(function(){
    shell = child_process.spawn ('./file-browser-go', [] , {
      stdio: ['pipe', 'pipe', 'ignore']
    });
		browser = new FileBrowser(shell);
	});

  it('can list contents of a directory', function(done) {
		browser.ls('/').then(res => {
			console.log(res);
			done();
		});
  });

});
