package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"time"

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
		StackName:    aws.String("CloudValidateTemp"),
		TemplateBody: aws.String(template),
		Parameters:   nil,
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
	for {
		out, err := cf.DescribeStacks(&cloudformation.DescribeStacksInput{
			NextToken: aws.String("1"),
			StackName: aws.String(stackID),
		})

		if err != nil {
			fmt.Println(err.Error())
			return err
		}
		if *out.Stacks[0].StackStatus != "CREATE_IN_PROGRESS" {
			return nil
		}
		time.Sleep(10 * time.Second)
	}

}

// deleteStack deletes the current stack and returns an error if failing to do so
func deleteStack(cf *cloudformation.CloudFormation, stackID string) error {
	_, err := cf.DeleteStack(&cloudformation.DeleteStackInput{
		StackName: aws.String(stackID),
	})
	return err
}
func main() {
	var templateLocation = flag.String("template", "", "the file location of the template")
	flag.Parse()
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
		color.Red("Failed to create stack")
		color.Red(err.Error())
		os.Exit(1)
	}
	err = describeStack(cf, stackID)
	for err != nil {
		color.Red("Failed to describe stack")
		color.Red(err.Error())
		err = describeStack(cf, stackID)
	}
	color.Yellow("Deleting stack")
	if err = deleteStack(cf, stackID); err != nil {
		color.Red(err.Error())
	}
}
