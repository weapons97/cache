package cache

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/require"
	"sort"
	"testing"
)

func TestSetUnion(t *testing.T) {
	s := NewSet[string]()
	s.Add(`a`)
	s.Add(`b`)
	s2 := NewSet[string]()
	s2.Add(`b`)
	s2.Add(`d`)
	s3 := s.Union(s2)
	wantS3 := []string{`a`, `b`, `d`}

	ans := s3.List()
	sort.Strings(ans)
	require.Equal(t, wantS3, ans)
	spew.Dump(s.List(), s2.List(), s3.List())
}

func TestSetJoin(t *testing.T) {
	s := NewSet[string]()
	s.Add(`a`)
	s.Add(`b`)
	s2 := NewSet[string]()
	s2.Add(`b`)
	s2.Add(`d`)
	s3 := s.Join(s2)
	wantS3 := []string{`b`}

	ans := s3.List()
	sort.Strings(ans)
	require.Equal(t, wantS3, ans)
	spew.Dump(s.List(), s2.List(), s3.List())
}

func TestSetJoinLeft(t *testing.T) {
	s := NewSet[string]()
	s.Add(`a`)
	s.Add(`b`)
	s2 := NewSet[string]()
	s2.Add(`b`)
	s2.Add(`d`)
	s3 := s.JoinLeft(s2)
	wantS3 := []string{`a`, `b`}
	ans := s3.List()
	sort.Strings(ans)

	require.Equal(t, wantS3, ans)
	spew.Dump(s.List(), s2.List(), s3.List())
}

func TestSetJoinRight(t *testing.T) {
	s := NewSet[string]()
	s.Add(`a`)
	s.Add(`b`)
	s2 := NewSet[string]()
	s2.Add(`b`)
	s2.Add(`d`)
	s3 := s.JoinRight(s2)
	wantS3 := []string{`b`, `d`}
	ans := s3.List()
	sort.Strings(ans)
	require.Equal(t, wantS3, ans)
	spew.Dump(s.List(), s2.List(), s3.List())
}

func TestSetSub(t *testing.T) {
	s := NewSet[string]()
	s.Add(`a`)
	s.Add(`b`)
	s2 := NewSet[string]()
	s2.Add(`b`)
	s2.Add(`d`)
	s3 := s.Sub(s2)
	wantS3 := []string{`a`}
	ans := s3.List()
	sort.Strings(ans)
	require.Equal(t, wantS3, ans)
	spew.Dump(s.List(), s2.List(), s3.List())
}
