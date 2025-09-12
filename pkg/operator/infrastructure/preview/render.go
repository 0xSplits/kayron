package preview

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

// TODO scanner needs Replace and Delete
func (p *Preview) render(res Resource, bra string) []byte {
	var out []byte

	var hsh string
	{
		hsh = Hash(bra)
	}

	{
		res.Ser = res.Ser.Append([]byte("  Service:"), []byte(hsh))
		res.Ser = res.Ser.Append([]byte("      ServiceName:"), []byte("-"+hsh))
		res.Ser = res.Ser.Append([]byte("      TaskDefinition:"), []byte(hsh))
		res.Ser = res.Ser.Append([]byte("        - TargetGroupArn:"), []byte(hsh))
	}

	{
		res.Tas = res.Tas.Append([]byte("  TaskDefinition:"), []byte(hsh))
		res.Tas = res.Tas.Append([]byte("      Family:"), []byte("-"+hsh))
	}

	{
		res.Dom = res.Dom.Append([]byte("  DomainRecord:"), []byte(hsh))
	}

	{
		res.Tar = res.Tar.Append([]byte("  TargetGroup:"), []byte(hsh))
	}

	{
		res.Lis = res.Lis.Append([]byte("  ListenerRule:"), []byte(hsh))
	}

	{
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
	}

	return out
}
