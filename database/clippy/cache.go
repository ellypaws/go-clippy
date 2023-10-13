package clippy

import (
	"fmt"
)

type CacheType struct {
	Config User
	Awards []*Award
}

type cacheMap map[string]*CacheType

var Cache = NewCache()

func NewCache() cacheMap {
	return make(cacheMap)
}

func (c cacheMap) Reset() {
	for k := range c {
		delete(c, k)
	}
}
func (c cacheMap) Private(snowflake string) bool {
	return c.GetConfig(snowflake).Private
}

func (c cacheMap) OptOut(snowflake string) bool {
	return c.GetConfig(snowflake).OptOut
}

func (c cacheMap) countAwards(snowflake string) int {
	if _, ok := c[snowflake]; !ok {
		c.allAwards(Request{})
	}
	return len(c[snowflake].Awards)
}

func (c cacheMap) updatePoints(snowflake string) {
	if _, ok := c[snowflake]; !ok {
		c.GetConfig(snowflake)
	} else {
		c[snowflake].Config.Points = c.countAwards(snowflake)
	}
	c[snowflake].Config.Record()
}

func (c cacheMap) QueryPoints(request Request) int {
	if c.OptOut(request.Snowflake) || c.Private(request.Snowflake) {
		return 0
	}
	if _, ok := c[request.Snowflake]; !ok {
		c.GetConfig(request.Snowflake)
	}
	return c[request.Snowflake].Config.Points
}

func (c cacheMap) precacheAwards(request Request) {
	users := getPublicUsers()
	for _, user := range users {
		c[user.Snowflake] = &CacheType{
			Config: *user,
		}
	}
	c.allAwards(request)
}

func (c cacheMap) LeaderboardPrecached(max int, request Request) string {
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
