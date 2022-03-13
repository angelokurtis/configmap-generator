package handler

import (
	"context"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/util/workqueue"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/angelokurtis/configmap-generator/api/v1beta1"
)

var enqueueLog = ctrl.Log.WithName("eventhandler").WithName("EnqueueRequestForSourceRef")

type EnqueueRequestForSourceRef struct {
	client client.Client
	scheme *runtime.Scheme
}

func NewEnqueueRequestForSourceRef(client client.Client, scheme *runtime.Scheme) *EnqueueRequestForSourceRef {
	return &EnqueueRequestForSourceRef{client: client, scheme: scheme}
}

func (e *EnqueueRequestForSourceRef) Create(evt event.CreateEvent, q workqueue.RateLimitingInterface) {
	if evt.Object == nil {
		enqueueLog.Error(nil, "CreateEvent received with no metadata", "event", evt)
		return
	}
	list, err := e.list()
	if err != nil {
		enqueueLog.Error(nil, "Failed to list all ConfigMapGenerator", "event", evt)
		return
	}
	for _, resource := range e.filter(list, evt.Object) {
		q.Add(reconcile.Request{NamespacedName: types.NamespacedName{
			Name:      resource.GetName(),
			Namespace: resource.GetNamespace(),
		}})
	}
}

func (e *EnqueueRequestForSourceRef) Update(evt event.UpdateEvent, q workqueue.RateLimitingInterface) {
	list, err := e.list()
	if err != nil {
		enqueueLog.Error(nil, "Failed to list all ConfigMapGenerator", "event", evt)
		return
	}
	switch {
	case evt.ObjectNew != nil:
		for _, resource := range e.filter(list, evt.ObjectNew) {
			q.Add(reconcile.Request{NamespacedName: types.NamespacedName{
				Name:      resource.GetName(),
				Namespace: resource.GetNamespace(),
			}})
		}
	case evt.ObjectOld != nil:
		for _, resource := range e.filter(list, evt.ObjectOld) {
			q.Add(reconcile.Request{NamespacedName: types.NamespacedName{
				Name:      resource.GetName(),
				Namespace: resource.GetNamespace(),
			}})
		}
	default:
		enqueueLog.Error(nil, "UpdateEvent received with no metadata", "event", evt)
	}
}

func (e *EnqueueRequestForSourceRef) Delete(evt event.DeleteEvent, q workqueue.RateLimitingInterface) {
	if evt.Object == nil {
		enqueueLog.Error(nil, "DeleteEvent received with no metadata", "event", evt)
		return
	}
	list, err := e.list()
	if err != nil {
		enqueueLog.Error(nil, "Failed to list all ConfigMapGenerator", "event", evt)
		return
	}
	for _, resource := range e.filter(list, evt.Object) {
		q.Add(reconcile.Request{NamespacedName: types.NamespacedName{
			Name:      resource.GetName(),
			Namespace: resource.GetNamespace(),
		}})
	}
}

func (e *EnqueueRequestForSourceRef) Generic(evt event.GenericEvent, q workqueue.RateLimitingInterface) {
	if evt.Object == nil {
		enqueueLog.Error(nil, "GenericEvent received with no metadata", "event", evt)
		return
	}
	list, err := e.list()
	if err != nil {
		enqueueLog.Error(nil, "Failed to list all ConfigMapGenerator", "event", evt)
		return
	}
	for _, resource := range e.filter(list, evt.Object) {
		q.Add(reconcile.Request{NamespacedName: types.NamespacedName{
			Name:      resource.GetName(),
			Namespace: resource.GetNamespace(),
		}})
	}
}

func (e *EnqueueRequestForSourceRef) list() (*v1beta1.ConfigMapGeneratorList, error) {
	list := &v1beta1.ConfigMapGeneratorList{}
	err := e.client.List(context.Background(), list)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (e *EnqueueRequestForSourceRef) filter(list *v1beta1.ConfigMapGeneratorList, obj client.Object) []v1beta1.ConfigMapGenerator {
	resources := make([]v1beta1.ConfigMapGenerator, 0, 0)
	if list == nil {
		return resources
	}
	gvks, _, _ := e.scheme.ObjectKinds(obj)
	for _, resource := range list.Items {
		ref := resource.GetSourceRef()
		for _, gvk := range gvks {
			version, kind := gvk.ToAPIVersionAndKind()
			sameVersion := ref.APIVersion == "" || ref.APIVersion == version
			sameKind := ref.Kind == kind
			sameNamespace := ref.Namespace == obj.GetNamespace()
			sameName := ref.Name == obj.GetName()
			if sameVersion && sameKind && sameNamespace && sameName {
				resources = append(resources, resource)
			}
		}
	}
	return resources
}
