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

// func TestToggleSettingString(t *testing.T) {
// 	toggle := toggleSetting("on")
// 	assert.Equal(t, "on", toggle.String())
// }

// func TestToggleSettingSet(t *testing.T) {
// 	tests := []struct {
// 		name     string
// 		input    string
// 		expected string
// 	}{
// 		{
// 			name:     "LowercaseOn",
// 			input:    "on",
// 			expected: "on",
// 		},
// 		{
// 			name:     "LowercaseOff",
// 			input:    "off",
// 			expected: "off",
// 		},
// 		{
// 			name:     "MixedCaseOn",
// 			input:    "oN",
// 			expected: "on",
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			var setting toggleSetting

// 			err := setting.Set(tt.input)
// 			require.NoError(t, err)

// 			require.Equal(t, tt.expected, string(setting))
// 		})
// 	}
// }

// func TestToggleSettingSetError(t *testing.T) {
// 	var setting toggleSetting

// 	err := setting.Set("not-supported")
// 	assert.EqualError(t, err, `supported values are "on" or "off" (case-insensitive)`)
// }

// func TestToggleSettingType(t *testing.T) {
// 	toggle := toggleSetting("on")
// 	assert.Equal(t, "string", toggle.Type())
// }

// // TODO: replace this with a mock/stub in a testing package (ec2/ec2test) ~> captures input
// type imdsAPI struct {
// 	mock.Mock
// }

// func (m *imdsAPI) GetMetadata(ctx context.Context, params *awsimds.GetMetadataInput, optFns ...func(*awsimds.Options)) (*awsimds.GetMetadataOutput, error) {
// 	args := m.Called(ctx, params, optFns)
// 	return args.Get(0).(*awsimds.GetMetadataOutput), args.Error(1)
// }

// func TestToggleMetadataTags(t *testing.T) {
// 	imds := &imdsAPI{}
// 	imds.On("GetMetadata", mock.Anything, mock.Anything, mock.Anything).Return(&awsec2.ModifyInstanceMetadataOptionsOutput{
// 		InstanceId: aws.String("12345"),
// 	}, nil)

// 	tests := []struct {
// 		name     string
// 		toggle   toggleSetting
// 		expected bool
// 	}{}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			ec2 := &ec2test.Client{}

// 			err := toggleMetadataTags(ec2, nil, tt.toggle)
// 			assert.NoError(t, err)
// 		})
// 	}
// }
