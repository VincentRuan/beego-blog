package admin

import (
	"github.com/astaxie/beego"
	"github.com/vincentruan/beego-blog/utilities/mongo"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type User struct {
	ID_      bson.ObjectId `bson:"_id,omitempty"`
	UserName string
	Pwd      string
	Portrait string
	Email    string
	Type     string
}

func Save(user *User) (*User, error) {
	session, err := mongo.CopyMasterSession()
	if err != nil {
		return nil, err
	}
	defer mongo.CloseSession(session)

	f := func(collection *mgo.Collection) error {
		return collection.Insert(user)
	}
	mongo.DoAction(session, "user", f)

	return user, nil
}

func CheckUser(userName, password string) bool {
	ct, err := CountUser(bson.M{"username": userName, "pwd": password})

	if err != nil {
		beego.Error(err)
		return false
	}

	if ct == 0 {
		beego.BeeLogger.Debug("Unable to find any user with user name [%s] and password [%s]", userName, password)
		return false
	}

	return true
}

func CheckAdminUser(userName, password, userType string) bool {
	ct, err := CountUser(bson.M{"username": userName, "pwd": password, "type": userType})

	if err != nil {
		beego.Error(err)
		return false
	}

	if ct == 0 {
		beego.BeeLogger.Debug("Unable to find any user with user name [%s] and password [%s] and type [%s]", userName, password, userType)
		return false
	}

	return true
}

func CountUser(query interface{}) (int, error) {
	session, err := mongo.CopyMasterSession()
	if err != nil {
		return 0, err
	}
	defer mongo.CloseSession(session)

	c := session.DB("").C("user")
	return c.Find(query).Count()
}

func QueryByUserName(userName string) (*User, error) {
	return QueryUser(bson.M{"username": userName})
}

func QueryAdmin() (*User, error) {
	return QueryUser(bson.M{"type": "host"})
}

func QueryUser(query interface{}) (*User, error) {
	session, err := mongo.CopyMasterSession()
	if err != nil {
		return nil, err
	}
	defer mongo.CloseSession(session)

	c := session.DB("").C("user")
	user := User{}

	err = c.Find(query).One(&user)
	if err != nil {
		beego.Error("Unable to query user by ", query)
		return nil, err
	}

	return &user, nil
}

func UpdateUser(user User) error {
	session, err := mongo.CopyMasterSession()
	if err != nil {
		return err
	}
	defer mongo.CloseSession(session)

	c := session.DB("").C("user")

	return c.UpdateId(user.ID_, &user)
}
