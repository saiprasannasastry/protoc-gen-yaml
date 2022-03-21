package main

import (
	"github.com/apex/log"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
	"gopkg.in/yaml.v2"
	"os"
	"sort"
)

var (
	fileName string
)

type Yaml struct {
	Messages []Message `yaml:"messages,omitempty`
	Services []Service `yaml:"services,omitempty"`
}

type Service struct {
	Name    string          `yaml:"name,omitempty"`
	Methods []MethodDetails `yaml:"methods,omitempty"`
}
type MethodDetails struct {
	Name       string `yaml:"name,omitempty"`
	InputType  string `yaml:"input_type,omitempty"`
	OutputType string `yaml:"output_type,omitempty"`
}
type Message struct {
	Name   string  `yaml:"name,omitempty"`
	Fields []Field `yaml:"fields,omitempty"`
}
type Field struct {
	Name   string `yaml:"name,omitempty"`
	Number int32  `yaml:"number,omitempty"`
}

func GenerateYaml(plugin *protogen.Plugin) (string, error) {
	var err error

	for _, file := range plugin.Files {
		packageName := *file.Proto.Package
		fileName = file.GeneratedFilenamePrefix + ".proto.yaml"
		f, err := os.Create("output/" + fileName)

		if err != nil {
			log.Fatalf("failed to create file %v", err)
		}

		genYaml := Yaml{}

		for _, msg := range file.Proto.MessageType {
			//fmt.Printf("The message is %v",msg)
			msgName := *msg.Name
			message := Message{}
			fullName := packageName + "." + msgName
			message.Name = fullName
			fields := genFields(msg.Field)
			message.Fields = fields
			genYaml.Messages = append(genYaml.Messages, message)

			if msg.NestedType != nil {
				for _, nestedMessage := range msg.NestedType {
					nestedMessageDetails := Message{}
					nestedMessageDetails.Name = fullName + "." + *nestedMessage.Name
					nestedMessageDetails.Fields = genFields(nestedMessage.Field)
					genYaml.Messages = append(genYaml.Messages, nestedMessageDetails)
				}
			}

		}
		sort.Slice(genYaml.Messages, func(i, j int) bool {
			return genYaml.Messages[i].Name < genYaml.Messages[j].Name
		})
		for _, srvc := range file.Services {
			srvcData := getSrvcData(srvc)
			genYaml.Services = append(genYaml.Services, srvcData)
		}
		sort.Slice(genYaml.Services, func(i, j int) bool {
			return genYaml.Services[i].Name < genYaml.Services[j].Name
		})

		yamlData, err := yaml.Marshal(&genYaml)
		if err != nil {
			log.WithError(err).Errorf("failed to marshal yaml data")
			return "", err
		}

		f.WriteString(string(yamlData))

	}
	stdout := plugin.Response()
	out, err := proto.Marshal(stdout)
	if err != nil {
		panic(err)
	}
	return string(out), err
}

//genFields returns field name and number for each message
func genFields(fields []*descriptorpb.FieldDescriptorProto) []Field {
	retFields := []Field{}
	for _, data := range fields {
		filed := Field{}
		filed.Name = data.GetName()
		filed.Number = data.GetNumber()
		retFields = append(retFields, filed)
	}
	//sort by field Number
	sort.Slice(retFields, func(i, j int) bool {
		return retFields[i].Number < retFields[j].Number
	})
	return retFields
}

//genSrvcData returns service data for each service
func getSrvcData(service *protogen.Service) Service {
	serviceName := service.Desc.FullName()
	methodDetails := getMethodDetails(service.Methods)
	return Service{string(serviceName), methodDetails}

}

//genMethodDetails gives all methods for each service
func getMethodDetails(methods []*protogen.Method) []MethodDetails {
	//Each Rpc method details
	methodsDetails := []MethodDetails{}
	for _, method := range methods {
		methodDetails := MethodDetails{}
		methodDetails.Name = method.GoName
		methodDetails.InputType = string(method.Desc.Input().FullName())
		methodDetails.OutputType = string(method.Desc.Output().FullName())
		methodsDetails = append(methodsDetails, methodDetails)
	}
	//sort by method Name
	sort.Slice(methodsDetails, func(i, j int) bool {
		return methodsDetails[i].Name < methodsDetails[j].Name
	})
	return methodsDetails
}
