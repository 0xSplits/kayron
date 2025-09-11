package preview

import (
	"bytes"
	"fmt"

	"github.com/iancoleman/strcase"
)

func (p *Preview) Render(bra []string) []byte {
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
		out = p.inp
	}

	for _, x := range bra {
		out = append(out, p.render(res, x)...)
	}

	return out
}

func (p *Preview) render(res Resource, bra string) []byte {
	var out []byte

	var cam string
	var keb string
	{
		cam = strcase.ToCamel(bra) // FancyFeatureBranch
		keb = strcase.ToKebab(bra) // fancy-feature-branch
	}

	fmt.Printf("%#v\n", cam)
	fmt.Printf("%#v\n", keb)

	{
		res.Ser = bytes.Replace(res.Ser, []byte("  Service:"), fmt.Appendf(nil, "  Service%s:", cam), 1)
		res.Tas = bytes.Replace(res.Tas, []byte("  TaskDefinition:"), fmt.Appendf(nil, "  TaskDefinition%s:", cam), 1)
		res.Dom = bytes.Replace(res.Dom, []byte("  DomainRecord:"), fmt.Appendf(nil, "  DomainRecord%s:", cam), 1)
		res.Tar = bytes.Replace(res.Tar, []byte("  TargetGroup:"), fmt.Appendf(nil, "  TargetGroup%s:", cam), 1)
		res.Lis = bytes.Replace(res.Lis, []byte("  ListenerRule:"), fmt.Appendf(nil, "  ListenerRule%s:", cam), 1)
	}

	{
		out = append(out, '\n')
		out = append(out, res.Ser...)
		out = append(out, '\n')
		out = append(out, res.Tas...)
		out = append(out, '\n')
		out = append(out, res.Dom...)
		out = append(out, '\n')
		out = append(out, res.Tar...)
		out = append(out, '\n')
		out = append(out, res.Lis...)
	}

	return out
}
