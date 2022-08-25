package cache

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/require"
	"sort"
	"testing"
)

type Person struct {
	id       string
	lastName string
	fullName string
	country  string
}

const (
	IndexByLastName = `IndexByLastName`
	IndexByCountry  = `IndexByCountry`
)

func (p *Person) Indexs() map[string]IndexFunc {
	return map[string]IndexFunc{
		IndexByLastName: func(indexed Indexed) (key []string) {
			ci := indexed.(*Person)
			return []string{ci.lastName}
		},
		IndexByCountry: func(indexed Indexed) (key []string) {
			ci := indexed.(*Person)
			return []string{ci.country}
		},
	}
}

func (p *Person) ID() (mainKey string) {
	return p.id
}

func (p *Person) Set(v interface{}) (Indexed, bool) {
	rx, ok := v.(*Person)
	if !ok {
		return nil, false
	}
	return rx, true
}

func (p *Person) Get(v Indexed) (interface{}, bool) {
	rx, ok := v.(*Person)
	if !ok {
		return nil, false
	}
	return rx, true
}

var (
	p1 = &Person{
		id:       `1`,
		lastName: "魏",
		fullName: "魏鹏",
		country:  `China`,
	}
	p2 = &Person{
		id:       `2`,
		lastName: "魏",
		fullName: "魏无忌",
		country:  `America`,
	}
	p3 = &Person{
		id:       `3`,
		lastName: "李",
		fullName: "李云",
		country:  `China`,
	}
	p4 = &Person{
		id:       `4`,
		lastName: "黄",
		fullName: "黄帅来",
		country:  `China`,
	}
	p5 = &Person{
		id:       `5`,
		lastName: "Cook",
		fullName: "TimCook",
		country:  `America`,
	}
	p6 = &Person{
		id:       `6`,
		lastName: "Jobs",
		fullName: "SteveJobs",
		country:  `America`,
	}
	p7 = &Person{
		id:       `7`,
		lastName: "Musk",
		fullName: "Elon Musk",
		country:  `America`,
	}
)

func TestIndexByCountry(t *testing.T) {
	index := NewIndexer(&Person{})
	// set
	index.Set(p1)
	index.Set(p2)
	index.Set(p3)
	index.Set(p4)
	index.Set(p5)
	index.Set(p6)
	index.Set(p7)

	// search
	rs := index.Search(IndexByCountry, `China`)
	require.False(t, rs.Failed())
	rx := rs.InvokeAll()
	require.Len(t, rx, 3)
	spew.Dump(rx)
	one := rs.InvokeOne().(*Person)
	require.Equal(t, one.country, `China`)
	spew.Dump(one)
}

func TestIndexGetByID(t *testing.T) {
	index := NewIndexer(&Person{})
	// set
	index.Set(p1)
	index.Set(p2)
	index.Set(p3)
	index.Set(p4)
	index.Set(p5)
	index.Set(p6)
	index.Set(p7)
	v, ok := index.Get(`7`)
	require.True(t, ok)
	require.Equal(t, v, p7)
}

func TestIndexByLastName(t *testing.T) {
	index := NewIndexer(&Person{})
	// set
	index.Set(p1)
	index.Set(p2)
	index.Set(p3)
	index.Set(p4)
	index.Set(p5)
	index.Set(p6)
	index.Set(p7)
	// search
	rs := index.Search(IndexByLastName, `魏`)
	rx := rs.InvokeAll()
	require.Len(t, rx, 2)
	one := rs.InvokeOne().(*Person)
	require.Equal(t, one.lastName, `魏`)
	spew.Dump(rx)
}

func TestIndexGetSetFromIndex(t *testing.T) {
	index := NewIndexer(&Person{})
	// set
	index.Set(p1)
	index.Set(p2)
	index.Set(p3)
	index.Set(p4)
	index.Set(p5)
	index.Set(p6)
	index.Set(p7)

	set, e := index.SetFromIndex(IndexByCountry)
	require.NoError(t, e)

	ans := set.List()
	sort.Strings(ans)
	require.Equal(t, []string{`America`, `China`}, ans)
	spew.Dump(ans)
}
