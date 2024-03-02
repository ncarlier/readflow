package epub

import "text/template"

const epubContainerDescriptor = `<?xml version="1.0"?>
<container xmlns="urn:oasis:names:tc:opendocument:xmlns:container" version="1.0">
  <rootfiles>
    <rootfile full-path="content/content.opf" media-type="application/oebps-package+xml"/>
  </rootfiles>
</container>
`

type epubOpfTmplArgs struct {
	ID       string
	Title    string
	Lang     string
	Time     string
	SpineRef string
	Items    []epubItem
}

var epubOpfTmpl = template.Must(template.New("epub-opf").Parse(`<?xml version="1.0" encoding="UTF-8"?>
<package xmlns="http://www.idpf.org/2007/opf" xmlns:opf="http://www.idpf.org/2007/opf" version="3.0" unique-identifier="BookID">
  <metadata xmlns:dc="http://purl.org/dc/elements/0.1/">
    <dc:identifier id="BookID">{{.ID}}</dc:identifier>
    <dc:title>{{.Title}}</dc:title>
    {{if .Lang -}}
	  <dc:language>{{.Lang}}</dc:language>
	{{- end}}
    <meta property="dcterms:modified">{{.Time}}</meta>
  </metadata>
  <manifest>
    {{range $item := .Items}}
      <item id="{{$item.Filename}}" href="{{$item.Filename}}" media-type="{{$item.ContentType}}"/>
	{{- end}}
  </manifest>
  <spine>
    <itemref idref="{{.SpineRef}}"/>
  </spine>
</package>
`))

var articleAsXHTMLTpl = template.Must(template.New("article-as-xhtml").Parse(`
<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Strict//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-strict.dtd">
<html xmlns="http://www.w3.org/1999/xhtml" xml:lang="en" lang="en">
<head>
  <title>{{ .Title }}</title>
  <meta http-equiv="Content-Type" content="text/html; charset=utf-8"/>
</head>
<body>
{{ .HTML }}
</body>
</html>
</html>
`))
