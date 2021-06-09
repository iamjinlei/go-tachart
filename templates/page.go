package templates

var PageTpl = `
{{- define "page" }}
<!DOCTYPE html>
<html>
    {{- template "header" . }}
<body style="background-color:{{ .PageBackgroundColor }}">

<div id="grid_container" style="display:grid;grid-template-columns:{{ .Layout.TemplateColumns }};justify-content:center;">

<div id="page_top" style="grid-column-start:1;grid-column-end:4;grid-row-start:1;grid-row-end:2;height:{{ .Layout.TopHeight }};">
	{{ .Layout.TopContent }}
</div>


<div id="page_left" style="grid-column-start:1;grid-column-end:2;grid-row-start:2;grid-row-end:3;">
	{{ .Layout.LeftContent }}
</div>

<div id="page_center" style="grid-column-start:2;grid-column-end:3;grid-row-start:2;grid-row-end:3;">

{{ if eq .ChartArea "none" }}
	<div style="widht:auto">
    	{{- range .Charts }} {{ template "base" . }} {{- end }}
	</div>
{{ end }}

{{ if eq .ChartArea "center" }}
	<style> .container {display: flex;justify-content: center;align-items: center;} .item {margin: auto;} </style>
	<div style="widht:auto">
    	{{- range .Charts }} {{ template "base" . }} {{- end }}
	</div>
{{ end }}

{{ if eq .ChartArea "flex" }}
	<style> .box { justify-content:center; display:flex; flex-wrap:wrap } </style>
	<div style="widht:auto">
		<div class="box"> {{- range .Charts }} {{ template "base" . }} {{- end }} </div>
	</div>
{{ end }}

</div>

<div id="page_right" style="grid-column-start:3;grid-column-end:4;grid-row-start:2;grid-row-end:3;">
	{{ .Layout.RightContent }}
</div>

<div id="page_bottom" style="grid-column-start:1;grid-column-end:4;grid-row-start:3;grid-row-end:4;height:{{ .Layout.BottomHeight }};">
	{{ .Layout.BottomContent }}
</div>

</div>

</body>
</html>
{{ end }}
`
