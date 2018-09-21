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

func (o *DiagnosticOptions) Run() error {
	// Run JX version to get the versions of the following: jx cli, helm, kubect
	output, err := o.getCommandOutput("", "jx", "version", "--no-version-check")
	if err != nil {
		return err
	}
	log.Info("*** Jenkins-X Version ***\n")
	log.Infof("%s\n", util.ColorInfo(output))

	// Print the Jenkins-X status
	output, err = o.getCommandOutput("", "jx", "status")
	if err != nil {
		log.Error("Unable to get the status of the Jenkins-X install.")
		return err
	}
	log.Info("\n*** Jenkins-X Status ***\n")
	log.Infof("%s\n", util.ColorInfo(output))

	// Get the PVCs in the current namespace
	output, err = o.getCommandOutput("", "kubectl", "get", "pvc")
	if err != nil {
		log.Error("Unable to get the Kubernetes PVCs")
		return err
	}
	log.Info("\n*** Kubernetes PVCs ***\n")
	log.Infof("%s\n", util.ColorInfo(output))

	// Get the pods in the current namespace
	output, err = o.getCommandOutput("", "kubectl", "get", "po")
	if err != nil {
		log.Error("Unable to get the Kubernetes pods")
		return err
	}
	log.Info("\n*** Kubernetes Pods ***\n")
	log.Infof("%s\n", util.ColorInfo(output))

	// Get the services in the current namespace
	output, err = o.getCommandOutput("", "kubectl", "get", "svc")
	if err != nil {
		log.Error("Unable to get the Kubernetes services")
		return err
	}
	log.Info("\n*** Kubernetes Services ***\n")
	log.Infof("%s\n", util.ColorInfo(output))

	log.Info("\nPlease visit https://jenkins-x.io/faq/issues/ for any known issues.")
	log.Info("\nFinished printing diagnostic information.")
	return nil
}
