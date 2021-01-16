/*
Copyright 2018 The Kubernetes Authors.

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

package hostpathProvisioner

import (
	"context"
	"errors"
	"os"
	"path"
	"time"

	"github.com/dennismaxjung/expanded-hostpath-provisioner/internal/config"
	"sigs.k8s.io/sig-storage-lib-external-provisioner/v6/controller"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	klog "k8s.io/klog/v2"
)

const (
	resyncPeriod              = 15 * time.Second
	exponentialBackOffOnError = false
	failedRetryThreshold      = 5
)

// NewHostPathProvisioner creates a new hostpath provisioner
func NewHostPathProvisioner(prov config.Provisioner) controller.Provisioner {
	var provisioner = hostPathProvisioner(prov)
	return &provisioner
}

// Provision creates a storage asset and returns a PV object representing it.
func (p *hostPathProvisioner) Provision(ctx context.Context, options controller.ProvisionOptions) (*v1.PersistentVolume, controller.ProvisioningState, error) {
	path := path.Join(p.Directory, options.PVC.Namespace+"-"+options.PVC.Name+"-"+options.PVName)
	klog.Infof("creating backing directory: %v", path)

	if err := os.MkdirAll(path, 0777); err != nil {
		return nil, controller.ProvisioningFinished, err
	}

	reclaimPolicy := *options.StorageClass.ReclaimPolicy
	if p.ReclaimPolicy != "" {
		reclaimPolicy = v1.PersistentVolumeReclaimPolicy(p.ReclaimPolicy)
	}

	pv := &v1.PersistentVolume{
		ObjectMeta: metav1.ObjectMeta{
			Name: options.PVName,
			Annotations: map[string]string{
				"hostPathProvisionerIdentity": nodeName,
			},
		},
		Spec: v1.PersistentVolumeSpec{
			PersistentVolumeReclaimPolicy: reclaimPolicy,
			AccessModes:                   options.PVC.Spec.AccessModes,
			Capacity: v1.ResourceList{
				v1.ResourceName(v1.ResourceStorage): options.PVC.Spec.Resources.Requests[v1.ResourceName(v1.ResourceStorage)],
			},
			PersistentVolumeSource: v1.PersistentVolumeSource{
				HostPath: &v1.HostPathVolumeSource{
					Path: path,
				},
			},
		},
	}

	return pv, controller.ProvisioningFinished, nil
}

// Delete removes the storage asset that was created by Provision represented
// by the given PV.
func (p *hostPathProvisioner) Delete(ctx context.Context, volume *v1.PersistentVolume) error {
	ann, ok := volume.Annotations["hostPathProvisionerIdentity"]
	if !ok {
		return errors.New("identity annotation not found on PV")
	}
	if ann != nodeName {
		return &controller.IgnoredError{Reason: "identity annotation on PV does not match ours"}
	}

	path := volume.Spec.PersistentVolumeSource.HostPath.Path
	klog.Info("removing backing directory: %v", path)
	if err := os.RemoveAll(path); err != nil {
		return err
	}

	return nil
}
