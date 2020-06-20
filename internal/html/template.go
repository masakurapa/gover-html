package html

import "html/template"

type TemplateData struct {
	Tree  template.HTML
	Files []TemplateFile
}

type TemplateFile struct {
	Name     string
	Body     template.HTML
	Coverage float64
}

var parsedTemplate = template.Must(template.New("html").Funcs(template.FuncMap{}).Parse(tmpl))

const tmpl = `
<!DOCTYPE html>
<html>
<head>
	<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
	<style></style>
</head>
<body>
	<div>{{.Tree}}</div>
	{{range $i, $f := .Files}}
		<div>
			<div>{{$f.Name}}</div>
			<div>{{printf "%.1f" $f.Coverage}}</div>
			<pre class="file" id="file{{$i}}">{{$f.Body}}</pre>
		</div>
	{{end}}
	<script></script>
</body>
</html>
`
