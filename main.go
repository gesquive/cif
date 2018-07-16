package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/gesquive/cig/certs"
	"github.com/gesquive/cli"
)

var version = "v0.1.0-git"
var dirty = ""

var cfgFile string
var displayVersion string

var showVersion bool
var verbose bool
var debug bool

func main() {
	displayVersion = fmt.Sprintf("cig %s%s",
		version,
		dirty)
	Execute(displayVersion)
}

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:              "cig [flags] <cert_path> [<cert_path>...]",
	Short:            "Generate certificate summary information",
	Long:             `Generate certificate information for multiple PEM formatted certificates`,
	ValidArgs:        []string{"cert_path"},
	PersistentPreRun: runRoot,
	Run:              run,
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(version string) {
	displayVersion = version
	RootCmd.SetHelpTemplate(fmt.Sprintf("%s\nVersion:\n  github.com/gesquive/%s\n",
		RootCmd.HelpTemplate(), displayVersion))
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	// flags:	-i replace in file
	// 			-o output

	RootCmd.PersistentFlags().BoolVar(&showVersion, "version", false,
		"Display the version number and exit")

	RootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false,
		"Print logs to stdout instead of file")

	RootCmd.PersistentFlags().BoolVarP(&debug, "debug", "D", false,
		"Include debug statements in log output")
	RootCmd.PersistentFlags().MarkHidden("debug")

	viper.SetEnvPrefix("cig")
	viper.AutomaticEnv()

}

func runRoot(cmd *cobra.Command, args []string) {
	if debug {
		cli.SetPrintLevel(cli.LevelDebug)
	}
	if showVersion {
		cli.Info(displayVersion)
		os.Exit(0)
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
