/*
Copyright 2021 Rancher Labs, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by main. DO NOT EDIT.

package v1alpha1

import (
	"context"
	"time"

	v1alpha1 "github.com/rancher/fleet/pkg/apis/fleet.cattle.io/v1alpha1"
	"github.com/rancher/lasso/pkg/client"
	"github.com/rancher/lasso/pkg/controller"
	"github.com/rancher/wrangler/pkg/apply"
	"github.com/rancher/wrangler/pkg/condition"
	"github.com/rancher/wrangler/pkg/generic"
	"github.com/rancher/wrangler/pkg/kv"
	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/cache"
)

type ClusterGroupHandler func(string, *v1alpha1.ClusterGroup) (*v1alpha1.ClusterGroup, error)

type ClusterGroupController interface {
	generic.ControllerMeta
	ClusterGroupClient

	OnChange(ctx context.Context, name string, sync ClusterGroupHandler)
	OnRemove(ctx context.Context, name string, sync ClusterGroupHandler)
	Enqueue(namespace, name string)
	EnqueueAfter(namespace, name string, duration time.Duration)

	Cache() ClusterGroupCache
}

type ClusterGroupClient interface {
	Create(*v1alpha1.ClusterGroup) (*v1alpha1.ClusterGroup, error)
	Update(*v1alpha1.ClusterGroup) (*v1alpha1.ClusterGroup, error)
	UpdateStatus(*v1alpha1.ClusterGroup) (*v1alpha1.ClusterGroup, error)
	Delete(namespace, name string, options *metav1.DeleteOptions) error
	Get(namespace, name string, options metav1.GetOptions) (*v1alpha1.ClusterGroup, error)
	List(namespace string, opts metav1.ListOptions) (*v1alpha1.ClusterGroupList, error)
	Watch(namespace string, opts metav1.ListOptions) (watch.Interface, error)
	Patch(namespace, name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.ClusterGroup, err error)
}

type ClusterGroupCache interface {
	Get(namespace, name string) (*v1alpha1.ClusterGroup, error)
	List(namespace string, selector labels.Selector) ([]*v1alpha1.ClusterGroup, error)

	AddIndexer(indexName string, indexer ClusterGroupIndexer)
	GetByIndex(indexName, key string) ([]*v1alpha1.ClusterGroup, error)
}

type ClusterGroupIndexer func(obj *v1alpha1.ClusterGroup) ([]string, error)

type clusterGroupController struct {
	controller    controller.SharedController
	client        *client.Client
	gvk           schema.GroupVersionKind
	groupResource schema.GroupResource
}

func NewClusterGroupController(gvk schema.GroupVersionKind, resource string, namespaced bool, controller controller.SharedControllerFactory) ClusterGroupController {
	c := controller.ForResourceKind(gvk.GroupVersion().WithResource(resource), gvk.Kind, namespaced)
	return &clusterGroupController{
		controller: c,
		client:     c.Client(),
		gvk:        gvk,
		groupResource: schema.GroupResource{
			Group:    gvk.Group,
			Resource: resource,
		},
	}
}

func FromClusterGroupHandlerToHandler(sync ClusterGroupHandler) generic.Handler {
	return func(key string, obj runtime.Object) (ret runtime.Object, err error) {
		var v *v1alpha1.ClusterGroup
		if obj == nil {
			v, err = sync(key, nil)
		} else {
			v, err = sync(key, obj.(*v1alpha1.ClusterGroup))
		}
		if v == nil {
			return nil, err
		}
		return v, err
	}
}

func (c *clusterGroupController) Updater() generic.Updater {
	return func(obj runtime.Object) (runtime.Object, error) {
		newObj, err := c.Update(obj.(*v1alpha1.ClusterGroup))
		if newObj == nil {
			return nil, err
		}
		return newObj, err
	}
}

func UpdateClusterGroupDeepCopyOnChange(client ClusterGroupClient, obj *v1alpha1.ClusterGroup, handler func(obj *v1alpha1.ClusterGroup) (*v1alpha1.ClusterGroup, error)) (*v1alpha1.ClusterGroup, error) {
	if obj == nil {
		return obj, nil
	}

	copyObj := obj.DeepCopy()
	newObj, err := handler(copyObj)
	if newObj != nil {
		copyObj = newObj
	}
	if obj.ResourceVersion == copyObj.ResourceVersion && !equality.Semantic.DeepEqual(obj, copyObj) {
		return client.Update(copyObj)
	}

	return copyObj, err
}

func (c *clusterGroupController) AddGenericHandler(ctx context.Context, name string, handler generic.Handler) {
	c.controller.RegisterHandler(ctx, name, controller.SharedControllerHandlerFunc(handler))
}

func (c *clusterGroupController) AddGenericRemoveHandler(ctx context.Context, name string, handler generic.Handler) {
	c.AddGenericHandler(ctx, name, generic.NewRemoveHandler(name, c.Updater(), handler))
}

func (c *clusterGroupController) OnChange(ctx context.Context, name string, sync ClusterGroupHandler) {
	c.AddGenericHandler(ctx, name, FromClusterGroupHandlerToHandler(sync))
}

func (c *clusterGroupController) OnRemove(ctx context.Context, name string, sync ClusterGroupHandler) {
	c.AddGenericHandler(ctx, name, generic.NewRemoveHandler(name, c.Updater(), FromClusterGroupHandlerToHandler(sync)))
}

func (c *clusterGroupController) Enqueue(namespace, name string) {
	c.controller.Enqueue(namespace, name)
}

func (c *clusterGroupController) EnqueueAfter(namespace, name string, duration time.Duration) {
	c.controller.EnqueueAfter(namespace, name, duration)
}

func (c *clusterGroupController) Informer() cache.SharedIndexInformer {
	return c.controller.Informer()
}

func (c *clusterGroupController) GroupVersionKind() schema.GroupVersionKind {
	return c.gvk
}

func (c *clusterGroupController) Cache() ClusterGroupCache {
	return &clusterGroupCache{
		indexer:  c.Informer().GetIndexer(),
		resource: c.groupResource,
	}
}

func (c *clusterGroupController) Create(obj *v1alpha1.ClusterGroup) (*v1alpha1.ClusterGroup, error) {
	result := &v1alpha1.ClusterGroup{}
	return result, c.client.Create(context.TODO(), obj.Namespace, obj, result, metav1.CreateOptions{})
}

func (c *clusterGroupController) Update(obj *v1alpha1.ClusterGroup) (*v1alpha1.ClusterGroup, error) {
	result := &v1alpha1.ClusterGroup{}
	return result, c.client.Update(context.TODO(), obj.Namespace, obj, result, metav1.UpdateOptions{})
}

func (c *clusterGroupController) UpdateStatus(obj *v1alpha1.ClusterGroup) (*v1alpha1.ClusterGroup, error) {
	result := &v1alpha1.ClusterGroup{}
	return result, c.client.UpdateStatus(context.TODO(), obj.Namespace, obj, result, metav1.UpdateOptions{})
}

func (c *clusterGroupController) Delete(namespace, name string, options *metav1.DeleteOptions) error {
	if options == nil {
		options = &metav1.DeleteOptions{}
	}
	return c.client.Delete(context.TODO(), namespace, name, *options)
}

func (c *clusterGroupController) Get(namespace, name string, options metav1.GetOptions) (*v1alpha1.ClusterGroup, error) {
	result := &v1alpha1.ClusterGroup{}
	return result, c.client.Get(context.TODO(), namespace, name, result, options)
}

func (c *clusterGroupController) List(namespace string, opts metav1.ListOptions) (*v1alpha1.ClusterGroupList, error) {
	result := &v1alpha1.ClusterGroupList{}
	return result, c.client.List(context.TODO(), namespace, result, opts)
}

func (c *clusterGroupController) Watch(namespace string, opts metav1.ListOptions) (watch.Interface, error) {
	return c.client.Watch(context.TODO(), namespace, opts)
}

func (c *clusterGroupController) Patch(namespace, name string, pt types.PatchType, data []byte, subresources ...string) (*v1alpha1.ClusterGroup, error) {
	result := &v1alpha1.ClusterGroup{}
	return result, c.client.Patch(context.TODO(), namespace, name, pt, data, result, metav1.PatchOptions{}, subresources...)
}

type clusterGroupCache struct {
	indexer  cache.Indexer
	resource schema.GroupResource
}

func (c *clusterGroupCache) Get(namespace, name string) (*v1alpha1.ClusterGroup, error) {
	obj, exists, err := c.indexer.GetByKey(namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(c.resource, name)
	}
	return obj.(*v1alpha1.ClusterGroup), nil
}

func (c *clusterGroupCache) List(namespace string, selector labels.Selector) (ret []*v1alpha1.ClusterGroup, err error) {

	err = cache.ListAllByNamespace(c.indexer, namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.ClusterGroup))
	})

	return ret, err
}

func (c *clusterGroupCache) AddIndexer(indexName string, indexer ClusterGroupIndexer) {
	utilruntime.Must(c.indexer.AddIndexers(map[string]cache.IndexFunc{
		indexName: func(obj interface{}) (strings []string, e error) {
			return indexer(obj.(*v1alpha1.ClusterGroup))
		},
	}))
}

func (c *clusterGroupCache) GetByIndex(indexName, key string) (result []*v1alpha1.ClusterGroup, err error) {
	objs, err := c.indexer.ByIndex(indexName, key)
	if err != nil {
		return nil, err
	}
	result = make([]*v1alpha1.ClusterGroup, 0, len(objs))
	for _, obj := range objs {
		result = append(result, obj.(*v1alpha1.ClusterGroup))
	}
	return result, nil
}

type ClusterGroupStatusHandler func(obj *v1alpha1.ClusterGroup, status v1alpha1.ClusterGroupStatus) (v1alpha1.ClusterGroupStatus, error)

type ClusterGroupGeneratingHandler func(obj *v1alpha1.ClusterGroup, status v1alpha1.ClusterGroupStatus) ([]runtime.Object, v1alpha1.ClusterGroupStatus, error)

func RegisterClusterGroupStatusHandler(ctx context.Context, controller ClusterGroupController, condition condition.Cond, name string, handler ClusterGroupStatusHandler) {
	statusHandler := &clusterGroupStatusHandler{
		client:    controller,
		condition: condition,
		handler:   handler,
	}
	controller.AddGenericHandler(ctx, name, FromClusterGroupHandlerToHandler(statusHandler.sync))
}

func RegisterClusterGroupGeneratingHandler(ctx context.Context, controller ClusterGroupController, apply apply.Apply,
	condition condition.Cond, name string, handler ClusterGroupGeneratingHandler, opts *generic.GeneratingHandlerOptions) {
	statusHandler := &clusterGroupGeneratingHandler{
		ClusterGroupGeneratingHandler: handler,
		apply:                         apply,
		name:                          name,
		gvk:                           controller.GroupVersionKind(),
	}
	if opts != nil {
		statusHandler.opts = *opts
	}
	controller.OnChange(ctx, name, statusHandler.Remove)
	RegisterClusterGroupStatusHandler(ctx, controller, condition, name, statusHandler.Handle)
}

type clusterGroupStatusHandler struct {
	client    ClusterGroupClient
	condition condition.Cond
	handler   ClusterGroupStatusHandler
}

func (a *clusterGroupStatusHandler) sync(key string, obj *v1alpha1.ClusterGroup) (*v1alpha1.ClusterGroup, error) {
	if obj == nil {
		return obj, nil
	}

	origStatus := obj.Status.DeepCopy()
	obj = obj.DeepCopy()
	newStatus, err := a.handler(obj, obj.Status)
	if err != nil {
		// Revert to old status on error
		newStatus = *origStatus.DeepCopy()
	}

	if a.condition != "" {
		if errors.IsConflict(err) {
			a.condition.SetError(&newStatus, "", nil)
		} else {
			a.condition.SetError(&newStatus, "", err)
		}
	}
	if !equality.Semantic.DeepEqual(origStatus, &newStatus) {
		if a.condition != "" {
			// Since status has changed, update the lastUpdatedTime
			a.condition.LastUpdated(&newStatus, time.Now().UTC().Format(time.RFC3339))
		}

		var newErr error
		obj.Status = newStatus
		newObj, newErr := a.client.UpdateStatus(obj)
		if err == nil {
			err = newErr
		}
		if newErr == nil {
			obj = newObj
		}
	}
	return obj, err
}

type clusterGroupGeneratingHandler struct {
	ClusterGroupGeneratingHandler
	apply apply.Apply
	opts  generic.GeneratingHandlerOptions
	gvk   schema.GroupVersionKind
	name  string
}

func (a *clusterGroupGeneratingHandler) Remove(key string, obj *v1alpha1.ClusterGroup) (*v1alpha1.ClusterGroup, error) {
	if obj != nil {
		return obj, nil
	}

	obj = &v1alpha1.ClusterGroup{}
	obj.Namespace, obj.Name = kv.RSplit(key, "/")
	obj.SetGroupVersionKind(a.gvk)

	return nil, generic.ConfigureApplyForObject(a.apply, obj, &a.opts).
		WithOwner(obj).
		WithSetID(a.name).
		ApplyObjects()
}

func (a *clusterGroupGeneratingHandler) Handle(obj *v1alpha1.ClusterGroup, status v1alpha1.ClusterGroupStatus) (v1alpha1.ClusterGroupStatus, error) {
	if !obj.DeletionTimestamp.IsZero() {
		return status, nil
	}

	objs, newStatus, err := a.ClusterGroupGeneratingHandler(obj, status)
	if err != nil {
		return newStatus, err
	}

	return newStatus, generic.ConfigureApplyForObject(a.apply, obj, &a.opts).
		WithOwner(obj).
		WithSetID(a.name).
		ApplyObjects(objs...)
}