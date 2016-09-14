package codegen

import (
	"fmt"
	"strings"

	"github.com/Jumpscale/go-raml/raml"
)

// API client definition
type clientDef struct {
	Name     string
	BaseURI  string
	Services map[string]ClientService
}

// create client definition from RAML API definition
func newClientDef(apiDef *raml.APIDefinition) clientDef {
	cd := clientDef{
		Name:     normalizeURI(apiDef.Title),
		BaseURI:  apiDef.BaseURI,
		Services: map[string]ClientService{},
	}
	if strings.Index(cd.BaseURI, "{version}") > 0 {
		cd.BaseURI = strings.Replace(cd.BaseURI, "{version}", apiDef.Version, -1)
	}
	return cd
}

// GenerateClient generates client library
func GenerateClient(apiDef *raml.APIDefinition, dir, packageName, lang, rootImportPath string) error {
	//check create dir
	if err := checkCreateDir(dir); err != nil {
		return err
	}

	// global variables
	globAPIDef = apiDef
	globRootImportPath = rootImportPath

	// creates base client struct
	cd := newClientDef(apiDef)

	services := map[string]*ClientService{}
	for k, v := range apiDef.Resources {
		rd := newResourceDef(apiDef, normalizeURITitle(apiDef.Title), packageName)
		rd.generateMethods(&v, lang)
		services[k] = &ClientService{
			lang:         lang,
			rootEndpoint: k,
			PackageName:  packageName,
			Methods:      rd.Methods,
		}
	}

	switch lang {
	case langGo:
		// rootImportPath only needed if we use libraries
		if rootImportPath == "" && len(apiDef.Libraries) > 0 {
			return fmt.Errorf("--import-path can't be empty when we use libraries")
		}

		gc := goClient{
			clientDef:      cd,
			libraries:      apiDef.Libraries,
			PackageName:    packageName,
			RootImportPath: rootImportPath,
			Services:       services,
		}
		return gc.generate(apiDef, dir)
	case langPython:
		pc := pythonClient{
			clientDef: cd,
			Services:  services,
		}
		return pc.generate(dir)
	}
	return errInvalidLang
}
