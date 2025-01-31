/*
Copyright 2022 The Kubernetes Authors All rights reserved.

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
	"os/exec"
	"time"

	"golang.org/x/mod/semver"
	"k8s.io/klog/v2"

	"k8s.io/minikube/hack/update"
)

const (
	// default context timeout
	cxTimeout = 5 * time.Minute
)

func main() {
	// set a context with defined timeout
	ctx, cancel := context.WithTimeout(context.Background(), cxTimeout)
	defer cancel()

	// get Docsy stable version
	stable, err := docsyVersion(ctx, "google", "docsy")
	if err != nil {
		klog.Fatalf("Unable to get Doscy stable version: %v", err)
	}
	klog.Infof("Doscy stable version: %s", stable)

	if err := exec.CommandContext(ctx, "./update_docsy_version.sh", stable).Run(); err != nil {
		klog.Fatalf("failed to update docsy commit: %v", err)
	}
}

// docsyVersion returns stable version in semver format.
func docsyVersion(ctx context.Context, owner, repo string) (stable string, err error) {
	// get Docsy version from GitHub Releases
	stable, _, _, err = update.GHReleases(ctx, owner, repo)
	if err != nil || !semver.IsValid(stable) {
		return "", err
	}
	return stable, nil
}
