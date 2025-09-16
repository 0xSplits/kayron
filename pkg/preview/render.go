package preview

import (
	"bytes"
	"fmt"

	"github.com/0xSplits/kayron/pkg/cache"
	"github.com/0xSplits/kayron/pkg/hash"
	"github.com/xh3b4sd/tracer"
)

var (
	header = bytes.Join(
		[][]byte{
			[]byte("  #"),
			[]byte("  # AUTO GENERATED PREVIEW DEPLOYMENT"),
			[]byte("  #"),
		},
		[]byte("\n"),
	)
)

func (p *Preview) Render(pre []cache.Object) ([]byte, error) {
	var err error

	// TODO write test
	if len(pre) == 0 {
		return p.inp, nil
	}

	var res Resource
	{
		res = Resource{
			Ser: p.sca.Search([]byte("  Service:")),
			Tas: p.sca.Search([]byte("  TaskDefinition:")),
			Dom: p.sca.Search([]byte("  DomainRecord:")),
			Tar: p.sca.Search([]byte("  TargetGroup:")),
			Lis: p.sca.Search([]byte("  ListenerRule:")),
		}
	}

	var out []byte
	{
		out = append(p.inp, '\n')
	}

	// Derive the base of our listener rule priority range from the provided
	// template defining a AWS::ElasticLoadBalancingV2::ListenerRule resource.

	var pri int
	{
		pri, err = lisPri(res.Lis.Search([]byte("      Priority:")).Bytes())
		if err != nil {
			return nil, tracer.Mask(err)
		}
	}

	for _, x := range pre {
		var hsh hash.Hash
		{
			hsh = hash.New(x.Release.Deploy.Branch.String())
		}

		var dom string
		{
			dom = fmt.Sprintf("%s.%s.${Environment}.splits.org", hsh.Hsh, x.Release.Docker.String())
		}

		var ima []byte
		{
			ima, err = repIma(res.Tas.Search([]byte("          Image:")).Bytes(), []byte(x.Artifact.Reference.Desired))
			if err != nil {
				return nil, tracer.Mask(err)
			}
		}

		// Increment the listener rule priority per preview release. This must be
		// done before the call to render(), because the base priority resolved
		// above is already taken by the main release defining our preview
		// deployments.

		{
			pri++
		}

		{
			out = append(out, p.render(res, dom, pri, hsh, ima)...)
		}
	}

	return out, nil
}

func (p *Preview) render(res Resource, dom string, pri int, hsh hash.Hash, ima []byte) []byte {
	{
		res.Ser = res.Ser.Append([]byte("  Service:"), hsh.Hsh)
		res.Ser = res.Ser.Append([]byte("      ServiceName:"), hsh.Dsh)
		res.Ser = res.Ser.Append([]byte("      TaskDefinition:"), hsh.Hsh)
		res.Ser = res.Ser.Append([]byte("        - TargetGroupArn:"), hsh.Hsh)
		res.Ser = res.Ser.Delete([]byte("      ServiceRegistries:"))
	}

	{
		res.Tas = res.Tas.Append([]byte("  TaskDefinition:"), hsh.Hsh)
		res.Tas = res.Tas.Append([]byte("      Family:"), hsh.Dsh)
		res.Tas = res.Tas.Delete([]byte("          Image:"), ima...)
	}

	{
		res.Dom = res.Dom.Append([]byte("  DomainRecord:"), hsh.Hsh)
		res.Dom = res.Dom.Delete([]byte("      Name:"), fmt.Appendf(nil, `      Name: !Sub "%s"`, dom)...)
	}

	{
		res.Tar = res.Tar.Append([]byte("  TargetGroup:"), hsh.Hsh)
	}

	{
		res.Lis = res.Lis.Append([]byte("  ListenerRule:"), hsh.Hsh)
		res.Lis = res.Lis.Append([]byte("          TargetGroupArn:"), hsh.Hsh)
		res.Lis = res.Lis.Delete([]byte("            Values:"), fmt.Appendf(nil, "            Values:\n              - !Sub \"%s\"", dom)...)
		res.Lis = res.Lis.Delete([]byte("      Priority:"), fmt.Appendf(nil, "      Priority: %d # Host header = %s", pri, dom)...)
	}

	var out []byte
	{
		out = append(out, header...)
		out = append(out, '\n', '\n')
		out = append(out, res.Ser.Bytes()...)
		out = append(out, '\n', '\n')
		out = append(out, res.Tas.Bytes()...)
		out = append(out, '\n', '\n')
		out = append(out, res.Dom.Bytes()...)
		out = append(out, '\n', '\n')
		out = append(out, res.Tar.Bytes()...)
		out = append(out, '\n', '\n')
		out = append(out, res.Lis.Bytes()...)
		out = append(out, '\n', '\n')
	}

	return out
}
