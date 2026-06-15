package versioncompare

import (
	"regexp"
	"strings"

	tupprv1alpha1 "github.com/home-operations/tuppr/api/v1alpha1"
)

var commitSuffixPattern = regexp.MustCompile(`-[0-9a-fA-F]{7,40}$`)

// Equivalent reports whether a current reported version satisfies the target
// version under the configured comparison policy. Invalid runtime policies fall
// back to exact comparison; webhooks reject invalid user-facing configs.
func Equivalent(current, target string, policy tupprv1alpha1.VersionComparisonSpec) bool {
	mode := policy.Mode
	if mode == "" {
		mode = tupprv1alpha1.VersionComparisonExact
	}

	switch mode {
	case tupprv1alpha1.VersionComparisonExact:
		return current == target
	case tupprv1alpha1.VersionComparisonIgnoreBuildMetadata:
		return stripBuildMetadata(current) == stripBuildMetadata(target)
	case tupprv1alpha1.VersionComparisonIgnoreCommitSuffix:
		return equivalentIgnoringSuffix(current, target, commitSuffixPattern)
	case tupprv1alpha1.VersionComparisonIgnoreMatchingSuffix:
		if policy.SuffixPattern == "" {
			return current == target
		}
		pattern, err := regexp.Compile(policy.SuffixPattern)
		if err != nil {
			return current == target
		}
		return equivalentIgnoringSuffix(current, target, pattern)
	default:
		return current == target
	}
}

func stripBuildMetadata(version string) string {
	base, _, _ := strings.Cut(version, "+")
	return base
}

func equivalentIgnoringSuffix(current, target string, pattern *regexp.Regexp) bool {
	if current == target {
		return true
	}
	if !strings.HasPrefix(current, target) {
		return false
	}
	suffix := strings.TrimPrefix(current, target)
	return suffix != "" && pattern.MatchString(suffix)
}
