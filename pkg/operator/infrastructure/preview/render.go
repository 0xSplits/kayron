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

func (p *Preview) Render(art []cache.Object) ([]byte, error) {
	var err error

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

	var pri int
	{
		pri = 30
	}

	for _, x := range art {
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

		{
			out = append(out, p.render(res, dom, pri, hsh, ima)...)
		}

		{
			pri++
		}
	}

	return out, nil
}

// TODO scanner needs Replace and Delete
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
