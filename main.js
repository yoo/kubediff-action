const spawnSync = require('child_process').spawnSync;

if (process.env.INPUT_CONTEXT != undefined) {
	args.push(`--context ${process.env.INPUT_CONTEXT}`);
}

const lines = process.env.INPUT_MANIFESTS.split('\n');

const manifests = lines.filter(function (el) {
	return el != "";
});

let args = ['diff'];
for (const index in manifests) {
	args.push('-f');
	args.push(manifests[index]);
};

console.log(args)
process.env.KUBECTL_EXTERNAL_DIFF = 'kubediff';
const proc = spawnSync('kubectl', args, {stdio: 'inherit'});
process.exit(proc.status);
