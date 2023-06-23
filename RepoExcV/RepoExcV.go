package RepoExcV

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type ExcValue struct {
	CurlFrom  string
	CurlTo    string
	Value     float64
	CreatedOn string
}

type Value struct {
	Disclaimer string
	License    string
	Timestamp  int
	Base       string
	Rates      rate
}

type rate struct {
	EUR float64
}

type Repo struct {
	db *gorm.DB
}

func New(host string, port int, user string, password string, dbname string) Repo {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	DB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return Repo{
		db: DB,
	}
}

func GetValue(jsonResp []byte) interface{} {
	var value Value
	err := json.Unmarshal(jsonResp, &value)
	if err != nil {
		panic(err)
	}
	fmt.Println(value)
	var v ExcValue
	v.CurlFrom = "USD"
	v.CurlTo = "EUR"
	v.Value = value.Rates.EUR
	ntime := int64(value.Timestamp)
	v.CreatedOn = time.Unix(ntime, 10000).Format("2006-01-02 15:04:05")
	fmt.Println("sono qui", v)
	return v
}

func (r Repo) GetToDB() []byte {
	var records []ExcValue
	if err := r.db.Find(&records).Error; err != nil {
		log.Fatal(err)
	}

	jsonData, err := json.Marshal(records)
	if err != nil {
		log.Fatal(err)
	}
	return jsonData
}

func (r Repo) AddToDB(val interface{}) {
	value, ok := val.(ExcValue)
	if !ok {
		fmt.Println("not a value")
		return
	}
	r.db.Create(value)
}
