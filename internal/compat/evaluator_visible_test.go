package compat

import (
	"path/filepath"
	"testing"

	"compatgate/internal/domain"
	"compatgate/internal/policy"
)

func TestVisibleSanityValidRegulatedBundle(t *testing.T) {
	catalog, err := policy.Load(filepath.Join("..", "..", "testdata", "policies"))
	if err != nil {
		t.Fatalf("load catalog: %v", err)
	}
	evaluator := NewEvaluator(catalog)
	resp, err := evaluator.Assess(domain.AssessmentRequest{
		ProfileID:          "regulated-container-host-a",
		BundleName:         "visible-regulated",
		SelectedComponents: []string{"kernel-trace-sensor", "attestation-seal", "content-sanitizer", "central-collector", "telemetry-relay", "egress-auditor"},
		DataPath:           []string{"kernel-trace-sensor", "content-sanitizer", "central-collector", "telemetry-relay", "egress-auditor"},
	})
	if err != nil {
		t.Fatalf("assess: %v", err)
	}
	if !resp.Compatible {
		t.Fatalf("expected compatible bundle, got blockers: %#v", resp.Blockers)
	}
}

func TestVisibleSanityBalancedRestrictedBundleDoesNotNeedAuditor(t *testing.T) {
	catalog, err := policy.Load(filepath.Join("..", "..", "testdata", "policies"))
	if err != nil {
		t.Fatalf("load catalog: %v", err)
	}
	evaluator := NewEvaluator(catalog)
	resp, err := evaluator.Assess(domain.AssessmentRequest{
		ProfileID:          "regulated-balanced-container-host-a",
		BundleName:         "visible-balanced",
		SelectedComponents: []string{"kernel-trace-sensor", "attestation-seal", "content-sanitizer", "central-collector", "telemetry-relay"},
		DataPath:           []string{"kernel-trace-sensor", "content-sanitizer", "central-collector", "telemetry-relay"},
	})
	if err != nil {
		t.Fatalf("assess: %v", err)
	}
	if !resp.Compatible {
		t.Fatalf("expected compatible balanced bundle, got blockers: %#v", resp.Blockers)
	}
}
