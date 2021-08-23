//+build amd64

package sm3

//go:noescape
func update(digest *[8]uint32, a []byte, b []byte)
