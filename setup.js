async function run() {
	const fs = require('fs');
	const tc = require('@actions/tool-cache');

	const kubeDiffVersion = fs.readFileSync('.version', 'utf8');

	let cachedKubeDiffPath = tc.find('kubediff', kubeDiffVersion);
	if (cachedKubeDiffPath) {
		core.addPath(cachedKubeDiffPath);
		return
	}

	const kubeDiffDownloadUrl = `https://github.com/yoo/kubediff-action/releases/download/v${kubeDiffVersion}/kubediff`;

	const kubeDiffPath = await tc.downloadTool(kubeDiffDownloadUrl);
	fs.chmodSync(kubeDiffPath, '775');
	cachedKubeDiffPath = await tc.cacheDir(kubeDiffPath, 'kubediff', kubeDiffVersion);
	core.addPath(cachedKubeDiffPath);
}

run();
