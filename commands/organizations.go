package commands

import (
	"context"

	"github.com/shurcooL/graphql"
	"github.com/urfave/cli"
)

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
