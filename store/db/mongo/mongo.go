package mongo

import (
	"errors"
	"log"
	"sync"
	"time"

	"github.com/Leon2012/goconfd/store"
	"github.com/Leon2012/goconfd/store/types"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func init() {
	store.Register("mongo", NewMongoAdapter())
}

type MongoAdapter struct {
	Conf    MongoConfig
	Session *mgo.Session
	isOpen  bool
	m       sync.Mutex
}

type MongoConfig struct {
	Url      string
	DbName   string
	Timeout  int
	Username string
	Password string
}

func NewMongoAdapter() *MongoAdapter {
	return &MongoAdapter{}
}

func (a *MongoAdapter) Open(c interface{}) error {
	config, ok := c.(MongoConfig)
	if !ok {
		return errors.New("Invaid Config")
	}
	timeout := time.Duration(config.Timeout) * time.Second
	session, err := mgo.DialWithTimeout(config.Url, timeout)
	if err != nil {
		return err
	}
	session.SetMode(mgo.Monotonic, true)

	a.Session = session
	a.Conf = config
	a.isOpen = true
	return nil
}

func (a *MongoAdapter) Close() error {
	a.Session.Close()
	a.isOpen = false
	return nil
}

func (a *MongoAdapter) IsOpen() bool {
	return a.isOpen
}

func (a *MongoAdapter) Online(agent *types.Agent) error {
	var err error
	now := time.Now()
	selector := types.Agent{}
	agent.Status = 1
	a.run("agents", func(c *mgo.Collection) {
		err = c.Find(bson.M{"hostname": agent.HostName, "keyprefix": agent.KeyPrefix}).One(&selector)
	})
	if err != nil { //未查到记录
		agent.CreatedAt = now
		a.run("agents", func(c *mgo.Collection) {
			err = c.Insert(agent)
		})
	} else {
		data := bson.M{"$set": bson.M{"status": 1, "updatedat": now}}
		a.run("agents", func(c *mgo.Collection) {
			err = c.Update(selector, data)
		})
	}
	return err
}

func (a *MongoAdapter) Offline(agent *types.Agent) error {
	var err error
	now := time.Now()
	selector := bson.M{"hostname": agent.HostName, "keyprefix": agent.KeyPrefix}
	data := bson.M{"$set": bson.M{"status": 0, "updatedat": now}}
	a.run("agents", func(c *mgo.Collection) {
		err = c.Update(selector, data)
	})
	return err
}

func (a *MongoAdapter) Heartbeat(log *types.Heartbeat) error {
	var err error
	now := time.Now()
	selector := bson.M{"hostname": log.HostName, "keyprefix": log.KeyPrefix}
	data := bson.M{"$set": bson.M{"heartbeattime": now}}
	a.run("agents", func(c *mgo.Collection) {
		err = c.Update(selector, data)
	})
	if err != nil {
		return err
	}
	log.CreatedAt = now
	a.run("heartbeats", func(c *mgo.Collection) {
		err = c.Insert(log)
	})
	return err
}

func (a *MongoAdapter) GetAgents() ([]*types.Agent, error) {
	var results []*types.Agent
	var err error
	a.run("agents", func(c *mgo.Collection) {
		err = c.Find(nil).All(&results)
	})
	if err != nil {
		return nil, err
	} else {
		return results, nil
	}
}

func (a *MongoAdapter) GetHeartbeats() ([]*types.Heartbeat, error) {
	var results []*types.Heartbeat
	var err error
	a.run("heartbeats", func(c *mgo.Collection) {
		err = c.Find(nil).Sort("-createdat").Limit(50).All(&results)
	})
	if err != nil {
		return nil, err
	} else {
		return results, nil
	}
}

func (a *MongoAdapter) GetHeartbeatsByAgent(agent *types.Agent) ([]*types.Heartbeat, error) {
	var results []*types.Heartbeat
	var err error
	selector := bson.M{"hostname": agent.HostName, "keyprefix": agent.KeyPrefix}
	a.run("heartbeats", func(c *mgo.Collection) {
		err = c.Find(selector).Sort("-createdat").Limit(50).All(&results)
	})
	if err != nil {
		return nil, err
	} else {
		return results, nil
	}
}

func (a *MongoAdapter) GetHeartbeatsByTime(time time.Time) ([]*types.Heartbeat, error) {
	var results []*types.Heartbeat
	var err error
	selector := bson.M{"CreatedAt": time}
	a.run("heartbeats", func(c *mgo.Collection) {
		err = c.Find(selector).All(&results)
	})
	if err != nil {
		return nil, err
	} else {
		return results, nil
	}
}

func (a *MongoAdapter) clone() *mgo.Session {
	a.m.Lock()
	defer a.m.Unlock()
	newSession := a.Session.Clone()
	newSession.Refresh()
	return newSession
}

func (a *MongoAdapter) run(collection string, f func(*mgo.Collection)) {
	session := a.clone()
	defer func() {
		session.Clone()
		if err := recover(); err != nil {
			log.Fatal(err)
		}
	}()

	db := session.DB(a.Conf.DbName)
	if a.Conf.Username != "" {
		err := db.Login(a.Conf.Username, a.Conf.Password)
		if err != nil {
			log.Fatal(err)
		}
	}
	c := db.C(collection)
	f(c)
}

func (a *MongoAdapter) search(collection string, q interface{}, skip, limit int, sorts []string) (searchResults []interface{}, err error) {
	var query *mgo.Query
	a.run(collection, func(c *mgo.Collection) {
		if q != nil {
			query = c.Find(q)
		} else {
			query = c.Find(nil)
		}
		if skip > 0 {
			query = query.Skip(skip)
		}
		if limit > 0 {
			query = query.Limit(limit)
		}
		if sorts != nil && len(sorts) > 0 {
			query = query.Sort(sorts...)
		}
		err = query.All(&searchResults)
	})
	return
}
