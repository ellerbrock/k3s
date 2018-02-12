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

package vsphere

import (
	"k8s.io/kubernetes/test/e2e/framework"
	"sync"
)

var once sync.Once
var waiting = make(chan bool)
var f *framework.Framework

// Bootstrap takes care of initializing necessary test context for vSphere tests
func Bootstrap(fw *framework.Framework) {
	done := make(chan bool)
	f = fw
	go func() {
		once.Do(bootstrapOnce)
		<-waiting
		done <- true
	}()
	<-done
}

func bootstrapOnce() {
	// 1. Read vSphere conf and get VSphere instances
	vsphereInstances, err := GetVSphereInstances()
	if err != nil {
		framework.Failf("Failed to bootstrap vSphere with error: %v", err)
	}
	// 2. Get all ready nodes
	nodeList := framework.GetReadySchedulableNodesOrDie(f.ClientSet)
	TestContext = VSphereContext{NodeMapper: &NodeMapper{}, VSphereInstances: vsphereInstances}

	// 3. Get Node to VSphere mapping
	err = TestContext.NodeMapper.GenerateNodeMap(vsphereInstances, *nodeList)
	if err != nil {
		framework.Failf("Failed to bootstrap vSphere with error: %v", err)
	}
	close(waiting)
}