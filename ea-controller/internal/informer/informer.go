package informer

import (
	"context"
	"log"

	"ea-controller/internal/model"
	"ea-controller/internal/store"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/tools/cache"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Start(dynamicClient dynamic.Interface, namespace string, st *store.Store) {
	gvr := schema.GroupVersionResource{
		Group:    "approximate.io",
		Version:  "v1",
		Resource: "approximateservices",
	}

	lw := &cache.ListWatch{
		ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
			return dynamicClient.Resource(gvr).Namespace(namespace).List(context.TODO(), options)
		},
		WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
			return dynamicClient.Resource(gvr).Namespace(namespace).Watch(context.TODO(), options)
		},
	}

	_, controller := cache.NewInformer(
		lw,
		&unstructured.Unstructured{},
		0,
		cache.ResourceEventHandlerFuncs{

			AddFunc: func(obj interface{}) {
				u := obj.(*unstructured.Unstructured)
				s := convert(u, namespace)
				st.Set(s.Name, s)
				log.Printf("[ADD] %s\n", s.Name)
			},

			UpdateFunc: func(_, newObj interface{}) {
				u := newObj.(*unstructured.Unstructured)
				s := convert(u, namespace)
				st.Set(s.Name, s)
				log.Printf("[UPDATE] %s (budget=%.2f)\n", s.Name, s.Budget)
			},

			DeleteFunc: func(obj interface{}) {
				u := obj.(*unstructured.Unstructured)
				st.Delete(u.GetName())
				log.Printf("[DELETE] %s\n", u.GetName())
			},
		},
	)

	go controller.Run(make(chan struct{}))
}

func convert(obj *unstructured.Unstructured, namespace string) *model.Service {
	spec, _ := obj.Object["spec"].(map[string]interface{})
	variantsRaw, _ := spec["variants"].([]interface{})

	var variants []model.Variant

	for _, v := range variantsRaw {
		vm := v.(map[string]interface{})
		name := vm["name"].(string)
		consumption := toFloat(vm["consumption"])

		variants = append(variants, model.Variant{
			Name:      name,
			Workload:  obj.GetName() + "-" + name,
			Utility:   toFloat(vm["utility"]),
			DefaultCj: consumption,
			PrevCij:   consumption,
		})
	}

	return &model.Service{
		Name:      obj.GetName(),
		Namespace: namespace,
		Variants:  variants,
		Budget:    toFloat(spec["energyBudget"]),
	}
}

func toFloat(v interface{}) float64 {
	switch val := v.(type) {
	case float64:
		return val
	case int64:
		return float64(val)
	case int:
		return float64(val)
	default:
		return 0
	}
}
