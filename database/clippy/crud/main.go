package main

import (
	"fmt"
	_ "go-clippy/database"
	"go-clippy/database/clippy"
	"log"
	"math/rand"
	"strconv"
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

func main() {
	//for i := 0; i < 500; i++ {
	//	// random chance
	//	switch rand.Intn(25) {
	//	case 0:
	//		testuser, testconfig := testUser(1)
	//		testuser.Record()
	//		testconfig.Record()
	//	case 1:
	//		testuser, testconfig := testUser(2)
	//		testuser.Record()
	//		testconfig.Record()
	//	case 2:
	//		testuser, testconfig := testUser(3)
	//		testuser.Record()
	//		testconfig.Record()
	//	}
	//}

	curTime := time.Now()
	fmt.Println(clippy.Leaderboard(5))
	fmt.Println("Non cached:", time.Since(curTime))

	curTime = time.Now()
	fmt.Println(clippy.Leaderboard(5))
	fmt.Println("Non cached:", time.Since(curTime))

	curTime = time.Now()
	fmt.Println(clippy.Leaderboard(5))
	fmt.Println("Non cached:", time.Since(curTime))

	curTime = time.Now()
	fmt.Println(clippy.LeaderboardCached(5))
	fmt.Println("Cached:", time.Since(curTime))
	curTime = time.Now()
	fmt.Println(clippy.LeaderboardCached(5))
	curTime = time.Now()
	fmt.Println("Cached:", time.Since(curTime))
	curTime = time.Now()
	fmt.Println(clippy.LeaderboardCached(5))
	fmt.Println("Cached:", time.Since(curTime))

	//clippy.CountTotalPoints("SNOWFLAKE1")
	//
	//clippy.CountTotalPointsCached("SNOWFLAKE1")

	timeTaken(func() {
		log.Println(clippy.CountTotalPoints("SNOWFLAKE1"))
	})

	timeTaken(func() {
		log.Println(clippy.CountTotalPointsCached("SNOWFLAKE1"))
	})
	timeTaken(func() {
		log.Println(clippy.CountTotalPoints("SNOWFLAKE1"))
	})

	timeTaken(func() {
		log.Println(clippy.CountTotalPointsCached("SNOWFLAKE1"))
	})
	timeTaken(func() {
		log.Println(clippy.CountTotalPoints("SNOWFLAKE1"))
	})

	timeTaken(func() {
		log.Println(clippy.CountTotalPointsCached("SNOWFLAKE1"))
	})

}

func timeTaken(f func()) {
	start := time.Now()
	f()
	log.Println("Took", time.Since(start))
}
