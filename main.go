package main

import (
	"encoding/json"
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
	Template   string // template body
	Name       string // name of stack
	ID         string // id of stack
	Parameters []*cloudformation.Parameter
}

// StackError is an error implementation that includes failure reason, status, and stack
type StackError struct {
	Stack       *Stack
	Reason      string
	Status      string
	StackEvents []*cloudformation.StackEvent
}

func (e StackError) Error() string {
	return fmt.Sprintf("%v failed\n Status: %v. Reason: %v.\n Events: %v", *e.Stack, e.Status, e.Reason, e.StackEvents)
}

/*
 * readFile reads in the template file specified in location
 * returns the template file as a string or an error
 */
func readFile(location string) (string, error) {
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
		Parameters:   stack.Parameters,
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

// describe checks the current status of the stack
// returns bool to indicate if successfully created and an error to describe the stack failure
func (stack *Stack) describe(cf *cloudformation.CloudFormation) (bool, error) {
	dsOut, err := cf.DescribeStacks(&cloudformation.DescribeStacksInput{
		NextToken: aws.String("1"),
		StackName: aws.String(stack.ID),
	})
	if err != nil {
		return false, StackError{
			Stack:       stack,
			Reason:      err.Error(),
			Status:      "AWS ERROR",
			StackEvents: nil,
		}
	}
	for i, cfStack := range dsOut.Stacks {
		if i > 0 {
			return false, StackError{
				Stack:       stack,
				Reason:      "Multiple stacks returned",
				Status:      "AWS ERROR",
				StackEvents: nil,
			}
		}
		fmt.Println(*cfStack.StackStatus)
		if *cfStack.StackStatus == "CREATE_COMPLETE" {
			return true, nil
		}
		if *cfStack.StackStatus == "ROLLBACK_INITIATED" {
			out, err := cf.DescribeStackEvents(&cloudformation.DescribeStackEventsInput{
				NextToken: aws.String("1"),
				StackName: aws.String(stack.ID),
			})
			if err != nil {
				return false, err
			}
			return false, StackError{
				Stack:       stack,
				Reason:      "Multiple stacks returned",
				Status:      *cfStack.StackStatus,
				StackEvents: out.StackEvents,
			}
		}
	}
	return false, nil
}

// delete deletes the current stack and returns an error if failing to do so
func (stack *Stack) delete(cf *cloudformation.CloudFormation) error {
	_, err := cf.DeleteStack(&cloudformation.DeleteStackInput{
		StackName: aws.String(stack.ID),
	})
	return err
}

// return parameters of cloudformation.Parameter type located in location file
// @param location is the paramters file location
// @returns slice of cloudformation.Parameter pointers and an error
func loadParameters(location string) ([]*cloudformation.Parameter, error) {
	fileBody, err := readFile(location)
	if err != nil {
		return nil, err
	}
	var parameters []*cloudformation.Parameter
	if err = json.Unmarshal([]byte(fileBody), &parameters); err != nil {
		return nil, err
	}
	return parameters, nil
}
func main() {
	var templateLocation = flag.String("template", "", "the file location of the template")
	var parameterLocation = flag.String("parameters", "", "the location of the JSON parameter file. Should contain a JSON array of cloudformation.Parameter objects. See examples/parameters.json for reference. ")
	var stackName = flag.String("name", "CloudValidate", "the name of the stack to create. Defaults to CloudValidate")
	var persist = flag.Bool("persist", false, "persist will persist the stack if successful. Defaults to false, deleting the stack after completion")
	flag.Parse()
	var parameters []*cloudformation.Parameter
	if *templateLocation == "" {
		color.Red("Please specify the location of the template file")
		os.Exit(1)
	}
	template, err := readFile(*templateLocation)
	if err != nil {
		color.Red(err.Error())
		os.Exit(1)
	}
	if *parameterLocation != "" {
		color.Yellow("Loading parameters")
		parameters, err = loadParameters(*parameterLocation)
		if err != nil {
			color.Red(err.Error())
			color.Yellow("Attempting to create stack without parameters")
			parameters = nil
		}
	}
	stack := &Stack{
		Template:   template,
		Name:       *stackName,
		ID:         "",
		Parameters: parameters,
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
	success, err := stack.describe(cf)
	for err != nil || success == false {
		if err != nil { // if it errors out don't retry, instead delete stack
			color.Red("Failed to describe stack")
			color.Red(err.Error())
			break
		}
		success, err = stack.describe(cf)
		time.Sleep(5 * time.Second)
	}
	if err != nil {
		color.Green("Successfully completed creating " + stack.Name)
	}
	// persist should only be considered if the template doesn't fail
	if *persist == true && err == nil {
		return
	}
	color.Yellow("Deleting stack")
	if err = stack.delete(cf); err != nil {
		color.Red(err.Error())
	}
}
