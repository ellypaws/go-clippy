package main

import (
	"fmt"
	_ "go-clippy/database"
	"go-clippy/database/clippy"
	"math/rand"
	"strconv"
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
	for i := 0; i < 500; i++ {
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

	fmt.Println(clippy.Leaderboard(5))

}
