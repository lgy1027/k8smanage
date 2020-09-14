// @Author : liguoyu
// @Date: 2019/10/29 15:42
package protocol

import (
	"errors"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"net/url"
)

type Service struct {
	Host     string
	Port     int
	Scheme   string
	Path     string
	Url      *url.URL
	Account  string
	Password string
}

func (h *Service) Init() (err error) {
	if h.Url == nil {
		if h.Scheme == "" {
			h.Scheme = "http"
		}
		if h.Port == 0 {
			h.Port = 80
		}
		u := fmt.Sprintf("%s://%s:%d%s", h.Scheme, h.Host, h.Port, h.Path)
		h.Url, err = url.Parse(u)
	}
	return
}

type Mysql struct {
	Address      string
	User         string
	Password     string
	DbName       string
	Timeout      string
	ReadTimeout  string
	WriteTimeout string
	db           *gorm.DB
}

func (m *Mysql) Init() (err error) {
	if m.User == "" || m.Address == "" || m.DbName == "" {
		return errors.New("mysql: missing user, address, db_name")
	}
	if m.db == nil {
		dsn := fmt.Sprintf(`%s:%s@tcp(%s)/%s?timeout=%s&readTimeout=%s&charset=utf8&parseTime=True&loc=Local`,
			m.User, m.Password, m.Address, m.DbName, m.Timeout, m.ReadTimeout,
		)
		m.db, err = gorm.Open("mysql", dsn)
		if err != nil {
			return err
		}
		m.db.SingularTable(true)
		//m.db.CreateTable(&User{},&Advice{})
	}
	return
}

func (m *Mysql) Db() *gorm.DB {
	if m.db == nil {
		m.Init()
	}
	return m.db
}

type Redis struct {
	Address string
	Db      int
	pool    *redis.Pool
}

func (r *Redis) Init() (err error) {
	if r.pool == nil {
		opts := make([]redis.DialOption, 0)
		if r.Db != 0 {
			opts = append(opts, redis.DialDatabase(r.Db))
		}
		r.pool = &redis.Pool{
			Dial: func() (redis.Conn, error) {
				c, err := redis.Dial("tcp", r.Address, opts...)
				if err != nil {
					return nil, err
				}
				return c, nil
			},
		}
		// have a try
		err = r.pool.Get().Close()
	}
	return
}

func (r *Redis) Pool() *redis.Pool {
	return r.pool
}
