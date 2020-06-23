{{ if .Versions -}}
{{   if .Unreleased.CommitGroups -}}
<a name="unreleased"></a>
## [Unreleased]

{{     range .Unreleased.CommitGroups -}}
### {{ .Title }}
{{       range .Commits -}}
- {{ if .Scope }}**{{ .Scope }}:** {{ end }}{{ .Subject }}
{{       end -}}
{{     end -}}
{{   end -}}

{{ range .Versions -}}
<a name="{{ .Tag.Name }}"></a>
## {{ if .Tag.Previous }}[{{ .Tag.Name }}]{{ else }}{{ .Tag.Name }}{{ end }} - {{ datetime "2006-01-02" .Tag.Date }}
{{   range .CommitGroups -}}
### {{ .Title }}
{{     range .Commits -}}
- {{ if .Scope }}**{{ .Scope }}:** {{ end }}{{ .Subject }}
{{     end }}
{{   end -}}

{{   if .NoteGroups -}}
{{     range .NoteGroups -}}
### {{ .Title }}
{{       range .Notes }}
{{         .Body }}
{{       end }}
{{     end -}}
{{   end -}}
{{ end -}}

[Unreleased]: {{ .Info.RepositoryURL }}/compare/{{ $latest := index .Versions 0 }}{{ $latest.Tag.Name }}...HEAD
{{   range .Versions -}}
{{     if .Tag.Previous -}}
[{{ .Tag.Name }}]: {{ $.Info.RepositoryURL }}/compare/{{ .Tag.Previous.Name }}...{{ .Tag.Name }}
{{     end -}}
{{   end -}}
{{ end -}}