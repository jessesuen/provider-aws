///*
//Copyright 2019 The Crossplane Authors.
//
//Licensed under the Apache License, Version 2.0 (the "License");
//you may not use this file except in compliance with the License.
//You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
//Unless required by applicable law or agreed to in writing, software
//distributed under the License is distributed on an "AS IS" BASIS,
//WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//See the License for the specific language governing permissions and
//limitations under the License.
//*/
//
package rds
//
//import (
//	"context"
//	"testing"
//
//	"github.com/crossplaneio/stack-aws/apis"
//
//	"github.com/google/go-cmp/cmp"
//	. "github.com/onsi/gomega"
//	"github.com/pkg/errors"
//	corev1 "k8s.io/api/core/v1"
//	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
//	"k8s.io/apimachinery/pkg/types"
//	"k8s.io/client-go/kubernetes/scheme"
//	. "sigs.k8s.io/controller-runtime/pkg/client/fake"
//	"sigs.k8s.io/controller-runtime/pkg/reconcile"
//
//	"github.com/crossplaneio/stack-aws/apis/database/v1alpha2"
//	. "github.com/crossplaneio/stack-aws/apis/database/v1alpha2"
//	"github.com/crossplaneio/stack-aws/pkg/clients/rds"
//	. "github.com/crossplaneio/stack-aws/pkg/clients/rds/fake"
//
//	runtimev1alpha1 "github.com/crossplaneio/crossplane-runtime/apis/core/v1alpha1"
//	"github.com/crossplaneio/crossplane-runtime/pkg/resource"
//	"github.com/crossplaneio/crossplane-runtime/pkg/test"
//)
//
//const (
//	namespace    = "default"
//	instanceName = "test-instance"
//
//	masterUserName = "testuser"
//	engine         = "mysql"
//	class          = "db.t2.small"
//	size           = int64(10)
//)
//
//var (
//	key = types.NamespacedName{
//		Name: instanceName,
//	}
//	request = reconcile.Request{
//		NamespacedName: key,
//	}
//)
//
//func init() {
//	if err := apis.AddToScheme(scheme.Scheme); err != nil {
//		panic(err)
//	}
//}
//
//func testResource() *RDSInstance {
//	return &RDSInstance{
//		ObjectMeta: metav1.ObjectMeta{
//			Name: instanceName,
//		},
//		Spec: RDSInstanceSpec{
//			ResourceSpec: runtimev1alpha1.ResourceSpec{
//				ProviderReference: &corev1.ObjectReference{},
//			},
//			RDSInstanceParameters: v1alpha2.RDSInstanceParameters{
//				MasterUsername: masterUserName,
//				Engine:         engine,
//				Class:          class,
//				Size:           size,
//			},
//		},
//	}
//}
//
//// assertResource a helper function to check on cluster and its status
//func assertResource(g *GomegaWithT, r *Reconciler, s runtimev1alpha1.ConditionedStatus) *RDSInstance {
//	resource := &RDSInstance{}
//	err := r.Get(ctx, key, resource)
//	g.Expect(err).To(BeNil())
//	g.Expect(cmp.Diff(s, resource.Status.ConditionedStatus, test.EquateConditions())).Should(BeZero())
//	return resource
//}
//
//func TestSyncClusterError(t *testing.T) {
//	g := NewGomegaWithT(t)
//
//	assert := func(instance *RDSInstance, client rds.Client, expectedResult reconcile.Result, expectedStatus runtimev1alpha1.ConditionedStatus) {
//		r := &Reconciler{
//			Client:                     NewFakeClient(instance),
//			ManagedReferenceResolver:   resource.NewAPIManagedReferenceResolver(NewFakeClient(instance)),
//			ManagedConnectionPublisher: resource.PublisherChain{}, // A no-op publisher.
//		}
//
//		rs, err := r._sync(instance, client)
//
//		g.Expect(err).NotTo(HaveOccurred())
//		g.Expect(rs).To(Equal(expectedResult))
//		assertResource(g, r, expectedStatus)
//	}
//
//	// get error
//	testError := errors.New("test-resource-retrieve-error")
//	cl := &MockRDSClient{
//		MockGetInstance: func(s string) (instance *rds.Instance, e error) {
//			return nil, testError
//		},
//	}
//	expectedStatus := runtimev1alpha1.ConditionedStatus{}
//	expectedStatus.SetConditions(runtimev1alpha1.ReconcileError(testError))
//	assert(testResource(), cl, resultRequeue, expectedStatus)
//
//	// instance is not ready
//	cl = &MockRDSClient{
//		MockGetInstance: func(s string) (instance *rds.Instance, e error) {
//			return &rds.Instance{
//				Status: string(RDSInstanceStateCreating),
//			}, nil
//		},
//	}
//	expectedStatus = runtimev1alpha1.ConditionedStatus{}
//	expectedStatus.SetConditions(runtimev1alpha1.Creating(), runtimev1alpha1.ReconcileSuccess())
//	assert(testResource(), cl, resultRequeue, expectedStatus)
//
//	// instance in failed state
//	cl = &MockRDSClient{
//		MockGetInstance: func(s string) (instance *rds.Instance, e error) {
//			return &rds.Instance{
//				Status: string(RDSInstanceStateFailed),
//			}, nil
//		},
//	}
//	expectedStatus = runtimev1alpha1.ConditionedStatus{}
//	expectedStatus.SetConditions(runtimev1alpha1.Unavailable(), runtimev1alpha1.ReconcileSuccess())
//	assert(testResource(), cl, result, expectedStatus)
//
//	// instance is in deleting state
//	cl = &MockRDSClient{
//		MockGetInstance: func(s string) (instance *rds.Instance, e error) {
//			return &rds.Instance{
//				Status: string(RDSInstanceStateDeleting),
//			}, nil
//		},
//	}
//
//	expectedStatus = runtimev1alpha1.ConditionedStatus{}
//	expectedStatus.SetConditions(
//		runtimev1alpha1.ReconcileError(errors.Errorf("unexpected resource status: %s", RDSInstanceStateDeleting)),
//	)
//	assert(testResource(), cl, resultRequeue, expectedStatus)
//
//	// failed to retrieve instance secret
//	cl = &MockRDSClient{
//		MockGetInstance: func(s string) (instance *rds.Instance, e error) {
//			return &rds.Instance{
//				Status: string(RDSInstanceStateAvailable),
//			}, nil
//		},
//	}
//
//	testError = errors.New("test-error-publish-secret")
//
//	assert = func(instance *RDSInstance, client rds.Client, expectedResult reconcile.Result, expectedStatus runtimev1alpha1.ConditionedStatus) {
//		r := &Reconciler{
//			Client:                   NewFakeClient(instance),
//			ManagedReferenceResolver: resource.NewAPIManagedReferenceResolver(NewFakeClient(instance)),
//			ManagedConnectionPublisher: resource.ManagedConnectionPublisherFns{
//				PublishConnectionFn: func(_ context.Context, _ resource.Managed, _ resource.ConnectionDetails) error {
//					return testError
//				},
//			},
//		}
//
//		rs, err := r._sync(instance, client)
//
//		g.Expect(err).NotTo(HaveOccurred())
//		g.Expect(rs).To(Equal(expectedResult))
//		assertResource(g, r, expectedStatus)
//	}
//
//	tr := testResource()
//	expectedStatus = runtimev1alpha1.ConditionedStatus{}
//	expectedStatus.SetConditions(runtimev1alpha1.Available(), runtimev1alpha1.ReconcileError(testError))
//	assert(tr, cl, resultRequeue, expectedStatus)
//
//}
//
//func TestSyncClusterUpdateSecretFailure(t *testing.T) {
//	g := NewGomegaWithT(t)
//
//	tr := testResource()
//
//	testError := errors.New("test-error-publish-secret")
//
//	r := &Reconciler{
//		Client:                   NewFakeClient(tr),
//		ManagedReferenceResolver: resource.NewAPIManagedReferenceResolver(NewFakeClient(tr)),
//		ManagedConnectionPublisher: resource.ManagedConnectionPublisherFns{
//			PublishConnectionFn: func(_ context.Context, _ resource.Managed, _ resource.ConnectionDetails) error {
//				return testError
//			},
//		},
//	}
//
//	called := false
//	cl := &MockRDSClient{
//		MockGetInstance: func(s string) (instance *rds.Instance, e error) {
//			called = true
//			return &rds.Instance{
//				Status: string(RDSInstanceStateAvailable),
//			}, nil
//		},
//	}
//
//	expectedStatus := runtimev1alpha1.ConditionedStatus{}
//	expectedStatus.SetConditions(runtimev1alpha1.Available(), runtimev1alpha1.ReconcileError(testError))
//
//	rs, err := r._sync(tr, cl)
//	g.Expect(rs).To(Equal(resultRequeue))
//	g.Expect(err).NotTo(HaveOccurred())
//	g.Expect(called).To(BeTrue())
//	assertResource(g, r, expectedStatus)
//}
//
//func TestSyncCluster(t *testing.T) {
//	g := NewGomegaWithT(t)
//
//	tr := testResource()
//
//	r := &Reconciler{
//		Client:                     NewFakeClient(tr),
//		ManagedReferenceResolver:   resource.NewAPIManagedReferenceResolver(NewFakeClient(tr)),
//		ManagedConnectionPublisher: resource.PublisherChain{}, // A no-op publisher.
//	}
//
//	called := false
//	cl := &MockRDSClient{
//		MockGetInstance: func(s string) (instance *rds.Instance, e error) {
//			called = true
//			return &rds.Instance{
//				Status: string(RDSInstanceStateAvailable),
//			}, nil
//		},
//	}
//
//	expectedStatus := runtimev1alpha1.ConditionedStatus{}
//	expectedStatus.SetConditions(runtimev1alpha1.Available(), runtimev1alpha1.ReconcileSuccess())
//
//	rs, err := r._sync(tr, cl)
//	g.Expect(rs).To(Equal(result))
//	g.Expect(err).NotTo(HaveOccurred())
//	g.Expect(called).To(BeTrue())
//	rr := assertResource(g, r, expectedStatus)
//	g.Expect(rr.Status.State).To(Equal(string(RDSInstanceStateAvailable)))
//}
//
//func TestDelete(t *testing.T) {
//	g := NewGomegaWithT(t)
//
//	tr := testResource()
//
//	r := &Reconciler{
//		Client:                   NewFakeClient(tr),
//		ManagedReferenceResolver: resource.NewAPIManagedReferenceResolver(NewFakeClient(tr)),
//	}
//
//	cl := &MockRDSClient{}
//
//	// test delete w/ reclaim policy
//	tr.Spec.ReclaimPolicy = runtimev1alpha1.ReclaimRetain
//	expectedStatus := runtimev1alpha1.ConditionedStatus{}
//	expectedStatus.SetConditions(runtimev1alpha1.Deleting(), runtimev1alpha1.ReconcileSuccess())
//
//	rs, err := r._delete(tr, cl)
//	g.Expect(rs).To(Equal(result))
//	g.Expect(err).NotTo(HaveOccurred())
//	assertResource(g, r, expectedStatus)
//
//	// test delete w/ delete policy
//	tr.Spec.ReclaimPolicy = runtimev1alpha1.ReclaimDelete
//	called := false
//	cl.MockDeleteInstance = func(name string) (instance *rds.Instance, e error) {
//		called = true
//		return nil, nil
//	}
//
//	rs, err = r._delete(tr, cl)
//	g.Expect(rs).To(Equal(result))
//	g.Expect(err).NotTo(HaveOccurred())
//	g.Expect(called).To(BeTrue())
//	assertResource(g, r, expectedStatus)
//
//	// test delete w/ delete policy and delete error
//	testError := errors.New("test-delete-error")
//	called = false
//	cl.MockDeleteInstance = func(name string) (instance *rds.Instance, e error) {
//		called = true
//		return nil, testError
//	}
//	expectedStatus = runtimev1alpha1.ConditionedStatus{}
//	expectedStatus.SetConditions(runtimev1alpha1.Deleting(), runtimev1alpha1.ReconcileError(testError))
//
//	rs, err = r._delete(tr, cl)
//	g.Expect(rs).To(Equal(resultRequeue))
//	g.Expect(err).NotTo(HaveOccurred())
//	g.Expect(called).To(BeTrue())
//	assertResource(g, r, expectedStatus)
//}
//
//func TestCreate(t *testing.T) {
//	g := NewGomegaWithT(t)
//
//	tr := testResource()
//
//	r := &Reconciler{
//		Client:                     NewFakeClient(tr),
//		ManagedReferenceResolver:   resource.NewAPIManagedReferenceResolver(NewFakeClient(tr)),
//		ManagedConnectionPublisher: resource.PublisherChain{}, // A no-op publisher.
//	}
//
//	called := false
//	cl := &MockRDSClient{
//		MockCreateInstance: func(s string, s2 string, spec *RDSInstanceSpec) (instance *rds.Instance, e error) {
//			called = true
//			return nil, nil
//		},
//	}
//
//	expectedStatus := runtimev1alpha1.ConditionedStatus{}
//	expectedStatus.SetConditions(runtimev1alpha1.Creating(), runtimev1alpha1.ReconcileSuccess())
//
//	rs, err := r._create(testResource(), cl)
//	g.Expect(rs).To(Equal(resultRequeue))
//	g.Expect(err).NotTo(HaveOccurred())
//	g.Expect(called).To(BeTrue())
//	assertResource(g, r, expectedStatus)
//}
//
//func TestCreateFail(t *testing.T) {
//	g := NewGomegaWithT(t)
//	tr := testResource()
//	cl := &MockRDSClient{}
//
//	called := false
//	cl.MockCreateInstance = func(s string, s2 string, spec *RDSInstanceSpec) (instance *rds.Instance, e error) {
//		called = true
//		return nil, nil
//	}
//
//	testError := errors.New("test-publish-secret-error")
//
//	r := &Reconciler{
//		Client:                   NewFakeClient(tr),
//		ManagedReferenceResolver: resource.NewAPIManagedReferenceResolver(NewFakeClient(tr)),
//		ManagedConnectionPublisher: resource.ManagedConnectionPublisherFns{
//			PublishConnectionFn: func(_ context.Context, _ resource.Managed, _ resource.ConnectionDetails) error {
//				return testError
//			},
//		},
//	}
//
//	expectedStatus := runtimev1alpha1.ConditionedStatus{}
//	expectedStatus.SetConditions(runtimev1alpha1.Creating(), runtimev1alpha1.ReconcileError(testError))
//
//	rs, err := r._create(tr, cl)
//	g.Expect(rs).To(Equal(resultRequeue))
//	g.Expect(err).NotTo(HaveOccurred())
//	assertResource(g, r, expectedStatus)
//
//	// test create resource error
//	tr = testResource()
//	called = false
//	testError = errors.New("test-create-error")
//	cl.MockCreateInstance = func(s string, s2 string, spec *RDSInstanceSpec) (instance *rds.Instance, e error) {
//		called = true
//		return nil, testError
//	}
//
//	expectedStatus = runtimev1alpha1.ConditionedStatus{}
//	expectedStatus.SetConditions(runtimev1alpha1.Creating(), runtimev1alpha1.ReconcileError(testError))
//
//	rs, err = r._create(tr, cl)
//	g.Expect(rs).To(Equal(resultRequeue))
//	g.Expect(err).NotTo(HaveOccurred())
//	g.Expect(called).To(BeTrue())
//	assertResource(g, r, expectedStatus)
//}
//
//func TestReconcile(t *testing.T) {
//	g := NewGomegaWithT(t)
//
//	tr := testResource()
//
//	r := &Reconciler{
//		Client:                   NewFakeClient(tr),
//		ManagedReferenceResolver: resource.NewAPIManagedReferenceResolver(NewFakeClient(tr)),
//	}
//
//	// test connect error
//	called := false
//	testError := errors.New("test-connect-error")
//	r.connect = func(instance *RDSInstance) (client rds.Client, e error) {
//		called = true
//		return nil, testError
//	}
//
//	expectedStatus := runtimev1alpha1.ConditionedStatus{}
//	expectedStatus.SetConditions(runtimev1alpha1.ReconcileError(testError))
//
//	rs, err := r.Reconcile(request)
//	g.Expect(rs).To(Equal(resultRequeue))
//	g.Expect(err).NotTo(HaveOccurred())
//	g.Expect(called).To(BeTrue())
//	assertResource(g, r, expectedStatus)
//
//	// test delete
//	r.connect = func(instance *RDSInstance) (client rds.Client, e error) {
//		t := metav1.Now()
//		instance.DeletionTimestamp = &t
//		return nil, nil
//	}
//	called = false
//	r.delete = func(instance *RDSInstance, client rds.Client) (i reconcile.Result, e error) {
//		called = true
//		return result, nil
//	}
//	r.Reconcile(request)
//	g.Expect(called).To(BeTrue())
//
//	// test create
//	r.connect = func(instance *RDSInstance) (client rds.Client, e error) {
//		return nil, nil
//	}
//	called = false
//	r.delete = r._delete
//	r.create = func(instance *RDSInstance, client rds.Client) (i reconcile.Result, e error) {
//		called = true
//		return result, nil
//	}
//	r.Reconcile(request)
//	g.Expect(called).To(BeTrue())
//
//	// test sync
//	r.connect = func(instance *RDSInstance) (client rds.Client, e error) {
//		instance.Status.InstanceName = "foo"
//		return nil, nil
//	}
//	called = false
//	r.create = r._create
//	r.sync = func(instance *RDSInstance, client rds.Client) (i reconcile.Result, e error) {
//		called = true
//		return result, nil
//	}
//	r.Reconcile(request)
//	g.Expect(called).To(BeTrue())
//
//}
