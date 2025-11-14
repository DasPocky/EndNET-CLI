package cloudinit

import "endnet-cli/pkg/models"

// Generator produces cloud-init user-data documents based on configuration and desired state.
type Generator interface {
	Generate(cfg *models.Config, plan *models.Plan) (string, error)
}

// TemplateGenerator is a placeholder implementation that returns canned data.
type TemplateGenerator struct{}

// NewGenerator creates a generator ready for use during early development.
func NewGenerator() Generator {
	return &TemplateGenerator{}
}

// Generate builds a trivial cloud-init document.
func (g *TemplateGenerator) Generate(cfg *models.Config, plan *models.Plan) (string, error) {
	return "#cloud-config\nusers:\n  - name: endnet", nil
}
