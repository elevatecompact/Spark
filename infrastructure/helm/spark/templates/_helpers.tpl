{{- define "spark.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{- define "spark.fullname" -}}
{{- if .Values.fullnameOverride }}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- $name := default .Chart.Name .Values.nameOverride }}
{{- if contains $name .Release.Name }}
{{- .Release.Name | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" }}
{{- end }}
{{- end }}
{{- end }}

{{- define "spark.labels" -}}
helm.sh/chart: {{ include "spark.name" . }}-{{ .Chart.Version }}
{{ include "spark.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{- define "spark.selectorLabels" -}}
app.kubernetes.io/name: {{ include "spark.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{- define "spark.databaseUrl" -}}
{{- $user := .Values.secrets.postgresUser }}
{{- $pass := .Values.secrets.postgresPassword }}
{{- $host := printf "%s-postgres" (include "spark.fullname" .) }}
postgres://{{ $user }}:{{ $pass }}@{{ $host }}:5432/{{ .dbName }}?sslmode=disable
{{- end }}
