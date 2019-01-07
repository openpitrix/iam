// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package service

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

func init() {
	if appPath, err := os.Executable(); err == nil {
		println(appPath)

		indexPath := filepath.Join(filepath.Dir(appPath), "public", "index.html")
		if data, err := ioutil.ReadFile(indexPath); err == nil {
			web_tmplHomepage = string(data)
		}
		readmeMdPath := filepath.Join(filepath.Dir(appPath), "public", "readme.md")
		if data, err := ioutil.ReadFile(readmeMdPath); err == nil {
			web_tmplReadme_md = string(data)
		}
	}
}

func web_readme_md() string {
	return web_tmplReadme_md
}

func web_homepage() string {
	return web_tmplHomepage
}

var web_tmplReadme_md = `
# OpenPitrix IAM Server

- [Swagger](/static/swagger)
- [Version](/v1.1/version:iam)

`
var web_tmplHomepage = `<!doctype html>
<html>
<head>
	<meta charset="utf-8"/>
	<title>OpenPitrix IAM Server</title>
</head>
<body>
	<div id="content"></div>
	<script src="/static/swagger/marked.min.js"></script>
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
