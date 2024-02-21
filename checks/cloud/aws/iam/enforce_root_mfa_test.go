package iam

import (
	"testing"

	trivyTypes "github.com/aquasecurity/trivy/pkg/iac/types"

	"github.com/aquasecurity/trivy/pkg/iac/state"

	"github.com/aquasecurity/trivy/pkg/iac/providers/aws/iam"
	"github.com/aquasecurity/trivy/pkg/iac/scan"

	"github.com/stretchr/testify/assert"
)

func TestCheckRootMFAEnabled(t *testing.T) {
	tests := []struct {
		name     string
		input    iam.IAM
		expected bool
	}{
		{
			name: "root user without mfa",
			input: iam.IAM{
				Users: []iam.User{
					{
						Metadata: trivyTypes.NewTestMetadata(),
						Name:     trivyTypes.String("root", trivyTypes.NewTestMetadata()),
					},
				},
			},
			expected: true,
		},
		{
			name: "other user without mfa",
			input: iam.IAM{
				Users: []iam.User{
					{
						Metadata: trivyTypes.NewTestMetadata(),
						Name:     trivyTypes.String("other", trivyTypes.NewTestMetadata()),
					},
				},
			},
			expected: false,
		},
		{
			name: "root user with mfa",
			input: iam.IAM{
				Users: []iam.User{
					{
						Metadata: trivyTypes.NewTestMetadata(),
						Name:     trivyTypes.String("root", trivyTypes.NewTestMetadata()),
						MFADevices: []iam.MFADevice{
							{
								Metadata:  trivyTypes.NewTestMetadata(),
								IsVirtual: trivyTypes.Bool(true, trivyTypes.NewTestMetadata()),
							},
						},
					},
				},
			},
			expected: false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var testState state.State
			testState.AWS.IAM = test.input
			results := checkRootMFAEnabled.Evaluate(&testState)
			var found bool
			for _, result := range results {
				if result.Status() == scan.StatusFailed && result.Rule().LongID() == checkRootMFAEnabled.LongID() {
					found = true
				}
			}
			if test.expected {
				assert.True(t, found, "Rule should have been found")
			} else {
				assert.False(t, found, "Rule should not have been found")
			}
		})
	}
}
