package commands

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strings"
	"text/tabwriter"

	"github.com/rezkit/cli/internal/config"
	"github.com/shurcooL/graphql"
	"github.com/urfave/cli"
	"golang.org/x/oauth2"
	"gopkg.in/yaml.v3"
)

const (
	DefaultGqlEndpoint = "https://api.staging.rezkit.app/api/graphql"
)

func getGqlEndpoint() string {
	return DefaultGqlEndpoint
}

func getGqlClient(ctx context.Context) *graphql.Client {

	auth := oauth2.StaticTokenSource(&oauth2.Token{
		AccessToken:  config.GetConfig().GetString("authentication.access_token"),
		RefreshToken: config.GetConfig().GetString("authentication.refresh_token"),
		Expiry:       config.GetConfig().GetTime("authentication.expires"),
		TokenType:    "bearer",
	})

	httpClient := oauth2.NewClient(ctx, auth)

	return graphql.NewClient(getGqlEndpoint(), httpClient)
}

func writeOutput(ctx *cli.Context, data interface{}) {
	format := ctx.GlobalString("format")

	switch format {
	case "json":
		encoder := json.NewEncoder(os.Stdout)
		encoder.SetIndent("", "\t")
		encoder.Encode(data)
		break
	case "yaml":
		encoder := yaml.NewEncoder(os.Stdout)
		encoder.SetIndent(4)
		encoder.Encode(data)
		break
	case "table":
		printTable(data)
	}
}

func printTable(data interface{}) {
	ref := reflect.TypeOf(data)

	tw := tabwriter.NewWriter(os.Stdout, 4, 2, 1, ' ', 0)

	fields := []string{}

	if ref.Kind() == reflect.Slice {
		// Generate header row from struct fields
		first := ref.Elem()

		if first.Kind() == reflect.Struct {
			fields = make([]string, first.NumField())

			for i := 0; i < first.NumField(); i++ {
				fields[i] = first.Field(i).Name
			}

			fmt.Fprintln(tw, strings.Join(fields, "\t"))
		}

		values := reflect.ValueOf(data)

		for i := 0; i < values.Len(); i++ {
			row := make([]string, len(fields))
			for f, field := range fields {
				row[f] = fmt.Sprintf("%v", values.Index(i).FieldByName(field))
			}

			fmt.Fprintln(tw, strings.Join(row, "\t"))
		}
	}

	tw.Flush()
}
