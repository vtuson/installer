package cmd

import (
	"errors"
	"github.com/spf13/cobra"
	"k8s.io/client-go/kubernetes"
	"strconv"
	"strings"
)

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "validate KubeApps components",
	Long:  `validate KubeApps components`,
	RunE:  validateRun,
}

type ConfigurationVal struct {
	Namespaces []string `json:"namespaces"`
	Endpoints  []string `json:"endpoints"`
}

var kpass_val bool
var passcount, passfail int

func getTestConfig() *ConfigurationVal {
	var testConf ConfigurationVal
	testConf.Namespaces = []string{"kubeless", "kubeapps", "kube-system"}
	testConf.Endpoints = []string{"/", "/api/v1/repos", "/kubeless"}
	return &testConf
}

func kubernetesClient() (*kubernetes.Clientset, error) {

	config, err := buildOutOfClusterConfig()
	if err != nil {
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return clientset, nil
}

func validateRun(cmd *cobra.Command, args []string) error {
	defer Report()
	kpass_val = true

	config := getTestConfig()

	client, err := kubernetesClient()
	if err != nil {
		return err
	}

	for _, n := range config.Namespaces {
		CheckPods(n, client)
		CheckEndpoints(n, client)
	}

	localPort, err := cmd.Flags().GetInt("port")
	if err != nil {
		return err
	}

	url, err := cmd.Flags().GetString("url")
	if err != nil {
		return err
	}
	if !strings.Contains(url, "http") {
		return errors.New("badly formed url")
	}

	for _, p := range config.Endpoints {
		PingPath(p, url+":"+strconv.Itoa(localPort))
	}

	return nil

}

func init() {
	RootCmd.AddCommand(validateCmd)
	validateCmd.Flags().Int("port", 8002, "local port")
	validateCmd.Flags().String("url", "http://localhost", "base url")

}