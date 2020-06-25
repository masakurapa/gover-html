package html

import "html/template"

type DefaultTemplateData struct {
	Files []TemplateFile
}

var parsedDefaultTemplate = template.Must(template.New("html").Funcs(template.FuncMap{}).Parse(defaultTmpl))

const defaultTmpl = `
<!DOCTYPE html>
<html>
<head>
	<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
	<style>
		.source {
			white-space: nowrap;
			overflow-x: scroll;
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
	<div>
		<select id="files" onchange="change();">
			{{range $i, $f := .Files}}
			<option value="file{{$i}}">{{$f.Name}} ({{printf "%.1f" $f.Coverage}}%)</option>
			{{end}}
		</select>
	</div>

	{{range $i, $f := .Files}}
		<div class="source" id="file{{$i}}" style="display: none">
			<pre>{{$f.Body}}</pre>
		</div>
	{{end}}

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
		function change(t) {
			select(document.getElementById('files').value);
			window.scrollTo(0, 0);
		}

		if (!current) {
			select("file0");
		}
	</script>
</body>
</html>
`
