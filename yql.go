package yql

import (
	"encoding/json"
	"errors"
	"exp/sql"
	"exp/sql/driver"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

func init() {
	sql.Register("yql", &YQLDriver{})
}

type YQLDriver struct {

}

type YQLConn struct {
	c *http.Client
}

func (d *YQLDriver) Open(dsn string) (driver.Conn, error) {
	return &YQLConn{http.DefaultClient}, nil
}

func (c *YQLConn) Close() error {
	c.c = nil
	return nil
}

type YQLStmt struct {
	c *YQLConn
	q string
}

func (c *YQLConn) Begin() (driver.Tx, error) {
	return nil, errors.New("Begin does not supported")
}

func (c *YQLConn) Prepare(query string) (driver.Stmt, error) {
	return &YQLStmt{c, query}, nil
}

func (s *YQLStmt) Close() error {
	return nil
}

func (s *YQLStmt) NumInput() int {
	return strings.Count(s.q, "?")
}

func (s *YQLStmt) bind(args []interface{}) error {
	b := s.q
	for _, v := range args {
		b = strings.Replace(b, "?", fmt.Sprintf("%q", v), 1)
	}
	s.q = b
	return nil
}

func (s *YQLStmt) Query(args []interface{}) (driver.Rows, error) {
	if err := s.bind(args); err != nil {
		return nil, err
	}
	res, err := http.Get(fmt.Sprintf("http://query.yahooapis.com/v1/public/yql?q=%s&format=json", url.QueryEscape(s.q)))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	var data interface{}
	b, err := ioutil.ReadAll(res.Body)
	//	println(string(b))
	err = json.Unmarshal(b, &data)
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, errors.New("Unsupported result")
	}
	var ok bool
	data = data.(map[string]interface{})["query"]
	if data == nil {
		return nil, errors.New("Unsupported result")
	}
	data = data.(map[string]interface{})["results"]
	if data == nil {
		return nil, errors.New("Unsupported result")
	}
	results, ok := data.(map[string]interface{})
	if !ok {
		return nil, errors.New("Unsupported result")
	}
	for _, v := range results {
		if vv, ok := v.([]interface{}); ok {
			return &YQLRows{s, 0, vv}, nil
		}
	}
	return nil, errors.New("Unsupported result")
}

type YQLResult struct {
	s *YQLStmt
}

func (s *YQLStmt) Exec(args []interface{}) (driver.Result, error) {
	return nil, errors.New("Exec does not supported")
}

type YQLRows struct {
	s *YQLStmt
	n int
	d []interface{}
}

func (rc *YQLRows) Close() error {
	return nil
}

func (rc *YQLRows) Columns() []string {
	return []string{"results"}
}

func (rc *YQLRows) Next(dest []interface{}) error {
	if rc.n == len(rc.d) {
		return errors.New("EOF")
	}
	if s, ok := rc.d[rc.n].(string); ok {
		dest[0] = s
	} else {
		dest[0] = rc.d[rc.n]
	}
	rc.n++
	return nil
}
