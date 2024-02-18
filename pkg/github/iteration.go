package github

import (
	"fmt"

	"github.com/cli/go-gh/v2/pkg/api"
	graphql "github.com/cli/shurcooL-graphql"
)

// ProjectV2IterationFieldIteration
// https://docs.github.com/en/graphql/reference/objects#projectv2iterationfielditeration
type ProjectV2IterationFieldIteration struct {
	Duration  int    `json:"duration"`
	ID        string `json:"id"`
	StartDate string `json:"startDate"`
	Title     string `json:"title"`
}

// ProjectV2IterationField
// https://docs.github.com/en/graphql/reference/objects#projectv2iterationfield
type ProjectV2IterationField struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	Configuration struct {
		CompletedIterations []ProjectV2IterationFieldIteration `json:"completedIterations"`
		Iterations          []ProjectV2IterationFieldIteration `json:"iterations"`
	} `json:"configuration"`
}

func FetchIterationFieldByName(projectID string, fieldName string) (*ProjectV2IterationField, error) {
	client, err := api.DefaultGraphQLClient()
	if err != nil {
		return nil, fmt.Errorf("failed to init GraphQL client: %w", err)
	}

	var query struct {
		Node struct {
			ProjectV2 struct {
				Field struct {
					ProjectV2IterationField ProjectV2IterationField `graphql:"... on ProjectV2IterationField"`
				} `graphql:"field(name: $field_name)"`
			} `graphql:"... on ProjectV2"`
		} `graphql:"node(id: $project_id)"`
	}
	variables := map[string]interface{}{
		"project_id": graphql.ID(projectID),
		"field_name": graphql.String(fieldName),
	}

	err = client.Query("IterationField", &query, variables)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve a iteration field: %w", err)
	}

	return &query.Node.ProjectV2.Field.ProjectV2IterationField, nil
}

func FetchIterationFieldByID(fieldID string) (*ProjectV2IterationField, error) {
	client, err := api.DefaultGraphQLClient()
	if err != nil {
		return nil, fmt.Errorf("failed to init GraphQL client: %w", err)
	}

	var query struct {
		Node struct {
			ProjectV2IterationField ProjectV2IterationField `graphql:"... on ProjectV2IterationField"`
		} `graphql:"node(id: $field_id)"`
	}
	variables := map[string]interface{}{
		"field_id": graphql.ID(fieldID),
	}

	err = client.Query("IterationField", &query, variables)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve a iteration field: %w", err)
	}

	return &query.Node.ProjectV2IterationField, nil
}

func UpdateIterationField(projectID string, fieldID string, itemID string, iterationID string) (string, error) {
	client, err := api.DefaultGraphQLClient()
	if err != nil {
		return "", fmt.Errorf("failed to init GraphQL client: %w", err)
	}
	type ProjectV2Item struct {
		ID string `graphql:"id"`
	}

	var mutation struct {
		UpdateProjectV2ItemFieldValue struct {
			ClientMutationID string        `graphql:"clientMutationId"`
			ProjectV2Item    ProjectV2Item `graphql:"projectV2Item"`
		} `graphql:"updateProjectV2ItemFieldValue(input: $input)"`
	}
	type ProjectV2FieldValue struct {
		IterationID string `json:"iterationId"`
	}
	// https://docs.github.com/en/graphql/reference/mutations#updateprojectv2itemfieldvalue
	type UpdateProjectV2ItemFieldValueInput struct {
		FieldID   string              `json:"fieldId"`
		ItemID    string              `json:"itemId"`
		ProjectID string              `json:"projectId"`
		Value     ProjectV2FieldValue `json:"value"`
	}

	variables := map[string]interface{}{
		"input": UpdateProjectV2ItemFieldValueInput{
			FieldID:   fieldID,
			ItemID:    itemID,
			ProjectID: projectID,
			Value:     ProjectV2FieldValue{IterationID: iterationID},
		},
	}
	err = client.Mutate("updateProjectV2ItemFieldValue", &mutation, variables)
	if err != nil {
		return "", fmt.Errorf("failed to update the iteration field: %w", err)
	}

	return mutation.UpdateProjectV2ItemFieldValue.ProjectV2Item.ID, nil
}
