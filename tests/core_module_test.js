var fs = require('fs');

var datastr = fs.readFileSync('tests/data.json');
data = JSON.parse(datastr);

console.log(datastr);

module.exports = data[0].name;