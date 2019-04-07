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

func createStack(session *session.Session, template string) error {
	svc := cloudformation.New(session, aws.NewConfig().WithRegion(*session.Config.Region))

	input := &cloudformation.CreateStackInput{
		StackName:    aws.String("CloudValidate_Temp"),
		TemplateBody: aws.String(template),
	}
	if err := input.Validate(); err != nil {
		return err
	}
	output, err := svc.CreateStack(input)
	if err != nil {
		return err
	}
	color.Green(output.String())
	return nil
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
	fmt.Println(template)
	fmt.Println("HELLO WORLD")
}
