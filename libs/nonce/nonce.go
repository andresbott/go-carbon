package nonce

import (
	"git.andresbott.com/Golang/carbon/libs/utils"
	"time"
)

// default length of the generated nonce
const defaultLen = 32 // same length as uuids

type Cfg struct {
	Duration time.Duration
	Len      int
	Store    ExpiryStore
}

type Nonce struct {
	dur   time.Duration
	store ExpiryStore
	len   int
}

// New returns a new instance of Nonce Manager
func New(cfg Cfg) *Nonce {

	var store ExpiryStore
	if cfg.Store != nil {
		store = cfg.Store
	} else {
		store = &memStore{
			expiration: cfg.Duration,
			db:         make(map[string]time.Time),
		}
	}

	if cfg.Len == 0 {
		cfg.Len = defaultLen
	}

	return &Nonce{
		dur:   cfg.Duration,
		store: store,
		len:   cfg.Len,
	}
}

func (n Nonce) Get() string {
	//uuid := uuid.New().String()
	s := utils.RandStr(n.len)
	n.store.Save(s)
	return s
}

func (n Nonce) Validate(in string) bool {
	return n.store.Valid(in)
}

// ExpiryStore is an interface for temporary keeping string that are going to expire eventually
type ExpiryStore interface {
	SetDuration(duration time.Duration)
	Save(in string)
	Valid(in string) bool
}

// memStore is a quick in memory temporary storage for expiration strings
type memStore struct {
	expiration time.Duration
	db         map[string]time.Time
}

func (m *memStore) SetDuration(d time.Duration) {
	m.expiration = d
}

func (m *memStore) Save(in string) {
	m.db[in] = time.Now().Add(m.expiration)
}

func (m *memStore) Valid(s string) bool {
	m.expire()
	if _, ok := m.db[s]; ok {
		return true
	} else {
		return false
	}
}

// expire will remove old items from the storage
func (m memStore) expire() {
	now := time.Now()
	for k, v := range m.db {
		if now.After(v) {
			delete(m.db, k)
		}
	}
}
