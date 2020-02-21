package mattrax

import "testing"

func NewMockServer(t *testing.T) *Server {
	// TODO: Make this fully functional. Init using other packages so it is functional and tests everything
	return &Server{
		Version: "0.0.0-test",
	}
}
