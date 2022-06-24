package informer

import (
	"github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
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
		DeleteFunc: func(obj interface{}) {
			podObj := obj.(*v1.Pod)
			logrus.Printf("deleted pod %+v", podObj.Name)
		},
	})
	informerFactory.Start(wait.NeverStop)
	informerFactory.WaitForCacheSync(wait.NeverStop)

	<-wait.NeverStop

	return nil
}
