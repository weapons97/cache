package cache

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/require"
	"reflect"
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
		IndexByLastName: func(indexed any) (key []string) {
			ci := indexed.(*Person)
			return []string{ci.lastName}
		},
		IndexByCountry: func(indexed any) (key []string) {
			ci := indexed.(*Person)
			return []string{ci.country}
		},
	}
}

func (p *Person) ID() (mainKey string) {
	return p.id
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
	index = NewIndexer[*Person]()
)

func init() {
	// set
	index.Set(p1)
	index.Set(p2)
	index.Set(p3)
	index.Set(p4)
	index.Set(p5)
	index.Set(p6)
	index.Set(p7)
}

func any2Slice[T any](ps []any) (rs []T) {
	for _, v := range ps {
		rs = append(rs, v.(T))
	}
	return rs
}

func TestIndexGetByID(t *testing.T) {
	tests := []struct {
		name string
		id   string
		res  *Person
	}{
		{`IndexGetByID.1`,
			`1`,
			p1,
		},
		{`IndexGetByID.2`,
			`2`,
			p2,
		},
		{`IndexGetByID.4`,
			`4`,
			p4,
		},
		{`IndexGetByID.5`,
			`5`,
			p5,
		},
		{`IndexGetByID.7`,
			`7`,
			p7,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rx, ok := index.Get(tt.id)
			require.True(t, ok)
			if !reflect.DeepEqual(rx, tt.res) {
				t.Errorf("got %v, want %v", rx, tt.res)
			}
		})
	}
}

func TestIndexByLastName(t *testing.T) {
	tests := []struct {
		name      string
		indexName string
		indexKey  string
		res       []*Person
	}{
		{`IndexByLastName.魏`,
			IndexByLastName,
			`魏`,
			[]*Person{p1, p2},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rs := index.Search(tt.indexName, tt.indexKey)
			rx := rs.InvokeAll()
			if !reflect.DeepEqual(rx, tt.res) {
				t.Errorf("got %v, want %v", rx, tt.res)
			}
		})
	}
}

func TestIndexByCountry(t *testing.T) {
	tests := []struct {
		name      string
		indexName string
		indexKey  string
		res       []*Person
	}{
		{`IndexByCountry`,
			IndexByCountry,
			`China`,
			[]*Person{p1, p3, p4},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rs := index.Search(tt.indexName, tt.indexKey)
			rx := rs.InvokeAll()
			require.Len(t, rx, len(tt.res))
			if !reflect.DeepEqual(rx, tt.res) {
				t.Errorf("got %v, want %v", rx, tt.res)
			}
		})
	}
}

func TestIndexGetSetFromIndex(t *testing.T) {

	set, e := index.SetFromIndex(IndexByCountry)
	require.NoError(t, e)

	ans := set.List()
	sort.Strings(ans)
	require.Equal(t, []string{`America`, `China`}, ans)
	spew.Dump(ans)
}
