package cache

import (
	"testing"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/require"
)

func TestNewCache(t *testing.T) {
	c := NewCache(WithTTL(time.Second))
	b := 1
	c.Set(`a`, b)
	d, ok := c.Get(`a`)
	require.True(t, ok)
	require.Equal(t, b, d)
	time.Sleep(time.Second)
	d, ok = c.Get(`a`)
	require.False(t, ok)
	require.Equal(t, d, nil)
}

func TestSetUnion(t *testing.T) {
	s := NewSet()
	s.Add(`a`)
	s.Add(`b`)
	s2 := NewSet()
	s2.Add(`b`)
	s2.Add(`d`)
	s3 := s.Union(s2)
	wantS3 := []string{`a`, `b`, `d`}
	require.Equal(t, s3.ListStrings(), wantS3)
	spew.Dump(s.List(), s2.List(), s3.List())
}

func TestSetJoin(t *testing.T) {
	s := NewSet()
	s.Add(`a`)
	s.Add(`b`)
	s2 := NewSet()
	s2.Add(`b`)
	s2.Add(`d`)
	s3 := s.Join(s2)
	wantS3 := []string{`b`}
	require.Equal(t, s3.ListStrings(), wantS3)
	spew.Dump(s.List(), s2.List(), s3.List())
}

func TestSetJoinLeft(t *testing.T) {
	s := NewSet()
	s.Add(`a`)
	s.Add(`b`)
	s2 := NewSet()
	s2.Add(`b`)
	s2.Add(`d`)
	s3 := s.JoinLeft(s2)
	wantS3 := []string{`a`, `b`}
	require.Equal(t, s3.ListStrings(), wantS3)
	spew.Dump(s.List(), s2.List(), s3.List())
}

func TestSetJoinRight(t *testing.T) {
	s := NewSet()
	s.Add(`a`)
	s.Add(`b`)
	s2 := NewSet()
	s2.Add(`b`)
	s2.Add(`d`)
	s3 := s.JoinRight(s2)
	wantS3 := []string{`b`, `d`}
	require.Equal(t, s3.ListStrings(), wantS3)
	spew.Dump(s.List(), s2.List(), s3.List())
}

func TestSetSub(t *testing.T) {
	s := NewSet()
	s.Add(`a`)
	s.Add(`b`)
	s2 := NewSet()
	s2.Add(`b`)
	s2.Add(`d`)
	s3 := s.Sub(s2)
	wantS3 := []string{`a`}
	require.Equal(t, s3.ListStrings(), wantS3)
	spew.Dump(s.List(), s2.List(), s3.List())
}

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
	one := rs.InvokeOne().(*Person)
	require.Equal(t, one.country, `China`)
	spew.Dump(rx)
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
	require.Equal(t, []string{`China`, `America`}, set.ListStrings())
	spew.Dump(set.ListStrings())
}
