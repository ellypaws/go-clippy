package clippy

import (
	"cmp"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/nokusukun/bingo"
	logger "go-clippy/gui/log"
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
		//log.Println(err)
		program.Send(logger.Message(fmt.Sprintf("Error recording award: %v", err)))
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
		//log.Println("Error recording user: ", err)
		program.Send(logger.Message(fmt.Sprintf("Error recording user: %v", err)))
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
		program.Send(logger.Message(fmt.Sprintf("User %v not in cache", user.Snowflake)))
		var exist bool
		user, exist = c.GetConfig(user.Snowflake)
		if !exist {
			log.Printf("User %v does not exist", user.Snowflake)
			return
		}
	}
	//log.Println("current points: ", user.Points)
	program.Send(logger.Message(fmt.Sprintf("current points: %v", user.Points)))
	old := user.Points
	user.Points++
	//log.Println("new points: ", user.Points)
	program.Send(logger.Message(fmt.Sprintf("adding one point to user: @%v", user.Username)))

	c.synchronizePoints(user.Snowflake)
	if user.Points == old {
		//log.Printf("points mismatch, user.Points: %v, u.Config.Points: %v", user.Points, u.Config.Points)
		program.Send(logger.Message(fmt.Sprintf("points mismatch, still at %v", user.Points)))
		return
	}
	program.Send(logger.Message(fmt.Sprintf("c.Map[%v].Config.Points = %v", user.Snowflake, user.Points)))

	c.Mutex.Lock()
	c.Map[user.Snowflake].Config.Points = user.Points
	c.Mutex.Unlock()
	//log.Println("recording user: ", user)
	program.Send(logger.Message(fmt.Sprintf("recording user: %v", user)))
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
		c.Map[snowflake] = &CachedUser{}
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
		c.Mutex.Lock()
		defer c.Mutex.Unlock()
		c.Map[award.Snowflake] = user
	} else {
		c.Mutex.Lock()
		defer c.Mutex.Unlock()
		c.Map[award.Snowflake] = &CachedUser{
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
		c.Mutex.Lock()
		defer c.Mutex.Unlock()
		c.Map[config.Snowflake] = user
	} else {
		c.Mutex.Lock()
		defer c.Mutex.Unlock()
		c.Map[config.Snowflake] = &CachedUser{
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

func getAllUsers() (users []*User) {
	result := Users.Query(bingo.Query[User]{
		Filter: func(user User) bool {
			return true
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

var program *tea.Program

func StoreProgram(p *tea.Program) {
	p.Send(logger.Message("Storing program"))
	program = p
	program.Send(logger.Message("Program stored"))
}
