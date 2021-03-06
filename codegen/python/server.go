package python

import (
	log "github.com/Sirupsen/logrus"

	"github.com/Jumpscale/go-raml/raml"
)

// Server represents a python server
type Server interface {
	Generate(dir string) error
}

// NewServer creates a new python server
func NewServer(kind string, apiDef *raml.APIDefinition, apiDocsDir string, withMain bool) Server {
	switch kind {
	case "", "flask":
		return NewFlaskServer(apiDef, apiDocsDir, withMain)
	case "sanic":
		return NewSanicServer(apiDef, apiDocsDir, withMain)
	default:
		log.Fatalf("Invalid kind of python server : %v", kind)
		return nil
	}
}
