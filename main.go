package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"runtime"

	"github.com/spf13/cobra"

	"github.com/gesquive/cig/certs"
	"github.com/gesquive/cli"
)

var (
	buildVersion = "v0.1.0-dev"
	buildCommit  = ""
	buildDate    = ""
)

var cfgFile string

var showVersion bool
var debug bool

func main() {
	Execute()
}

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:              "cig [flags] <cert_path> [<cert_path>...]",
	Short:            "Formats PEM certificates in a human readable (mkcert.org) format",
	Long:             `Generate certificate summary information for PEM certificates and output (in mkcert.org format)`,
	ValidArgs:        []string{"cert_path"},
	PersistentPreRun: preRun,
	Run:              run,
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	RootCmd.SetHelpTemplate(fmt.Sprintf("%s\nVersion:\n  github.com/gesquive/cig %s\n",
		RootCmd.HelpTemplate(), buildVersion))
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	// TODO: future flags to implement
	// 		-i replace in file
	// 		-o output
	RootCmd.PersistentFlags().BoolVar(&showVersion, "version", false,
		"Display the version number and exit")

	RootCmd.PersistentFlags().BoolVarP(&debug, "debug", "D", false,
		"Include debug statements in log output")
	RootCmd.PersistentFlags().MarkHidden("debug")
}

func preRun(cmd *cobra.Command, args []string) {
	if showVersion {
		fmt.Printf("github.com/gesquive/cig\n")
		fmt.Printf(" Version:    %s\n", buildVersion)
		if len(buildCommit) > 6 {
			fmt.Printf(" Git Commit: %s\n", buildCommit[:7])
		}
		if buildDate != "" {
			fmt.Printf(" Build Date: %s\n", buildDate)
		}
		fmt.Printf(" Go Version: %s\n", runtime.Version())
		fmt.Printf(" OS/Arch:    %s/%s\n", runtime.GOOS, runtime.GOARCH)
		os.Exit(0)
	}
	if debug {
		cli.SetPrintLevel(cli.LevelDebug)
	}
	cli.Debug("Running with debug turned on")
}

func run(cmd *cobra.Command, args []string) {
	certPaths := args[0:]

	cli.Debug("Reading cert files '%v'", certPaths)

	// Detect if user is piping in text
	fileInput, err := os.Stdin.Stat()
	if err != nil {
		cli.Error(err.Error())
		os.Exit(2)
	}

	annotatedCerts := ""

	pipeFound := fileInput.Mode()&os.ModeNamedPipe != 0
	if pipeFound {
		cli.Debug("Pipe input detected, reading from pipe")
		annotatedCerts, err = parseCertFile(os.Stdin, annotatedCerts)
		if err != nil {
			cli.Error("Failed to parse piped data")
			cli.Error(err.Error())
		}
	}

	errCount := 0
	for _, f := range certPaths {
		certFile, err := os.Open(f)
		if err != nil {
			cli.Errorf("Failed to send '%s'\n\r", f)
			cli.Error(err.Error())
			errCount++
			continue
		}
		defer certFile.Close()
		annotatedCerts, err = parseCertFile(certFile, annotatedCerts)
		if err != nil {
			cli.Errorf("Failed to parse '%s'\n", f)
			cli.Error(err.Error())
			errCount++
		}
	}

	if !pipeFound && len(certPaths) == 0 {
		cli.Warn("No data was piped into or specified on the command line.\n")
		cmd.Usage()
		os.Exit(1)
	} else if errCount == 0 {
		cli.Debug("All files successfully parsed.")
	} else {
		cli.Warn("There were some errors while parsing files.")
	}

	fmt.Fprint(os.Stdout, annotatedCerts)
}

func parseCertFile(certFile *os.File, annotations string) (string, error) {
	certData, err := ioutil.ReadAll(certFile)
	if err != nil {
		return "", err
	}

	certificates, err := certs.DecodePEMBlock(certData)
	if err != nil {
		return "", err
	}

	for _, cert := range certificates {
		annotations += certs.WriteCert(cert)
	}

	return annotations, nil
}
