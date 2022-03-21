## protoc-gen-yaml

protoc-gen-yaml is a proto buff plugin that converts proto format to yaml format

## Installation guide
clone the source code / extract the source code to gopath and run ``go install .``
## How to Run
`` protoc \   
--proto_path . \
-I=test/input \
test/input/proto/echo.proto \
--yaml_out=./output``