package html

import (
	"fmt"
	"html/template"
)

var parsedTreeTemplate = template.
	Must(template.New("html").
		Funcs(template.FuncMap{
			"indent": func(i int) int { return i*24 + 12 },
			"coverageColor": func(coverage float64) template.CSS {
				// Professional color gradient with proper opacity
				var h float64
				if coverage < 50 {
					// Red to orange gradient
					h = coverage * 0.6 // 0-30 degrees (red to orange)
				} else if coverage < 80 {
					// Orange to yellow gradient
					h = 30 + ((coverage - 50) / 30 * 30) // 30-60 degrees
				} else {
					// Yellow to green gradient
					h = 60 + ((coverage - 80) / 20 * 60) // 60-120 degrees
				}
				return template.CSS(fmt.Sprintf("hsla(%.0f, 70%%, 50%%, 0.4)", h))
			},
		}).
		Parse(treeTemplate))

const treeTemplate = `<!DOCTYPE html>
<html>
	<head>
		<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
		<meta name="viewport" content="width=device-width, initial-scale=1">
		<title>Coverage Report</title>
		<style>
			* {
				box-sizing: border-box;
			}
			body {
				margin: 0;
				font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, sans-serif;
				font-size: 14px;
				line-height: 1.6;
			}
			
			/* Custom scrollbar */
			::-webkit-scrollbar {
				width: 8px;
				height: 8px;
			}
			::-webkit-scrollbar-track {
				background: transparent;
			}
			::-webkit-scrollbar-thumb {
				background: rgba(156, 163, 175, 0.3);
				border-radius: 4px;
			}
			::-webkit-scrollbar-thumb:hover {
				background: rgba(156, 163, 175, 0.5);
			}
			.dark ::-webkit-scrollbar-thumb {
				background: rgba(71, 85, 105, 0.5);
			}
			.dark ::-webkit-scrollbar-thumb:hover {
				background: rgba(71, 85, 105, 0.7);
			}
			.main {
				width: 100%;
				min-height: 100vh;
				display: flex;
			}
			.light {
				background: #f8fafc;
				color: #1e293b;
			}
			.dark {
				background: #0f172a;
				color: #e2e8f0;
			}

			#tree {
				width: 280px;
				min-width: 280px;
				height: 100vh;
				padding: 24px 16px;
				white-space: nowrap;
				overflow-x: hidden;
				overflow-y: auto;
				position: sticky;
				top: 0;
				left: 0;
				background: #ffffff;
				box-shadow: 1px 0 0 rgba(0, 0, 0, 0.05);
			}
			.dark #tree {
				background: #1e293b;
				box-shadow: 1px 0 0 rgba(255, 255, 255, 0.05);
			}
			#tree > div {
				padding: 0;
				position: relative;
				margin: 4px 0;
			}
			.tree-item {
				position: relative;
				padding: 8px 12px;
				border-radius: 6px;
				transition: all 0.2s ease;
				background: transparent;
			}
			.tree-item:hover {
				background-color: rgba(59, 130, 246, 0.05);
			}
			.dark .tree-item:hover {
				background-color: rgba(96, 165, 250, 0.1);
			}
			.tree-item-bg {
				position: absolute;
				left: 0;
				top: 0;
				bottom: 0;
				z-index: 0;
				border-radius: 6px;
			}
			.tree-item-content {
				position: relative;
				z-index: 1;
				font-weight: 500;
				display: flex;
				justify-content: space-between;
				align-items: center;
			}
			.tree-name {
				flex: 1;
				overflow: hidden;
				text-overflow: ellipsis;
			}
			.tree-coverage {
				margin-left: 12px;
				font-size: 12px;
				font-weight: 600;
				color: #64748b;
			}
			.dark .tree-coverage {
				color: #94a3b8;
			}
			.clickable {
				cursor: pointer;
			}
			.clickable .tree-item-content {
				color: inherit;
				text-decoration: none;
			}

			.current .tree-item {
				background-color: rgba(59, 130, 246, 0.1);
				box-shadow: inset 3px 0 0 #3b82f6;
			}
			.dark .current .tree-item {
				background-color: rgba(96, 165, 250, 0.15);
				box-shadow: inset 3px 0 0 #60a5fa;
			}

			#coverage {
				flex: 1;
				padding: 24px;
				overflow-y: auto;
			}

			.source {
				white-space: nowrap;
				background: #ffffff;
				border-radius: 8px;
				padding: 24px;
				box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
				overflow-x: auto;
				margin-top: 24px;
			}
			.dark .source {
				background: #1e293b;
				box-shadow: 0 1px 3px rgba(0, 0, 0, 0.3);
			}
			pre {
				counter-reset: line;
				font-family: "SF Mono", Monaco, "Cascadia Code", "Roboto Mono", Consolas, "Courier New", monospace;
				font-size: 13px;
				line-height: 1.6;
				margin: 0;
			}
			ol {
				list-style: none;
				counter-reset: number;
				margin: 0;
				padding: 0;
			}
			li {
				padding: 2px 0;
				transition: background-color 0.2s ease;
			}
			li:hover {
				background-color: rgba(59, 130, 246, 0.05);
			}
			.dark li:hover {
				background-color: rgba(96, 165, 250, 0.05);
			}
			li:before {
				counter-increment: number;
				content: counter(number);
				margin-right: 24px;
				display: inline-block;
				width: 50px;
				text-align: right;
				color: #94a3b8;
				font-weight: normal;
				cursor: pointer;
				transition: color 0.2s ease;
			}
			li:hover:before {
				color: #3b82f6;
			}
			.dark li:hover:before {
				color: #60a5fa;
			}
			.dark li:before {
				color: #64748b;
			}
			.range-highlight {
				background-color: rgba(59, 130, 246, 0.15) !important;
			}
			.dark .range-highlight {
				background-color: rgba(96, 165, 250, 0.15) !important;
			}

			.cov0 {
				background-color: rgba(239, 68, 68, 0.1);
				font-weight: 600;
			}
			.cov1 {
				background-color: rgba(34, 197, 94, 0.1);
				font-weight: normal;
			}
			table {
				width: 100%;
				background: #ffffff;
				border-radius: 8px;
				overflow: hidden;
				box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
			}
			.dark table {
				background: #1e293b;
				box-shadow: 0 1px 3px rgba(0, 0, 0, 0.3);
			}
			tr {
				border-bottom: 1px solid rgba(148, 163, 184, 0.1);
			}
			.dark tr {
				border-bottom: 1px solid rgba(51, 65, 85, 0.5);
			}
			tr:last-child {
				border-bottom: none;
			}
			th {
				background: #f1f5f9;
				font-weight: 600;
				text-align: left;
				padding: 10px 16px;
				font-size: 13px;
				text-transform: uppercase;
				letter-spacing: 0.05em;
				color: #64748b;
			}
			.dark th {
				background: #0f172a;
				color: #94a3b8;
			}
			td {
				padding: 8px 16px;
			}
			table .total {
				font-weight: 600;
			}
			table .fnc {
				font-family: "SF Mono", Monaco, "Cascadia Code", "Roboto Mono", Consolas, "Courier New", monospace;
				font-size: 13px;
			}
			table .fnc .clickable {
				color: #3b82f6;
				text-decoration: none;
				transition: color 0.2s ease;
			}
			table .fnc .clickable:hover {
				color: #2563eb;
				text-decoration: underline;
			}
			.dark table .fnc .clickable {
				color: #60a5fa;
			}
			.dark table .fnc .clickable:hover {
				color: #93bbfc;
			}
			table .cov {
				width: 120px;
				text-align: right;
				padding: 0;
			}
			.cov-cell {
				position: relative;
				padding: 8px 16px;
				background-color: transparent;
			}
			.cov-text {
				position: relative;
				z-index: 1;
				font-weight: 600;
				font-size: 13px;
			}
			.cov-bg {
				position: absolute;
				left: 0;
				top: 0;
				bottom: 0;
				z-index: 0;
			}
		</style>
	</head>
	<body>
		<div class="main {{.Theme}}">
			<div id="tree">
			{{range $i, $t := .Tree}}
				{{if $t.IsFile}}
				<div class="file" id="tree{{$t.ID}}" onclick="change({{$t.ID}}, {{$t.Indent}});">
					<div class="tree-item clickable" style="margin-inline-start: {{indent $t.Indent}}px;">
						<div class="tree-item-bg" style="width: {{if eq $t.Coverage 0.0}}100{{else}}{{$t.Coverage}}{{end}}%; background-color: {{coverageColor $t.Coverage}}"></div>
						<div class="tree-item-content">
							<span class="tree-name">{{$t.Name}}</span>
							<span class="tree-coverage">{{$t.Coverage}}%</span>
						</div>
					</div>
				</div>
				{{else}}
				<div>
					<div class="tree-item" style="margin-inline-start: {{indent $t.Indent}}px">
						<div class="tree-item-bg" style="width: {{if eq $t.Coverage 0.0}}100{{else}}{{$t.Coverage}}{{end}}%; background-color: {{coverageColor $t.Coverage}}"></div>
						<div class="tree-item-content">
							<span class="tree-name">{{$t.Name}}/</span>
							<span class="tree-coverage">{{$t.Coverage}}%</span>
						</div>
					</div>
				</div>
				{{end}}
			{{end}}
			</div>
			<div id="coverage">
				{{range $i, $f := .Files}}
				<div id="file{{$f.ID}}"  style="display: none">
					<table>
						<tr><th>Function</th><th style="text-align: right">Coverage</th></tr>
						<tr>
							<td class="total">Total Coverage</td>
							<td class="cov">
								<div class="cov-cell">
									<div class="cov-bg" style="width: {{if eq $f.Coverage 0.0}}100{{else}}{{$f.Coverage}}{{end}}%; background-color: {{coverageColor $f.Coverage}}"></div>
									<div class="cov-text">{{$f.Coverage}}%</div>
								</div>
							</td>
						</tr>
						{{range $j, $fn := .Functions}}
						<tr>
							<td class="fnc"><span class="clickable" onclick="scrollToLine({{$f.ID}}, {{$fn.Line}});">{{$fn.Name}}</span></td>
							<td class="cov">
								<div class="cov-cell">
									<div class="cov-bg" style="width: {{if eq $fn.Coverage 0.0}}100{{else}}{{$fn.Coverage}}{{end}}%; background-color: {{coverageColor $fn.Coverage}}"></div>
									<div class="cov-text">{{$fn.Coverage}}%</div>
								</div>
							</td>
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
			let current;
			let currentTree;
			updateByQuery();

			window.addEventListener('popstate', function(e) {
				updateByQuery();
			})

			window.addEventListener('hashchange', function(e) {
				handleHashChange();
			})

			window.addEventListener('load', function() {
				handleHashChange();
				setupLineNumberClicks();
			})

			function updateByQuery() {
				const searchParams = new URLSearchParams(window.location.search);
				const n = searchParams.get('n');
				const i = searchParams.get('i');
				if (n && i) {
					change(n, i);
				}
				// Check for hash after loading the file
				setTimeout(handleHashChange, 100);
			}

			function handleHashChange() {
				const hash = window.location.hash;
				if (hash) {
					// Check if it's a range (e.g., #file1-L10-L20)
					const rangeMatch = hash.match(/#file(\d+)-L(\d+)-L(\d+)/);
					if (rangeMatch) {
						const fileId = rangeMatch[1];
						const startLine = parseInt(rangeMatch[2]);
						const endLine = parseInt(rangeMatch[3]);
						
						// Highlight the range
						highlightLineRange(fileId, startLine, endLine);
						
						// Scroll to start of range
						const startElm = document.getElementById('file' + fileId + '-' + startLine);
						if (startElm) {
							startElm.scrollIntoView({ behavior: 'smooth', block: 'center' });
						}
					} else {
						// Single line
						const elm = document.getElementById(hash.substring(1));
						if (elm) {
							elm.scrollIntoView({ behavior: 'smooth', block: 'center' });
							// Clear previous highlights
							clearHighlights();
							// Highlight the line briefly
							elm.style.backgroundColor = 'rgba(59, 130, 246, 0.2)';
							setTimeout(() => {
								elm.style.backgroundColor = '';
							}, 1500);
						}
					}
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
				// Setup line number clicks for the newly displayed file
				setTimeout(setupLineNumberClicks, 100);
			}

			function selectTree(n, indent) {
				if (currentTree) {
					currentTree.classList.remove('current');
				}

				currentTree = document.getElementById('tree' + n);
				if (!currentTree) {
					return;
				}
				currentTree.classList.add('current');
			}

			function scrollById(id) {
				const elm = document.getElementById(id);
				if (elm) {
					elm.scrollIntoView({ behavior: 'smooth', block: 'center' });
				}
			}

			function scrollToLine(fileId, lineNum) {
				const lineId = 'file' + fileId + '-' + lineNum;
				// Update URL with hash
				const url = new URL(window.location.href);
				url.hash = lineId;
				history.pushState("", "", url);
				
				// Scroll to the line
				scrollById(lineId);
				
				// Highlight the line
				const elm = document.getElementById(lineId);
				if (elm) {
					elm.style.backgroundColor = 'rgba(59, 130, 246, 0.2)';
					setTimeout(() => {
						elm.style.backgroundColor = '';
					}, 1500);
				}
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

			let lastClickedLine = null;

			function setupLineNumberClicks() {
				// Wait for content to be loaded
				setTimeout(() => {
					document.querySelectorAll('li').forEach(li => {
						const id = li.id;
						if (id && id.includes('-')) {
							// Create clickable line number
							li.addEventListener('click', function(e) {
								if (e.target === li || window.getComputedStyle(e.target, ':before').getPropertyValue('content')) {
									const parts = id.split('-');
									const fileId = parts[0].replace('file', '');
									const lineNum = parseInt(parts[1]);
									
									if (e.shiftKey && lastClickedLine && lastClickedLine.fileId === fileId) {
										// Range selection
										const startLine = Math.min(lastClickedLine.lineNum, lineNum);
										const endLine = Math.max(lastClickedLine.lineNum, lineNum);
										
										// Update URL with range
										const url = new URL(window.location.href);
										url.hash = 'file' + fileId + '-L' + startLine + '-L' + endLine;
										history.pushState("", "", url);
										
										// Highlight range
										highlightLineRange(fileId, startLine, endLine);
									} else {
										// Single line selection
										const url = new URL(window.location.href);
										url.hash = id;
										history.pushState("", "", url);
										
										// Clear previous highlights
										clearHighlights();
										
										// Highlight the line
										li.style.backgroundColor = 'rgba(59, 130, 246, 0.2)';
										setTimeout(() => {
											if (!li.classList.contains('range-highlight')) {
												li.style.backgroundColor = '';
											}
										}, 1500);
										
										lastClickedLine = { fileId, lineNum };
									}
								}
							});
						}
					});
				}, 100);
			}

			function highlightLineRange(fileId, startLine, endLine) {
				clearHighlights();
				for (let i = startLine; i <= endLine; i++) {
					const li = document.getElementById('file' + fileId + '-' + i);
					if (li) {
						li.classList.add('range-highlight');
						li.style.backgroundColor = 'rgba(59, 130, 246, 0.15)';
					}
				}
			}

			function clearHighlights() {
				document.querySelectorAll('.range-highlight').forEach(li => {
					li.classList.remove('range-highlight');
					li.style.backgroundColor = '';
				});
			}
		</script>
	</body>
</html>
`
