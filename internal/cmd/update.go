package cmd

import (
	"fmt"
	"os"
	"slices"

	"github.com/urfave/cli/v2"
	s "github.com/yitsushi/totp-cli/internal/storage"
	"github.com/yitsushi/totp-cli/internal/terminal"
)

// UpdateCommand is the subcommand to update a TOTP token options.
func UpdateCommand() *cli.Command {
	return &cli.Command{
		Name:  "update",
		Usage: "Update the TOTP token options.",
		Flags: []cli.Flag{
			flagLength(),
			flagPrefix(),
			flagAlgorithm(),
			flagTimePeriod(),
			flagFavorite(),
			flagNoFavorite(),
		},
		ArgsUsage: "[namespace] [account]",
		Action: func(ctx *cli.Context) error {
			namespaceName := ctx.Args().Get(argSetPrefixPositionNamespace)
			if len(namespaceName) < 1 {
				return CommandError{Message: errMsgNamespaceNotDefined}
			}

			accountName := ctx.Args().Get(argSetPrefixPositionAccount)
			if len(accountName) < 1 {
				return CommandError{Message: "account is not defined"}
			}

			term := terminal.New(os.Stdin, os.Stdout, os.Stderr)
			storage := prepareStorage(term)

			err := storage.Prepare()
			if err != nil {
				return err
			}

			opts := AccountOptions{}

			if slices.Contains(ctx.LocalFlagNames(), "algorithm") {
				opts.Algorithm = ctx.String("algorithm")
			}

			if slices.Contains(ctx.LocalFlagNames(), "length") {
				opts.Length = ctx.Uint("length")
			}

			if slices.Contains(ctx.LocalFlagNames(), "time-period") {
				opts.TimePeriod = ctx.Int64("time-period")
			}

			if slices.Contains(ctx.LocalFlagNames(), "prefix") {
				opts.Prefix = ctx.String("prefix")
			}

			setFavorite := false
			unsetFavorite := false

			if slices.Contains(ctx.LocalFlagNames(), "favorite") {
				setFavorite = ctx.Bool("favorite")
			}

			if slices.Contains(ctx.LocalFlagNames(), "no-favorite") {
				unsetFavorite = ctx.Bool("no-favorite")
			}

			if setFavorite && unsetFavorite {
				return CommandError{Message: "--favorite and --no-favorite cannot be used together"}
			}

			return executeUpdate(storage, namespaceName, accountName, opts, ctx.LocalFlagNames(), setFavorite, unsetFavorite)
		},
	}
}

func executeUpdate(storage s.Storage, nsName, accName string, opts AccountOptions, setFlags []string, setFavorite, unsetFavorite bool) error {
	account, err := getAccount(storage, nsName, accName)
	if err != nil {
		return err
	}

	if slices.Contains(setFlags, "algorithm") {
		account.Algorithm = opts.Algorithm
	}

	if slices.Contains(setFlags, "length") {
		account.Length = opts.Length
	}

	if slices.Contains(setFlags, "time-period") {
		account.TimePeriod = opts.TimePeriod
	}

	if slices.Contains(setFlags, "prefix") {
		account.Prefix = opts.Prefix
	}

	if setFavorite {
		account.Favorite = true
	}

	if unsetFavorite {
		account.Favorite = false
	}

	err = storage.Save()
	if err != nil {
		return fmt.Errorf("failed to save the storage: %w", err)
	}

	return nil
}
