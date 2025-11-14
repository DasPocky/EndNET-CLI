package cloudinit

import (
	"bytes"
	"text/template"

	"endnet-cli/pkg/models"
)

var (
	edgeTemplate = template.Must(template.New("edge").Parse(`#cloud-config
hostname: {{ .Roles.Edge.Name }}
write_files:
  - path: /etc/endnet/edge.info
    content: |
      project={{ .Project }}
      forgejo={{ .DNS.ForgejoHost }}
`))

	wgTemplate = template.Must(template.New("wg").Parse(`#cloud-config
hostname: {{ .Roles.WG.Name }}
write_files:
  - path: /etc/endnet/wg.info
    content: |
      subnet={{ .Network.CIDR }}
`))

	forgeTemplate = template.Must(template.New("forge").Parse(`#cloud-config
hostname: {{ .Roles.Forge.Name }}
write_files:
  - path: /etc/endnet/forgejo.info
    content: |
      forgejo_url=http://{{ .DNS.ForgejoHost }}
`))
)

// RenderEdgeCloudInit renders the edge cloud-init template.
func RenderEdgeCloudInit(spec models.EndnetSpec) (string, error) {
	return render(spec, edgeTemplate)
}

// RenderWGCloudInit renders the WireGuard cloud-init template.
func RenderWGCloudInit(spec models.EndnetSpec) (string, error) {
	return render(spec, wgTemplate)
}

// RenderGitCloudInit renders the Forgejo cloud-init template.
func RenderGitCloudInit(spec models.EndnetSpec) (string, error) {
	return render(spec, forgeTemplate)
}

func render(spec models.EndnetSpec, tmpl *template.Template) (string, error) {
	buf := bytes.NewBuffer(nil)
	if err := tmpl.Execute(buf, spec); err != nil {
		return "", err
	}
	return buf.String(), nil
}
