package cmd

import (
	"io"

	"github.com/jenkins-x/jx/pkg/kube"
	"github.com/jenkins-x/jx/pkg/log"
	"github.com/jenkins-x/jx/pkg/util"
	"github.com/spf13/cobra"
	"gopkg.in/AlecAivazis/survey.v1/terminal"
)

type DiagnosticOptions struct {
	CommonOptions
}

// NewCmdDiagnostic creates a new command "jx diagnostic"
func NewCmdDiagnostic(f Factory, in terminal.FileReader, out terminal.FileWriter, errOut io.Writer) *cobra.Command {
	options := &DiagnosticOptions{
		CommonOptions: CommonOptions{
			Factory: f,
			In:      in,
			Out:     out,
			Err:     errOut,
		},
	}

	cmd := &cobra.Command{
		Use:   "diagnostic",
		Short: "Print diagnostic information about the Jenkins-X installation",
		Run: func(cmd *cobra.Command, args []string) {
			options.Cmd = cmd
			options.Args = args
			err := options.Run()
			CheckErr(err)
		},
	}
	options.addCommonFlags(cmd)
	return cmd
}

func (o *DiagnosticOptions) Run() error {
	// Display the current namespace
	config, _, err := kube.LoadConfig()
	if err != nil {
		return err
	}
	currentNS := kube.CurrentNamespace(config)
	log.Infof("Running in namespace: %s", util.ColorInfo(currentNS))

	err = printStatus(o, "Jenkins-X Version", "jx", "version", "--no-version-check")
	if err != nil {
		return err
	}

	err = printStatus(o, "Jenkins-X Status", "jx", "status")
	if err != nil {
		return err
	}

	err = printStatus(o, "Kubernetes PVCs", "kubectl", "get", "pvc")
	if err != nil {
		return err
	}

	err = printStatus(o, "Kubernetes Pods", "kubectl", "get", "po")
	if err != nil {
		return err
	}

	err = printStatus(o, "Kubernetes Services", "kubectl", "get", "svc")
	if err != nil {
		return err
	}

	log.Info("\nPlease visit https://jenkins-x.io/faq/issues/ for any known issues.")
	log.Info("\nFinished printing diagnostic information.\n")
	return nil
}

// Print the the status of the resource by running the specified command
func printStatus(o *DiagnosticOptions, resource string, command string, options ...string) error {
	output, err := o.getCommandOutput("", command, options...)
	if err != nil {
		log.Errorf("Unable to get the %s", resource)
		return err
	}
	log.Infof("\n%s:\n", resource)
	log.Infof("%s\n", util.ColorInfo(output))
	return nil
}
