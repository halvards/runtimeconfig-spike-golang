// Copyright 2018 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/runtimeconfig/v1beta1"

	"github.com/google/uuid"
)

func main() {
	ctx := context.Background()
	fmt.Println("Creating Runtime Configuration API client")
	applicationDefaultCredentials, err := createRuntimeConfigClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create Runtime Configuration API client.\n")
		os.Exit(1)
	}
	runtimeconfigClient := applicationDefaultCredentials.Client
	projectID := applicationDefaultCredentials.ProjectID
	fmt.Printf("Project ID: %s\n", projectID)
	projectPath := "projects/" + projectID

	uuidVal, err := uuid.NewRandom()
	if err != nil {
		log.Fatalf("Failed to create UUID for config name: %v", err)
		os.Exit(1)
	}
	configPath := projectPath + "/configs/config-" + uuidVal.String()
	fmt.Printf("Creating config %s\n", configPath)
	config, err := runtimeconfigClient.Projects.Configs.Create(projectPath, &runtimeconfig.RuntimeConfig{
		Name:        configPath,
		Description: "Configuration created via API",
	}).Do()
	if err != nil {
		log.Fatalf("Failed to create config: %v", err)
		os.Exit(1)
	}
	fmt.Printf("Config: %v\n", config.Name)

	fmt.Printf("Listing existing configs for project ID %s:\n", projectID)
	listConfigsResponse, err := runtimeconfigClient.Projects.Configs.List(projectPath).Do()
	if err != nil {
		log.Fatalf("Failed to retrieve list of configs: %v", err)
		os.Exit(1)
	}
	for _, savedConfig := range listConfigsResponse.Configs {
		fmt.Printf("Saved config: %s\n", savedConfig.Name)
	}

	role := "roles/viewer"
	serviceAccountName := "runtimeconfig-spike@" + projectID + ".iam.gserviceaccount.com"
	fmt.Printf("Setting IAM policy with role \"%s\" for serviceAccount:%s\n", role, serviceAccountName)
	policy, err := runtimeconfigClient.Projects.Configs.SetIamPolicy(configPath, &runtimeconfig.SetIamPolicyRequest{
		Policy: &runtimeconfig.Policy{
			Bindings: []*runtimeconfig.Binding{
				&runtimeconfig.Binding{
					Role:    role,
					Members: []string{"serviceAccount:" + serviceAccountName},
				},
			},
		},
	}).Do()
	if err != nil {
		log.Fatalf("Failed to set config IAM policy: %v", err)
		os.Exit(1)
	}
	policyJSON, err := policy.MarshalJSON()
	if err != nil {
		log.Fatalf("Failed to marshall policy as JSON: %v", err)
	}
	fmt.Printf("Policy:\n%v\n", string(policyJSON))

	variablePath := configPath + "/variables/myvar1"
	fmt.Printf("Creating config variable %s\n", variablePath)
	createdVariable, err := runtimeconfigClient.Projects.Configs.Variables.Create(configPath, &runtimeconfig.Variable{
		Name: variablePath,
		Text: "mysecret1",
	}).Do()
	if err != nil {
		log.Fatalf("Failed to create config variable: %v", err)
	}
	fmt.Printf("Variable name: %s\n", createdVariable.Name)

	fmt.Printf("Getting config variable %s\n", variablePath)
	variable, err := runtimeconfigClient.Projects.Configs.Variables.Get(variablePath).Do()
	if err != nil {
		log.Fatalf("Failed to get config variable: %v", err)
	}
	fmt.Printf("Variable text: %s\n", variable.Text) // Value for base64 encoded binary value

	fmt.Printf("Deleting config variable %s\n", variablePath)
	_, err = runtimeconfigClient.Projects.Configs.Variables.Delete(variablePath).Do()
	if err != nil {
		log.Fatalf("Failed to delete config variable: %v", err)
	}

	fmt.Printf("Deleting config %s\n", configPath)
	_, err = runtimeconfigClient.Projects.Configs.Delete(configPath).Do()
	if err != nil {
		log.Fatalf("Failed to delete config: %v", err)
	}
}

// ApplicationDefaultCredentials contains an authenticated API client
// and the detected project ID
type ApplicationDefaultCredentials struct {
	Client    *runtimeconfig.Service
	ProjectID string
}

func createRuntimeConfigClient(ctx context.Context) (*ApplicationDefaultCredentials, error) {
	credentials, err := authenticate(ctx,
		"https://www.googleapis.com/auth/cloud-platform", // this scope may not be required
		"https://www.googleapis.com/auth/cloudruntimeconfig",
	)
	if err != nil {
		log.Fatalf("Failed to authenticate: %v", err)
		return nil, err
	}
	client := oauth2.NewClient(ctx, credentials.TokenSource)
	runtimeconfigClient, err := runtimeconfig.New(client)
	if err != nil {
		log.Fatalf("Failed to create Runtime Configuration API client: %v", err)
		return nil, err
	}
	adc := &ApplicationDefaultCredentials{
		Client:    runtimeconfigClient,
		ProjectID: credentials.ProjectID,
	}
	return adc, nil
}

func authenticate(ctx context.Context, scopes ...string) (*google.Credentials, error) {
	credentials, err := google.FindDefaultCredentials(ctx, scopes...)
	if err != nil {
		log.Fatalf("Failed to find default credentials: %v", err)
		return nil, err
	}
	if len(credentials.ProjectID) == 0 {
		log.Printf("Could not find project ID from runtime environment, trying the GOOGLE_CLOUD_PROJECT environment variable instead.\n")
		credentials.ProjectID = os.Getenv("GOOGLE_CLOUD_PROJECT")
	}
	if len(credentials.ProjectID) == 0 {
		log.Fatalf("Environment variable GOOGLE_CLOUD_PROJECT must contain your project ID.\n")
		return nil, err
	}
	return credentials, nil
}
