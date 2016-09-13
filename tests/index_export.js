//console.log("load index.js");
var data = require('./data');
var sort = require('./helper.js').sort;

function search(i){
	return sort(data)[i].name;
}
exports.search=search
