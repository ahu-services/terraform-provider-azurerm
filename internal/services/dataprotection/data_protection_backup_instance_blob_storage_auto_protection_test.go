// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package dataprotection

import (
	"reflect"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
)

func TestExpandBlobBackupAutoProtection(t *testing.T) {
	withoutRules := expandBlobBackupAutoProtection([]interface{}{})
	settings, ok := withoutRules.AutoProtectionSettings.(BlobBackupRuleBasedAutoProtectionSettings)
	if !ok {
		t.Fatalf("expected BlobBackupRuleBasedAutoProtectionSettings, got %T", withoutRules.AutoProtectionSettings)
	}
	if !settings.Enabled {
		t.Fatalf("expected auto protection to be enabled")
	}
	if settings.Rules != nil {
		t.Fatalf("expected no rules when no prefixes are excluded, got %+v", *settings.Rules)
	}

	withRules := expandBlobBackupAutoProtection([]interface{}{"temp-", "test-"})
	settings = withRules.AutoProtectionSettings.(BlobBackupRuleBasedAutoProtectionSettings)
	expected := []BlobBackupAutoProtectionRule{
		{ObjectType: "BlobBackupAutoProtectionRule", Mode: BlobBackupRuleModeExclude, Type: BlobBackupPatternTypePrefix, Pattern: "temp-"},
		{ObjectType: "BlobBackupAutoProtectionRule", Mode: BlobBackupRuleModeExclude, Type: BlobBackupPatternTypePrefix, Pattern: "test-"},
	}
	if !reflect.DeepEqual(pointer.From(settings.Rules), expected) {
		t.Fatalf("unexpected rules:\n want: %+v\n got:  %+v", expected, pointer.From(settings.Rules))
	}
}

func TestFlattenBlobBackupAutoProtection(t *testing.T) {
	testCases := []struct {
		name             string
		input            BlobBackupDatasourceParametersForAutoProtection
		expectedEnabled  bool
		expectedPrefixes []string
	}{
		{
			name:             "no settings",
			input:            BlobBackupDatasourceParametersForAutoProtection{},
			expectedEnabled:  false,
			expectedPrefixes: []string{},
		},
		{
			name: "enabled without rules",
			input: BlobBackupDatasourceParametersForAutoProtection{
				AutoProtectionSettings: BlobBackupRuleBasedAutoProtectionSettings{Enabled: true},
			},
			expectedEnabled:  true,
			expectedPrefixes: []string{},
		},
		{
			name: "enabled with rules keeps order and ignores unknown rule kinds",
			input: BlobBackupDatasourceParametersForAutoProtection{
				AutoProtectionSettings: BlobBackupRuleBasedAutoProtectionSettings{
					Enabled: true,
					Rules: &[]BlobBackupAutoProtectionRule{
						{Mode: BlobBackupRuleModeExclude, Type: BlobBackupPatternTypePrefix, Pattern: "temp-"},
						{Mode: "Include", Type: BlobBackupPatternTypePrefix, Pattern: "keep-"},
						{Mode: BlobBackupRuleModeExclude, Type: BlobBackupPatternTypePrefix, Pattern: "test-"},
					},
				},
			},
			expectedEnabled:  true,
			expectedPrefixes: []string{"temp-", "test-"},
		},
		{
			name: "unknown settings implementation only exposes enabled",
			input: BlobBackupDatasourceParametersForAutoProtection{
				AutoProtectionSettings: RawBlobBackupAutoProtectionSettingsImpl{
					blobBackupAutoProtectionSettings: BaseBlobBackupAutoProtectionSettingsImpl{Enabled: true},
				},
			},
			expectedEnabled:  true,
			expectedPrefixes: []string{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			enabled, prefixes := flattenBlobBackupAutoProtection(tc.input)
			if enabled != tc.expectedEnabled {
				t.Fatalf("expected enabled = %t, got %t", tc.expectedEnabled, enabled)
			}
			if !reflect.DeepEqual(prefixes, tc.expectedPrefixes) {
				t.Fatalf("expected prefixes %v, got %v", tc.expectedPrefixes, prefixes)
			}
		})
	}
}
