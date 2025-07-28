package artifact

import (
	"github.com/0xSplits/kayron/pkg/release/artifact"
	"github.com/0xSplits/kayron/pkg/release/artifact/container"
	"github.com/0xSplits/kayron/pkg/release/artifact/identity"
	"github.com/0xSplits/kayron/pkg/release/artifact/reference"
	"github.com/0xSplits/kayron/pkg/release/schema/service"
	"github.com/xh3b4sd/choreo/parallel"
	"github.com/xh3b4sd/tracer"
	"golang.org/x/sync/errgroup"
)

func (h *Handler) Ensure() error {
	var err error

	var fil []service.Service
	for i := range h.rel.Length() {
		var ser service.Service
		var exi bool
		{
			ser, exi = h.rel.Search(i)
			if !exi {
				return tracer.Mask(serviceNotCachedError)
			}
			if ser.Provider == "cloudformation" {
				continue
			}
			if ser.Deploy.Suspend {
				continue
			}
		}

		{
			fil = append(fil, ser)
		}
	}

	var art []artifact.Artifact
	{
		art = make([]artifact.Artifact, len(fil))
	}

	fnc := func(i int, x service.Service) error {
		con, ref, ide, err := h.ociRef(x)
		if err != nil {
			return tracer.Mask(err)
		}

		art[i] = artifact.Artifact{
			Container: con,
			Reference: ref,
			Identity:  ide,
		}

		return nil
	}

	{
		err = parallel.Slice(fil, fnc)
		if err != nil {
			return tracer.Mask(err)
		}
	}

	return nil
}

func (h *Handler) ociRef(s service.Service) (container.Struct, reference.Struct, identity.String, error) {
	var err error

	var grp errgroup.Group
	{
		grp = errgroup.Group{}
	}

	var con container.Struct

	{
		grp.Go(func() error {
			var tag string
			{
				// TODO always lookup docker image tag from remote registry
			}

			con = container.Struct{
				Tag:        tag,
				Repository: s.Docker.String(),
			}

			return nil
		})
	}

	var ref reference.Struct

	if !s.Deploy.Branch.Empty() {
		grp.Go(func() error {
			var val string
			{
				val = s.Deploy.Release.String()
			}

			{
				// TODO lookup github commit sha from remote repository if setting is a branch
				//
				//     !s.Deploy.Branch.Empty()     s.Deploy.String() is a branch
				//     !s.Deploy.Release.Empty()    s.Deploy.String() is a tag
				//
			}

			ref = reference.Struct{
				Repository: s.Github.String(),
				Value:      val,
			}

			return nil
		})
	}

	var ide identity.String
	{
		ide = identity.String(s.Docker) // use docker for now, should be more reliable
	}

	{
		err = grp.Wait()
		if err != nil {
			return container.Struct{}, reference.Struct{}, identity.String(""), tracer.Mask(err)
		}
	}

	return con, ref, ide, nil
}
