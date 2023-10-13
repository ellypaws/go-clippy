package main

import (
	"fmt"
	_ "go-clippy/database"
	"go-clippy/database/clippy"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

func testUser(i int) (clippy.Award, clippy.User) {
	return clippy.Award{
			Username:        "TestUser" + strconv.Itoa(i),
			Snowflake:       "SNOWFLAKE" + strconv.Itoa(i),
			GuildName:       "Discord Server",
			GuildID:         "randomGuildID",
			Channel:         "test-forum",
			ChannelID:       "channelID",
			MessageID:       "messageID",
			OriginUsername:  "TestUser1" + strconv.Itoa(i),
			OriginSnowflake: "SNOWFLAKE1" + strconv.Itoa(i),
			InteractionID:   strconv.Itoa(rand.Int()),
		},
		clippy.User{
			Username:  "TestUser" + strconv.Itoa(i),
			Snowflake: "SNOWFLAKE" + strconv.Itoa(i),
			OptOut:    false,
			Private:   false,
		}
}

var out = strings.Builder{}

func recordAwards(runs int) {
	for i := 0; i < runs; i++ {
		// random chance
		switch rand.Intn(25) {
		case 0:
			testuser, testconfig := testUser(1)
			testuser.Record()
			testconfig.Record()
		case 1:
			testuser, testconfig := testUser(2)
			testuser.Record()
			testconfig.Record()
		case 2:
			testuser, testconfig := testUser(3)
			testuser.Record()
			testconfig.Record()
		}
	}
}

func precachedLeaderboard(runs int) {
	timeTaken("PRECACHED (first run)", func() {
		fmt.Println(clippy.Cache.LeaderboardPrecached(5, clippy.Request{}))
	})

	for i := 0; i < runs; i++ {
		timeTaken(fmt.Sprintf("PRECACHED (%v)", i), func() {
			fmt.Println(clippy.Cache.LeaderboardPrecached(5, clippy.Request{}))
		})
	}
}

func singleUserPrecached(runs int) {
	for i := 0; i < runs; i++ {
		timeTaken(fmt.Sprintf("SINGLE USER PRECACHED (%v)", i), func() {
			fmt.Println(clippy.Cache.QueryPoints(clippy.Request{Snowflake: "SNOWFLAKE1"}))
		})
	}
}

func main() {
	clippy.Cache.Reset()
	precachedLeaderboard(3)

	//clippy.Cache.Reset()
	singleUserPrecached(3)

	fmt.Println(out.String())

}

func timeTaken(text string, f func()) {
	start := time.Now()
	f()
	toPrint := fmt.Sprintln(text, "TOOK:", time.Since(start))
	log.Println(toPrint)
	out.WriteString(toPrint)
}
