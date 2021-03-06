// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2020 Intel Corporation

package assets

import (
	"context"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/klog/v2/klogr"
)

//  runtime.Object implementation
type InvalidRuntimeType struct {
}

func (i *InvalidRuntimeType) GetObjectKind() schema.ObjectKind {
	return schema.EmptyObjectKind
}
func (i *InvalidRuntimeType) DeepCopyObject() runtime.Object {
	return i
}

var _ = Describe("Asset Tests", func() {

	log := klogr.New()

	var _ = Describe("Manager", func() {
		var _ = It("Run Manager with no assets (setKernel false)", func() {
			var err error
			log = klogr.New().WithName("N3000Assets-Test")

			manager := Manager{Client: k8sClient, Log: log}

			err = manager.LoadAndDeploy(context.TODO(), false)
			Expect(err).ToNot(HaveOccurred())
		})
		var _ = It("Run Manager with no assets (setKernel true)", func() {
			var err error
			log = klogr.New().WithName("N3000Assets-Test")

			manager := Manager{Client: k8sClient, Log: log}

			err = manager.LoadAndDeploy(context.TODO(), true)
			Expect(err).To(HaveOccurred())
		})
		var _ = It("Run Manager (setKernel true)", func() {
			var err error
			log = klogr.New().WithName("N3000Assets-Test")

			assets := []Asset{
				{
					log:  log,
					Path: "/tmp/dummy.bin",
				},
			}

			manager := Manager{Client: k8sClient,
				Log:    log,
				Assets: assets}

			err = manager.LoadAndDeploy(context.TODO(), true)
			Expect(err).To(HaveOccurred())
		})
		var _ = It("Run Manager loadFromDir (setKernel false)", func() {
			var err error
			log = klogr.New().WithName("N3000Assets-Test")

			assets := []Asset{
				{
					log:  log,
					Path: "/tmp/",
				},
			}

			manager := Manager{Client: k8sClient,
				Log:    log,
				Assets: assets}

			err = manager.LoadAndDeploy(context.TODO(), false)
			Expect(err).To(HaveOccurred())
		})
		var _ = It("Run Manager loadFromFile (setKernel false, no retries)", func() {
			var err error
			log = klogr.New().WithName("N3000Assets-Test")

			assets := []Asset{
				{
					log:           log,
					Path:          fakeAssetFile,
					substitutions: map[string]string{"one": "two"},
					BlockingReadiness: ReadinessPollConfig{
						Retries: 1,
					},
				},
			}

			manager := Manager{Client: k8sClient,
				Log:    log,
				Assets: assets,
				Owner:  fakeOwner,
				Scheme: scheme.Scheme}

			err = manager.LoadAndDeploy(context.TODO(), false)
			Expect(err).ToNot(HaveOccurred())
		})
		var _ = It("Run Manager loadFromFile (setKernel false)", func() {
			var err error
			log = klogr.New().WithName("N3000Assets-Test")

			assets := []Asset{
				{
					log:           log,
					Path:          fakeAssetFile,
					substitutions: map[string]string{"one": "two"},
				},
			}

			manager := Manager{Client: k8sClient,
				Log:    log,
				Assets: assets,
				Owner:  fakeOwner,
				Scheme: scheme.Scheme}

			err = manager.LoadAndDeploy(context.TODO(), false)
			Expect(err).ToNot(HaveOccurred())
		})
		var _ = It("Run LoadAndDeploy (fail setting Owner)", func() {
			var err error
			log = klogr.New().WithName("N3000Assets-Test")

			var invalidObject InvalidRuntimeType

			assets := []Asset{
				{
					log:           log,
					Path:          fakeAssetFile,
					substitutions: map[string]string{"one": "two"},
					objects: []runtime.Object{
						&invalidObject},
				},
			}

			manager := Manager{Client: k8sClient,
				Log:    log,
				Assets: assets,
				Owner:  fakeOwner,
				Scheme: scheme.Scheme}

			Expect(manager).ToNot(Equal(nil))

			// Create a Node
			node := &corev1.Node{
				ObjectMeta: v1.ObjectMeta{
					Name: "dummy",
					Labels: map[string]string{
						"fpga.intel.com/intel-accelerator-present": "",
					},
				},
			}

			err = k8sClient.Create(context.Background(), node)
			Expect(err).ToNot(HaveOccurred())

			err = manager.LoadAndDeploy(context.TODO(), true)
			Expect(err).To(HaveOccurred())

			// Cleanup
			err = k8sClient.Delete(context.TODO(), node)
			Expect(err).ToNot(HaveOccurred())
		})
		var _ = It("Run Manager loadFromFile (bad file)", func() {
			var err error
			log = klogr.New().WithName("N3000Assets-Test")

			assets := []Asset{
				{
					log:           log,
					Path:          "/dev/null",
					substitutions: map[string]string{"one": "two"},
					BlockingReadiness: ReadinessPollConfig{
						Retries: 1,
					},
				},
			}

			manager := Manager{Client: k8sClient,
				Log:    log,
				Assets: assets,
				Owner:  fakeOwner,
				Scheme: scheme.Scheme}

			err = manager.LoadAndDeploy(context.TODO(), false)
			Expect(err).To(HaveOccurred())
		})
		var _ = It("Run Manager loadFromFile (missing file)", func() {
			var err error
			log = klogr.New().WithName("N3000Assets-Test")

			assets := []Asset{
				{
					log:           log,
					Path:          "/dev/null_fake",
					substitutions: map[string]string{"one": "two"},
					BlockingReadiness: ReadinessPollConfig{
						Retries: 1,
					},
				},
			}

			manager := Manager{Client: k8sClient,
				Log:    log,
				Assets: assets,
				Owner:  fakeOwner,
				Scheme: scheme.Scheme}

			err = manager.LoadAndDeploy(context.TODO(), false)
			Expect(err).To(HaveOccurred())
		})
		var _ = It("Run Manager loadFromFile (invalid retries count)", func() {
			var err error
			log = klogr.New().WithName("N3000Assets-Test")

			assets := []Asset{
				{
					log:           log,
					Path:          fakeAssetFile,
					substitutions: map[string]string{"one": "two"},
					BlockingReadiness: ReadinessPollConfig{
						Retries: -1,
					},
				},
			}

			manager := Manager{Client: k8sClient,
				Log:    log,
				Assets: assets,
				Owner:  fakeOwner,
				Scheme: scheme.Scheme}

			err = manager.LoadAndDeploy(context.TODO(), false)
			Expect(err).To(HaveOccurred())
		})
	})
})
