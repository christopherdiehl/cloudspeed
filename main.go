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

// Stack is a helper struct to store  the template body, stack id, and stack name in one place
type Stack struct {
	Template string // template body
	Name     string // name of stack
	ID       string // id of stack
}

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

func (stack *Stack) create(cf *cloudformation.CloudFormation) (string, error) {
	input := &cloudformation.CreateStackInput{
		StackName:    aws.String(stack.Name),
		TemplateBody: aws.String(stack.Template),
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

// describe checks the current status of the stack and returns an error if the stack
func (stack *Stack) describe(cf *cloudformation.CloudFormation) error {
	for {
		out, err := cf.DescribeStacks(&cloudformation.DescribeStacksInput{
			NextToken: aws.String("1"),
			StackName: aws.String(stack.ID),
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

// delete deletes the current stack and returns an error if failing to do so
func (stack *Stack) delete(cf *cloudformation.CloudFormation) error {
	_, err := cf.DeleteStack(&cloudformation.DeleteStackInput{
		StackName: aws.String(stack.ID),
	})
	return err
}
func main() {
	var templateLocation = flag.String("template", "", "the file location of the template")
	var stackName = flag.String("name", "CloudValidate", "the name of the stack to create. Defaults to CloudValidate")
	var persist = *flag.Bool("persist", false, "persist will persist the stack if successful. Defaults to false, deleting the stack after completion")
	flag.Parse()
	fmt.Println(persist)
	if *templateLocation == "" {
		color.Red("Please specify the location of the template file")
		os.Exit(1)
	}
	template, err := getTemplate(*templateLocation)
	if err != nil {
		color.Red(err.Error())
		os.Exit(1)
	}
	stack := &Stack{
		Template: template,
		Name:     *stackName,
		ID:       "",
	}
	session, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
	})
	if err != nil {
		color.Red(err.Error())
		os.Exit(1)
	}
	cf := cloudformation.New(session, aws.NewConfig().WithRegion(*session.Config.Region))
	stackID, err := stack.create(cf)
	if err != nil {
		color.Red("Failed to create stack")
		color.Red(err.Error())
		os.Exit(1)
	}
	stack.ID = stackID
	err = stack.describe(cf)
	for err != nil {
		color.Red("Failed to describe stack")
		color.Red(err.Error())
		err = stack.describe(cf)
	}
	if persist == true {
		return
	}
	color.Yellow("Deleting stack")
	if err = stack.delete(cf); err != nil {
		color.Red(err.Error())
	}
}
