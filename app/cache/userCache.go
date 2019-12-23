package cache

import (
	. "asyncMessageSystem/app/middleware/redis"
)

type UserCache struct {}

func(uc UserCache) Demo()bool{
	rep,errRep := Cache.SetEx("demo",200,1)
	if rep == "OK" && errRep == nil {
		return true
	}else{
		return false
	}
}