package istio

import (
	"context"
	"log"

	"ea-controller/internal/constants"
	"ea-controller/internal/model"

	istioclient "istio.io/client-go/pkg/clientset/versioned"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Apply updates Istio VirtualService weights
func Apply(s *model.Service, weights []int, client *istioclient.Clientset) {
	ctx := context.TODO()
	name := s.Name + constants.VirtualServiceSuffix

	vs, err := client.NetworkingV1alpha3().
		VirtualServices(s.Namespace).
		Get(ctx, name, metav1.GetOptions{})

	if err != nil {
		log.Println(err)
		return
	}

	routes := vs.Spec.Http[0].Route
	changed := false
	updated := map[string]int{}

	for i := range routes {
		for j, v := range s.Variants {
			if routes[i].Destination.Subset == v.Name {
				newW := int32(weights[j])
				if routes[i].Weight != newW {
					routes[i].Weight = newW
					changed = true
				}
				updated[v.Name] = int(newW)
			}
		}
	}

	if !changed {
		log.Printf("[SERVICE %s] No change in weights\n", s.Name)
		return
	}

	vs.Spec.Http[0].Route = routes

	_, err = client.NetworkingV1alpha3().
		VirtualServices(s.Namespace).
		Update(ctx, vs, metav1.UpdateOptions{})

	if err != nil {
		log.Println(err)
		return
	}

	log.Printf("[SERVICE %s] Updated weights:\n", s.Name)
	for _, v := range s.Variants {
		if w, ok := updated[v.Name]; ok {
			log.Printf("  - %s: %d", v.Name, w)
		}
	}
}
