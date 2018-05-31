package database

import (
	"github.com/globalsign/mgo"
	log "github.com/sirupsen/logrus"
)

var (
	mgoSession *mgo.Session
)

// GetSession either creates or returns a clone of a single session
// This can then be used for operations on the database
func GetSession(dbUrl string) (*mgo.Session, error) {
	if mgoSession == nil {
		log.WithFields(log.Fields{
			"URL": dbUrl,
		}).Debug("Creating mgoSession")

		dialInfo, err := mgo.ParseURL(dbUrl)
		dialInfo.Direct = true
		dialInfo.FailFast = true

		mgoSession, err = mgo.DialWithInfo(dialInfo)
		if err != nil {
			log.WithFields(log.Fields{
				"error": err.Error(),
			}).Error("Failed to connect to DB")
			return nil, err
		}

	}
	return mgoSession.Clone(), nil
}
