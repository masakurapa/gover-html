package html

import "html/template"

var parsedTreeTemplate = template.
	Must(template.New("html").
		Funcs(template.FuncMap{
			"indent": func(i int) int { return i*30 + 8 },
		}).
		Parse(treeTemplate))

const treeTemplate = `<!DOCTYPE html>
<html>
	<head>
		<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
		<style>
			body {
				margin: 0;
			}
			.main {
				width: 100%;
				display: flex;
			}
			#tree {
				width: 25%;
				height: 100vh;
				padding: 8px 0;
				white-space: nowrap;
				overflow: scroll;
				position: sticky;
				position: -webkit-sticky;
				top: 0;
				left: 0;
				background: #FFFFFF;
				border-right: 1px solid #000000;
			}
			#tree div {
				padding: 4px 0;
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
				margin-left: 16px;
				margin-right: 32px;
			}
			.source {
				white-space: nowrap;
			}
			pre {
				counter-reset: line;
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
			table, tr, td, th {
				border-collapse: collapse;
				border:1px solid #BBBBBB;
			}
			table {
				margin: 16px 0 32px 74px;
			}
			table .total {
				min-width: 300px;
				text-align: left;
				padding-left: 8px;
			}
			table .fnc {
				min-width: 300px;
				text-align: left;
				padding: 0 20px 0 20px;
			}
			table .cov {
				width: 70px;
				text-align: right;
				padding-right: 8px;
			}
		</style>
	</head>
	<body>
		<div class="main">
			<div id="tree">
			{{range $i, $t := .Tree}}
				{{if $t.IsFile}}
				<div class="file" style="padding-inline-start: {{indent $t.Indent}}px;" id="tree{{$t.ID}}" onclick="change({{$t.ID}}, {{$t.Indent}});">
					{{$t.Name}} ({{$t.Coverage}}%)
				</div>
				{{else}}
				<div style="padding-inline-start: {{indent $t.Indent}}px">{{$t.Name}}/</div>
				{{end}}
			{{end}}
			</div>

			<div id="cov" class="content">
				{{range $i, $f := .Files}}
				<div id="file{{$f.ID}}"  style="display: none">
					<table>
						<tr><th colspan="2">Coverages</th></tr>
						<tr><td class="total">Total</td><td class="cov">{{$f.Coverage}}%</td></tr>
						{{range $j, $fn := .Functions}}
						<tr><td class="fnc">{{$fn.Name}}</td><td class="cov">{{$fn.Coverage}}%</td></tr>
						{{end}}
					</table>

					<div class="source">
						<pre>{{$f.Body}}</pre>
					</div>
				</div>
				{{end}}
			</div>
		</div>

		<script>
			// tree max width
			const scrollWidth = document.getElementById('tree').scrollWidth;

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
			function selectTree(n, indent) {
				if (currentTree) {
					currentTree.classList.remove('current');
				}

				currentTree = document.getElementById('tree' + n);
				if (!current) {
					return;
				}
				currentTree.classList.add('current');
				currentTree.style.width = scrollWidth - (indent * 30 + 8) + 'px';
			}
			function change(n, i) {
				select(n);
				selectTree(n, i);
			}
		</script>
	</body>
</html>
`
