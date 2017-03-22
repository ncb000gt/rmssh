const exec = require('ssh-exec');
const Stream = require('stream');
const colors = require('colors');
const read = require('read');
const argv = require('yargs').argv;

let user = argv.user;
let password = argv.password;
const filter = argv.filter;
const commands = argv._;

function genStream(server) {
	const writer = new Stream();
	writer.writable = true;

	writer.write = function(buf) {
		const  bufStr = buf.toString();
		const snl = bufStr.split('\n');

		snl.forEach(line => {
			if ((filter && line.indexOf(filter) >= 0) || !filter) {
				console.log(colors.green(server), ' : ', line);
			}
		});
	}

	return writer;
}


function listenToServers() {
	commands.forEach(command => {
		var ssplit = command.split(' ');
		var server = ssplit[0];
		exec(ssplit.slice(1).join(' '), {host: server, user: user, password: password}).pipe(genStream(server));
	});
}

read({prompt: 'Username: '}, (_, un) => {
	user = un;

	read({prompt: 'Password: ', silent: true}, (_, pw) => {
		password = pw;

		listenToServers();
	});
});
