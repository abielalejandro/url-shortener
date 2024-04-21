package services

import (
	"crypto/rand"
	"math"
	rnd "math/rand"
	"strconv"

	"github.com/abielalejandro/shortener-service/config"
	"github.com/abielalejandro/shortener-service/internals/storage"
	"github.com/abielalejandro/shortener-service/pkg/logger"
	"github.com/abielalejandro/shortener-service/pkg/utils"
	"go.uber.org/atomic"
)

type Service interface {
	GenerateToken() (string, error)
}

type Range struct {
	Min int
	Max int
}

type TgsService struct {
	config       *config.Config
	log          *logger.Logger
	rng          *Range
	atom         atomic.Uint32
	storage      storage.Storage
	cachestorage storage.CacheStorage
}

func NewTgsService(config *config.Config, storage storage.Storage, cachestorage storage.CacheStorage) *TgsService {
	return &TgsService{
		config:       config,
		log:          logger.New(config.Log.Level),
		rng:          &Range{Min: 0, Max: 0},
		storage:      storage,
		cachestorage: cachestorage,
	}
}

func (svc *TgsService) GenerateRange() {
	next, err := svc.storage.GetNext("default")
	if err != nil {
		svc.log.Error(err)
		panic(err)
	}

	factor := math.Pow10(svc.config.Token.Range)
	minRange := next * int(factor)
	maxRange := ((next+1)*int(factor) - 1)
	svc.rng.Max = maxRange
	svc.rng.Min = minRange
}

func (svc *TgsService) GenerateToken() (string, error) {
	counter := svc.atom.Inc()
	if int(counter) > (svc.rng.Max - svc.rng.Min) {
		svc.atom.Store(1)
		svc.GenerateRange()
	}

	randomBytesToAdd := make([]byte, svc.config.Range)
	rand.Read(randomBytesToAdd)

	nextSeed := (rnd.Intn(svc.rng.Max - svc.rng.Min)) + svc.rng.Min
	seedBytes := []byte(strconv.Itoa(nextSeed))
	toTokenizer := make([]byte, svc.config.Range)

	for i := 0; i < len(randomBytesToAdd); i++ {
		toTokenizer[i] = seedBytes[i] + randomBytesToAdd[i]
	}
	next := string(toTokenizer[:])

	return string(utils.ToBase62(next)[:(svc.config.Range + 1)]), nil
}
