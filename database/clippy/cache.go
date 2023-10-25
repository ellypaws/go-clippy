package clippy

import (
	"fmt"
	"github.com/charmbracelet/bubbles/table"
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

func (c Cache) Reset() *Cache {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()
	for k := range c.Map {
		delete(c.Map, k)
	}
	return &c
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
		program.Send(fmt.Sprintf("User %v has %v points, which is correct.", snowflake, configPoints))
		return
	} else {
		program.Send(fmt.Sprintf("User %v has %v points, but %v awards.", snowflake, configPoints, countedAwards))
	}
	c.Mutex.Lock()
	c.Map[snowflake].Config.Points = countedAwards
	c.Mutex.Unlock()

	c.Mutex.RLock()
	defer c.Mutex.RUnlock()
	c.Map[snowflake].Config.Record()
}

func (c Cache) SynchronizeAllPoints() *Cache {
	users := getAllUsers()
	for _, user := range users {
		c.synchronizePoints(user.Snowflake)
	}
	return &c
}

func (c Cache) QueryPoints(request Request) int {
	if c.OptOut(request.Snowflake) || c.Private(request.Snowflake) {
		return 0
	}
	c.Mutex.RLock()
	_, ok := c.Map[request.Snowflake]
	c.Mutex.RUnlock()
	// TODO Check if this is necessary
	if !ok {
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

//	columns := []table.Column{
//		{Title: "#", Width: 4},
//		{Title: "Username", Width: 12},
//		{Title: "Snowflake", Width: 12},
//		{Title: "Points", Width: 8},
//		{Title: "Private", Width: 8},
//	}
//
//	rows := []table.Row{
//		{"1", "Username", "000000", "15", "false"},
//	}

func (c Cache) LeaderboardTable(max int, request Request) []table.Row {
	users := getAllUsers()
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

	var rows []table.Row
	for i, user := range userPointsSlice {
		rows = append(rows, table.Row{
			fmt.Sprintf("%d", i+1),
			user.Username,
			user.Snowflake,
			fmt.Sprintf("%d", user.Points),
			fmt.Sprintf("%v", c.Private(user.Snowflake)),
		})
	}

	return rows
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
