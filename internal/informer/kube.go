package informer

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"os"
	"path/filepath"
)

// KubernetesClient represents the Kubernetes configuration of the project.
type KubernetesClient struct {
	KubeConfig      string
	Client          kubernetes.Interface
	DynClient       dynamic.Interface
	DiscoveryClient discovery.DiscoveryInterface
	config          *rest.Config
}

func NewClient(kubeconfig string) *KubernetesClient {
	return &KubernetesClient{
		KubeConfig: kubeconfig,
	}
}

// Init initializes the Kubernetes client-go.
func (kc *KubernetesClient) Init() error {
	var (
		config *rest.Config
		err    error
	)

	if kc.KubeConfig != "" {
		config, err = kc.getConfigFromFile(kc.KubeConfig)
	} else {
		// if no kubeconfigfile is provided creates the in-cluster config
		config, err = kc.inClusterConfig()

		if err != nil {
			logrus.Errorf("In-cluster config failed with err: %s", err)

			// If inCluster config does not work, try with the default kube config path
			logrus.Info("Trying to connect to cluster using default kubeconfig(`.kube/config`)")
			config, err = kc.extClusterConfig()
		}
	}

	if err != nil {
		return err
	}

	kc.config = config

	// return k8s client and err
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}

	kc.Client = client

	dyn, err := dynamic.NewForConfig(config)
	if err != nil {
		return err
	}

	kc.DynClient = dyn

	dis, err := discovery.NewDiscoveryClientForConfig(config)
	if err != nil {
		return err
	}

	kc.DiscoveryClient = dis

	return nil
}

// Config returns the rest configuration to the Kubernetes API.
func (kc *KubernetesClient) Config() *rest.Config {
	return kc.config
}

func (kc *KubernetesClient) inClusterConfig() (*rest.Config, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}

	return config, nil
}

func (kc *KubernetesClient) extClusterConfig() (*rest.Config, error) {
	if home := os.Getenv("HOME"); home != "" {
		kubeConfigPath := filepath.Join(home, ".kube", "config")
		config, err := kc.getConfigFromFile(kubeConfigPath)

		if err != nil {
			return nil, err
		}

		return config, nil
	}

	return nil, fmt.Errorf("can not configure external cluster configuration from the default $HOME/.kube/config path")
}

func (kc *KubernetesClient) getConfigFromFile(kubeConfigPath string) (*rest.Config, error) {
	kubeConfigContent, err := ioutil.ReadFile(kubeConfigPath)
	if err != nil {
		return nil, err
	}

	return clientcmd.RESTConfigFromKubeConfig(kubeConfigContent)
}
