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

type NameContainerIndexed struct {
	shortName string
	fullName  string
	setID     string
}

const (
	IndexByShortName = `IndexByShortName`
	IndexBySetID     = `IndexBySetID`
)

func (nci *NameContainerIndexed) Indexs() map[string]IndexFunc {
	return map[string]IndexFunc{
		IndexByShortName: func(indexed Indexed) (key []string) {
			ci := indexed.(*NameContainerIndexed)
			return []string{ci.shortName}
		},
		IndexBySetID: func(indexed Indexed) (key []string) {
			ci := indexed.(*NameContainerIndexed)
			return []string{ci.setID}
		},
	}
}

func (nci *NameContainerIndexed) Id() (mainKey string) {
	return nci.fullName
}

func (nci *NameContainerIndexed) Set(v interface{}) (Indexed, bool) {
	rx, ok := v.(*NameContainerIndexed)
	if !ok {
		return nil, false
	}
	return rx, true
}

func (nci *NameContainerIndexed) Get(v Indexed) (interface{}, bool) {
	rx, ok := v.(*NameContainerIndexed)
	if !ok {
		return nil, false
	}
	return rx, true
}

var (
	testC1 = &NameContainerIndexed{
		shortName: "uhost-access",
		fullName:  "/NS/region10027/set11/uhost/access",
		setID:     `11`,
	}
	testC1_ = &NameContainerIndexed{
		shortName: "uhost-access",
		fullName:  "/NS/region10027/set11/uhost/access",
	}
	testC2 = &NameContainerIndexed{
		shortName: "uhost-manager",
		fullName:  "/NS/region10027/set11/uhost/manager",
		setID:     `11`,
	}
	testC3 = &NameContainerIndexed{
		shortName: "uimage3-access",
		fullName:  "/NS/region10027/set11/uimage3/access",
		setID:     `11`,
	}
	testC4 = &NameContainerIndexed{
		shortName: "uhost-access",
		fullName:  "/NS/region10027/set10/uhost/access",
		setID:     `10`,
	}
	testC5 = &NameContainerIndexed{
		shortName: "uhost-manager",
		fullName:  "/NS/region10027/set10/uhost/manager",
		setID:     `10`,
	}
	testC6 = &NameContainerIndexed{
		shortName: "uimage3-access",
		fullName:  "/NS/region10027/set10/uimage3/access",
		setID:     `10`,
	}
)

func TestIndex(t *testing.T) {
	index := NewIndexer(&NameContainerIndexed{})
	// set
	index.Set(testC1)
	index.Set(testC2)
	index.Set(testC3)
	index.Set(testC4)
	index.Set(testC5)
	index.Set(testC6)
	// search
	rs := index.Search(IndexBySetID, `10`)
	require.False(t, rs.Failed())
	rx := rs.InvokeAll()
	require.Len(t, rx, 3)
	one := rs.InvokeOne().(*NameContainerIndexed)
	require.Equal(t, one.setID, `10`)
	spew.Dump(rx)
	// search
	rs = index.Search(IndexByShortName, `uhost-access`)
	rx = rs.InvokeAll()
	require.Len(t, rx, 2)
	one = rs.InvokeOne().(*NameContainerIndexed)
	require.Equal(t, one.shortName, `uhost-access`)
	spew.Dump(rx)
	// search
	v, ok := index.Get(`/NS/region10027/set10/uhost/access`)
	require.True(t, ok)
	require.Equal(t, v, testC4)
	index.Del(testC4)
	// search
	rs = index.Search(IndexByShortName, `uhost-access`)
	rx = rs.InvokeAll()
	require.Len(t, rx, 1)
	one = rs.InvokeOne().(*NameContainerIndexed)
	require.Equal(t, one, testC1)
	spew.Dump(rx)
}

func TestIndex2(t *testing.T) {
	index := NewIndexer(&NameContainerIndexed{})
	// set
	index.Set(testC1)
	index.Set(testC1_)
	rs := index.Search(IndexBySetID, `11`)
	x := rs.InvokeAll()
	spew.Dump(x)
}