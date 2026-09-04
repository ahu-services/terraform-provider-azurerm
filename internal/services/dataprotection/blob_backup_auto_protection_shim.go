// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package dataprotection

// TEMPORARY: this file mirrors the models that `hashicorp/go-azure-sdk` will generate for
// `resource-manager/dataprotection/2026-03-01/backupinstanceresources` once that API version
// is published (the Pandora config landed in hashicorp/pandora#5490, the SDK regeneration
// is still pending). The type and constant names intentionally match the generator output so
// that, once the SDK is bumped, this file can be deleted and the usages switched to the
// `backupinstanceresources` package with no further changes.
//
// Schema reference: Microsoft.DataProtection/backupVaults/backupInstances @ 2026-03-01
//   BlobBackupDatasourceParametersForAutoProtection
//   BlobBackupRuleBasedAutoProtectionSettings
//   BlobBackupAutoProtectionRule

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2025-07-01/backupinstanceresources"
)

type BlobBackupRuleMode string

const (
	BlobBackupRuleModeExclude BlobBackupRuleMode = "Exclude"
)

func PossibleValuesForBlobBackupRuleMode() []string {
	return []string{
		string(BlobBackupRuleModeExclude),
	}
}

type BlobBackupPatternType string

const (
	BlobBackupPatternTypePrefix BlobBackupPatternType = "Prefix"
)

func PossibleValuesForBlobBackupPatternType() []string {
	return []string{
		string(BlobBackupPatternTypePrefix),
	}
}

// BlobBackupAutoProtectionRule - `objectType` is a required plain string (not a discriminator)
// and must be set to "BlobBackupAutoProtectionRule" by the caller.
type BlobBackupAutoProtectionRule struct {
	Mode       BlobBackupRuleMode    `json:"mode"`
	ObjectType string                `json:"objectType"`
	Pattern    string                `json:"pattern"`
	Type       BlobBackupPatternType `json:"type"`
}

type BlobBackupAutoProtectionSettings interface {
	BlobBackupAutoProtectionSettings() BaseBlobBackupAutoProtectionSettingsImpl
}

var _ BlobBackupAutoProtectionSettings = BaseBlobBackupAutoProtectionSettingsImpl{}

type BaseBlobBackupAutoProtectionSettingsImpl struct {
	Enabled    bool   `json:"enabled"`
	ObjectType string `json:"objectType"`
}

func (s BaseBlobBackupAutoProtectionSettingsImpl) BlobBackupAutoProtectionSettings() BaseBlobBackupAutoProtectionSettingsImpl {
	return s
}

var _ BlobBackupAutoProtectionSettings = RawBlobBackupAutoProtectionSettingsImpl{}

// RawBlobBackupAutoProtectionSettingsImpl is returned when the Discriminated Value doesn't match any of the defined types
type RawBlobBackupAutoProtectionSettingsImpl struct {
	blobBackupAutoProtectionSettings BaseBlobBackupAutoProtectionSettingsImpl
	Type                             string
	Values                           map[string]interface{}
}

func (s RawBlobBackupAutoProtectionSettingsImpl) BlobBackupAutoProtectionSettings() BaseBlobBackupAutoProtectionSettingsImpl {
	return s.blobBackupAutoProtectionSettings
}

func UnmarshalBlobBackupAutoProtectionSettingsImplementation(input []byte) (BlobBackupAutoProtectionSettings, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling BlobBackupAutoProtectionSettings into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["objectType"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "BlobBackupRuleBasedAutoProtectionSettings") {
		var out BlobBackupRuleBasedAutoProtectionSettings
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into BlobBackupRuleBasedAutoProtectionSettings: %+v", err)
		}
		return out, nil
	}

	var parent BaseBlobBackupAutoProtectionSettingsImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseBlobBackupAutoProtectionSettingsImpl: %+v", err)
	}

	return RawBlobBackupAutoProtectionSettingsImpl{
		blobBackupAutoProtectionSettings: parent,
		Type:                             value,
		Values:                           temp,
	}, nil
}

var _ BlobBackupAutoProtectionSettings = BlobBackupRuleBasedAutoProtectionSettings{}

type BlobBackupRuleBasedAutoProtectionSettings struct {
	Rules *[]BlobBackupAutoProtectionRule `json:"rules,omitempty"`

	// Fields inherited from BlobBackupAutoProtectionSettings

	Enabled    bool   `json:"enabled"`
	ObjectType string `json:"objectType"`
}

func (s BlobBackupRuleBasedAutoProtectionSettings) BlobBackupAutoProtectionSettings() BaseBlobBackupAutoProtectionSettingsImpl {
	return BaseBlobBackupAutoProtectionSettingsImpl{
		Enabled:    s.Enabled,
		ObjectType: s.ObjectType,
	}
}

var _ json.Marshaler = BlobBackupRuleBasedAutoProtectionSettings{}

func (s BlobBackupRuleBasedAutoProtectionSettings) MarshalJSON() ([]byte, error) {
	type wrapper BlobBackupRuleBasedAutoProtectionSettings
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling BlobBackupRuleBasedAutoProtectionSettings: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling BlobBackupRuleBasedAutoProtectionSettings: %+v", err)
	}

	decoded["objectType"] = "BlobBackupRuleBasedAutoProtectionSettings"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling BlobBackupRuleBasedAutoProtectionSettings: %+v", err)
	}

	return encoded, nil
}

var _ backupinstanceresources.BackupDatasourceParameters = BlobBackupDatasourceParametersForAutoProtection{}

type BlobBackupDatasourceParametersForAutoProtection struct {
	AutoProtectionSettings BlobBackupAutoProtectionSettings `json:"autoProtectionSettings"`

	// Fields inherited from BackupDatasourceParameters

	ObjectType string `json:"objectType"`
}

func (s BlobBackupDatasourceParametersForAutoProtection) BackupDatasourceParameters() backupinstanceresources.BaseBackupDatasourceParametersImpl {
	return backupinstanceresources.BaseBackupDatasourceParametersImpl{
		ObjectType: s.ObjectType,
	}
}

var _ json.Marshaler = BlobBackupDatasourceParametersForAutoProtection{}

func (s BlobBackupDatasourceParametersForAutoProtection) MarshalJSON() ([]byte, error) {
	type wrapper BlobBackupDatasourceParametersForAutoProtection
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling BlobBackupDatasourceParametersForAutoProtection: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling BlobBackupDatasourceParametersForAutoProtection: %+v", err)
	}

	decoded["objectType"] = "BlobBackupDatasourceParametersForAutoProtection"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling BlobBackupDatasourceParametersForAutoProtection: %+v", err)
	}

	return encoded, nil
}

var _ json.Unmarshaler = &BlobBackupDatasourceParametersForAutoProtection{}

func (s *BlobBackupDatasourceParametersForAutoProtection) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		ObjectType string `json:"objectType"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.ObjectType = decoded.ObjectType

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling BlobBackupDatasourceParametersForAutoProtection into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["autoProtectionSettings"]; ok {
		impl, err := UnmarshalBlobBackupAutoProtectionSettingsImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'AutoProtectionSettings' for 'BlobBackupDatasourceParametersForAutoProtection': %+v", err)
		}
		s.AutoProtectionSettings = impl
	}

	return nil
}

// asBlobBackupDatasourceParametersForAutoProtection bridges the gap while the vendored SDK does not know the
// `BlobBackupDatasourceParametersForAutoProtection` discriminator: the SDK hands such entries back as
// `RawBackupDatasourceParametersImpl`, which is re-decoded here. Once the SDK is bumped this collapses to a
// plain type assertion on `backupinstanceresources.BlobBackupDatasourceParametersForAutoProtection`.
func asBlobBackupDatasourceParametersForAutoProtection(input backupinstanceresources.BackupDatasourceParameters) (*BlobBackupDatasourceParametersForAutoProtection, error) {
	switch v := input.(type) {
	case BlobBackupDatasourceParametersForAutoProtection:
		return &v, nil
	case *BlobBackupDatasourceParametersForAutoProtection:
		return v, nil
	case backupinstanceresources.RawBackupDatasourceParametersImpl:
		if !strings.EqualFold(v.Type, "BlobBackupDatasourceParametersForAutoProtection") {
			return nil, nil
		}
		encoded, err := json.Marshal(v.Values)
		if err != nil {
			return nil, fmt.Errorf("re-marshaling raw BackupDatasourceParameters: %+v", err)
		}
		var out BlobBackupDatasourceParametersForAutoProtection
		if err := json.Unmarshal(encoded, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into BlobBackupDatasourceParametersForAutoProtection: %+v", err)
		}
		return &out, nil
	}
	return nil, nil
}
