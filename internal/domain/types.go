package domain

import "time"

type Budget struct {
	CPUMilli  int `json:"cpu_milli"`
	MemoryMB  int `json:"memory_mb"`
	HookUnits int `json:"hook_units"`
}

type Condition struct {
	SecureBoot   *bool    `json:"secure_boot,omitempty"`
	Zones        []string `json:"zones,omitempty"`
	EgressModes  []string `json:"egress_modes,omitempty"`
	Incident     []string `json:"incident_states,omitempty"`
	Runtimes     []string `json:"runtimes,omitempty"`
	FIPSModes    []string `json:"fips_modes,omitempty"`
	KernelLT     string   `json:"kernel_lt,omitempty"`
	KernelGTE    string   `json:"kernel_gte,omitempty"`
	RequireRaw   *bool    `json:"require_raw_payload,omitempty"`
}

type Requirement struct {
	ID     string    `json:"id"`
	Reason string    `json:"reason"`
	When   Condition `json:"when"`
}

type AlternativeRequirement struct {
	IDs    []string  `json:"ids"`
	Reason string    `json:"reason"`
	When   Condition `json:"when"`
}

type Conflict struct {
	ID       string    `json:"id"`
	Code     string    `json:"code"`
	Message  string    `json:"message"`
	Waivable bool      `json:"waivable"`
	When     Condition `json:"when"`
}

type VersionRange struct {
	Min          string `json:"min"`
	MaxExclusive string `json:"max_exclusive"`
}

type Component struct {
	ID               string                   `json:"id"`
	Class            string                   `json:"class"`
	Roles            []string                 `json:"roles"`
	RawPayload       bool                     `json:"raw_payload"`
	SupportedOS      []string                 `json:"supported_os"`
	SupportedRuntime []string                 `json:"supported_runtimes"`
	KernelRanges     []VersionRange           `json:"kernel_ranges"`
	Overhead         Budget                   `json:"overhead"`
	Requires         []Requirement            `json:"requires"`
	RequiresAny      []AlternativeRequirement `json:"requires_any"`
	Conflicts        []Conflict               `json:"conflicts"`
}

type ComponentsCatalog struct {
	Components []Component `json:"components"`
}

type Profile struct {
	ID            string    `json:"id"`
	OS            string    `json:"os"`
	KernelVersion string    `json:"kernel_version"`
	Runtime       string    `json:"runtime"`
	Zone          string    `json:"zone"`
	EgressMode    string    `json:"egress_mode"`
	IncidentState string    `json:"incident_state"`
	FIPSMode      string    `json:"fips_mode"`
	SecureBoot    bool      `json:"secure_boot"`
	Budgets       Budget    `json:"budgets"`
	AsOf          time.Time `json:"as_of"`
}

type ProfilesCatalog struct {
	Profiles []Profile `json:"profiles"`
}

// FlowPolicy defines the global data-path and budget rules loaded from
// compatibility_rules.json.
type FlowPolicy struct {
	CollectorID        string              `json:"collector_id"`
	RelayID            string              `json:"relay_id"`
	SanitizersByState  map[string][]string `json:"sanitizers_by_state"`
	ProtectedZones     []string            `json:"protected_zones"`
	WarningUtilization float64             `json:"warning_utilization"`
	// Note: auditor_id and auditor_required_fips_modes are defined in the
	// policy contract but are not yet wired into this struct.
}

type GlobalRules struct {
	Flow FlowPolicy `json:"flow"`
}

// Approval represents a waiver ticket loaded from approval_register.json.
type Approval struct {
	TicketID       string    `json:"ticket_id"`
	RuleCode       string    `json:"rule_code"`
	Components     []string  `json:"components"`
	ProfileIDs     []string  `json:"profile_ids"`
	Zones          []string  `json:"zones"`
	IncidentStates []string  `json:"incident_states"`
	// fips_modes and egress_modes scope fields are present in the JSON but not
	// mapped here yet.
	ExpiresAt time.Time `json:"expires_at"`
}

type ApprovalsCatalog struct {
	Approvals []Approval `json:"approvals"`
}

type AssessmentRequest struct {
	ProfileID          string   `json:"profile_id"`
	BundleName         string   `json:"bundle_name"`
	SelectedComponents []string `json:"selected_components"`
	DataPath           []string `json:"data_path"`
	Approvals          []string `json:"approvals"`
}

type Blocker struct {
	Code       string   `json:"code"`
	Axis       string   `json:"axis"`
	Message    string   `json:"message"`
	Components []string `json:"components"`
}

type Trace struct {
	EvaluatedAxes []string `json:"evaluated_axes"`
	Totals        Budget   `json:"totals"`
}

type AssessmentResponse struct {
	ProfileID       string    `json:"profile_id"`
	BundleName      string    `json:"bundle_name"`
	Decision        string    `json:"decision"`
	Compatible      bool      `json:"compatible"`
	Score           float64   `json:"score"`
	Blockers        []Blocker `json:"blockers"`
	RequiredActions []string  `json:"required_actions"`
	Trace           Trace     `json:"trace"`
}
