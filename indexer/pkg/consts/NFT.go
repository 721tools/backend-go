package consts

// NFT Plateform
type Plateform uint8

const (
	OPENSEA Plateform = iota
	Blur
	KIWISWAP
)

// NFT Behavior
type Behavior int

const (
	NFT_TRANSFORM Behavior = iota
	NFT_MINT
	NFT_LIST
	NFT_SALE
)

const (
	NFT_BEHAVIOR = "NFT-BEHAVIOR"
)

func ToString(b Behavior) string {
	switch b {
	case NFT_TRANSFORM:
		return "NFT-TRANSFORM"
	case NFT_MINT:
		return "NFT-MINT"
	case NFT_LIST:
		return "NFT-LIST"
	case NFT_SALE:
		return "NFT-SALE"
	default:
		return ""
	}
}

// NFT keys in MQ
const (
	OPENSEA_ETH_ORDER_LISTING    = "OPENSEA-ETH-ORDER-LISTING"
	OPENSEA_ETH_ORDER_OFFER      = "OPENSEA-ETH-ORDER-OFFER"
	OPENSEA_ETH_COLLECTION_OFFER = "OPENSEA-ETH-COLLECTION-OFFER"
	GAS_PRICE_NOW_MQ             = "GAS-PRICE-NOW"
	GAS_PRICE_NOW_KEY            = "GAS-PRICE-NOW-KEY"
)
