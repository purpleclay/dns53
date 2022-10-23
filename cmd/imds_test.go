/*
Copyright (c) 2022 Purple Clay

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestToggleSettingString(t *testing.T) {
	toggle := toggleSetting("on")
	assert.Equal(t, "on", toggle.String())
}

func TestToggleSettingSet(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "LowercaseOn",
			input:    "on",
			expected: "on",
		},
		{
			name:     "LowercaseOff",
			input:    "off",
			expected: "off",
		},
		{
			name:     "MixedCaseOn",
			input:    "oN",
			expected: "on",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var setting toggleSetting

			err := setting.Set(tt.input)
			require.NoError(t, err)

			require.Equal(t, tt.expected, string(setting))
		})
	}
}

func TestToggleSettingSetError(t *testing.T) {
	var setting toggleSetting

	err := setting.Set("not-supported")
	assert.EqualError(t, err, `supported values are "on" or "off" (case-insensitive)`)
}

func TestToggleSettingType(t *testing.T) {
	toggle := toggleSetting("on")
	assert.Equal(t, "string", toggle.Type())
}

// TODO: test through the command directly

// func TestToggleMetadataTags(t *testing.T) {
// 	tests := []struct {
// 		name     string
// 		toggle   toggleSetting
// 		expected string
// 	}{
// 		{
// 			name:     "On",
// 			toggle:   toggleSettingOn,
// 			expected: "enabled",
// 		},
// 		{
// 			name:     "Off",
// 			toggle:   toggleSettingOff,
// 			expected: "disabled",
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			mockEC2 := ec2mock.New(t)
// 			mockEC2.On("ModifyInstanceMetadataOptions", mock.Anything, mock.MatchedBy(func(req *ec2.ModifyInstanceMetadataOptionsInput) bool {
// 				return req.InstanceMetadataTags == types.InstanceMetadataTagsState(tt.expected)
// 			}), mock.Anything).Return(&ec2.ModifyInstanceMetadataOptionsOutput{}, nil)

// 			err := toggleMetadataTags(mockEC2, imdsstub.New(t), tt.toggle)
// 			assert.NoError(t, err)
// 		})
// 	}
// }
