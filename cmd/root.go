package cmd

import (
	"fmt"
	"github.com/Al-Pragliola/k8s-informers/internal/informer"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "inform",
	Short: "informers example from Programming Kubernetes",
	Long:  "informers example from Programming Kubernetes",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		logrus.SetLevel(logrus.DebugLevel)
	},
}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "run the example",
	Long:  "run the example",
	RunE: func(cmd *cobra.Command, args []string) error {
		logrus.Printf("kubeconfig is %s \n", kubeconfig)

		if kubeconfig == "" {
			logrus.Errorln("missing kubeconfig path")
			return fmt.Errorf("missing kubeconfig path")
		}

		logrus.Println("creating kubeclient...")
		kc := informer.NewClient(kubeconfig)

		logrus.Println("initializing kubeclient...")
		err := kc.Init()
		if err != nil {
			return err
		}

		return informer.StartPodInformer(kc)
	},
}

var kubeconfig string

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		logrus.Errorln(err)
		os.Exit(1)
	}
}

func init() {
	runCmd.PersistentFlags().StringVar(
		&kubeconfig,
		"kubeconfig",
		"",
		"path to the kubeconfig, if needed",
	)
	rootCmd.AddCommand(runCmd)
}
