package framework

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/binary"
	"time"

	"github.com/go-redis/redis"
)

type Session struct {
	manager *SessionManager
	id      string
	Values  map[string]string
}

type SessionManager struct {
	client *redis.Client
}

func NewSessionManager(redisURL string) (*SessionManager, error) {
	url, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, err
	}

	manager := new(SessionManager)
	manager.client = redis.NewClient(url)

	return manager, nil
}

func (sess *Session) Save() error {
	client := sess.manager.client

	for key, value := range sess.Values {
		if value == "" {
			status := client.HDel(sess.id, key)
			err := status.Err()
			if err != nil {
				return err
			}
		} else {
			status := client.HSet(sess.id, key, value)
			err := status.Err()
			if err != nil {
				return err
			}
		}
	}

	client.Expire(sess.id, 24*time.Hour)

	return nil
}

func (sess *Session) Id() string {
	return sess.id
}

func (mgr *SessionManager) Load(id string) (*Session, error) {
	client := mgr.client
	mapping := client.HGetAll(id)
	err := mapping.Err()
	if err != nil {
		return nil, err
	}

	session := new(Session)
	session.Values = mapping.Val()
	session.id = id
	session.manager = mgr

	return session, nil
}

func (mgr *SessionManager) Create() (*Session, error) {
	id := make([]byte, 16)
	now := time.Now().UnixNano()
	binary.LittleEndian.PutUint64(id, uint64(now))
	_, err := rand.Reader.Read(id[8:])
	if err != nil {
		return nil, err
	}
	idBase64 := base64.RawURLEncoding.EncodeToString(id)
	session := new(Session)
	session.Values = make(map[string]string)
	session.id = idBase64
	session.manager = mgr

	return session, nil
}
