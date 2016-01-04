// +build amd64,!darwin,!netbsd,!openbsd

package goselect

// darwin, netbsd and openbsd uses uint32 on both amd64 and 386

const (
	// NFDBITS is the amount of bits per mask
	NFDBITS = 8 * 8
)
