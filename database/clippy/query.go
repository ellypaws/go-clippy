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
// it calls [Cache.updateAwards] to update the Cache for the user's awards
// then [Cache.addPointRecord] will be called to increment the user's points
// it will also call [User.Record] to record the new points in the database,
// and then it will call [Cache.updateCachedConfig] to update the Cache
func (point Award) Record() {
	_, err := Awards.Insert(point, bingo.Upsert)
	if err != nil {
		log.Println(err)
	}
	user, ok := GetCache().GetConfig(point.Snowflake)
	if !ok {
		return
	}
	GetCache().updateAwards(&point)
	GetCache().addPointRecord(user)
}

// Record records the user's config to the database
// it first checks if the user is in the cache for the user points
// then it will update the cache as well
func (config User) Record() {
	// We shouldn't grab from cache again since we're updating with new values
	//user, exist := GetCache().GetConfig(config.Snowflake)
	//if exist {
	//	config.Points = user.Points
	//}
	_, err := Users.Insert(config, bingo.Upsert)
	if err != nil {
		log.Println("Error recording user: ", err)
	}
	GetCache().updateCachedConfig(config)
}

func (moderator Moderator) Record() {
	_, err := Moderators.Insert(moderator, bingo.Upsert)
	if err != nil {
		panic(err)
	}
}

// addPointRecord adds a point to the user's record
// it first checks if the user is in the Cache for the user points
// then [user.Record] will update the Cache as well
func (c Cache) addPointRecord(user *User) {
	c.Mutex.RLock()
	_, ok := c.Map[user.Snowflake]
	c.Mutex.RUnlock()
	if !ok {
		var exist bool
		user, exist = c.GetConfig(user.Snowflake)
		if !exist {
			log.Printf("User %v does not exist", user.Snowflake)
			return
		}
	}
	log.Println("current points: ", user.Points)
	user.Points++
	log.Println("new points: ", user.Points)
	c.Mutex.Lock()
	c.Map[user.Snowflake].Config.Points = user.Points
	c.Mutex.Unlock()
	log.Println("recording user: ", user)
	user.Record()
}

func (c Cache) allAwards(request Request) *bingo.QueryResult[Award] {
	result := Awards.Query(bingo.Query[Award]{
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

func (c Cache) queryConfig(snowflake string) *bingo.QueryResult[User] {
	result := Users.Query(bingo.Query[User]{
		//Filter: func(user User) bool {
		//	return user.Snowflake == snowflake
		//},
		Keys: [][]byte{
			[]byte(snowflake),
		},
	})
	if !result.Any() {
		c.Mutex.Lock()
		c.Map[snowflake] = &CacheType{}
		c.Mutex.Unlock()
	}
	return result
}

// updateAwards updates the Cache for the user's awards
// if a user is not in the Cache, it will create a new entry and store the award
func (c Cache) updateAwards(award *Award) {
	c.Mutex.RLock()
	user, ok := c.Map[award.Snowflake]
	c.Mutex.RUnlock()

	if ok {
		user.Awards = append(user.Awards, award)
	} else {
		c.Mutex.Lock()
		defer c.Mutex.Unlock()
		c.Map[award.Snowflake] = &CacheType{
			Awards: []*Award{award},
		}
	}
}

// Cache.updateCachedConfig updates the cached config for a user
// if it's not in the Cache, it will store the incoming config
func (c Cache) updateCachedConfig(config User) {
	c.Mutex.RLock()
	user, ok := c.Map[config.Snowflake]
	c.Mutex.RUnlock()

	if ok {
		user.Config = config
	} else {
		c.Mutex.Lock()
		defer c.Mutex.Unlock()
		c.Map[config.Snowflake] = &CacheType{
			Config: config,
		}
	}
}

// GetConfig returns the cached config for a user
// if it's not in the cache, it will query the database
func (c Cache) GetConfig(snowflake string) (user *User, exist bool) {
	c.Mutex.RLock()
	_, ok := c.Map[snowflake]
	c.Mutex.RUnlock()

	if !ok {
		config := c.queryConfig(snowflake)
		if !config.Any() || config.Error != nil {
			return &User{}, false
		}
		c.updateCachedConfig(*config.First())
	}
	c.Mutex.RLock()
	defer c.Mutex.RUnlock()
	return &c.Map[snowflake].Config, true
}

func getPublicUsers() (users []*User) {
	result := Users.Query(bingo.Query[User]{
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
