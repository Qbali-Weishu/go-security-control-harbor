package policy

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"compatgate/internal/domain"
)

// Catalog 包含所有策略数据
type Catalog struct {
	Components map[string]domain.Component
	Profiles   map[string]domain.Profile
	Rules      domain.GlobalRules
	Approvals  []domain.Approval
}

// Load 从指定目录加载策略目录
func Load(root string) (Catalog, error) {
	components, err := loadComponents(filepath.Join(root, "component_catalog.json"))
	if err != nil {
		return Catalog{}, err
	}
	profiles, err := loadProfiles(filepath.Join(root, "profile_catalog.json"))
	if err != nil {
		return Catalog{}, err
	}
	rules, err := loadRules(filepath.Join(root, "compatibility_rules.json"))
	if err != nil {
		return Catalog{}, err
	}
	approvals, err := loadApprovals(filepath.Join(root, "approval_register.json"))
	if err != nil {
		return Catalog{}, err
	}
	return Catalog{Components: components, Profiles: profiles, Rules: rules, Approvals: approvals}, nil
}

// loadComponents 加载组件目录
func loadComponents(path string) (map[string]domain.Component, error) {
	var catalog domain.ComponentsCatalog
	if err := decode(path, &catalog); err != nil {
		return nil, err
	}
	items := make(map[string]domain.Component, len(catalog.Components))
	for _, component := range catalog.Components {
		items[component.ID] = component
	}
	return items, nil
}

// loadProfiles 加载配置文件目录
func loadProfiles(path string) (map[string]domain.Profile, error) {
	var catalog domain.ProfilesCatalog
	if err := decode(path, &catalog); err != nil {
		return nil, err
	}
	items := make(map[string]domain.Profile, len(catalog.Profiles))
	for _, profile := range catalog.Profiles {
		items[profile.ID] = profile
	}
	return items, nil
}

// loadRules 加载全局规则
func loadRules(path string) (domain.GlobalRules, error) {
	var rules domain.GlobalRules
	if err := decode(path, &rules); err != nil {
		return domain.GlobalRules{}, err
	}
	return rules, nil
}

// loadApprovals 加载批准记录
func loadApprovals(path string) ([]domain.Approval, error) {
	var catalog domain.ApprovalsCatalog
	if err := decode(path, &catalog); err != nil {
		return nil, err
	}
	return catalog.Approvals, nil
}

// decode 从文件读取并解析 JSON
func decode(path string, target any) error {
	payload, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("读取 %s: %w", path, err)
	}
	if err := json.Unmarshal(payload, target); err != nil {
		return fmt.Errorf("解码 %s: %w", path, err)
	}
	return nil
}
