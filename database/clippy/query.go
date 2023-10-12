package clippy

import (
	"cmp"
	"fmt"
	"github.com/nokusukun/bingo"
	"slices"
)

func (point Award) Record() {
	_, err := Collection.Insert(point, bingo.Upsert)
	if err != nil {
		panic(err)
	}
	cached.updateAwards(&point)
	//log.Println("Inserted", id)
}

func RecordMany(awards []Award) {
	id, err := Collection.InsertMany(awards, bingo.Upsert)
	if err != nil {
		panic(err)
	}
	for _, i := range id {
		fmt.Println("Inserted", i)
	}
}

func (user Config) Record() {
	_, err := UserSettings.Insert(user, bingo.Upsert)
	if err != nil {
		panic(err)
	}
	//fmt.Println("Inserted", id)
}

func QueryClippy(snowflake string) *bingo.QueryResult[Award] {
	return Collection.Query(bingo.Query[Award]{
		Filter: func(point Award) bool {
			return point.Snowflake == snowflake
		},
		// we cannot use keys because we're just storing the awards
		//Keys: [][]byte{
		//	[]byte(snowflake),
		//},
	})
}

func QueryConfig(snowflake string) *bingo.QueryResult[Config] {
	return UserSettings.Query(bingo.Query[Config]{
		//Filter: func(user Config) bool {
		//	return user.Snowflake == snowflake
		//},
		Keys: [][]byte{
			[]byte(snowflake),
		},
	})
}

func awardsSnowflake(snowflake string) []*Award {
	q := QueryClippy(snowflake)
	if q.Error != nil || !q.Any() {
		return nil
	}
	return q.Items
}

func countTotalPoints(snowflake string) int {
	return len(awardsSnowflake(snowflake))
}

func awardsSnowflakeGuild(snowflake string, guild string) (point []*Award) {
	result := Collection.Query(bingo.Query[Award]{
		Filter: func(point Award) bool {
			return point.Snowflake == snowflake && point.GuildID == guild
		},
	})
	return result.Items
}

func countTotalPointsGuild(snowflake string, guild string) int {
	return len(awardsSnowflakeGuild(snowflake, guild))
}

func GetOptedInConfigs() (users []*Config) {
	result := UserSettings.Query(bingo.Query[Config]{
		Filter: func(user Config) bool {
			return !user.OptOut
		},
	})
	for _, user := range result.Items {
		users = append(users, user)
	}
	return users
}

func getPointsShowConfig() (users []*Config) {
	result := UserSettings.Query(bingo.Query[Config]{
		Filter: func(user Config) bool {
			return !user.Private
		},
	})
	for _, user := range result.Items {
		users = append(users, user)
	}
	return users
}

func slicePointersToValues(users []*Config) []Config {
	var usersdereferenced []Config
	for _, user := range users {
		usersdereferenced = append(usersdereferenced, *user)
	}
	return usersdereferenced
}

type userPoints struct {
	Snowflake string
	Username  string
	Points    int
}

var userPointsSlice []userPoints

func sortUsersByPoints(users []*Config) {
	usersToSort := slicePointersToValues(users)
	userPointsSlice = []userPoints{}
	for _, user := range usersToSort {
		userPointsSlice = append(userPointsSlice, userPoints{
			Snowflake: user.Snowflake,
			Username:  user.Username,
			Points:    countTotalPoints(user.Snowflake),
		})
	}
}

func sortUserPoints(users []userPoints) []userPoints {
	slices.SortStableFunc(users, func(a, b userPoints) int {
		return cmp.Compare(b.Points, a.Points)
	})
	return users
}

func Leaderboard(max int, guild ...string) string {
	users := getPointsShowConfig()

	// get all awards depending on guild
	userPointsSlice = []userPoints{}
	if len(guild) > 0 {
		for _, u := range users {
			userPointsSlice = append(userPointsSlice, userPoints{
				Snowflake: u.Snowflake,
				Username:  u.Username,
				Points:    countTotalPointsGuild(u.Snowflake, guild[0]),
			})
		}
	} else {
		for _, u := range users {
			userPointsSlice = append(userPointsSlice, userPoints{
				Snowflake: u.Snowflake,
				Username:  u.Username,
				Points:    countTotalPoints(u.Snowflake),
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

	return leaderboard
}
