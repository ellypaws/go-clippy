package clippy

import "fmt"

type Cache struct {
	Config Config
	Awards []*Award
}

type cacheMap map[string]*Cache

var cached = make(cacheMap)

func (c cacheMap) queryDB(snowflake string) {
	c[snowflake] = &Cache{
		Config: *QueryConfig(snowflake).First(),
		Awards: QueryClippy(snowflake).Items,
	}
}

func (c cacheMap) updateAwards(award *Award) {
	if user, ok := c[award.Snowflake]; ok {
		user.Awards = append(user.Awards, award)
	} else {
		c[award.Snowflake] = &Cache{
			Awards: []*Award{award},
		}
	}
}

func (c cacheMap) updateConfig(config Config) {
	if user, ok := c[config.Snowflake]; ok {
		user.Config = config
	} else {
		c[config.Snowflake] = &Cache{
			Config: config,
		}
	}
}

func (c cacheMap) getAward(snowflake string) ([]*Award, error) {
	if user, ok := c[snowflake]; ok && user.Awards != nil {
		return user.Awards, nil
	}

	q := QueryClippy(snowflake)

	if !q.Any() {
		return nil, fmt.Errorf("no awards found for %s", snowflake)
	}

	c[snowflake] = &Cache{
		Awards: q.Items,
	}

	return q.Items, q.Error
}

func (c cacheMap) getConfig(snowflake string) (*Config, error) {
	if user, ok := c[snowflake]; ok {
		return &user.Config, nil
	}

	q := QueryConfig(snowflake)

	if !q.Any() {
		return nil, fmt.Errorf("no config found for %s", snowflake)
	}

	c[snowflake] = &Cache{
		Config: *q.First(),
	}

	return q.First(), q.Error
}
