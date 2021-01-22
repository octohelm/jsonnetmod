package jsonnetmod

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"

	"github.com/google/go-jsonnet"
	"github.com/pkg/errors"
)

type Resolver interface {
	Resolve(ctx context.Context, importPath string, importedFrom string) (fullDir string, err error)
}

func NewImporter(ctx context.Context, resolver Resolver) jsonnet.Importer {
	return &Importer{ctx: ctx, resolver: resolver}
}

type Importer struct {
	ctx      context.Context
	resolver Resolver
	caches   map[string]jsonnet.Contents
}

func (i *Importer) Import(importedFrom, importedPath string) (jsonnet.Contents, string, error) {
	fullImportedPath, err := i.resolve(importedFrom, importedPath)
	if err != nil {
		return jsonnet.Contents{}, "", err
	}
	contents, err := i.load(fullImportedPath)
	if err != nil {
		return jsonnet.Contents{}, "", errors.Wrapf(err, "load %s failed", fullImportedPath)
	}
	return contents, fullImportedPath, nil
}

func (i *Importer) resolve(importedFrom, importedPath string) (string, error) {
	// hack the `k.libsonnet`
	if importedPath == "k.libsonnet" {
		importedPath = "k/main.libsonnet"
	}

	abs := importedPath

	if !filepath.IsAbs(importedPath) {
		abs = filepath.Join(filepath.Dir(importedFrom), importedPath)
	}

	if _, err := os.Stat(abs); err != nil {
		if os.IsNotExist(err) && importedPath[0] == '.' {
			return "", err
		}

		p := filepath.Dir(importedPath)
		filename := filepath.Base(importedPath)

		dir, err := i.resolver.Resolve(i.ctx, p, importedFrom)
		if err != nil {
			return "", errors.Wrapf(err, "resolve failed `%s`", importedPath)
		}

		abs = filepath.Join(dir, filename)
	}

	return abs, nil
}

func (i *Importer) load(file string) (jsonnet.Contents, error) {
	if i.caches == nil {
		i.caches = map[string]jsonnet.Contents{}
	}

	if contents, ok := i.caches[file]; ok {
		return contents, nil
	}

	data, err := ioutil.ReadFile(file)
	if err != nil {
		return jsonnet.Contents{}, err
	}

	ext := filepath.Ext(file)

	switch ext {
	case ".yaml", ".yml":
		d, err := yamlToJSON(data)
		if err != nil {
			return jsonnet.Contents{}, err
		}
		data = patchFileContext(d, file)
	}

	contents := jsonnet.MakeContents(string(data))

	i.caches[file] = contents

	return contents, nil
}

func yamlToJSON(data []byte) ([]byte, error) {
	v := map[string]interface{}{}
	if err := yaml.Unmarshal(data, &v); err != nil {
		return nil, err
	}
	return json.MarshalIndent(v, "", "  ")
}

func patchFileContext(d []byte, file string) []byte {
	meta := fmt.Sprintf(" + { __dirname:: '%s', __filename:: '%s'}", filepath.Dir(file), file)
	return append(d, []byte(meta)...)
}
