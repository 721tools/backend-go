package hotness

import (
	"context"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/721tools/backend-go/indexer/pkg/consts"
	"github.com/721tools/backend-go/indexer/pkg/mq"
	"github.com/721tools/backend-go/indexer/pkg/utils/log16"
	"github.com/721tools/backend-go/indexer/pkg/utils/quit"
	"github.com/721tools/backend-go/indexer/pkg/utils/timedset"
)

var (
	log        = log16.NewLogger("module", "hotness")
	hotnessMtx sync.RWMutex
	hotness    []string
)

type NFT struct {
	Contract string
	TokenID  string
	Behavior int
	Time     time.Time
}

func parseKEY(key string) (contract, token string, behavior int) {
	str := strings.Split(key, "/")

	if len(str) == 3 {
		i, _ := strconv.Atoi(str[2])
		return str[0], str[1], i
	}

	return "", "", 0
}

func calculateHotness(s *timedset.TimedSet, weights map[int]float64, n int) []string {
	for {
		items := s.GetAll()
		nfts := make([]NFT, 0)
		for _, key := range items {
			nft := key.(NFT)
			nfts = append(nfts, nft)
		}

		// 计算 top N NFT
		topN := calculateHotnessInner(nfts, weights, n)

		// 更新 hotnessMap
		hotnessMtx.Lock()
		hotness = topN
		hotnessMtx.Unlock()

		fmt.Println(hotness)
		time.Sleep(1 * time.Minute)
	}
}

func calculateHotnessInner(nfts []NFT, weights map[int]float64, n int) []string {
	// 构建每个 behavior 对应的权重
	behaviorWeights := make(map[int]float64)
	for behavior, weight := range weights {
		behaviorWeights[behavior] = weight
	}

	// 统计每个 NFT 的得分
	scores := make(map[string]float64)
	for _, nft := range nfts {
		key := nft.Contract
		weight, ok := behaviorWeights[nft.Behavior]
		if ok {
			// 计算时间衰减权重，衰减因子为 0.9，衰减间隔为 1 小时
			decayFactor := 1.0
			now := time.Now()
			hours := now.Sub(nft.Time).Hours()
			for i := 0; i < int(hours); i++ {
				decayFactor *= 0.9
			}
			weight *= decayFactor

			scores[key] += weight
		}
	}

	// 构建分数与 key 的对应关系，用于排序
	var pairs []pair
	for key, score := range scores {
		pairs = append(pairs, pair{key, score})
	}

	// 根据分数排序
	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].score > pairs[j].score
	})

	// 取出前 N 个 key
	topN := make([]string, 0, n)
	for i := 0; i < n && i < len(pairs); i++ {
		key := pairs[i].key
		topN = append(topN, key)
	}

	return topN
}

type pair struct {
	key   string
	score float64
}

func Init() {

	ctx := context.Background()
	s := timedset.NewTimedSet(time.Hour * 4)
	q := quit.NewQuit()
	q.WatchOsSignal()
	mq := mq.GetMQ()
	go mq.Subscribe(ctx, consts.NFT_BEHAVIOR)

	go func() {
		for {
			pop := mq.Pop()
			contract, token, behavior := parseKEY(pop.Payload)
			if contract == "" {
				continue
			}

			nft := NFT{
				Contract: contract,
				TokenID:  token,
				Behavior: behavior,
				Time:     time.Now(),
			}

			s.Add(nft)
		}
	}()

	// behavior 权重
	weights := map[int]float64{
		int(consts.NFT_TRANSFORM): 0.1,
		int(consts.NFT_LIST):      0.1,
		int(consts.NFT_MINT):      0.4,
		int(consts.NFT_SALE):      1.6,
	}

	// 每隔一分钟，计算 Top 100
	go calculateHotness(s, weights, 100)

	if q.IsQuit() {
		os.Exit(0)
	}

	fmt.Println("hshshshshhss")
}

func GetHotness() []string {
	fmt.Println(hotness)
	return hotness
}
