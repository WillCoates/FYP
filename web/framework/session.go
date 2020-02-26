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
	client redis.Client
}

func (sess *Session) Save() error {
	client := sess.manager.client
	var toDelete = make([]string, 0)
	var toSet = make(map[string]interface{})

	for key, value := range sess.Values {
		if value == "" {
			toDelete = append(toDelete, key)
		} else {
			toSet[key] = value
		}
	}

	if len(toDelete) > 0 {
		status := client.HDel(sess.id, toDelete...)
		err := status.Err()
		if err != nil {
			return err
		}
	}
	status := client.HMSet(sess.id, toSet)
	err := status.Err()
	if err != nil {
		return err
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
	session.id = idBase64
	session.manager = mgr

	return session, nil
}
