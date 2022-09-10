package commands

import (
	"context"

	"github.com/shurcooL/graphql"
	"github.com/urfave/cli"
)

type Product struct {
	ID         graphql.ID
	Name       graphql.String
	Category   graphql.String
	ExternalId *graphql.String
	Dates      *struct {
		Start graphql.String
		End   graphql.String
		Kind  graphql.String
	}
}

func ListProducts(cli *cli.Context) error {
	ctx := context.Background()

	var query struct {
		Products struct {
			Products []Product `graphql:"products(organizationId: $orgId, limit: 500)"`
		}
	}

	vars := map[string]interface{}{
		"orgId": graphql.ID(cli.Parent().String("org")),
	}

	if err := getGqlClient(ctx).Query(ctx, &query, vars); err != nil {
		return err
	}

	writeOutput(cli, query.Products.Products)
	return nil
}
