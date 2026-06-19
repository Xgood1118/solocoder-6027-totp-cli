package cmd

import (
	"github.com/urfave/cli/v2"
	"github.com/yitsushi/totp-cli/internal/security"
)

func flagAlgorithm() *cli.StringFlag {
	return &cli.StringFlag{
		Name:  "algorithm",
		Value: defaultAlgorithm,
		Usage: "Algorithm to use for HMAC (sha1, sha256, sha512).",
		Action: func(_ *cli.Context, value string) error {
			if value != defaultAlgorithm && value != "sha256" && value != "sha512" {
				return invalidAlgorithmError(value)
			}

			return nil
		},
	}
}

func flagShowRemaining() *cli.BoolFlag {
	return &cli.BoolFlag{
		Name:  "show-remaining",
		Value: false,
		Usage: "Show how much time left until the code will be invalid.",
	}
}

func flagLength() *cli.UintFlag {
	return &cli.UintFlag{
		Name:  "length",
		Value: security.DefaultLength,
		Usage: "Length of the generated token.",
	}
}

func flagPrefix() *cli.StringFlag {
	return &cli.StringFlag{
		Name:  "prefix",
		Value: "",
		Usage: "Prefix for the token.",
	}
}

func flagYesPlease(usage string) *cli.BoolFlag {
	return &cli.BoolFlag{
		Name:  "yes-please",
		Value: false,
		Usage: usage,
	}
}

func flagOutput() *cli.StringFlag {
	return &cli.StringFlag{
		Name:     "output",
		Usage:    "Output file. (REQUIRED)",
		Required: true,
	}
}

func flagInput() *cli.StringFlag {
	return &cli.StringFlag{
		Name:     "input",
		Usage:    "Input YAML file. (REQUIRED)",
		Required: true,
	}
}

func flagFollow() *cli.BoolFlag {
	return &cli.BoolFlag{
		Name:  "follow",
		Value: false,
		Usage: "Generate codes continuously.",
	}
}

func flagClearPrefix() *cli.BoolFlag {
	return &cli.BoolFlag{
		Name:  "clear",
		Value: false,
		Usage: "Clear prefix from account.",
	}
}

func flagTimePeriod() *cli.Int64Flag {
	return &cli.Int64Flag{
		Name:  "time-period",
		Value: security.DefaultTimePeriod,
		Usage: "Time period in seconds.",
	}
}

func flagSearch() *cli.StringFlag {
	return &cli.StringFlag{
		Name:  "search",
		Usage: "Search query for account names.",
	}
}

func flagRegex() *cli.BoolFlag {
	return &cli.BoolFlag{
		Name:  "regex",
		Value: false,
		Usage: "Use regex matching for search.",
	}
}

func flagFavorite() *cli.BoolFlag {
	return &cli.BoolFlag{
		Name:  "favorite",
		Value: false,
		Usage: "Mark account as favorite, or generate all favorite accounts.",
	}
}

func flagNoFavorite() *cli.BoolFlag {
	return &cli.BoolFlag{
		Name:  "no-favorite",
		Value: false,
		Usage: "Remove favorite mark from account.",
	}
}
