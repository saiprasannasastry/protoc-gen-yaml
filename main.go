package main

import (
	"fmt"
	"github.com/apex/log"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/pluginpb"
	"io/ioutil"
	"os"
)

func main() {
	input, _ := ioutil.ReadAll(os.Stdin)
	var req pluginpb.CodeGeneratorRequest
	proto.Unmarshal(input, &req)

	// Initialise our plugin with default options
	opts := protogen.Options{}
	plugin, err := opts.New(&req)
	if err != nil {
		log.Fatalf("failed to initialise protoc-gen-yaml %v", err)
	}
	var out string
	out, err = GenerateYaml(plugin)
	if err != nil {
		log.WithError(err).Error("failed to generate yaml")
	}
	fmt.Fprintf(os.Stdout, out)
}
