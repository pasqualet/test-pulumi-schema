package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/go-yaml/yaml"
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi/pkg/v3/codegen/schema"
	"github.com/pulumi/pulumi/sdk/go/common/util/contract"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("Usage: %s <out-dir>\n", os.Args[0])
		os.Exit(1)
	}

	outdir := os.Args[1]

	pkgSpec := generateSchema()
	mustWritePulumiSchema(pkgSpec, outdir)
}

func mustWritePulumiSchema(pkgSpec schema.PackageSpec, outdir string) {
	schemaJSON, err := yaml.Marshal(pkgSpec)
	if err != nil {
		panic(errors.Wrap(err, "marshaling Pulumi schema"))
	}
	mustWriteFile(outdir, "schema.yaml", schemaJSON)
}

func mustWriteFile(rootDir, filename string, contents []byte) {
	outPath := filepath.Join(rootDir, filename)

	if err := os.MkdirAll(filepath.Dir(outPath), 0755); err != nil {
		panic(err)
	}
	err := ioutil.WriteFile(outPath, contents, 0600)
	if err != nil {
		panic(err)
	}
}

func generateSchema() schema.PackageSpec {
	return schema.PackageSpec{
		Name:       "test-pulumi-schema",
		License:    "Apache-2.0",
		Keywords:   []string{"pulumi", "aws", "eks"},
		Repository: "https://github.com/pasqualet/test-pulumi-schema",
		Resources: map[string]schema.ResourceSpec{
			"test-pulumi-schema:kubernetes:NodeGroup": {
				IsComponent: true,
				ObjectTypeSpec: schema.ObjectTypeSpec{
					Description: "Kubernetes cluster production ready.",
					Properties: map[string]schema.PropertySpec{
						"name": {
							TypeSpec:    schema.TypeSpec{Type: "string"},
							Description: "Kubeconfig",
						},
					},
				},
			},
			"test-pulumi-schema:kubernetes:Cluster": {
				IsComponent: true,
				ObjectTypeSpec: schema.ObjectTypeSpec{
					Description: "Kubernetes cluster production ready.",
					Properties: map[string]schema.PropertySpec{
						"kubeconfig": {
							TypeSpec:    schema.TypeSpec{Type: "string"},
							Description: "Kubeconfig",
						},
						"nodegroup": {
							TypeSpec:    schema.TypeSpec{Ref: "#/resources/test-pulumi-schema:kubernetes:NodeGroup"},
							Description: "Node group",
						},
					},
				},
			},
		},
		Language: map[string]schema.RawMessage{
			"nodejs": rawMessage(map[string]interface{}{
				"respectSchemaVersion": true,
			}),
		},
	}
}

func rawMessage(v interface{}) schema.RawMessage {
	bytes, err := json.Marshal(v)
	contract.Assert(err == nil)
	return bytes
}
