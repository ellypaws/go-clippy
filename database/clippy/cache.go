package clippy

import (
	"fmt"
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
		c.allAwards(Request{})
	}
	return len(c[snowflake].Awards)
}

func (c Cache) updatePoints(snowflake string) {
		c.GetConfig(snowflake)
	} else {
		c[snowflake].Config.Points = c.countAwards(snowflake)
	}
	c[snowflake].Config.Record()
}

func (c Cache) QueryPoints(request Request) int {
	if c.OptOut(request.Snowflake) || c.Private(request.Snowflake) {
		return 0
	}
	if _, ok := c[request.Snowflake]; !ok {
		c.GetConfig(request.Snowflake)
	}
	return c[request.Snowflake].Config.Points
}

func (c Cache) precacheAwards(request Request) {
	users := getPublicUsers()
	for _, user := range users {
		c[user.Snowflake] = &CacheType{
			Config: *user,
		}
	}
	c.allAwards(request)
}

func (c Cache) LeaderboardPrecached(max int, request Request) string {
	users := getPublicUsers()
	if len(Cache) == 0 {
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
