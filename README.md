# Iroworks

Iroworks aims to alleviate the frustration of debugging cloudformation templates. Currently the `aws cloudformation validate-template` tool just validates the file is either valid JSON or valid YAML. It doesn't validate the properties and configuration. This tool ensures the template is valid by actually creating the stack with the associated infrastructure. 

Iroworks is more than a simple validation tool. It was specifically created to reduce stack debugging and creation time. Therefore, if the stack fails to create, the tool then outputs all of the stack events to the console before deletion. If you would like to keep the stack around after it succeeds, pass in the persist flag to ensure the infrastructure remains untouched. 

## Goals
`This project is in active(ish) development and breaking changes will most likely not occur. Using the project may incur charges against your AWS account`


## Usage 
```
Usage of iroworks
  -name string (optional)
        the name of the stack to create. Defaults to Iroworks (default "Iroworks")
  -parameters string (optional)
        the location of the JSON parameter file. Should contain a JSON array of cloudformation.Parameter objects. See examples/parameters.json for reference. 
  -persist boolean (optional)
        persist will persist the stack if successful. Defaults to false, deleting the stack after completion
  -template string (required)
        the file location of the template
```