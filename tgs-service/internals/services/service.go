package services

import (
	"math/rand"
	"strconv"
	"time"

	"github.com/abielalejandro/tgs-service/config"
	"github.com/abielalejandro/tgs-service/internals/storage"
	"github.com/abielalejandro/tgs-service/pkg/logger"
	"github.com/abielalejandro/tgs-service/pkg/utils"
	"go.uber.org/atomic"
)

type Range struct {
	Min int
	Max int
}

type TgsService struct {
	config *config.Config
	log    *logger.Logger
	rng    *Range
	atom   atomic.Uint32
}

func NewTgsService(config *config.Config) *TgsService {
	return &TgsService{
		config: config,
		log:    logger.New(config.Log.Level),
		rng:    &Range{Min: 0, Max: 0},
	}
}

func (svc *TgsService) GenerateRange(storage storage.Storage) {
	next, _ := storage.GetNext("default")
	minRange, _ := strconv.Atoi(utils.PaddingRight(strconv.Itoa(next), "0", svc.config.Token.Range))
	maxRange, _ := strconv.Atoi(utils.PaddingRight(strconv.Itoa(next), "9", svc.config.Token.Range))
	svc.rng.Max = maxRange
	svc.rng.Min = minRange
}

func (svc *TgsService) GenerateToken() (string, error) {
	counter := svc.atom.Inc()
	if int(counter) > (svc.rng.Max - svc.rng.Min) {
		svc.atom.Store(1)
	}
	rand.Seed(time.Now().UnixNano())

	randomBytesToAdd := make([]byte, svc.config.Range)
	rand.Read(randomBytesToAdd)

	nextSeed := (rand.Intn(svc.rng.Max - svc.rng.Min)) + svc.rng.Min
	seedBytes := []byte(strconv.Itoa(nextSeed))
	toTokenizer := make([]byte, svc.config.Range)

	for i := 0; i < len(seedBytes); i++ {
		toTokenizer[i] = seedBytes[i] + randomBytesToAdd[i]
	}
	next := string(toTokenizer[:])

	return utils.ToBase62(next), nil
}
