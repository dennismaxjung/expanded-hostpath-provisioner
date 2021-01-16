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

package main

import (
	"context"
	"flag"
	"os"
	"sync"
	"syscall"

	conf "github.com/dennismaxjung/expanded-hostpath-provisioner/internal/config"
	"github.com/dennismaxjung/expanded-hostpath-provisioner/internal/hostpathProvisioner"

	klog "k8s.io/klog/v2"
)

var (
	config conf.Config
)

func getDefaultProvisioner() conf.Provisioner {
	// The provisioner name is "microk8s.io/hostpath" for compatibility reasonmust be the one used in the storage class manifest
	provisionerName := "microk8s.io/hostpath"

	pvdir := os.Getenv("PV_DIR")
	if pvdir == "" {
		klog.Fatal("The env variable PV_DIR has to be set.")
	}
	pvreclaim := os.Getenv("PV_RECLAIM_POLICY")
	if pvreclaim == "" {
		klog.Info("The env variable PV_RECLAIM_POLICY is not set. The ReclaimPolicy of the StorageClass will be used")
	}

	return conf.Provisioner{Name: provisionerName, Directory: pvdir, ReclaimPolicy: pvreclaim}
}

func main() {
	syscall.Umask(0)

	if !flag.Parsed() {
		flag.Parse()
	}
	flag.Set("logtostderr", "true")

	config, err := conf.LoadConfiguration("config.json")

	if err != nil {
		klog.Fatalf("Failed to get 'config.json' : %v", err)
	}

	// Check if the default config should be used.
	if config.DefaultProvisioner == true {
		klog.Info("DefaultProvisioner is activated.")
		defaultProvisioner := getDefaultProvisioner()
		config.Provisioners = append(config.Provisioners, defaultProvisioner)
	}

	var wg sync.WaitGroup

	for _, provisioner := range config.Provisioners {
		// Start the provision controller which will dynamically provision hostPath PVs
		provisionController := hostpathProvisioner.NewProvisionerController(provisioner)
		wg.Add(1)
		go func() {
			defer wg.Done()
			provisionController.Run(context.Background())
		}()
	}

	wg.Wait()
}
