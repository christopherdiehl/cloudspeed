# CloudValidate

CloudValidate aims to alleviate the frustration of debugging cloudformation templates. Currently the `aws cloudformation validate-template` tool just validates the file is either valid JSON or valid YAML. It doesn't validate the properties and configuration.

This tool accepts a path to a cloudformation template file location, and attempts to create the cloudformation stack with project name: CloudValidate_Temp. If the stack fails to create, the tool then outputs the errors to the console. CloudValidate will delete the stack once created _or_ failed.

## Goals

- Allow custom stack project names.
- Allow additional parameters to the template file.
- Pass in region to create the stack.

`This project is in active development and breaking changes will occur. Additionally, using the project may incur charges against your AWS account`
