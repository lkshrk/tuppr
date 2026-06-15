package versioncompare

import (
	"testing"

	tupprv1alpha1 "github.com/home-operations/tuppr/api/v1alpha1"
)

const (
	target       = "v1.34.0"
	currentBuild = "v1.34.0+deadbee"
)

func TestEquivalent(t *testing.T) {
	tests := []struct {
		name    string
		current string
		target  string
		policy  tupprv1alpha1.VersionComparisonSpec
		want    bool
	}{
		{name: "empty policy exact match", current: target, target: target, want: true},
		{name: "empty policy exact mismatch", current: "v1.34.0-deadbee", target: target, want: false},
		{
			name:    "explicit exact rejects build metadata",
			current: currentBuild,
			target:  target,
			policy:  tupprv1alpha1.VersionComparisonSpec{Mode: tupprv1alpha1.VersionComparisonExact},
			want:    false,
		},
		{
			name:    "ignore build metadata accepts plus suffix",
			current: currentBuild,
			target:  target,
			policy:  tupprv1alpha1.VersionComparisonSpec{Mode: tupprv1alpha1.VersionComparisonIgnoreBuildMetadata},
			want:    true,
		},
		{
			name:    "ignore build metadata strips current and target",
			current: "v1.34.0+nodebuild",
			target:  "v1.34.0+gitopsbuild",
			policy:  tupprv1alpha1.VersionComparisonSpec{Mode: tupprv1alpha1.VersionComparisonIgnoreBuildMetadata},
			want:    true,
		},
		{
			name:    "ignore build metadata rejects prerelease",
			current: "v1.34.0-rc.1",
			target:  target,
			policy:  tupprv1alpha1.VersionComparisonSpec{Mode: tupprv1alpha1.VersionComparisonIgnoreBuildMetadata},
			want:    false,
		},
		{
			name:    "ignore commit suffix accepts lowercase hex",
			current: "v1.34.0-deadbee",
			target:  target,
			policy:  tupprv1alpha1.VersionComparisonSpec{Mode: tupprv1alpha1.VersionComparisonIgnoreCommitSuffix},
			want:    true,
		},
		{
			name:    "ignore commit suffix accepts uppercase hex",
			current: "v1.34.0-DEADBEE",
			target:  target,
			policy:  tupprv1alpha1.VersionComparisonSpec{Mode: tupprv1alpha1.VersionComparisonIgnoreCommitSuffix},
			want:    true,
		},
		{
			name:    "ignore commit suffix accepts forty hex chars",
			current: "v1.34.0-0123456789abcdef0123456789abcdef01234567",
			target:  target,
			policy:  tupprv1alpha1.VersionComparisonSpec{Mode: tupprv1alpha1.VersionComparisonIgnoreCommitSuffix},
			want:    true,
		},
		{
			name:    "ignore commit suffix rejects short hex",
			current: "v1.34.0-deadbe",
			target:  target,
			policy:  tupprv1alpha1.VersionComparisonSpec{Mode: tupprv1alpha1.VersionComparisonIgnoreCommitSuffix},
			want:    false,
		},
		{
			name:    "ignore commit suffix rejects prerelease",
			current: "v1.34.0-rc.1",
			target:  target,
			policy:  tupprv1alpha1.VersionComparisonSpec{Mode: tupprv1alpha1.VersionComparisonIgnoreCommitSuffix},
			want:    false,
		},
		{
			name:    "ignore commit suffix rejects word suffix",
			current: "v1.34.0-talos",
			target:  target,
			policy:  tupprv1alpha1.VersionComparisonSpec{Mode: tupprv1alpha1.VersionComparisonIgnoreCommitSuffix},
			want:    false,
		},
		{
			name:    "ignore matching suffix accepts configured vendor suffix",
			current: "v1.34.0-hcloud.20260614",
			target:  target,
			policy: tupprv1alpha1.VersionComparisonSpec{
				Mode:          tupprv1alpha1.VersionComparisonIgnoreMatchingSuffix,
				SuffixPattern: "-hcloud\\.[0-9]{8}$",
			},
			want: true,
		},
		{
			name:    "ignore matching suffix requires target prefix",
			current: "v1.35.0-hcloud.20260614",
			target:  target,
			policy: tupprv1alpha1.VersionComparisonSpec{
				Mode:          tupprv1alpha1.VersionComparisonIgnoreMatchingSuffix,
				SuffixPattern: "-hcloud\\.[0-9]{8}$",
			},
			want: false,
		},
		{
			name:    "unknown runtime mode falls back to exact",
			current: currentBuild,
			target:  target,
			policy:  tupprv1alpha1.VersionComparisonSpec{Mode: "Loose"},
			want:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Equivalent(tt.current, tt.target, tt.policy)
			if got != tt.want {
				t.Fatalf("Equivalent(%q, %q, %+v) = %v, want %v", tt.current, tt.target, tt.policy, got, tt.want)
			}
		})
	}
}
