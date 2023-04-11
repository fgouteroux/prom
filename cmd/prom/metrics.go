package main

import (
	"bytes"
	"fmt"
	"github.com/urfave/cli/v2"
	"os"

	checkMetrics "github.com/fgouteroux/prom/pkg/metrics/check"
	fmtMetrics "github.com/fgouteroux/prom/pkg/metrics/format"
	pushClient "github.com/fgouteroux/prom/pkg/metrics/push/client"
	pushMetrics "github.com/fgouteroux/prom/pkg/metrics/push/metrics"
	"github.com/fgouteroux/prom/pkg/utils"
)

var metricsCommonFlags = []cli.Flag{
	&cli.StringSliceFlag{
		Name:     "files",
		Aliases:  []string{"f"},
		Usage:    "prometheus metrics files.",
		EnvVars:  []string{"PROM_METRICS_FILES"},
		Required: false,
		Value:    nil,
	},
}

var metricsPushFlags = []cli.Flag{
	&cli.StringFlag{
		Name:     "url",
		Aliases:  []string{"u"},
		Usage:    "remote write url",
		Required: true,
	},
	&cli.StringFlag{
		Name:     "job",
		Aliases:  []string{"j"},
		Usage:    "job label to attach to metrics",
		Required: false,
		Value:    "prom/metrics",
	},
	&cli.StringFlag{
		Name:     "tls-ca-file",
		Usage:    "remote write tls ca file.",
		Required: false,
	},
	&cli.StringFlag{
		Name:     "tls-key-file",
		Usage:    "remote write tls key file.",
		Required: false,
	},
	&cli.StringFlag{
		Name:     "tls-cert-file",
		Usage:    "remote write tls cert file.",
		Required: false,
	},
	&cli.BoolFlag{
		Name:     "tls-skip-verify",
		Usage:    "remote write tls skip verify.",
		Required: false,
	},
	&cli.BoolFlag{
		Name:     "enable-HTTP2",
		Usage:    "remote write HTTP2.",
		Required: false,
	},
	&cli.IntFlag{
		Name:     "timeout",
		Aliases:  []string{"t"},
		Usage:    "remote write timeout.",
		Required: false,
		Value:    30,
	},
}

var metricsCheckFlags = []cli.Flag{
	&cli.BoolFlag{
		Name:     "silent",
		Aliases:  []string{"s"},
		Usage:    "silent error(s)",
		Required: false,
	},
}

var metricsFmtFlags = []cli.Flag{
	&cli.StringFlag{
		Name:     "output",
		Aliases:  []string{"o"},
		Usage:    "save output to file (if using STDIN)",
		Required: false,
	},
	&cli.BoolFlag{
		Name:     "write",
		Usage:    "Don't write to source files --write=false (always disabled if using STDIN)",
		Required: false,
		Value:    true,
	},
	&cli.BoolFlag{
		Name:     "diff",
		Usage:    "Display diffs of formatting changes",
		Required: false,
	},
}

func metricsCommand() *cli.Command {
	return &cli.Command{
		Name:  "metrics",
		Usage: "metrics operations",

		Subcommands: []*cli.Command{
			{
				Name:   "push",
				Usage:  "push prometheus metrics to a remote write.",
				Action: pushMetricsAction,
				Flags:  append(metricsCommonFlags, metricsPushFlags...),
			},
			{
				Name:   "check",
				Usage:  "check prometheus metrics.",
				Action: checkMetricsAction,
				Flags:  append(metricsCommonFlags, metricsCheckFlags...),
			},
			{
				Name:   "fmt",
				Usage:  "format prometheus metrics.",
				Action: fmtMetricsAction,
				Flags:  append(metricsCommonFlags, metricsFmtFlags...),
			},
		},
	}
}

func pushMetricsAction(c *cli.Context) error {
	debug := c.Bool("debug")
	url := c.String("url")
	jobLabel := c.String("job")
	tlsCAFile := c.String("tlsCAFile")
	tlsKeyFile := c.String("tlsKeyFile")
	tlsCertFile := c.String("tlsCertFile")
	tlsSkipVerify := c.Bool("tlsSkipVerify")
	enableHTTP2 := c.Bool("enableHTTP2")
	timeout := c.Int("timeout")
	files := c.StringSlice("files")
	headers := c.StringSlice("headers")

	var pushMetricsClient *pushClient.Client
	if tlsCAFile == "" && tlsKeyFile == "" && tlsCertFile == "" && !tlsSkipVerify {
		pushMetricsClient = pushClient.Configure(url, debug, timeout, headers)
	} else {
		pushMetricsClient = pushClient.ConfigureWithTLS(url, tlsCAFile, tlsKeyFile, tlsCertFile, enableHTTP2, debug, tlsSkipVerify, timeout, headers)
	}

	// add empty string to avoid matching filename
	if len(files) == 0 {
		files = append(files, "")
	}

	for _, f := range files {
		data, err := utils.ReadFile(f)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}

		metricsData, err := pushMetrics.ParseTextAndFormat(bytes.NewReader(data), jobLabel)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		err = pushMetricsClient.PushWithRetries(metricsData)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}

	return nil
}

func checkMetricsAction(c *cli.Context) error {
	files := c.StringSlice("files")
	silent := c.Bool("silent")

	// add empty string to avoid matching filename
	if len(files) == 0 {
		files = append(files, "")
	}

	for _, f := range files {
		data, err := utils.ReadFile(f)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}

		exitCode, errs := checkMetrics.Check(bytes.NewReader(data))
		if !silent && exitCode != 0 {
			for _, error := range errs {
				fmt.Fprintln(os.Stderr, error)
			}
		}
		os.Exit(exitCode)
	}

	return nil
}

func fmtMetricsAction(c *cli.Context) error {
	files := c.StringSlice("files")
	output := c.String("output")
	write := c.Bool("write")
	showDiff := c.Bool("diff")

	if len(files) == 0 { // Assuming stdin, then.
		files = append(files, "")
		write = false
	} else if output != "" {
		fmt.Fprintln(os.Stderr, fmt.Errorf("Option --output cannot be used when reading from files"))
		os.Exit(1)
	}

	for _, f := range files {
		data, err := utils.ReadFile(f)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}

		metricsFmt, err := fmtMetrics.Format(string(data))
		if err != nil {
			fmt.Fprintln(os.Stderr, fmt.Errorf("Failed to generate diff for %s: %s", f, err))
			os.Exit(1)
		}

		if showDiff {
			diff, err := utils.BytesDiff(data, []byte(metricsFmt), f)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			if len(diff) > 0 {
				fmt.Println(string(diff))
			}
		}

		if output != "" {
			utils.WritetoFile(output, metricsFmt)
		} else if write {
			utils.WritetoFile(f, metricsFmt)
		}
	}

	return nil
}
