const core = require('@actions/core');
const exec = require('@actions/exec');
const fs = require('fs');
const hasbin = require('hasbin');
const tc = require('@actions/tool-cache');

async function setup() {

	if (hasbin.sync('kubediff')) {
		return;
	}
	const kubeDiffVersion = fs.readFileSync('.version', 'utf8');

	let cachedKubeDiffPath = tc.find('kubediff', kubeDiffVersion);
	if (cachedKubeDiffPath) {
		core.addPath(cachedKubeDiffPath);
		return;
	}

	const kubeDiffDownloadUrl = `https://github.com/yoo/kubediff-action/releases/download/v${kubeDiffVersion}/kubediff`;

	const kubeDiffPath = await tc.downloadTool(kubeDiffDownloadUrl);
	fs.chmodSync(kubeDiffPath, '775');
	cachedKubeDiffPath = await tc.cacheDir(kubeDiffPath, 'kubediff', kubeDiffVersion);
	core.addPath(cachedKubeDiffPath);
}

async function run() {
	await setup();

	let args = []
	const context = core.getInput('context');
	if (context != undefined) {
		args.push(`--context ${process.env.INPUT_CONTEXT}`);
	}

	const lines = core.getInput('manifests')
	const manifests = lines.split('\n').filter(function (el) {
		return el != "";
	});

	args.push('diff');
	for (const index in manifests) {
		args.push('-f');
		args.push(manifests[index]);
	};

	console.log(args)
	core.exportVariable('KUBECTL_EXTERNAL_DIFF', 'kubediff');
	await exec.exec('kubectl', args);
}

run();
