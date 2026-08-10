package object_test

import (
	"encoding/json"
	"reflect"
	"testing"

	. "github.com/scaleway/scaleway-cli/v2/internal/namespaces/object/v1"
)

func Test_UnmarshalStringList(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    StringList
		wantErr bool
	}{
		{
			name:    "Single string input",
			input:   `"s3:GetObject"`,
			want:    StringList{"s3:GetObject"},
			wantErr: false,
		},

		{
			name:    "Array of strings input",
			input:   `["s3:GetObject", "s3:PutObject"]`,
			want:    StringList{"s3:GetObject", "s3:PutObject"},
			wantErr: false,
		},

		{
			name:    "Empty array",
			input:   `[]`,
			want:    StringList{},
			wantErr: false,
		},

		{
			name:    "Invalid JSON type (number)",
			input:   `123`,
			want:    nil,
			wantErr: true,
		},

		{
			name:    "Invalid JSON syntax",
			input:   `["s3:GetObject"`,
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got StringList
			err := got.UnmarshalJSON([]byte(tt.input))
			if (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UnmarshalJSON() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPrincipal_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    BucketPolicyPrincipal
		wantErr bool
	}{
		{
			name:  "Wildcard string principal",
			input: `"*"`,
			want: BucketPolicyPrincipal{
				Raw: "*",
			},
			wantErr: false,
		},

		{
			name:  "Structured principal with single string values",
			input: `{"SCW": "arn:aws:iam::123456789012:root"}`,
			want: BucketPolicyPrincipal{
				SCW: StringList{"arn:aws:iam::123456789012:root"},
			},
			wantErr: false,
		},

		{
			name:  "Structured principal with array values",
			input: `{"SCW": ["arn:aws:iam::11111:root", "arn:aws:iam::22222:root"]}`,
			want: BucketPolicyPrincipal{
				SCW: StringList{"arn:aws:iam::11111:root", "arn:aws:iam::22222:root"},
			},
			wantErr: false,
		},

		{
			name:    "Invalid principal payload (number)",
			input:   `123`,
			want:    BucketPolicyPrincipal{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got BucketPolicyPrincipal
			err := json.Unmarshal([]byte(tt.input), &got)

			if (err != nil) != tt.wantErr {
				t.Errorf("Principal.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Principal.UnmarshalJSON() got = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestPrincipal_MarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		input   BucketPolicyPrincipal
		want    string
		wantErr bool
	}{
		{
			name:  "Marshal wildcard principal",
			input: BucketPolicyPrincipal{Raw: "*"},
			want:  `"*"`,
		},

		{
			name: "Marshal structured principal",
			input: BucketPolicyPrincipal{
				SCW: StringList{"arn:aws:iam::123456789012:root"},
			},
			want: `{"SCW":["arn:aws:iam::123456789012:root"]}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotBytes, err := tt.input.MarshalJSON()

			if (err != nil) != tt.wantErr {
				t.Errorf("Principal.MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && string(gotBytes) != tt.want {
				t.Errorf("Principal.MarshalJSON() got = %s, want %s", string(gotBytes), tt.want)
			}
		})
	}
}

func TestBucketPolicy_RoundTrip(t *testing.T) {
	jsonInput := `{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "PublicReadGetObject",
      "Effect": "Allow",
      "Principal": "*",
      "Action": [
        "s3:GetObject"
      ],
      "Resource": [
        "arn:aws:s3:::my-bucket/*"
      ]
    },
    {
      "Sid": "IPAllow",
      "Effect": "Deny",
      "Principal": {
        "SCW": [
          "arn:aws:iam::123456789012:root"
        ]
      },
      "Action": [
        "s3:PutObject",
        "s3:DeleteObject"
      ],
      "Resource": [
        "arn:aws:s3:::my-bucket/*"
      ],
      "Condition": {
        "NotIpAddress": {
          "aws:SourceIp": "192.0.2.0/24"
        }
      }
    }
  ]
}`

	expectedStruct := BucketPolicy{
		Version: "2012-10-17",
		Statement: []*BucketPolicyStatement{
			{
				Sid:       "PublicReadGetObject",
				Effect:    "Allow",
				Principal: BucketPolicyPrincipal{Raw: "*"},
				Action:    StringList{"s3:GetObject"},
				Resource:  StringList{"arn:aws:s3:::my-bucket/*"},
			},
			{
				Sid:       "IPAllow",
				Effect:    "Deny",
				Principal: BucketPolicyPrincipal{SCW: StringList{"arn:aws:iam::123456789012:root"}},
				Action:    StringList{"s3:PutObject", "s3:DeleteObject"},
				Resource:  StringList{"arn:aws:s3:::my-bucket/*"},
				Condition: map[string]map[string]any{
					"NotIpAddress": {
						"aws:SourceIp": "192.0.2.0/24",
					},
				},
			},
		},
	}

	var parsedBucketPolicy BucketPolicy
	if err := json.Unmarshal([]byte(jsonInput), &parsedBucketPolicy); err != nil {
		t.Fatalf("Failed to unmarshal valid policy JSON: %v", err)
	}

	if !reflect.DeepEqual(parsedBucketPolicy, expectedStruct) {
		t.Errorf("Parsed struct mismatch.\nGot:  %+v\nWant: %+v", parsedBucketPolicy, expectedStruct)
	}

	marshaledBytes, err := json.MarshalIndent(parsedBucketPolicy, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal BucketPolicy back to JSON: %v", err)
	}

	if string(marshaledBytes) != jsonInput {
		t.Errorf("Marshaled output mismatch.\nGot:\n%s\n\nWant:\n%s", string(marshaledBytes), jsonInput)
	}
}
