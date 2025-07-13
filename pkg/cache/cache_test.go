package cache

import "testing"

func Test_Cache_Lifecycle(t *testing.T) {
	var cac Interface[string, int]
	{
		cac = New[string, int]()
	}

	var foo string
	var bar string
	var baz string
	{
		foo = "foo"
		bar = "bar"
		baz = "baz"
	}

	{
		siz := cac.Length()
		if siz != 0 {
			t.Fatal("expected", 0, "got", siz)
		}
	}

	{
		cac.Update(foo, 33)
		cac.Update(bar, 47)
	}

	{
		siz := cac.Length()
		if siz != 2 {
			t.Fatal("expected", 2, "got", siz)
		}
	}

	{
		exi := cac.Exists(foo)
		if !exi {
			t.Fatal("expected", true, "got", false)
		}
	}

	{
		exi := cac.Exists(bar)
		if !exi {
			t.Fatal("expected", true, "got", false)
		}
	}

	{
		exi := cac.Exists(baz)
		if exi {
			t.Fatal("expected", false, "got", true)
		}
	}

	{
		val, _ := cac.Search(foo)
		if val != 33 {
			t.Fatal("expected", 33, "got", val)
		}
	}

	{
		val, _ := cac.Search(bar)
		if val != 47 {
			t.Fatal("expected", 47, "got", val)
		}
	}

	{
		val, _ := cac.Search(baz)
		if val != 0 {
			t.Fatal("expected", 0, "got", val)
		}
	}

	{
		cac.Update(foo, 99)
	}

	{
		siz := cac.Length()
		if siz != 2 {
			t.Fatal("expected", 2, "got", siz)
		}
	}

	{
		exi := cac.Exists(foo)
		if !exi {
			t.Fatal("expected", true, "got", false)
		}
	}

	{
		exi := cac.Exists(bar)
		if !exi {
			t.Fatal("expected", true, "got", false)
		}
	}

	{
		exi := cac.Exists(baz)
		if exi {
			t.Fatal("expected", false, "got", true)
		}
	}

	{
		val, _ := cac.Search(foo)
		if val != 99 {
			t.Fatal("expected", 99, "got", val)
		}
	}

	{
		val, _ := cac.Search(bar)
		if val != 47 {
			t.Fatal("expected", 47, "got", val)
		}
	}

	{
		val, _ := cac.Search(baz)
		if val != 0 {
			t.Fatal("expected", 0, "got", val)
		}
	}

	{
		cac.Delete(bar)
	}

	{
		siz := cac.Length()
		if siz != 1 {
			t.Fatal("expected", 1, "got", siz)
		}
	}

	{
		exi := cac.Exists(foo)
		if !exi {
			t.Fatal("expected", true, "got", false)
		}
	}

	{
		exi := cac.Exists(bar)
		if exi {
			t.Fatal("expected", false, "got", true)
		}
	}

	{
		exi := cac.Exists(baz)
		if exi {
			t.Fatal("expected", false, "got", true)
		}
	}

	{
		val, _ := cac.Search(foo)
		if val != 99 {
			t.Fatal("expected", 99, "got", val)
		}
	}

	{
		val, _ := cac.Search(bar)
		if val != 0 {
			t.Fatal("expected", 0, "got", val)
		}
	}

	{
		val, _ := cac.Search(baz)
		if val != 0 {
			t.Fatal("expected", 0, "got", val)
		}
	}

	{
		siz := cac.Length()
		if siz != 1 {
			t.Fatal("expected", 1, "got", siz)
		}
	}

	{
		cac.Delete(foo)
	}

	{
		siz := cac.Length()
		if siz != 0 {
			t.Fatal("expected", 0, "got", siz)
		}
	}
}
