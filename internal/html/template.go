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
		<title>Coverage Report</title>
		<style>
			body {
				font-family: "SF Mono", "Menlo", "Monaco", "Consolas", "Liberation Mono", "Courier New", monospace;
				margin: 0;
			}
			.main {
				width: 100%;
				display: flex;
			}
			.light {
				background: #FFFFFF;
				color: rgb(80, 80, 80);
			}
			.dark {
				background: #000000;
				color: rgb(160, 160, 160);
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
				border-right: 1px solid rgb(160, 160, 160);
			}
			#tree div {
				padding: 4px 0;
			}
			.clickable {
				cursor: pointer;
				color: #3EA6FF;
				text-decoration: underline;
			}

			.current {
				font-weight: bold;
			}
			.light .current {
				background-color: #E6F0FF;
			}
			.dark .current {
				background-color: #555555;
			}

			#coverage {
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
			}
			.light li:before {
				color: rgb(200, 200, 200);
			}
			.dark li:before {
				color: rgb(80, 80, 80);
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
		<div class="main {{.Theme}}">
			<div id="tree">
			{{range $i, $t := .Tree}}
				{{if $t.IsFile}}
				<div class="file clickable" style="padding-inline-start: {{indent $t.Indent}}px;" id="tree{{$t.ID}}" onclick="change({{$t.ID}}, {{$t.Indent}});">
					{{$t.Name}} ({{$t.Coverage}}%)
				</div>
				{{else}}
				<div style="padding-inline-start: {{indent $t.Indent}}px">{{$t.Name}}/ ({{$t.Coverage}}%)</div>
				{{end}}
			{{end}}
			</div>
			<div id="coverage">
				{{range $i, $f := .Files}}
				<div id="file{{$f.ID}}"  style="display: none">
					<table>
						<tr><th colspan="2">Coverages</th></tr>
						<tr><td class="total">Total</td><td class="cov">{{$f.Coverage}}%</td></tr>
						{{range $j, $fn := .Functions}}
						<tr>
							<td class="fnc"><span class="clickable" onclick="scrollById('file{{$f.ID}}-{{$fn.Line}}');">{{$fn.Name}}</span></td>
							<td class="cov">{{$fn.Coverage}}%</td>
						</tr>
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
			updateByQuery();

			window.addEventListener('popstate', function(e) {
				updateByQuery();
			})

			function updateByQuery() {
				const searchParams = new URLSearchParams(window.location.search);
				const n = searchParams.get('n');
				const i = searchParams.get('i');
				if (n && i) {
					change(n, i);
				} 
			}

			function select(n) {
				if (current) {
					current.style.display = 'none';
				}

				current = document.getElementById('file' + n);
				if (!current) {
					return;
				}
				current.style.display = 'block';
				scrollById('coverage');
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
			function scrollById(id) {
				const elm = document.getElementById(id);
				const rect = elm.getBoundingClientRect();
				document.documentElement.scrollTop = rect.top + window.pageYOffset;
			}
			function change(n, i) {
				select(n);
				selectTree(n, i);
				updateUrl(n, i)
			}
			function updateUrl(n, i) {
				const url = new URL(window.location.href);
				if( !url.searchParams.get('n') ) {
					url.searchParams.append('n',n);
					url.searchParams.append('i',i);
					location.href = url;
				} else {
					if (url.searchParams.get('n') != n || url.searchParams.get('i') != i) {
						url.searchParams.set('n',n);
						url.searchParams.set('i',i);
						history.pushState("", "", url);
					}
				}
			}
		</script>
	</body>
</html>
`
