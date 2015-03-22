// Copyright 2013 Ardan Studios. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package mongo provides mongo connectivity support.
package mongo

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/vincent3i/beego-blog/g"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"os"
	"strings"
	"time"
)

const (
	// MasterSession provides direct access to master database.
	MasterSession = "master"

	// MonotonicSession provides reads to slaves.
	MonotonicSession = "monotonic"
)

var (
	// Reference to the singleton.
	singleton mongoManager
)

type (
	// mongoConfiguration contains settings for initialization.
	mongoConfiguration struct {
		Hosts    string
		Database string
		UserName string
		Password string
	}

	// mongoManager contains dial and session information.
	mongoSession struct {
		mongoDBDialInfo *mgo.DialInfo
		mongoSession    *mgo.Session
	}

	// mongoManager manages a map of session.
	mongoManager struct {
		sessions map[string]mongoSession
	}

	// DBCall defines a type of function that can be used
	// to excecute code against MongoDB.
	DBCall func(*mgo.Collection) error
)

// Startup brings the manager to a running state.
func Startup() error {
	// If the system has already been started ignore the call.
	if singleton.sessions != nil {
		return nil
	}

	beego.Debug("Mongo session startup...")

	// Pull in the configuration.
	config := mongoConfiguration{
		g.Cfg.String("mongo_addresses"),
		g.Cfg.String("mongo_database"),
		g.Cfg.String("mongo_username"),
		g.Cfg.String("mongo_password"),
	}

	// Create the Mongo Manager.
	singleton = mongoManager{
		sessions: make(map[string]mongoSession),
	}

	// Log the mongodb connection straps.
	beego.BeeLogger.Debug("Startup MongoDB : Hosts[%s], Database[%s], Username[%s]", config.Hosts, config.Database, config.UserName)

	hosts := strings.Split(config.Hosts, ",")

	// Create the strong session.
	if err := CreateSession("strong", MasterSession, hosts, config.Database, config.UserName, config.Password); err != nil {
		beego.BeeLogger.Error("Create strong session error, %s", err.Error())
		return err
	}

	// Create the monotonic session.
	if err := CreateSession("monotonic", MonotonicSession, hosts, config.Database, config.UserName, config.Password); err != nil {
		beego.BeeLogger.Error("Create monotonic session error, %s", err.Error())
		return err
	}

	mgo.SetDebug(false)
	mgo.SetLogger(log.New(os.Stdout, "", log.LstdFlags))

	beego.Debug("Mongo session startup completed!")
	return nil
}

// Shutdown systematically brings the manager down gracefully.
func Shutdown() error {
	beego.Debug("Shuting down mongo sessions...")

	// Close the databases
	for _, session := range singleton.sessions {
		CloseSession(session.mongoSession)
	}

	beego.Debug("Shutdown mongo sessions completed...")
	return nil
}

// CreateSession creates a connection pool for use.
func CreateSession(mode string, sessionName string, hosts []string, databaseName string, username string, password string) error {
	beego.BeeLogger.Debug("CreateSession Mode[%s] SessionName[%s] Hosts[%s] DatabaseName[%s] Username[%s]", mode, sessionName, hosts, databaseName, username)

	// Create the database object
	mongoSession := mongoSession{
		mongoDBDialInfo: &mgo.DialInfo{
			Addrs:    hosts,
			Timeout:  60 * time.Second,
			Database: databaseName,
			Username: username,
			Password: password,
		},
	}

	// Establish the master session.
	var err error
	mongoSession.mongoSession, err = mgo.DialWithInfo(mongoSession.mongoDBDialInfo)
	if err != nil {
		beego.BeeLogger.Error("CreateSession error, %s", err.Error())
		return err
	}

	switch mode {
	case "strong":
		// Reads and writes will always be made to the master server using a
		// unique connection so that reads and writes are fully consistent,
		// ordered, and observing the most up-to-date data.
		// http://godoc.org/github.com/finapps/mgo#Session.SetMode
		mongoSession.mongoSession.SetMode(mgo.Strong, true)
		break

	case "monotonic":
		// Reads may not be entirely up-to-date, but they will always see the
		// history of changes moving forward, the data read will be consistent
		// across sequential queries in the same session, and modifications made
		// within the session will be observed in following queries (read-your-writes).
		// http://godoc.org/github.com/finapps/mgo#Session.SetMode
		mongoSession.mongoSession.SetMode(mgo.Monotonic, true)
	}

	// Have the session check for errors.
	// http://godoc.org/github.com/finapps/mgo#Session.SetSafe
	mongoSession.mongoSession.SetSafe(&mgo.Safe{})

	// Add the database to the map.
	singleton.sessions[sessionName] = mongoSession

	beego.Debug("Create Session Completed!")
	return nil
}

// CopyMasterSession makes a copy of the master session for client use.
func CopyMasterSession() (*mgo.Session, error) {
	return CopySession(MasterSession)
}

// CopyMonotonicSession makes a copy of the monotonic session for client use.
func CopyMonotonicSession() (*mgo.Session, error) {
	return CopySession(MonotonicSession)
}

// CopySession makes a copy of the specified session for client use.
func CopySession(useSession string) (*mgo.Session, error) {
	beego.BeeLogger.Debug("Copy session use [%s].", useSession)

	// Find the session object.
	session := singleton.sessions[useSession]

	if session.mongoSession == nil {
		err := fmt.Errorf("Unable To Locate Session %s", useSession)
		beego.BeeLogger.Error("Copy session error %s.", err.Error())
		return nil, err
	}

	// Copy the master session.
	mongoSession := session.mongoSession.Copy()

	beego.Debug("Copy session completed!")
	return mongoSession, nil
}

// CloneMasterSession makes a clone of the master session for client use.
func CloneMasterSession(sessionID string) (*mgo.Session, error) {
	return CloneSession(sessionID, MasterSession)
}

// CloneMonotonicSession makes a clone of the monotinic session for client use.
func CloneMonotonicSession(sessionID string) (*mgo.Session, error) {
	return CloneSession(sessionID, MonotonicSession)
}

// CloneSession makes a clone of the specified session for client use.
func CloneSession(sessionID string, useSession string) (*mgo.Session, error) {
	beego.BeeLogger.Debug("Clone session use %s.", useSession)

	// Find the session object.
	session := singleton.sessions[useSession]

	if session.mongoSession == nil {
		err := fmt.Errorf("Unable To Locate Session %s", useSession)
		beego.BeeLogger.Error("Clone session error %s", err.Error())
		return nil, err
	}

	// Clone the master session.
	mongoSession := session.mongoSession.Clone()

	beego.Debug("Clone session completed!")
	return mongoSession, nil
}

// CloseSession puts the connection back into the pool.
func CloseSession(mongoSession *mgo.Session) {
	mongoSession.Close()
	beego.Debug("Close session completed!")
}

// GetDatabase returns a reference to the specified database.
func GetDatabase(mongoSession *mgo.Session, useDatabase string) *mgo.Database {
	return mongoSession.DB(useDatabase)
}

// GetCollection returns a reference to a collection for the specified database and collection name.
func GetCollection(mongoSession *mgo.Session, useDatabase string, useCollection string) *mgo.Collection {
	return mongoSession.DB(useDatabase).C(useCollection)
}

// CollectionExists returns true if the collection name exists in the specified database.
func CollectionExists(sessionID string, mongoSession *mgo.Session, useDatabase string, useCollection string) bool {
	database := mongoSession.DB(useDatabase)
	collections, err := database.CollectionNames()

	if err != nil {
		return false
	}

	for _, collection := range collections {
		if collection == useCollection {
			return true
		}
	}

	return false
}

// ToString converts the quer map to a string.
func ToString(queryMap interface{}) string {
	json, err := json.Marshal(queryMap)
	if err != nil {
		return ""
	}

	return string(json)
}

// ToStringD converts bson.D to a string.
func ToStringD(queryMap bson.D) string {
	json, err := json.Marshal(queryMap)
	if err != nil {
		return ""
	}

	return string(json)
}

// Execute the MongoDB literal function.
func Execute(mongoSession *mgo.Session, databaseName string, collectionName string, dbCall DBCall) error {
	beego.BeeLogger.Debug("Execute Database[%s] Collection[%s]", databaseName, collectionName)

	// Capture the specified collection.
	collection := GetCollection(mongoSession, databaseName, collectionName)
	if collection == nil {
		err := fmt.Errorf("Collection %s does not exist", collectionName)
		beego.BeeLogger.Error("Execute error %s", err.Error())
		return err
	}

	// Execute the MongoDB call.
	err := dbCall(collection)
	if err != nil {
		beego.BeeLogger.Error("Execute error %s", err.Error())
		return err
	}

	return nil
}

func DoAction(mongoSession *mgo.Session, collectionName string, dbCall DBCall) error {
	return Execute(mongoSession, "", collectionName, dbCall)
}
