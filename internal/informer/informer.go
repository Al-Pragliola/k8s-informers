package informer

import (
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/tools/cache"
	"time"
)

func StartPodInformer(client *KubernetesClient) error {
	informerFactory := informers.NewSharedInformerFactory(client.Client, time.Second*30)
	podInformer := informerFactory.Core().V1().Pods()
	podInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    func(new interface{}) {},
		UpdateFunc: func(old, new interface{}) {},
		DeleteFunc: func(obj interface{}) {},
	})
	informerFactory.Start(wait.NeverStop)
	informerFactory.WaitForCacheSync(wait.NeverStop)
	pod, err := podInformer.Lister().Pods("informers-test").Get("prova")
	if err != nil {
		return err
	}

	logrus.Printf("pods: %+v", pod)

	return nil
}
