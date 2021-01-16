/*
   Copyright (C) 2021 Dennis Jung

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
	"os"

	"github.com/dennismaxjung/expanded-hostpath-provisioner/internal/config"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	klog "k8s.io/klog/v2"
	"sigs.k8s.io/sig-storage-lib-external-provisioner/v6/controller"
)

type hostPathProvisioner config.Provisioner

var (
	nodeName string = getNodeName()
)

func getNodeName() string {
	name := os.Getenv("NODE_NAME")
	if name == "" {
		//klog.Fatal()
	}
	return name
}

func getClientset() *kubernetes.Clientset {
	kubeConfig, err := rest.InClusterConfig()
	if err != nil {
		klog.Fatalf("Failed to get kubeConfig: %v", err)
	}

	clientset, err := kubernetes.NewForConfig(kubeConfig)
	if err != nil {
		klog.Fatalf("Failed to get clientset; %v", err)
	}

	return clientset
}

func NewProvisionerController(provisioner config.Provisioner) *controller.ProvisionController {
	p := NewHostPathProvisioner(provisioner)
	clientset := getClientset()

	serverVersion, err := clientset.Discovery().ServerVersion()
	if err != nil {
		klog.Fatalf("Failed to get serverVersion: %v", err)
	}

	return controller.NewProvisionController(
		clientset,
		provisioner.Name,
		p,
		serverVersion.GitVersion,
		controller.ExponentialBackOffOnError(exponentialBackOffOnError),
		controller.ResyncPeriod(resyncPeriod),
		controller.FailedProvisionThreshold(failedRetryThreshold),
		controller.FailedDeleteThreshold(failedRetryThreshold),
	)
}
