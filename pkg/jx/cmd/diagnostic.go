package cmd

import (
	"io"

	"github.com/jenkins-x/jx/pkg/log"
	"github.com/jenkins-x/jx/pkg/util"
	"github.com/spf13/cobra"
)

type DiagnosticOptions struct {
	CommonOptions
}

func NewCmdDiagnostic(f Factory, out io.Writer, errOut io.Writer) *cobra.Command {
	options := &DiagnosticOptions{
		CommonOptions: CommonOptions{
			Factory: f,
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
	// Get the JX version
	output, err := o.getCommandOutput("", "jx", "version")
	if err != nil {
		return err
	}
	log.Info("*** JX Version ***\n")
	log.Infof("%s\n", util.ColorInfo(output))

	// Get the PVCs in the current namespace
	output, err = o.getCommandOutput("", "kubectl", "get", "pvc")
	if err != nil {
		return err
	}
	log.Info("\n*** Kubernetes PVCs ***\n")
	log.Infof("%s\n", util.ColorInfo(output))

	// Get the pods in the current namespace
	output, err = o.getCommandOutput("", "kubectl", "get", "po")
	if err != nil {
		return err
	}
	log.Info("\n*** Kubernetes Pods ***\n")
	log.Infof("%s\n", util.ColorInfo(output))

	output, err = o.getCommandOutput("", "kubectl", "get", "svc")
	if err != nil {
		return err
	}
	log.Info("\n*** Kubernetes Services ***\n")
	log.Infof("%s\n", util.ColorInfo(output))
	return nil
}
