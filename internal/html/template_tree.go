package html

import "html/template"

type TreeTemplateData struct {
	Tree  template.HTML
	Files []TemplateFile
}

var parsedTreeTemplate = template.Must(template.New("html").Funcs(template.FuncMap{}).Parse(treeTmpl))

const treeTmpl = `
<!DOCTYPE html>
<html>
<head>
	<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
	<style>
		.content {
			display: flex;
			padding: 8px;
		}
		.tree {
			width: 30%;
			white-space: nowrap;
			overflow-x: scroll;
		}
		.cover {
			width: 70%;
			margin-left: 32px;
		}
		.source {
			white-space: nowrap;
			overflow-x: scroll;
		}
		ul {
			list-style: none;
			margin-block-start: unset;
			margin-block-end: unset;
		}
		li .file {
			cursor: pointer;
		}
		pre {
			font-family: Menlo, monospace;
			font-weight: bold;
			color: rgb(80, 80, 80);
		}
		.uncov {
			color: rgb(192, 0, 0);
		}
		.cov {
			color: rgb(44, 212, 149);
		}
	</style>
</head>
<body>
	<div class="content">
		<div class="tree">{{.Tree}}</div>
		<div class="cover">
			{{range $i, $f := .Files}}
				<div id="file{{$i}}" style="display: none">
					<div>{{$f.Name}}</div>
					<div>{{printf "%.1f" $f.Coverage}}%</div>
					<div class="source">
						<pre>{{$f.Body}}</pre>
					</div>
				</div>
			{{end}}
		</div>
	</div>

	<script>
		let current;

		function select(f) {
			if (current) {
				current.style.display = 'none';
			}

			current = document.getElementById(f)
			if (!current) {
				return;
			}

			current.style.display = 'block';
		}
		function change(f) {
			select(f);
			window.scrollTo(0, 0);
		}
	</script>
</body>
</html>
`
