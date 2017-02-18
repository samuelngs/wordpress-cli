package app

import (
	"bytes"
	"html/template"
)

// Service struct
type Service struct {
	Name        string
	Image       string
	Version     string
	DependsOn   []string
	Volumes     map[string]string
	Environment map[string]string
	Ports       map[uint]uint
	Links       map[string]string
	Networks    []string
}

// Configuration struct
type Configuration struct {
	Services []*Service
	Volumes  []string
	Networks []string
}

// Template for docker compose
const Template = `
version: '2'
{{- $services := .Services }}
{{- $volumes := .Volumes }}
{{- $networks := .Networks }}

{{- if gt (len $services) 0 }}
services:
{{- range $service := $services }}
  {{ $service.Name }}:
    image: {{ $service.Image }}:{{ $service.Version }}

    {{- if gt (len $service.DependsOn) 0 }}
    depends_on:
    {{- range $dependency := $service.DependsOn }}
      - {{ $dependency }}
    {{- end }}
    {{- end }}

    {{- if gt (len $service.Volumes) 0 }}
    volumes:
    {{- range $target, $destination := $service.Volumes }}
      - {{ $target }}:{{ $destination }}
    {{- end }}
    {{- end }}

    {{- if gt (len $service.Environment) 0 }}
    environment:
    {{- range $key, $val := $service.Environment }}
      {{ $key }}: {{ $val }}
    {{- end }}
    {{- end }}

    {{- if gt (len $service.Ports) 0 }}
    ports:
    {{- range $target, $destination := $service.Ports }}
      - {{ $target }}:{{ $destination }}
    {{- end }}
    {{- end }}

    {{- if gt (len $service.Links) 0 }}
    links:
    {{- range $target, $endpoint := $service.Links }}
	{{- if eq $target $endpoint }}
      - {{ $target }}
	{{- else }}
      - {{ $target }}:{{ $endpoint }}
	{{- end }}
    {{- end }}
    {{- end }}

    {{- if gt (len $service.Networks) 0 }}
    networks:
    {{- range $network := $service.Networks }}
      - {{ $network }}
    {{- end }}
    {{- end }}

{{- end }}
{{- end }}

{{- if gt (len $volumes) 0 }}
volumes:
  {{- range $volume := $volumes }}
  {{ $volume }}: {}
  {{- end }}
{{- end }}

{{- if gt (len $networks) 0 }}
networks:
  {{- range $network := $networks }}
  {{ $network }}: {}
  {{- end }}
{{- end }}

`

// buildConfig to build docker compose configuration
func buildConfig(hash, dir string) *Configuration {
	return &Configuration{
		Services: []*Service{
			&Service{
				Name:    "db",
				Image:   "mariadb",
				Version: "latest",
				Volumes: map[string]string{
					"db": "/var/lib/mysql",
				},
				Environment: map[string]string{
					"MYSQL_ROOT_PASSWORD": hash,
				},
			},
			&Service{
				Name:    "wordpress",
				Image:   "wordpress",
				Version: "4.7.2",
				DependsOn: []string{
					"db",
				},
				Volumes: map[string]string{
					dir + "/plugins": "/var/www/html/wp-content/plugins",
					"theme":          "/var/www/html/wp-content/themes/default",
				},
				Environment: map[string]string{
					"WORDPRESS_DB_HOST":     hash + "_db_1:3306",
					"WORDPRESS_DB_PASSWORD": hash,
				},
			},
			&Service{
				Name:    "webpack",
				Image:   "samuelngs/wordpress-sucks",
				Version: "latest",
				Volumes: map[string]string{
					dir + "/theme": "/theme",
					"theme":        "/wp-content/themes/default",
				},
				Environment: map[string]string{
					"WP_TARGET_HOST": hash + "_wordpress_1",
					"WP_TARGET_PORT": "80",
					"WP_PROXY_HOST":  "localhost",
					"WP_PROXY_ADDR":  "0.0.0.0",
					"WP_PROXY_PORT":  "5001",
				},
				Ports: map[uint]uint{
					5001: 5001,
				},
			},
		},
		Volumes: []string{
			"db",
			"theme",
		},
	}
}

// compile to compile docker compose configuration
func compile(conf *Configuration) ([]byte, error) {
	var doc bytes.Buffer
	t := template.New("docker-compose.yaml")
	t, _ = t.Parse(Template)
	if err := t.Execute(&doc, conf); err != nil {
		return nil, err
	}
	return doc.Bytes(), nil
}
