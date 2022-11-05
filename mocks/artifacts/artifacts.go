package artifacts

import (
	_ "embed"
)

//go:embed github-0.0.0.mock.tar.gz
var MockGitHubPluginBundle []byte

//go:embed github-2.0.0.tar.gz.sig
var MockGitHubPluginBundleSig []byte
