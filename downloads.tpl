---
layout: default
title: Downloads
---

NOTE: You can get latest version of scaleway-cli using `go get -u github.com/scaleway/scaleway-cli/cmd/scw`.

{{.AppName}} downloads (version {{.Version}})

{{range $k, $v := .Categories}}### {{$k}}

{{range $v}} * [{{.Text}}]({{.RelativeLink}})
{{end}}
{{end}}

---
{{.ExtraVars.footer}}
