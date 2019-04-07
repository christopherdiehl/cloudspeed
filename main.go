package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/fatih/color"
)

/*
 * getTemplate reads in the template file specified in location
 * returns the template file as a string or an error
 */
func getTemplate(location string) (string, error) {
	b, err := ioutil.ReadFile(location)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func createStack(cf *cloudformation.CloudFormation, template string) (string, error) {
	input := &cloudformation.CreateStackInput{
		StackName:    aws.String("CloudValidate_Temp"),
		TemplateBody: aws.String(template),
	}
	if err := input.Validate(); err != nil {
		return "", err
	}
	output, err := cf.CreateStack(input)
	if err != nil {
		return "", err
	}
	color.Green(output.String())
	return *output.StackId, nil
}

// describeStack checks the current status of the stack and returns an error if the stack
func describeStack(cf *cloudformation.CloudFormation, stackID string) error {
	out, err := cf.DescribeStacks(&cloudformation.DescribeStacksInput{
		NextToken: aws.String("1"),
		StackName: aws.String(stackID),
	})
	fmt.Printf("%v\n", *out)
	fmt.Println(err.Error())
	return err
}

// deleteStack deletes the current stack and returns an error if failing to do so
func deleteStack(cf *cloudformation.CloudFormation, stackID string) error {
	out, err := cf.DeleteStack(&cloudformation.DeleteStackInput{
		StackName: aws.String(stackID),
	})
	fmt.Printf("%v\n", out)
	return err
}
func main() {
	var templateLocation = flag.String("template-location", "", "the template file location")
	if *templateLocation == "" {
		color.Red("Please specify the location of the template file")
		os.Exit(1)
	}
	template, err := getTemplate(*templateLocation)
	if err != nil {
		color.Red(err.Error())
		os.Exit(1)
	}
	session, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
	})
	if err != nil {
		color.Red(err.Error())
		os.Exit(1)
	}
	cf := cloudformation.New(session, aws.NewConfig().WithRegion(*session.Config.Region))
	stackID, err := createStack(cf, template)
	if err != nil {
		color.Red(err.Error())
		os.Exit(1)
	}
	err = describeStack(cf, stackID)
	for err != nil {
		fmt.Println(err.Error())
		err = describeStack(cf, stackID)
	}
	fmt.Println("Deleting stack")
	if err = deleteStack(cf, stackID); err != nil {
		fmt.Println("stackId ", stackID, " failed to delete with error: ", err.Error())
	}
}
