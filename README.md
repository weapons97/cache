### cache 是一个带索引带超时的缓存库
目的在于优化代码结构，提供了若干实践。
#### example
### 实践1 定义泛型函数
1.18 已经发布一段实践了。通过泛型函数。我们可以减少循环的使用，优化代码结构。下面分享几个泛型函数和代码上的实践。
#### Filter 函数
```go
// Filter filter one slice
func Filter[T any](objs []T, filter func(obj T) bool) []T {
	res := make([]T, 0, len(objs))
	for i := range objs {
		ok := filter(objs[i])
		if ok {
			res = append(res, objs[i])
		}
	}
	return res
}
```

```go
// 测试[]int
func TestFilter(t *testing.T) {
	ans := []int{2, 4, 6}
	a := []int{1, 2, 3, 4, 5, 6}
	b := Filter(a, func(i int) bool {
        return i%2 == 0
	})
	require.Equal(t, ans, b)
	spew.Dump(b)
}
// 结果
=== RUN   TestFilter
([]int) (len=3 cap=6) {
 (int) 2,
 (int) 4,
 (int) 6
}
--- PASS: TestFilter (0.00s)
PASS

// NoSpace is filter func for strings
func NoSpace(s string) bool {
	return strings.TrimSpace(s) != ""
}
// 测试[]sting
func TestFilterNoSpace(t *testing.T) {
	ans1 := []string{"1", "2", "3"}
	a := []string{"", "1", "", "2", "", "3", ""}
	b := Filter(a, NoSpace)
	require.Equal(t, ans1, b)
	spew.Dump(b)
}
// 结果
=== RUN   TestFilterNoSpace
([]string) (len=3 cap=7) {
 (string) (len=1) "1",
 (string) (len=1) "2",
 (string) (len=1) "3"
}
--- PASS: TestFilterNoSpace (0.00s)
PASS
```
#### Map 函数
```go
// Map one slice
func Map[T any, K any](objs []T, mapper func(obj T) ([]K, bool)) []K {
	res := make([]K, 0, len(objs))
	for i := range objs {
		others, ok := mapper(objs[i])
		if ok {
			res = append(res, others...)
		}
	}
	return res
}
// 测试 []int -> []string
func TestMap(t *testing.T) {
	ans := []string{"2", "4", "6", "end"}
	a := []int{1, 2, 3, 4, 5, 6}
	b := Map(a, func(i int) ([]string, bool) {
		if i == 6 {
			return []string{fmt.Sprintf(`%v`, i), `end`}, true
		}
		if i%2 == 0 {
			return []string{fmt.Sprintf(`%v`, i)}, true
		} else {
			return nil, false
		}
	})
	require.Equal(t, ans, b)
	spew.Dump(b)
}
// 结果
=== RUN   TestMap
([]string) (len=4 cap=6) {
 (string) (len=1) "2",
 (string) (len=1) "4",
 (string) (len=1) "6",
 (string) (len=3) "end"
}
--- PASS: TestMap (0.00s)
PASS
```
### First 函数
```go
// First make return first for slice
func First[T any](objs []T) (T, bool) {
	if len(objs) > 0 {
		return objs[0], true
	}
	return *new(T), false
}

func TestFirstInt(t *testing.T) {
	ans1, ans2 := 1, 0
	a := []int{1, 2, 3, 4, 5, 6}
	b, ok := First(a)
	require.True(t, ok)
	require.Equal(t, ans1, b)
	spew.Dump(b)
	c := []int{}
	d, ok := First(c)
	require.False(t, ok)
	require.Equal(t, ans2, d)
	spew.Dump(d)
}
// result
=== RUN   TestFirstInt
(int) 1
(int) 0
--- PASS: TestFirstInt (0.00s)
PASS

func TestFirstString(t *testing.T) {
	ans1, ans2 := "1", ""
	a := []string{"1", "2", "3", "4", "5", "6"}
	b, ok := First(a)
	require.True(t, ok)
	require.Equal(t, ans1, b)
	spew.Dump(b)
	c := []string{}
	d, ok := First(c)
	require.False(t, ok)
	require.Equal(t, ans2, d)
	spew.Dump(d)
}
// result
=== RUN   TestFirstString
(string) (len=1) "1"
(string) ""
--- PASS: TestFirstString (0.00s)
PASS
```
### 实践2 带超时的cache
- 某些情况下，我们删除过期的cache， 通过利用带超时的cache，简化代码
- cache 结构 github.com/weapons97/cache/cache.go
```go
// 用辅助map删除
if apiRet.TotalCount > 0 {
 var hc sync.Map
 for _, h := range apiRet.Hcis {
  hc.Store(h.HostID, h)
  hostCpu.Store(h.HostID, h)
 }
 hostCpu.Range(func(key, _ interface{}) bool {
  _, ok := hc.Load(key)
  if !ok {
   hostCpu.Delete(key)
  }
  return true
 })
}
// 直接设置，过期的key 会删除
for _, h := range apiRet.Hcis {
	hostCpu.Set(h.HostID, h)
}
```
```go
func TestNewCache(t *testing.T) {
	c := NewCache(WithTTL[string, int](time.Second))
	b := 1
	c.Set(`a`, b)
	d, ok := c.Get(`a`)
	require.True(t, ok)
	require.Equal(t, b, d)
	time.Sleep(time.Second)
	d, ok = c.Get(`a`)
	require.False(t, ok)
	// 超时返回0值
	require.Equal(t, d, 0)
}
```
### 实践3 集合操作
通过 set 做集合，可以给集合去重。可以给结合相并，想交，等操作。
set 结构 github.com/weapons97/cache/set.go
```go
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

```
- 通过set 去重
```go
// ShowImageInManifest 抓取 manifest 中imgs
func ShowImageInManifest(manifest string) (imgs []string) {
	rx := regImages.FindAllStringSubmatch(manifest, -1)
	set := cache.NewSet[string]()
	for i := range rx {
		for j := range rx[i] {
			if strings.HasPrefix(rx[i][j], `image:`) {
				continue
			}
			tx0 := strings.TrimSpace(rx[i][j])
			tx1 := strings.Trim(tx0, `'`)
			tx2 := strings.Trim(tx1, `"`)
			set.Add(tx2)
		}
	}
	imgs = set.List()
	return imgs
}
```
#### 4 带索引的cache
- 某些情况下，我们可能根据cache 的某个元素对cache进行遍历，这时候如果给cache 加上索引结构，可以对遍历加速。
- index 结构 github.com/weapons97/cache/index.go
```go
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

```

```go
// 测试数据
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
```

```go 
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
// result
=== RUN   TestIndexByCountry
([]interface {}) (len=3 cap=3) {
 (*cache.Person)(0x14139c0)({
  id: (string) (len=1) "3",
  lastName: (string) (len=3) "李",
  fullName: (string) (len=6) "李云",
  country: (string) (len=5) "China"
 }),
 (*cache.Person)(0x1413a00)({
  id: (string) (len=1) "4",
  lastName: (string) (len=3) "黄",
  fullName: (string) (len=9) "黄帅来",
  country: (string) (len=5) "China"
 }),
 (*cache.Person)(0x1413940)({
  id: (string) (len=1) "1",
  lastName: (string) (len=3) "魏",
  fullName: (string) (len=6) "魏鹏",
  country: (string) (len=5) "China"
 })
}
(*cache.Person)(0x14139c0)({
 id: (string) (len=1) "3",
 lastName: (string) (len=3) "李",
 fullName: (string) (len=6) "李云",
 country: (string) (len=5) "China"
})
--- PASS: TestIndexByCountry (0.00s)
PASS

```