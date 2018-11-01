package kmgStringMutexMap

import "sync"

// a golang data race free map[string]*sync.mutex implement.
// it will free memory after lock and unlock a lot of strings.
type StringMutexMap struct{
	locker sync.RWMutex
	m map[string]*mutexItem
	tryGcNum int
	itemPool []*mutexItem
}

type mutexItem struct{
	mutex sync.Mutex
	lockerNumber int
}

// use sync.Pool make the function slower
//var gStringMutexMapPool sync.Pool

func (m *StringMutexMap) LockByString(s string){
	m.locker.Lock()
	if m.m==nil{
		m.m = map[string]*mutexItem{}
	}
	item :=m.m[s]
	if item ==nil{
		//itemObj:=gStringMutexMapPool.Get()
		//if itemObj!=nil{
		//	item = itemObj.(*mutexItem)
		//}else{
		//	item = &mutexItem{}
		//}
		if len(m.itemPool)>0{
			item = m.itemPool[len(m.itemPool)-1]
			m.itemPool = m.itemPool[:len(m.itemPool)-1]
		}else{
			item = &mutexItem{}
		}
		m.m[s] = item
	}
	item.lockerNumber++
	m.locker.Unlock()
	item.mutex.Lock()
}

func (m *StringMutexMap) UnlockByString(s string){
	m.locker.Lock()
	if m.m==nil{
		m.locker.Unlock()
		panic("StringMutexMap: unlock of unlocked mutex 1")
	}
	item :=m.m[s]
	if item ==nil{
		m.locker.Unlock()
		panic("StringMutexMap: unlock of unlocked mutex 2")
	}
	item.lockerNumber--
	if item.lockerNumber==0 && len(m.m)>= 1024{
		m.tryGcNum++
		if m.tryGcNum>=1024{
			m.tryGcNum = 0
			for k,v:=range m.m{
				if v.lockerNumber==0{
					delete(m.m,k)
					if len(m.itemPool)<1024*2{
						m.itemPool = append(m.itemPool,v)
					}
					//gStringMutexMapPool.Put(v)
				}
			}
		}
	}
	m.locker.Unlock()
	item.mutex.Unlock()
}