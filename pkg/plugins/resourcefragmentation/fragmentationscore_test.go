/*
Copyright 2026 KubeNexus Authors.

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

package resourcefragmentation

import (
	"testing"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

func TestName(t *testing.T) {
	plugin := &ResourceFragmentationScore{}
	if got := plugin.Name(); got != Name {
		t.Errorf("Name() = %v, want %v", got, Name)
	}
}

func TestScoreExtensions(t *testing.T) {
	plugin := &ResourceFragmentationScore{}
	if got := plugin.ScoreExtensions(); got != nil {
		t.Errorf("ScoreExtensions() = %v, want nil", got)
	}
}

func TestConstants(t *testing.T) {
	tests := []struct {
		name string
		got  interface{}
		want interface{}
	}{
		{"Name", Name, "ResourceFragmentationScore"},
		{"LargeIslandThreshold", LargeIslandThreshold, 4},
		{"SmallRequestThreshold", SmallRequestThreshold, 2},
		{"PenaltyFragmentPristineIsland", int64(PenaltyFragmentPristineIsland), int64(0)},
		{"BonusPerfectFit", int64(BonusPerfectFit), int64(90)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.want {
				t.Errorf("%s = %v, want %v", tt.name, tt.got, tt.want)
			}
		})
	}
}

func TestGetGPURequest(t *testing.T) {
	tests := []struct {
		name        string
		pod         *v1.Pod
		expectedGPU int
	}{
		{
			name: "single container with GPU",
			pod: &v1.Pod{
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Resources: v1.ResourceRequirements{
								Requests: v1.ResourceList{
									ResourceGPU: resource.MustParse("4"),
								},
							},
						},
					},
				},
			},
			expectedGPU: 4,
		},
		{
			name: "no GPU request",
			pod: &v1.Pod{
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Resources: v1.ResourceRequirements{
								Requests: v1.ResourceList{
									v1.ResourceCPU: resource.MustParse("2"),
								},
							},
						},
					},
				},
			},
			expectedGPU: 0,
		},
		{
			name: "multiple containers with GPUs",
			pod: &v1.Pod{
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Resources: v1.ResourceRequirements{
								Requests: v1.ResourceList{
									ResourceGPU: resource.MustParse("2"),
								},
							},
						},
						{
							Resources: v1.ResourceRequirements{
								Requests: v1.ResourceList{
									ResourceGPU: resource.MustParse("2"),
								},
							},
						},
					},
				},
			},
			expectedGPU: 4,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gpuCount := getGPURequest(tt.pod)
			if gpuCount != tt.expectedGPU {
				t.Errorf("Expected %d GPUs, got %d", tt.expectedGPU, gpuCount)
			}
		})
	}
}
