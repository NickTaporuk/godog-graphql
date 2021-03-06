package godogsql

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
	"github.com/go-test/deep"
	"github.com/jakubknejzlik/godog-graphql/graphql"
)

type gqlFeature struct {
	client    *graphql.Client
	query     string
	variables map[string]interface{}
	response  interface{}
	error     *string
}

func (f *gqlFeature) iSendQuery(arg1 *gherkin.DocString) error {
	f.query = arg1.Content
	f.response = nil
	f.error = nil

	ctx := context.Background()
	err := f.client.SendQuery(ctx, f.query, f.variables, &f.response)
	if err != nil {
		_err := err.Error()
		f.error = &_err
	}
	return nil
}

func (f *gqlFeature) iHaveVariables(arg1 *gherkin.DocString) error {
	return json.Unmarshal([]byte(arg1.Content), &f.variables)
}

func (f *gqlFeature) theResponseShouldBe(arg1 *gherkin.DocString) (err error) {
	var expected interface{}
	err = json.Unmarshal([]byte(arg1.Content), &expected)
	if err != nil {
		return
	}

	if diff := deep.Equal(expected, f.response); diff != nil {
		text1, _ := json.MarshalIndent(expected, "", " ")
		text2, _ := json.MarshalIndent(f.response, "", " ")
		err = errors.New(fmt.Sprintf("Expected response: %s \n\nActual response: %s\n", text1, text2))
	}
	return
}

func (f *gqlFeature) theErrorShouldBe(arg1 *gherkin.DocString) (err error) {
	expected := arg1.Content

	if f.error != nil && *f.error != expected {
		err = errors.New(fmt.Sprintf("Expected error: %s \n\nActual error: %s\n", expected, *f.error))
	}
	return
}
func (f *gqlFeature) theErrorShouldBeEmpty() (err error) {
	if f.error != nil {
		err = errors.New(fmt.Sprintf("Expected error to be empty\n\nActual error: %s\n", *f.error))
	}
	return
}
func (f *gqlFeature) theErrorShouldNotBeEmpty() (err error) {
	if f.error == nil {
		err = errors.New(fmt.Sprintf("Expected error to not be empty, but it is nil\n"))
	}
	return
}

func FeatureContext(s *godog.Suite) {
	URL := os.Getenv("GRAPHQL_URL")
	if URL == "" {
		panic(fmt.Errorf("Missing required environment variable GRAPHQL_URL"))
	}

	c, err := graphql.NewClient(URL)
	if err != nil {
		panic(err)
	}
	feature := &gqlFeature{client: c}

	s.Step(`^I send query:$`, feature.iSendQuery)
	s.Step(`^I have variables:$`, feature.iHaveVariables)
	s.Step(`^the response should be:$`, feature.theResponseShouldBe)
	s.Step(`^the error should be:$`, feature.theErrorShouldBe)
	s.Step(`^the error should be empty$`, feature.theErrorShouldBeEmpty)
	s.Step(`^the error should not be empty$`, feature.theErrorShouldNotBeEmpty)
}
