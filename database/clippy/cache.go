package clippy

import (
	"fmt"
	"github.com/nokusukun/bingo"
)

type Cache struct {
	Config Config
	Awards []*Award
}

type cacheMap map[string]*Cache

var CacheMap = make(cacheMap)

func NewCache() cacheMap {
	return make(cacheMap)
}

func (c cacheMap) Reset() {
	for k := range c {
		delete(c, k)
	}
}

func (c cacheMap) queryDB(snowflake string) {
	c[snowflake] = &Cache{
		Config: *QueryConfig(snowflake).First(),
		Awards: QueryClippy(snowflake).Items,
	}
}

func (c cacheMap) updateAwards(award *Award) {
	if user, ok := c[award.Snowflake]; ok {
		user.Awards = append(user.Awards, award)
	} else {
		c[award.Snowflake] = &Cache{
			Awards: []*Award{award},
		}
	}
}

func (c cacheMap) updateConfig(config Config) {
	if user, ok := c[config.Snowflake]; ok {
		user.Config = config
	} else {
		c[config.Snowflake] = &Cache{
			Config: config,
		}
	}
}

func (c cacheMap) getAward(snowflake string) ([]*Award, error) {
	if user, ok := c[snowflake]; ok && user.Awards != nil {
		return user.Awards, nil
	}

	q := QueryClippy(snowflake)

	if !q.Any() {
		return nil, fmt.Errorf("no awards found for %s", snowflake)
	}

	c[snowflake] = &Cache{
		Awards: q.Items,
	}

	return q.Items, q.Error
}

func (c cacheMap) awardsSnowflakeGuildCached(snowflake string, guild string) []*Award {
	if user, ok := c[snowflake]; ok && user.Awards != nil {
		return user.Awards
	}

	q := awardsSnowflakeGuild(snowflake, guild)

	if len(q) == 0 {
		return nil
	}

	c[snowflake] = &Cache{
		Awards: q,
	}

	return q
}

func (c cacheMap) getConfig(snowflake string) (*Config, error) {
	if user, ok := c[snowflake]; ok {
		return &user.Config, nil
	}

	q := QueryConfig(snowflake)

	if !q.Any() {
		return nil, fmt.Errorf("no config found for %s", snowflake)
	}

	c[snowflake] = &Cache{
		Config: *q.First(),
	}

	return q.First(), q.Error
}

func CountTotalPointsCached(snowflake string) int {
	award, err := CacheMap.getAward(snowflake)
	if err != nil {
		return 0
	}
	return len(award)
}

func CountTotalPointsGuildCached(snowflake string, guild string) int {
	award := CacheMap.awardsSnowflakeGuildCached(snowflake, guild)
	return len(award)
}

func CountTotalPointsGuildPrecached(snowflake string, guild ...string) int {
	if _, ok := CacheMap[snowflake]; !ok {
		PrecacheAwards(CacheMap, guild...)
	}
	award := CacheMap[snowflake].Awards
	return len(award)
}

func PrecacheAwards(cache cacheMap, guild ...string) {
	users := getPointsShowConfig()
	for _, user := range users {
		cache[user.Snowflake] = &Cache{
			Config: *user,
		}
	}
	var result *bingo.QueryResult[Award]

	result = Collection.Query(bingo.Query[Award]{
		Filter: func(award Award) bool {
			if len(guild) > 0 {
				return award.GuildID == guild[0]
			} else {
				return true
			}
		},
	})

	for _, award := range result.Items {
		if user, ok := cache[award.Snowflake]; ok {
			user.Awards = append(user.Awards, award)
		} else {
			cache[award.Snowflake] = &Cache{
				Awards: []*Award{award},
			}
		}
	}
}

func LeaderboardCached(max int, guild ...string) string {
	users := getPointsShowConfig()

	// get all awards depending on guild
	userPointsSlice = []userPoints{}
	if len(guild) > 0 {
		for _, u := range users {
			userPointsSlice = append(userPointsSlice, userPoints{
				Snowflake: u.Snowflake,
				Username:  u.Username,
				Points:    CountTotalPointsGuildCached(u.Snowflake, guild[0]),
			})
		}
	} else {
		for _, u := range users {
			userPointsSlice = append(userPointsSlice, userPoints{
				Snowflake: u.Snowflake,
				Username:  u.Username,
				Points:    CountTotalPointsCached(u.Snowflake),
			})
		}
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

func LeaderboardPrecached(max int, guild ...string) string {
	users := getPointsShowConfig()
	if len(CacheMap) == 0 {
		PrecacheAwards(CacheMap, guild...)
	}

	//log.Printf("Precached %v users", precache)

	//for k, v := range precache {
	//	log.Printf("%v: %v", k, v)
	//}

	// get all awards depending on guild
	userPointsSlice = []userPoints{}
	for _, u := range users {
		userPointsSlice = append(userPointsSlice, userPoints{
			Snowflake: u.Snowflake,
			Username:  u.Username,
			Points:    len(CacheMap[u.Snowflake].Awards),
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
