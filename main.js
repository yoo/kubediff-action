const core = require('@actions/core');
const exec = require('@actions/exec');
const fs = require('fs');
const hasbin = require('hasbin');
const tc = require('@actions/tool-cache');

async function setup() {
	if (hasbin.sync('kubediff')) {
		return;
	}
	const kubeDiffVersion = fs.readFileSync('.version', 'utf8').trim();

	let cachedKubeDiffPath = tc.find('kubediff', kubeDiffVersion);
	if (cachedKubeDiffPath) {
		core.addPath(cachedKubeDiffPath);
		return;
	}

	const kubeDiffDownloadUrl = `https://github.com/yoo/kubediff-action/releases/download/v${kubeDiffVersion}/kubediff`;
	console.log(`download kubediff from ${kubeDiffDownloadUrl}`)

	let kubeDiffPath = '';
	try {
		kubeDiffPath = await tc.downloadTool(kubeDiffDownloadUrl);
	} catch (error) {
		console.log(error);
		return core.setFailed('failed to download kubediff');
	}
	fs.chmodSync(kubeDiffPath, '775');

	try {
		cachedKubeDiffPath = await tc.cacheFile(kubeDiffPath, 'kubediff', 'kubediff', kubeDiffVersion);
	} catch (error) {
		console.log(error)
		return core.setFailed('failed to cache kubediff');
	}
	core.addPath(cachedKubeDiffPath);
}

async function run() {
	await setup();

	let args = []
	const context = core.getInput('context');
	if (context != '') {
		args.push(`--context ${context}`);
	}

	const lines = core.getInput('manifests');
	const manifests = lines.split('\n').filter(function (el) {
		return el != "";
	});

	args.push('diff');
	for (const index in manifests) {
		args.push('-f');
		args.push(manifests[index]);
	};

	console.log(args);
	core.exportVariable('KUBECTL_EXTERNAL_DIFF', 'kubediff');
	try {
		await exec.exec('kubectl', args);
	} catch (error) {
		console.log(error)
		core.setFailed('failed to run kubectl diff');
	}
}

run();
