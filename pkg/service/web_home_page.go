// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package service

func web_readme_md() string {
	return web_tmplReadme_md
}

func web_homepage() string {
	return web_tmplHomepage
}

const web_tmplReadme_md = `
# IAM帐号和权限管理服务

- [Swagger页面](/swagger)
- [版本信息](/static/version)

`
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
