package cmd

import (
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"

	"github.com/urfave/cli/v2"
	s "github.com/yitsushi/totp-cli/internal/storage"
	"github.com/yitsushi/totp-cli/internal/terminal"
)

const maxTokenDisplayLength = 100

// ListCommand is the list subcommand.
func ListCommand() *cli.Command {
	return &cli.Command{
		Name:      "list",
		Usage:     "List all available namespaces or accounts under a namespace.",
		ArgsUsage: "[namespace]",
		Flags: []cli.Flag{
			flagSearch(),
			flagRegex(),
		},
		Action: func(ctx *cli.Context) error {
			term := terminal.New(os.Stdin, os.Stdout, os.Stderr)
			storage := prepareStorage(term)

			err := storage.Prepare()
			if err != nil {
				return err
			}

			searchQuery := ctx.String("search")
			useRegex := ctx.Bool("regex")
			ns := ctx.Args().Get(argSetPrefixPositionNamespace)

			if searchQuery != "" {
				return executeSearch(storage, ns, searchQuery, useRegex)
			}

			return executeList(storage, ns)
		},
	}
}

func executeList(storage s.Storage, ns string) error {
	if len(ns) < 1 {
		for _, namespace := range storage.ListNamespaces() {
			fmt.Printf("%s (Number of accounts: %d)\n", namespace.Name, len(namespace.Accounts))
		}

		return nil
	}

	namespace, err := storage.FindNamespace(ns)
	if err != nil {
		return err
	}

	sort.Slice(namespace.Accounts, func(i, j int) bool {
		if namespace.Accounts[i].Favorite != namespace.Accounts[j].Favorite {
			return namespace.Accounts[i].Favorite
		}

		return namespace.Accounts[i].Name < namespace.Accounts[j].Name
	})

	for _, account := range namespace.Accounts {
		prefix := ""
		if account.Favorite {
			prefix = "★ "
		}

		fmt.Printf("%s%s.%s\n", prefix, namespace.Name, account.Name)
	}

	return nil
}

func executeSearch(storage s.Storage, nsFilter, query string, useRegex bool) error {
	var matchFunc func(string) bool

	if useRegex {
		re, err := regexp.Compile(query)
		if err != nil {
			fmt.Fprintf(os.Stderr, " !!! regex compile error: %s\n", err.Error())
			fmt.Fprintln(os.Stderr, "Falling back to substring matching.")

			useRegex = false
		} else {
			matchFunc = func(s string) bool {
				return re.MatchString(s)
			}
		}
	}

	if !useRegex {
		lowerQuery := strings.ToLower(query)
		matchFunc = func(s string) bool {
			return strings.Contains(strings.ToLower(s), lowerQuery)
		}
	}

	type searchResult struct {
		namespace string
		account   *s.Account
	}

	var results []searchResult

	namespaces := storage.ListNamespaces()
	if nsFilter != "" {
		ns, err := storage.FindNamespace(nsFilter)
		if err != nil {
			return err
		}

		namespaces = []*s.Namespace{ns}
	}

	for _, ns := range namespaces {
		for _, account := range ns.Accounts {
			if matchFunc(account.Name) {
				results = append(results, searchResult{
					namespace: ns.Name,
					account:   account,
				})
			}
		}
	}

	sort.Slice(results, func(i, j int) bool {
		if results[i].namespace != results[j].namespace {
			return results[i].namespace < results[j].namespace
		}

		return results[i].account.Name < results[j].account.Name
	})

	for _, r := range results {
		tokenDisplay := r.account.Token
		if len(r.account.Token) > maxTokenDisplayLength {
			tokenDisplay = fmt.Sprintf("[token length: %d chars]", len(r.account.Token))
		}

		fmt.Printf("%s.%s: %s\n", r.namespace, r.account.Name, tokenDisplay)
	}

	return nil
}
