# Cloudspeed

*Cloudspeed is a cli tool to quickly debug and create cloudformation stacks.*

Cloudspeed aims to alleviate the frustration of debugging cloudformation templates. Currently the `aws cloudformation validate-template` tool validates if the file is either valid JSON or valid YAML. It doesn't validate the properties and configuration. This tool ensures the template is valid by creating the stack with the associated infrastructure and deleting the stack after creation unless the persist flag is set to true. If the stack fails to create, the stack's events are displayed in the terminal, allowing for more rapid iterations over a troublesome template and quicker debugging. If you would like to keep the stack around if successfully created, pass in the persist flag to ensure the infrastructure remains untouched. Parameters can be passed in via an optional parameters file, containing a JSON array of [cloudformation.Parameter](https://docs.aws.amazon.com/sdk-for-go/api/service/cloudformation/#Parameter) objects

## Usage 
```
Usage of cloudspeed
  -name string (optional)
        the name of the stack to create. Defaults to cloudspeed (default "cloudspeed")
  -parameters string (optional)
        the location of the JSON parameter file. Should contain a JSON array of cloudformation.Parameter objects. See examples/parameters.json for reference. 
  -persist boolean (optional)
        persist will persist the stack if successful. Defaults to false, deleting the stack after completion
  -template string (required)
        the file location of the template
```


` cloudspeed -template templates/test-s3website.yaml -name mywebsite -persist=true -parameters templates/test-s3website-parameters.json`
