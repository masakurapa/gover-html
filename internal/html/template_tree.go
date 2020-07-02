package html

import "html/template"

var parsedTreeTemplate = template.Must(template.New("html").Funcs(template.FuncMap{}).Parse(treeTemplate))

const treeTemplate = `<!DOCTYPE html>
<html>
	<head>
		<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
		<style>
			.main {
				height: 100%;
				width: 100%;
				display: flex;
				padding: 8px;
			}
			.tree {
				width: 25%;
				height: 95vh;
				white-space: nowrap;
				overflow: scroll;
			}
			.tree div {
				padding: 3px 0;
			}
			.file {
				cursor: pointer;
			}
			.current {
				font-weight: bold;
				background-color: #E1F5FF;
			}
			.content {
				width: 70%;
				margin-left: 32px;
				margin-right: 32px;
			}
			.source {
				white-space: nowrap;
				overflow-x: scroll;
			}
			pre {
				counter-reset: line;
				height: 90vh;
				font-family: Menlo, monospace;
				font-weight: bold;
				color: rgb(80, 80, 80);
			}
			ol {
				list-style: none;
				counter-reset: number;
				margin: 0;
				padding: 0;
			}
			li:before {
				counter-increment: number;
				content: counter(number);
				margin-right: 24px;
				display: inline-block;
				width: 50px;
				text-align: right;
				color: rgb(200, 200, 200);
			}
			.cov0 {
				color: rgb(192, 0, 0);
			}
			.cov1 {
				color: rgb(44, 212, 149);
			}
		</style>
	</head>
	<body>
		<div class="main">
			<div class="tree">{{.Tree}}</div>
			<div id="cov" class="content">
				{{range $i, $f := .Files}}
				<div id="file{{$i}}" class="source" style="display: none">
					<pre>{{$f.Body}}</pre>
				</div>
				{{end}}
			</div>
		</div>

		<script>
			let current;
			let currentTree;

			function select(n) {
				if (current) {
					current.style.display = 'none';
				}

				current = document.getElementById('file' + n);
				if (!current) {
					return;
				}
				current.style.display = 'block';
				current.scrollLeft = 0;
				current.scrollTop = 0;
			}
			function selectTree(n) {
				if (currentTree) {
					currentTree.classList.remove('current');
				}

				currentTree = document.getElementById('tree' + n);
				if (!current) {
					return;
				}
				currentTree.classList.add('current');
			}

			function change(n) {
				select(n);
				selectTree(n);
			}
		</script>
	</body>
</html>
`
