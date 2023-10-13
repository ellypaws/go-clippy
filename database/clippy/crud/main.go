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

func testUser(i int) (clippy.Award, clippy.Config) {
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
		clippy.Config{
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

func nonCached(runs int) {
	for i := 0; i < runs; i++ {
		timeTaken(fmt.Sprintf("NON CACHED (%v)", i), func() {
			fmt.Println(clippy.Leaderboard(5))
		})
	}
}

func cached(runs int) {
	timeTaken("CACHED (first run)", func() {
		fmt.Println(clippy.LeaderboardCached(5))
	})

	for i := 0; i < runs; i++ {
		timeTaken(fmt.Sprintf("CACHED (%v)", i), func() {
			fmt.Println(clippy.LeaderboardCached(5))
		})
	}
}

func precached(runs int) {
	timeTaken("PRECACHED (first run)", func() {
		fmt.Println(clippy.LeaderboardPrecached(5))
	})

	for i := 0; i < runs; i++ {
		timeTaken(fmt.Sprintf("PRECACHED (%v)", i), func() {
			fmt.Println(clippy.LeaderboardPrecached(5))
		})
	}
}

func singleUser(runs int) {
	for i := 0; i < runs; i++ {
		timeTaken(fmt.Sprintf("SINGLE USER (%v)", i), func() {
			fmt.Println(clippy.CountTotalPoints("SNOWFLAKE1"))
		})
	}
}

func singleUserCached(runs int) {
	for i := 0; i < runs; i++ {
		timeTaken(fmt.Sprintf("SINGLE USER CACHED (%v)", i), func() {
			fmt.Println(clippy.CountTotalPointsCached("SNOWFLAKE1"))
		})
	}
}

func singleUserCached2(runs int) {
	for i := 0; i < runs; i++ {
		timeTaken(fmt.Sprintf("SINGLE USER CACHED (%v)", i), func() {
			fmt.Println(clippy.CountTotalPointsCached("SNOWFLAKE2"))
		})
	}
}

func singleUserPrecached(runs int) {
	for i := 0; i < runs; i++ {
		timeTaken(fmt.Sprintf("SINGLE USER PRECACHED (%v)", i), func() {
			fmt.Println(clippy.CountTotalPointsGuildPrecached("SNOWFLAKE1"))
		})
	}
}

func singleUserPrecached2(runs int) {
	for i := 0; i < runs; i++ {
		timeTaken(fmt.Sprintf("SINGLE USER PRECACHED (%v)", i), func() {
			fmt.Println(clippy.CountTotalPointsGuildPrecached("SNOWFLAKE2"))
		})
	}
}

func main() {
	nonCached(1)
	clippy.CacheMap.Reset()
	cached(3)
	clippy.CacheMap.Reset()
	precached(3)

	singleUser(3)
	clippy.CacheMap.Reset()
	singleUserCached(3)
	singleUserCached2(3)
	clippy.CacheMap.Reset()
	singleUserPrecached(3)
	singleUserPrecached2(3)

	fmt.Println(out.String())

}

func timeTaken(text string, f func()) {
	start := time.Now()
	f()
	toPrint := fmt.Sprintln(text, "TOOK:", time.Since(start))
	log.Println(toPrint)
	out.WriteString(toPrint)
}
