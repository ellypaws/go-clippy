package clippy

import (
	"cmp"
	"github.com/nokusukun/bingo"
	"log"
	"slices"
)

type Request struct {
	Snowflake string
	GuildID   string
}

// Record records the user's point to the database
// it first checks if the user is in the cache for the user points
// it calls [cacheMap.updateAwards] to update the cache for the user's awards
// then [cacheMap.addPointRecord] will be called to increment the user's points
// it will also call [User.Record] to record the new points in the database,
// and then it will call [cacheMap.updateCachedConfig] to update the cache
func (point Award) Record() {
	_, err := Collection.Insert(point, bingo.Upsert)
	if err != nil {
		log.Println(err)
	}
	user, ok := Cache.GetConfig(point.Snowflake)
	if !ok {
		return
	}
	Cache.updateAwards(&point)
	Cache.addPointRecord(user)
}

// Record records the user's config to the database
// it first checks if the user is in the cache for the user points
// then it will update the cache as well
func (config User) Record() {
	user, exist := Cache.GetConfig(config.Snowflake)
	if exist {
		config.Points = user.Points
	}
	_, err := UserSettings.Insert(config, bingo.Upsert)
	if err != nil {
		log.Println("Error recording user: ", err)
	}
	Cache.updateCachedConfig(config)
}

func (moderator Moderator) Record() {
	_, err := Moderators.Insert(moderator, bingo.Upsert)
	if err != nil {
		panic(err)
	}
}

// addPointRecord adds a point to the user's record
// it first checks if the user is in the cache for the user points
// then [user.Record] will update the cache as well
func (c cacheMap) addPointRecord(user User) {
	if config := c[user.Snowflake]; config == nil {
		var exist bool
		user, exist = c.GetConfig(user.Snowflake)
		if !exist {
			log.Printf("User %v does not exist", user.Snowflake)
			return
		}
	}
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

// updateAwards updates the cache for the user's awards
// if a user is not in the cache, it will create a new entry and store the award
func (c cacheMap) updateAwards(award *Award) {
	if user, ok := c[award.Snowflake]; ok {
		user.Awards = append(user.Awards, award)
	} else {
		c[award.Snowflake] = &CacheType{
			Awards: []*Award{award},
		}
	}
}

// cacheMap.updateCachedConfig updates the cached config for a user
// if it's not in the cache, it will store the incoming config
func (c cacheMap) updateCachedConfig(config User) {
	if user, ok := c[config.Snowflake]; ok {
		user.Config = config
	} else {
		c[config.Snowflake] = &CacheType{
			Config: config,
		}
	}
}

// GetConfig returns the cached config for a user
// if it's not in the cache, it will query the database
func (c cacheMap) GetConfig(snowflake string) (user User, exist bool) {
	if _, ok := c[snowflake]; !ok {
		config := c.queryConfig(snowflake)
		if !config.Any() || config.Error != nil {
			return User{}, false
		}
		c.updateCachedConfig(*config.First())
	}
	return c[snowflake].Config, true
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
