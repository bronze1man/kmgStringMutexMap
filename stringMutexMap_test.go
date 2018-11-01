package kmgStringMutexMap

import (
	"testing"
	"strconv"
	"sync"
)

func TestStringMutexMap(ot *testing.T){
	m:=StringMutexMap{}
	m.LockByString("1")
	m.LockByString("2")
	m.UnlockByString("2")
	m.UnlockByString("1")

	wg:=sync.WaitGroup{}
	num:=0
	wg.Add(10)
	for i:=0;i<10;i++ {
		go func() {
			m.LockByString("1")
			m.LockByString("2")
			num++
			m.UnlockByString("2")
			m.UnlockByString("1")
			wg.Done()
		}()
	}
	wg.Wait()
	if num!=10{
		panic("fail "+strconv.Itoa(num))
	}

	num=0
	wg.Add(10)
	for i:=0;i<10;i++ {
		go func() {
			m.LockByString("1")
			num++
			m.UnlockByString("1")
			wg.Done()
		}()
	}
	wg.Wait()
	if num!=10{
		panic("fail "+strconv.Itoa(num))
	}

	for i:=0;i<1024*2;i++{
		iS:=strconv.Itoa(i)
		m.LockByString(iS)
	}
	for i:=0;i<1024*2;i++{
		iS:=strconv.Itoa(i)
		m.UnlockByString(iS)
	}
	// 1024*2 测试数据量。
	// Benchmark [244.23ns/op] [3.9049MBop/s] duration:[500.178µs] allocNum:[1.00000B/op] allocSize:[16.0000B/op]
	// 使用sync.Pool 减掉一个alloc Benchmark [268.83ns/op] [3.5476MBop/s] duration:[550.555µs] allocNum:[0.00000B/op] allocSize:[0.00000B/op]
	// 使用 数组池子进行缓存 Benchmark [219.91ns/op] [4.3367MBop/s] duration:[450.375µs] allocNum:[0.00000B/op] allocSize:[0.00000B/op]

	// 1024*1024 测试数据量。
	// 使用sync.Pool Benchmark [306.80ns/op] [3.1084MBop/s] duration:[321.707726ms] allocNum:[0.00287B/op] allocSize:[0.44386B/op]
	// 不使用sync.Pool Benchmark [298.83ns/op] [3.1914MBop/s] duration:[313.343433ms] allocNum:[1.00105B/op] allocSize:[16.4377B/op]
	// 使用 数组池子进行缓存 Benchmark [259.22ns/op] [3.6789MBop/s] duration:[271.816786ms] allocNum:[0.00095B/op] allocSize:[0.41692B/op]
	//const benchNum = 1024*2
	//iList:=[]string{}
	//for i:=0;i<benchNum;i++ {
	//	iList = append(iList,strconv.Itoa(i))
	//}
	//kmgTest.Benchmark(func(){
	//	kmgTest.BenchmarkSetNum(benchNum)
	//	for i:=0;i<benchNum;i++{
	//		m.LockByString(iList[i])
	//		m.UnlockByString(iList[i])
	//	}
	//})
}
