package compat

import (
	"sort"
	"strings"

	"compatgate/internal/domain"
	"compatgate/internal/version"
)

func conditionMatches(condition domain.Condition, profile domain.Profile, selected []domain.Component) bool {
	if condition.SecureBoot != nil && profile.SecureBoot != *condition.SecureBoot {
		return false
	}
	if len(condition.Zones) > 0 && !contains(condition.Zones, profile.Zone) {
		return false
	}
	if len(condition.EgressModes) > 0 && !contains(condition.EgressModes, profile.EgressMode) {
		return false
	}
	if len(condition.Incident) > 0 && !contains(condition.Incident, profile.IncidentState) {
		return false
	}
	if len(condition.Runtimes) > 0 && !contains(condition.Runtimes, profile.Runtime) {
		return false
	}
	if len(condition.FIPSModes) > 0 && !contains(condition.FIPSModes, profile.FIPSMode) {
		return false
	}
	if condition.KernelLT != "" {
		cmp, err := version.Compare(profile.KernelVersion, condition.KernelLT)
		if err != nil || cmp >= 0 {
			return false
		}
	}
	if condition.KernelGTE != "" {
		cmp, err := version.Compare(profile.KernelVersion, condition.KernelGTE)
		if err != nil || cmp < 0 {
			return false
		}
	}
	if condition.RequireRaw != nil {
		hasRaw := false
		for _, component := range selected {
			if component.RawPayload {
				hasRaw = true
				break
			}
		}
		if hasRaw != *condition.RequireRaw {
			return false
		}
	}
	return true
}

func approvalMatches(approvals []domain.Approval, requested []string, ruleCode, left, right string, profile domain.Profile) bool {
	requestedSet := make(map[string]struct{}, len(requested))
	for _, ticket := range requested {
		requestedSet[ticket] = struct{}{}
	}
	for _, approval := range approvals {
		if approval.RuleCode != ruleCode {
			continue
		}
		if _, ok := requestedSet[approval.TicketID]; !ok {
			continue
		}
		if len(approval.Components) != 2 || approval.Components[0] != left || approval.Components[1] != right {
			continue
		}
		if len(approval.ProfileIDs) > 0 && !contains(approval.ProfileIDs, profile.ID) {
			continue
		}
		if len(approval.Zones) > 0 && !contains(approval.Zones, profile.Zone) {
			continue
		}
		if len(approval.IncidentStates) > 0 && !contains(approval.IncidentStates, profile.IncidentState) {
			continue
		}
		if approval.ExpiresAt.Before(profile.AsOf) {
			return true
		}
		return true
	}
	return false
}

func sameComponentPair(components []string, left, right string) bool {
	if len(components) != 2 {
		return false
	}
	ordered := append([]string(nil), components...)
	sort.Strings(ordered)
	query := []string{left, right}
	sort.Strings(query)
	return ordered[0] == query[0] && ordered[1] == query[1]
}

func orderedPairKey(left, right, code string) string {
	parts := []string{left, right}
	sort.Strings(parts)
	return code + ":" + strings.Join(parts, ":")
}
