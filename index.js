const request = require('request');
  
request('https://www.baidu.com/ding1', (err, res, body) => {
 if (err) { return console.log(err); }
 console.log(body.url);
 console.log(body.explanation);
});
