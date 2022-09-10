package commands

import (
	"context"

	"github.com/shurcooL/graphql"
	"github.com/urfave/cli"
)

type getOrganizationQuery struct {
	Organization struct {
		Id        graphql.ID
		Name      graphql.String
		CreatedAt graphql.String
		UpdatedAt graphql.String

		Users []struct {
			Id    graphql.ID
			Name  graphql.String
			Email graphql.String
			Roles []graphql.String
		}

		Invites []struct {
			Id        graphql.ID
			Email     graphql.String
			CreatedAt graphql.String
			UpdatedAt graphql.String
		}
	} `graphql:"organization(id: $id)"`
}

func ListOrganizations(cli *cli.Context) error {
	ctx := context.Background()

	gql := getGqlClient(ctx)

	var query struct {
		Organizations []struct {
			Id   graphql.ID
			Name graphql.String
		} `graphql:"organizations"`
	}

	if err := gql.Query(ctx, &query, nil); err != nil {
		return err
	}

	writeOutput(cli, query.Organizations)

	return nil
}

func ShowOrganization(cli *cli.Context) error {
	id := cli.Args().First()

	query, err := getOrganization(id)

	if err != nil {
		return err
	}

	writeOutput(cli, query.Organization)
	return nil
}

func getOrganization(id string) (*getOrganizationQuery, error) {
	ctx := context.Background()

	var query getOrganizationQuery
	vars := map[string]interface{}{
		"id": graphql.ID(id),
	}

	if err := getGqlClient(ctx).Query(ctx, &query, vars); err != nil {
		return nil, err
	}

	return &query, nil
}

func ListUsers(cli *cli.Context) error {
	id := cli.Args().First()

	query, err := getOrganization(id)

	if err != nil {
		return err
	}

	writeOutput(cli, query.Organization.Users)
	return nil
}
