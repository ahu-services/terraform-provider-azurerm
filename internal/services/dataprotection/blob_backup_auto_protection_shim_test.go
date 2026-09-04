// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package dataprotection

import (
	"encoding/json"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dataprotection/2025-07-01/backupinstanceresources"
)

// payload taken from the 2026-03-01 example `PutBackupInstance_BlobBackupAutoProtection.json`
const blobAutoProtectionExample = `{
  "backupDatasourceParametersList": [
    {
      "autoProtectionSettings": {
        "enabled": true,
        "objectType": "BlobBackupRuleBasedAutoProtectionSettings",
        "rules": [
          {
            "objectType": "BlobBackupAutoProtectionRule",
            "mode": "Exclude",
            "type": "Prefix",
            "pattern": "temp-"
          },
          {
            "objectType": "BlobBackupAutoProtectionRule",
            "mode": "Exclude",
            "type": "Prefix",
            "pattern": "test-"
          }
        ]
      },
      "objectType": "BlobBackupDatasourceParametersForAutoProtection"
    }
  ]
}`

func TestBlobBackupAutoProtection_RoundTripThroughSDK(t *testing.T) {
	// the vendored SDK does not know the discriminator, so it must hand back a Raw impl which the shim re-decodes
	var params backupinstanceresources.PolicyParameters
	if err := json.Unmarshal([]byte(blobAutoProtectionExample), &params); err != nil {
		t.Fatalf("unmarshaling example: %+v", err)
	}
	list := pointer.From(params.BackupDatasourceParametersList)
	if len(list) != 1 {
		t.Fatalf("expected 1 datasource parameter, got %d", len(list))
	}
	if _, ok := list[0].(backupinstanceresources.RawBackupDatasourceParametersImpl); !ok {
		t.Fatalf("expected the SDK to return a RawBackupDatasourceParametersImpl, got %T", list[0])
	}

	decoded, err := asBlobBackupDatasourceParametersForAutoProtection(list[0])
	if err != nil {
		t.Fatalf("bridging raw impl: %+v", err)
	}
	if decoded == nil {
		t.Fatalf("expected auto protection parameters, got nil")
	}
	settings, ok := decoded.AutoProtectionSettings.(BlobBackupRuleBasedAutoProtectionSettings)
	if !ok {
		t.Fatalf("expected BlobBackupRuleBasedAutoProtectionSettings, got %T", decoded.AutoProtectionSettings)
	}
	if !settings.Enabled {
		t.Fatalf("expected enabled = true")
	}
	rules := pointer.From(settings.Rules)
	if len(rules) != 2 || rules[0].Pattern != "temp-" || rules[1].Pattern != "test-" {
		t.Fatalf("unexpected rules: %+v", rules)
	}
	if rules[0].Mode != BlobBackupRuleModeExclude || rules[0].Type != BlobBackupPatternTypePrefix || rules[0].ObjectType != "BlobBackupAutoProtectionRule" {
		t.Fatalf("unexpected rule shape: %+v", rules[0])
	}

	// the non-matching discriminator must be ignored, not error
	other, err := asBlobBackupDatasourceParametersForAutoProtection(backupinstanceresources.BlobBackupDatasourceParameters{ContainersList: []string{"a"}})
	if err != nil || other != nil {
		t.Fatalf("expected (nil, nil) for BlobBackupDatasourceParameters, got (%+v, %+v)", other, err)
	}

	// marshaling through the SDK's PolicyParameters must re-produce the discriminators of the example
	out := backupinstanceresources.PolicyParameters{
		BackupDatasourceParametersList: &[]backupinstanceresources.BackupDatasourceParameters{*decoded},
	}
	encoded, err := json.Marshal(out)
	if err != nil {
		t.Fatalf("marshaling: %+v", err)
	}
	var want, got interface{}
	if err := json.Unmarshal([]byte(blobAutoProtectionExample), &want); err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(encoded, &got); err != nil {
		t.Fatal(err)
	}
	wantBytes, _ := json.Marshal(want)
	gotBytes, _ := json.Marshal(got)
	if string(wantBytes) != string(gotBytes) {
		t.Fatalf("round trip mismatch:\n want: %s\n got:  %s", wantBytes, gotBytes)
	}
}
