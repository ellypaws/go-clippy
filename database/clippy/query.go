package clippy

import (
	"cmp"
	"github.com/nokusukun/bingo"
	"slices"
)

type Request struct {
	Snowflake string
	GuildID   string
}

func (point Award) Record() {
	_, err := Collection.Insert(point, bingo.Upsert)
	if err != nil {
		panic(err)
	}
	Cache.updateAwards(&point)
	Cache.addPointRecord(Cache.GetConfig(point.Snowflake))
	//log.Println("Inserted", id)
}

func (config User) Record() {
	_, err := UserSettings.Insert(config, bingo.Upsert)
	if err != nil {
		panic(err)
	}
	Cache.updateConfig(config)
	//fmt.Println("Inserted", id)
}

func (c cacheMap) addPointRecord(user User) {
	user.Points++
	c[user.Snowflake].Config.Points = user.Points
	user.Record()
}

func (c cacheMap) allAwards(request Request) *bingo.QueryResult[Award] {
	result := Collection.Query(bingo.Query[Award]{
		Filter: func(point Award) bool {
			snowflakeMatch := request.Snowflake == "" || point.Snowflake == request.Snowflake
			guildIDMatch := request.GuildID == "" || point.GuildID == request.GuildID
			return snowflakeMatch && guildIDMatch
		},
		// we cannot use keys because we're just storing the awards
		//Keys: [][]byte{
		//	[]byte(snowflake),
		//},
	})
	for _, award := range result.Items {
		c.updateAwards(award)
	}
	return result
}

func (c cacheMap) queryConfig(snowflake string) *bingo.QueryResult[User] {
	result := UserSettings.Query(bingo.Query[User]{
		//Filter: func(user User) bool {
		//	return user.Snowflake == snowflake
		//},
		Keys: [][]byte{
			[]byte(snowflake),
		},
	})
	if !result.Any() {
		c[snowflake] = &CacheType{}
	}
	return result
}

func (c cacheMap) updateAwards(award *Award) {
	if user, ok := c[award.Snowflake]; ok {
		user.Awards = append(user.Awards, award)
	} else {
		c[award.Snowflake] = &CacheType{
			Awards: []*Award{award},
		}
	}
}

func (c cacheMap) updateConfig(config User) {
	if user, ok := c[config.Snowflake]; ok {
		user.Config = config
	} else {
		c[config.Snowflake] = &CacheType{
			Config: config,
		}
	}
}

func (c cacheMap) GetConfig(snowflake string) User {
	if _, ok := c[snowflake]; !ok {
		config := c.queryConfig(snowflake)
		if !config.Any() {
			return User{}
		}
		c.updateConfig(*config.First())
	}
	return c[snowflake].Config
}

func getPublicUsers() (users []*User) {
	result := UserSettings.Query(bingo.Query[User]{
		Filter: func(user User) bool {
			return !user.Private
		},
	})
	for _, user := range result.Items {
		users = append(users, user)
	}
	return users
}

type userPoints struct {
	Snowflake string
	Username  string
	Points    int
}

var userPointsSlice []userPoints

func sortUserPoints(users []userPoints) []userPoints {
	slices.SortStableFunc(users, func(a, b userPoints) int {
		return cmp.Compare(b.Points, a.Points)
	})
	return users
}
