package cmd

import (
	"io"

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

func (o *DiagnosticOptions) Run() {
	// Get the JX version
	log.Info("*** Jenkins-X Version ***\n")
	output, err := o.getCommandOutput("", "jx", "version", "--no-version-check")
	if err != nil {
		log.Error("Unable to get Jenkins-X version info.\n")
	} else {
		log.Infof("%s\n", util.ColorInfo(output))
	}

	// Print the Jenkins-X status
	log.Info("\n*** Jenkins-X Status ***\n")
	output, err = o.getCommandOutput("", "jx", "status")
	if err != nil {
		log.Error("Failed to get the status of the Jenkins-X installation.")
	} else {
		log.Infof("%s\n", util.ColorInfo(output))
	}

	// Get the PVCs in the current namespace
	log.Info("\n*** Kubernetes PVCs ***\n")
	output, err = o.getCommandOutput("", "kubectl", "get", "pvc")
	if err != nil {
		log.Error("Failed to get the Kubernetes PVCs.\n")
	} else {
		log.Infof("%s\n", util.ColorInfo(output))
	}

	// Get the pods in the current namespace
	log.Info("\n*** Kubernetes Pods ***\n")
	output, err = o.getCommandOutput("", "kubectl", "get", "po")
	if err != nil {
		log.Error("Unable to get the Kubernetes pods.\n")
	} else {
		log.Infof("%s\n", util.ColorInfo(output))
	}

	log.Info("\n*** Kubernetes Services ***\n")
	output, err = o.getCommandOutput("", "kubectl", "get", "svc")
	if err != nil {
		log.Error("Unable to get the Kubernetes services.\n")
	} else {
		log.Infof("%s\n", util.ColorInfo(output))
	}

	log.Info("\nPlease visit https://jenkins-x.io/faq/issues/ for any known issues.")
	log.Info("\nFinished printing diagnostic information.")
}
