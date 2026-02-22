package scoring

import (
	"testing"
	"context"
	"k8s.io/api/core/v1"
	framework "k8s.io/kube-scheduler/framework"
)

func TestTopologySpreadScorePlugin_Name(t *testing.T) {
	plugin := &TopologySpreadScorePlugin{}
	if plugin.Name() != TopologyScoringName {
		t.Errorf("expected plugin name %s, got %s", TopologyScoringName, plugin.Name())
	}
}

func TestTopologySpreadScorePlugin_Score(t *testing.T) {
	plugin := &TopologySpreadScorePlugin{}
	pod := &v1.Pod{}
	nodeInfo := framework.NewNodeInfo()
	ctx := context.Background()
	_, status := plugin.Score(ctx, nil, pod, nodeInfo)
	if status == nil {
		t.Errorf("expected non-nil status")
	}
}
