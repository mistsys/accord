package accord

// lets figure out the methods we need first and what we'll do with it
// the idea is that the clients will lookup the PSK with different
// mechanisms depending on the deployment requirements
// the same code is shared with the server but it looks them up from different sources
type PSKStore interface {
	GetPSK([]byte) ([]byte, error)
}
