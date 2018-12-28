// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package service

import (
	"fmt"
	"io/ioutil"
	"sync"
)

var (
	web_readme_once    sync.Once
	web_readme_content string
)

func web_readme_md(readmeFile string) string {
	web_readme_once.Do(func() {
		if data, err := ioutil.ReadFile(readmeFile); err == nil {
			web_readme_content = string(data)
		} else {
			web_readme_content = fmt.Sprintf("%q file missing", readmeFile)
		}
	})
	return web_readme_content
}

func web_homepage() string {
	return web_tmplHomepage
}

const web_tmplHomepage = `<!doctype html>
<html>
<head>
	<meta charset="utf-8"/>
	<title>OpenPitrix IAM Server</title>
</head>
<body>
	<div id="content"></div>
	<script src="/swagger/marked.min.js"></script>
	<script>
		fetch('readme.md').then(response => {
			return response.text()
		}).then(text => {
			document.getElementById('content').innerHTML = marked(text)
		})
	</script>
</body>
</html>
`
