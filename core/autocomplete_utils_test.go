package core_test

import (
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/stretchr/testify/assert"
)

func TestProfileFlagFromWords(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		leftWords  []string
		rightWords []string
		want       string
	}{
		{
			name:      "short flag",
			leftWords: []string{"scw", "k8s", "cluster", "get", "-p", "myprofile"},
			want:      "myprofile",
		},
		{
			name:      "long flag",
			leftWords: []string{"scw", "instance", "server", "get", "--profile", "myprofile"},
			want:      "myprofile",
		},
		{
			name:      "long flag equals",
			leftWords: []string{"scw", "k8s", "cluster", "get", "--profile=myprofile"},
			want:      "myprofile",
		},
		{
			name:      "short flag equals",
			leftWords: []string{"scw", "k8s", "cluster", "get", "-p=myprofile"},
			want:      "myprofile",
		},
		{
			name:      "dangling short flag while completing value",
			leftWords: []string{"scw", "k8s", "cluster", "get", "-p"},
			want:      "",
		},
		{
			name:      "dangling long flag while completing value",
			leftWords: []string{"scw", "k8s", "cluster", "get", "--profile"},
			want:      "",
		},
		{
			name:       "profile on the right of the cursor",
			leftWords:  []string{"scw", "k8s", "cluster", "get"},
			rightWords: []string{"-p", "myprofile"},
			want:       "myprofile",
		},
		{
			name:      "last profile wins",
			leftWords: []string{"scw", "-p", "first", "k8s", "get", "-p", "second"},
			want:      "second",
		},
		{
			name:      "no profile",
			leftWords: []string{"scw", "k8s", "cluster", "get"},
			want:      "",
		},
		{
			name:      "dangling flag with next word looking like a flag",
			leftWords: []string{"scw", "get", "-p", "-D"},
			want:      "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.want, core.ProfileFlagFromWords(tt.leftWords, tt.rightWords))
		})
	}
}
