package clippy

import (
	"fmt"
	"sync"
)

type CacheType struct {
	Config User
	Awards []*Award
}

type Cache struct {
	Map   map[string]*CacheType
	Mutex *sync.RWMutex
}

var cache *Cache

func GetCache() *Cache {
	if cache != nil {
		return cache
	}
	return &Cache{
		Map:   map[string]*CacheType{},
		Mutex: &sync.RWMutex{},
	}
}

func (c Cache) Reset() {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()
	for k := range c.Map {
		delete(c.Map, k)
	}
}

func (c Cache) Private(snowflake string) bool {
	user, exist := c.GetConfig(snowflake)
	if !exist {
		return false
	}
	return user.Private
}

func (c Cache) OptOut(snowflake string) bool {
	user, exist := c.GetConfig(snowflake)
	if !exist {
		return false
	}
	return user.OptOut
}

func (c Cache) countAwards(snowflake string) int {
	c.Mutex.RLock()
	_, ok := c.Map[snowflake]
	c.Mutex.RUnlock()
	if !ok || c.Map[snowflake].Awards == nil {
		c.allAwards(Request{})
	}
	c.Mutex.RLock()
	defer c.Mutex.RUnlock()
	return len(c.Map[snowflake].Awards)
}

func (c Cache) synchronizePoints(snowflake string) {
	c.Mutex.RLock()
	_, ok := c.Map[snowflake]
	c.Mutex.RUnlock()
	if !ok {
		c.GetConfig(snowflake)
	}
	c.Mutex.RLock()
	configPoints := c.Map[snowflake].Config.Points
	c.Mutex.RUnlock()
	countedAwards := c.countAwards(snowflake)
	if configPoints == countedAwards {
		return
	}
	c.Mutex.Lock()
	c.Map[snowflake].Config.Points = countedAwards
	c.Mutex.Unlock()

	c.Mutex.RLock()
	defer c.Mutex.RUnlock()
	c.Map[snowflake].Config.Record()
}

func (c Cache) SynchronizeAllPoints() {
	users := getAllUsers()
	for _, user := range users {
		c.synchronizePoints(user.Snowflake)
	}
}

func (c Cache) QueryPoints(request Request) int {
	if c.OptOut(request.Snowflake) || c.Private(request.Snowflake) {
		return 0
	}
	c.Mutex.RLock()
	_, ok := c.Map[request.Snowflake]
	c.Mutex.RUnlock()
	if ok {
		c.GetConfig(request.Snowflake)
	}
	c.Mutex.RLock()
	defer c.Mutex.RUnlock()
	return c.Map[request.Snowflake].Config.Points
}

func (c Cache) precacheAwards(request Request) {
	users := getPublicUsers()
	for _, user := range users {
		c.Mutex.Lock()
		c.Map[user.Snowflake] = &CacheType{
			Config: *user,
		}
		c.Mutex.Unlock()
	}
	c.allAwards(request)
}

func (c Cache) LeaderboardPrecached(max int, request Request) string {
	users := getPublicUsers()
	if len(c.Map) == 0 {
		c.precacheAwards(request)
	}

	userPointsSlice = []userPoints{}
	for _, u := range users {
		userPointsSlice = append(userPointsSlice, userPoints{
			Snowflake: u.Snowflake,
			Username:  u.Username,
			Points:    u.Points,
		})
	}

	sortUserPoints(userPointsSlice)
	if max == 0 {
		max = len(userPointsSlice)
	}
	userPointsSlice = userPointsSlice[:min(max, len(userPointsSlice))]

	var leaderboard string = "Leaderboard:\n"
	for i, user := range userPointsSlice {
		leaderboard += fmt.Sprintf("%d. %v (%v)\n", i+1, user.Username, user.Points)
	}

	return leaderboard[:len(leaderboard)-1]
}
