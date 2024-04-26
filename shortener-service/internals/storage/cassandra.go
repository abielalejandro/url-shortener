package storage

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/abielalejandro/shortener-service/config"
	"github.com/gocql/gocql"
)

type CUrl struct {
	Short       string    `cql:"short"`
	Long        string    `cql:"long"`
	LastVisited time.Time `cql:"last_visited"`
	CreatedAt   time.Time `cql:"created_at"`
	ExpiresAt   time.Time `cql:"expires_at"`
}

type CassandraStorage struct {
	*gocql.ClusterConfig
	*gocql.Session
}

func NewCassandraStorage(config *config.Config) *CassandraStorage {
	addr := config.Storage.Addr
	hosts := strings.Split(addr, ",")
	clusterConfig := gocql.NewCluster(hosts...)
	clusterConfig.Keyspace = config.Storage.Db
	clusterConfig.Consistency = gocql.Quorum
	clusterConfig.ProtoVersion = 4

	session, err := clusterConfig.CreateSession()

	if err != nil {
		panic(err)
	}

	return &CassandraStorage{
		ClusterConfig: clusterConfig,
		Session:       session,
	}
}

func (storage *CassandraStorage) GetUrlByLong(ctx context.Context, longUrl string) (*Url, error) {
	var record Url
	var (
		short       string
		long        string
		lastVisited time.Time
		createdAt   time.Time
		expiresAt   time.Time
	)
	scanner := storage.Query(`SELECT short, long, created_at, expires_at, last_visited FROM urls WHERE long=?`,
		longUrl).WithContext(ctx).Iter().Scanner()
	for scanner.Next() {

		err := scanner.Scan(&short, &long, &createdAt, &expiresAt, &lastVisited)
		if err != nil {
			log.Fatal(err)
			return nil, &NotFoundError{}
		}
		record.CreatedAt = createdAt
		record.ExpiresAt = expiresAt
		record.LastVisited = lastVisited
		record.Long = long
		record.Short = short
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
		return nil, &NotFoundError{}
	}

	return &record, nil
}

func (storage *CassandraStorage) ExistsByShort(ctx context.Context, shortUrl string) (bool, error) {
	url, err := storage.GetUrlByShort(ctx, shortUrl)
	if err != nil {
		return false, err
	}

	return url != nil, nil

}

func (storage *CassandraStorage) Create(ctx context.Context, url *Url) (bool, error) {
	command := "INSERT INTO urls (short, long,last_visited,created_at,expires_at) VALUES(?,?,?,?,?)"
	err := storage.Query(command, url.Short, url.Long, time.Now(), time.Now(), time.Now()).WithContext(ctx).Exec()
	if err != nil {
		return false, err
	}

	return true, nil
}

func (storage *CassandraStorage) Update(ctx context.Context, url *Url) (bool, error) {
	command := "UPDATE urls SET last_visited=?, created_at =?,expires_at=? WHERE short='?'"

	err := storage.Query(command, url.LastVisited, url.CreatedAt, url.ExpiresAt, url.Short).WithContext(ctx).Exec()
	if err != nil {
		return false, err
	}

	return true, nil
}

func (storage *CassandraStorage) GetUrlByShort(ctx context.Context, shortUrl string) (*Url, error) {
	var record Url
	var (
		short       string
		long        string
		lastVisited time.Time
		createdAt   time.Time
		expiresAt   time.Time
	)
	scanner := storage.Query(`SELECT short, long, created_at, expires_at, last_visited FROM urls WHERE short=?`,
		shortUrl).WithContext(ctx).Iter().Scanner()
	for scanner.Next() {

		err := scanner.Scan(&short, &long, &createdAt, &expiresAt, &lastVisited)
		if err != nil {
			log.Fatal(err)
			return nil, &NotFoundError{}
		}
		record.CreatedAt = createdAt
		record.ExpiresAt = expiresAt
		record.LastVisited = lastVisited
		record.Long = long
		record.Short = short
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
		return nil, &NotFoundError{}
	}

	return &record, nil
}
