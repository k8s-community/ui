package storage

import (
	"sync"

	"github.com/AlekSi/pointer"
	"github.com/Sirupsen/logrus"
	"github.com/icza/session"
	"github.com/k8s-community/k8s-community/models"
	"gopkg.in/reform.v1"
	"encoding/json"
)

type SessionAttrs struct {
	Activated bool `json:"Activated"`
	HasError  bool `json:"HasError"`
}

type DB struct {
	db     *reform.DB
	logger logrus.FieldLogger
	mux    *sync.RWMutex // mutex to synchronize access to sessions
}

func NewDB(db *reform.DB, logger logrus.FieldLogger) *DB {
	s := &DB{
		db:     db,
		logger: logger,
		mux:    &sync.RWMutex{},
	}

	return s
}

// Get returns the session specified by its id.
// The returned session will have an updated access time (set to the current time).
// nil is returned if this store does not contain a session with the specified id.
func (s *DB) Get(id string) session.Session {
	logger := s.logger.WithField("session_id", id)
	logger.Infof("Get session")

	s.mux.RLock()
	defer s.mux.RUnlock()

	st, err := s.db.FindOneFrom(models.UserTable, "session_id", id)
	if err == reform.ErrNoRows {
		logger.Infof("Session is not found")
		return nil
	}

	if err != nil {
		logger.Errorf("Couldn't get session from DB: %+v", err)
		return nil
	}

	user := st.(*models.User)
	if user.SessionData == nil {
		logger.Infof("Session data is empty")
		return nil
	}

	var data SessionAttrs
	err = json.Unmarshal([]byte(*user.SessionData), &data)
	if err != nil {
		logger.Errorf("Couldn't unmarshal session %+v: %+v", user.SessionData, err)
		return nil
	}

	sessionData := session.NewSessionOptions(&session.SessOptions{
		CAttrs: map[string]interface{}{"Login": user.Name},
		Attrs:  map[string]interface{}{"Activated": data.Activated, "HasError": data.HasError},
	})

	logger.Info("Session was found")

	return sessionData
}

// Add adds a new session to the store.
func (s *DB) Add(sess session.Session) {
	s.mux.RLock()
	defer s.mux.RUnlock()

	sessID := sess.ID()
	logger := s.logger.WithField("session_id", sessID)
	logger.Infof("Add session...")

	user := &models.User{}

	login := sess.CAttr("Login").(string)
	source := sess.CAttr("Source").(string)

	st, err := s.db.SelectOneFrom(models.UserTable, "WHERE source = $1 AND name = $2", source, login)

	if err != nil && err != reform.ErrNoRows {
		logger.Errorf("Couldn't get user data (%s) from DB: %+v", login, err)
		return
	} else if err != reform.ErrNoRows {
		user = st.(*models.User)
	}

	user.Source = source
	user.Name = login
	user.SessionID = pointer.ToString(sessID)

	data := &SessionAttrs{
		Activated: sess.Attr("Activated").(bool),
		HasError:  sess.Attr("HasError").(bool),
	}

	jsData, err := json.Marshal(data)
	if err != nil {
		logger.Errorf("Couldn't marshal data: %+v, %+v", data, err)
		return
	}

	user.SessionData = pointer.ToString(string(jsData))

	err = s.db.Save(user)
	if err != nil {
		logger.Errorf("Couldn't save session in database %+v: %+v", user, err)
	} else {
		logger.Info("Session data was saved")
	}
}

// Remove removes a session from the store.
func (s *DB) Remove(sess session.Session) {
	s.logger.Infof("Remove session %s", sess.ID())

	s.mux.RLock()
	defer s.mux.RUnlock()

	st, err := s.db.FindOneFrom(models.UserTable, "session_id", sess.ID())
	if err == reform.ErrNoRows {
		return
	}

	if err != nil {
		s.logger.Errorf("Couldn't get session %s from DB: %+v", sess.ID(), err)
		return
	}

	user := st.(*models.User)
	user.SessionID = nil
	user.SessionData = nil
	err = s.db.Save(user)

	if err != nil {
		s.logger.Errorf("Couldn't save session in database %+v: %+v", user, err)
	}
}

// Close closes the session store, releasing any resources that were allocated.
func (s *DB) Close() {
	return
}
