// Generates docs for kr8+ code and commands.
//
//go:generate go run ./docs.go
package main

import (
	"github.com/ice-bergtech/kr8/pkg/kr8_types"
	"github.com/invopop/jsonschema"
	"github.com/rs/zerolog/log"
)

type Kr8Cluster kr8_types.Kr8Cluster

type Kr8Component kr8_types.Kr8Cluster

func main() {
	r := new(jsonschema.Reflector)
	if err := r.AddGoComments("github.com/ice-bergtech/kr8/pkg/kr8_types", "../pkg/kr8_types"); err != nil {
		// deal with error
	}
	s := r.Reflect(&Kr8Component{})
	// output
	outfile, err := s.MarshalJSON()
	if err != nil {
		log.Error().Err(err).Msg("issue marshaling jsonschema")
	}
}
