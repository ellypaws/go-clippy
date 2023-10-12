package clippy

import (
	"cmp"
	"fmt"
	"github.com/nokusukun/bingo"
	"slices"
)

func (user Clippy) Record() {
	id, err := Collection.Insert(user, bingo.Upsert)
	if err != nil {
		panic(err)
	}
	fmt.Println("Inserted", id)
}

func RecordMany(c []Clippy) {
	id, err := Collection.InsertMany(c, bingo.Upsert)
	if err != nil {
		panic(err)
	}
	for _, i := range id {
		fmt.Println("Inserted", i)
	}
}

func GetUser(s string) (*Clippy, error) {
	q := QueryClippy(s)
	if q.Error != nil {
		return nil, q.Error
	}
	return q.First(), nil
}

func QueryClippy(s string) *bingo.QueryResult[Clippy] {
	return Collection.Query(bingo.Query[Clippy]{
		Filter: func(doc Clippy) bool {
			return doc.Snowflake == s
		},
	})
}

func (user Clippy) countTotalPoints(collection *bingo.Collection[Clippy]) int {
	var count int
	for _, point := range user.Points {
		count += len(point.Awards)
	}
	return count
}

func PointsOfSnowflake(s string) int {
	user, err := GetUser(s)
	if err != nil || user == nil {
		return 0
	}
	return user.countTotalPoints(Collection)
}

func PointsOfSnowflakeWithGuild(snowflake string, guild string) int {
	user, err := GetUser(snowflake)
	if err != nil || user == nil {
		return 0
	}
	for _, point := range user.Points {
		if point.GuildID == guild {
			return len(point.Awards)
		}
	}
	return 0
}

func GetUsersWithAwardsFromGuild(guild string) (users []*Clippy) {
	res := Collection.Query(bingo.Query[Clippy]{
		Filter: func(user Clippy) bool {
			for _, point := range user.Points {
				if point.GuildID == guild && len(point.Awards) > 0 {
					return true
				}
			}
			return false
		},
	})
	return res.Items
}

func GetUsersWithAwards() (users []*Clippy) {
	res := Collection.Query(bingo.Query[Clippy]{
		Filter: func(doc Clippy) bool {
			return len(doc.Points) > 0
		},
	})
	return res.Items
}

func slicePointersToValues(users []*Clippy) []Clippy {
	var usersdereferenced []Clippy
	for _, user := range users {
		usersdereferenced = append(usersdereferenced, *user)
	}
	return usersdereferenced
}

func sortUserByPoints(users []*Clippy) []Clippy {
	usersToSort := slicePointersToValues(users)
	slices.SortStableFunc(usersToSort, func(a, b Clippy) int {
		return cmp.Compare(a.countTotalPoints(Collection), b.countTotalPoints(Collection))
	})
	return usersToSort
}

func Leaderboard(max int, guild ...string) string {
	var users []*Clippy
	if len(guild) > 0 {
		for _, g := range guild {
			users = append(users, GetUsersWithAwardsFromGuild(g)...)
		}
	} else {
		users = GetUsersWithAwards()
	}
	if len(users) == 0 {
		return "No users found"
	}
	sorted := sortUserByPoints(users)
	if max == 0 {
		max = len(sorted)
	}
	sorted = sorted[:min(max, len(sorted))-1]

	var leaderboard string
	for i, user := range sorted {
		leaderboard += fmt.Sprintf("%d. %s (%d)\n", i+1, user.Username, user.countTotalPoints(Collection))
	}
	return leaderboard
}
