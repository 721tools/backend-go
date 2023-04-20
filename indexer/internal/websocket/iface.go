package websocket

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/721tools/backend-go/indexer/pkg/utils/log16"
	"github.com/721tools/stream-api-go/sdk"
)

var log = log16.NewLogger("module", "websocket")

// EventType
const (
	ITEM_METADATA_UPDATED string = "item_metadata_updated"
	ITEM_LISTED                  = "item_listed"
	ITEM_SOLD                    = "item_sold"
	ITEM_TRANSFERRED             = "item_transferred"
	ITEM_RECEIVED_OFFER          = "item_received_offer"
	ITEM_RECEIVED_BID            = "item_received_bid"
	ITEM_CANCELLED               = "item_cancelled"
	ITEM_FEACH_ALL               = "*"
)

const (
	MAIN_NET = iota
	TEST_NET
)

type ItemListedRes struct {
	EventTimestamp string `json:"event_timestamp"`
	BasePrice      string `json:"base_price"`
	Collection     struct {
		Slug string `json:"slug"`
	} `json:"collection"`
	ExpirationDate string `json:"expiration_date"`
	IsPrivate      bool   `json:"is_private"`
	ListingDate    string `json:"listing_date"`
	ListingType    string `json:"listing_type"`
	Maker          struct {
		Address string `json:"address"`
	} `json:"maker"`
	PaymentToken struct {
		Address  string `json:"address"`
		Decimals int    `json:"decimals"`
		EthPrice int    `json:"eth_price"`
		Name     string `json:"name"`
		Symbol   string `json:"Symbol"`
		UsdPrice string `json:"usd_price"`
	} `json:"payment_token"`
	Item struct {
		Chain struct {
			Name string `json:"name"`
		} `json:"chain"`
		Metadata  ItemListRes `json:"metadata"`
		NFTId     string      `json:"nft_id"`
		Permalink string      `json:"permalink"`
	} `json:"item"`
	Quantity int    `json:"quantity"`
	Taker    string `json:"taker"`
}

type ItemSoldRes struct {
	EventTimestamp string `json:"event_timestamp"`
	ClosingDate    string `json:"closing_date"`
	IsPrivate      bool   `json:"is_private"`
	ListingDate    string `json:"listing_date"`
	ListingType    string `json:"listing_type"`
	Maker          struct {
		Address string `json:"address"`
	} `json:"maker"`
	PaymentToken struct {
		Address  string `json:"address"`
		Decimals int    `json:"decimals"`
		EthPrice int    `json:"eth_price"`
		Name     string `json:"name"`
		Symbol   string `json:"Symbol"`
		UsdPrice string `json:"usd_price"`
	} `json:"payment_token"`
	Quantity int `json:"quantity"`
	Taker    struct {
		Address string `json:"address"`
	} `json:"taker"`
	Transaction struct {
		Timestamp string `json:"Timestamp"`
		Hash      string `json:"hash"`
	} `json:"transaction"`
}

type ItemTransferredRes struct {
	EventTimestamp string `json:"event_timestamp"`
	Transaction    struct {
		Timestamp string `json:"Timestamp"`
		Hash      string `json:"hash"`
	} `json:"transaction"`
	FromAccount struct {
		Address string `json:"address"`
	} `json:"from_account"`
	ToAccount struct {
		Address string `json:"address"`
	} `json:"to_account"`
	Quantity int `json:"quantity"`
}

type ItemMetadataUpdatedRes struct {
	Name            string `json:"name"`
	Description     string `json:"description"`
	ImagePreviewUrl string `json:"image_preview_url"`
	AnimationUrl    string `json:"animation_url"`
	BackgroundColor string `json:"background_color"`
	MetadataUrl     string `json:"metadata_url"`
}

type ItemListRes struct {
	AnimationUrl string `json:"animation_url"`
	ImageUrl     string `json:"image_url"`
	MetadataUrl  string `json:"metadata_url"`
	Name         string `json:"name"`
	Rank         uint64 `json:"rank"`
	FloorPrice   string `json:"floor_price"`
	Verified     uint64 `json:"verified"`
	TotalSupply  uint64 `json:"total_supply"`
}

type ItemCancelledRes struct {
	EventTimestamp string `json:"event_timestamp"`
	ListingType    string `json:"listing_type"`
	PaymentToken   struct {
		Address  string `json:"address"`
		Decimals int    `json:"decimals"`
		EthPrice int    `json:"eth_price"`
		Name     string `json:"name"`
		Symbol   string `json:"Symbol"`
		UsdPrice string `json:"usd_price"`
	} `json:"payment_token"`
	Quantity    int `json:"quantity"`
	Transaction struct {
		Timestamp string `json:"Timestamp"`
		Hash      string `json:"hash"`
	} `json:"transaction"`
}

type ItemReceivedOfferRes struct {
	EventTimestamp string `json:"event_timestamp"`
	BasePrice      string `json:"base_price"`
	CreatedDate    string `json:"created_date"`
	ExpirationDate string `json:"expiration_date"`
	Item           struct {
		Chain struct {
			Name string `json:"name"`
		} `json:"chain"`
		Metadata  ItemMetadataUpdatedRes `json:"metadata"`
		NFTId     string                 `json:"nft_id"`
		Permalink string                 `json:"permalink"`
	} `json:"item"`
	Maker struct {
		Address string `json:"address"`
	} `json:"maker"`
	PaymentToken struct {
		Address  string `json:"address"`
		Decimals int    `json:"decimals"`
		EthPrice int    `json:"eth_price"`
		Name     string `json:"name"`
		Symbol   string `json:"Symbol"`
		UsdPrice string `json:"usd_price"`
	} `json:"payment_token"`
	Quantity int `json:"quantity"`
	Taker    struct {
		Address string `json:"address"`
	} `json:"taker"`
}

type ItemReceivedBidRes struct {
	EventTimestamp string `json:"event_timestamp"`
	BasePrice      string `json:"base_price"`
	CreatedDate    string `json:"created_date"`
	ExpirationDate string `json:"expiration_date"`
	Maker          struct {
		Address string `json:"address"`
	} `json:"maker"`
	Item struct {
		Chain struct {
			Name string `json:"name"`
		} `json:"chain"`
		Metadata  ItemMetadataUpdatedRes `json:"metadata"`
		NFTId     string                 `json:"nft_id"`
		Permalink string                 `json:"permalink"`
	} `json:"item"`
	PaymentToken struct {
		Address  string `json:"address"`
		Decimals int    `json:"decimals"`
		EthPrice int    `json:"eth_price"`
		Name     string `json:"name"`
		Symbol   string `json:"Symbol"`
		UsdPrice string `json:"usd_price"`
	} `json:"payment_token"`
	Quantity int `json:"quantity"`
	Taker    struct {
		Address string `json:"address"`
	} `json:"taker"`
}

type PayloadJson struct {
	EventType string      `json:"event_type,omitempty"`
	Payload   interface{} `json:"payload,omitempty"`
	SentAt    string      `json:"sent_at,omitempty"`
}

type Message struct {
	Topic   string      `json:"topic"`
	Event   string      `json:"event"`
	Payload PayloadJson `json:"payload"`
	Ref     int         `json:"ref"`
}

func (m *Message) UnmarshalJSON(data []byte) error {

	type alice Message
	var temp alice

	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	m.Event = temp.Event
	m.Topic = temp.Topic
	m.Ref = temp.Ref
	m.Payload = PayloadJson{}

	payload := PayloadJson{}
	switch temp.Event {
	case "phx_reply":
		return nil
	case "item_metadata_update":
		payload.Payload = &ItemMetadataUpdatedRes{}
	case "item_listed":
		payload.Payload = &ItemListedRes{}
	case "item_sold":
		payload.Payload = &ItemSoldRes{}
	case "item_transferred":
		payload.Payload = &ItemTransferredRes{}
	case "item_metadata_updated":
		payload.Payload = &ItemMetadataUpdatedRes{}
	case "item_cancelled":
		payload.Payload = &ItemCancelledRes{}
	case "item_received_offer":
		payload.Payload = &ItemReceivedOfferRes{}
	case "item_received_bid":
		payload.Payload = &ItemReceivedBidRes{}
	}

	jsonObj, _ := json.Marshal(temp.Payload)
	json.Unmarshal(jsonObj, &payload)
	m.Payload = payload
	return nil
}

func (m Message) getSlugNameFromTopic() string {
	// get slug from payload.payload.collection.slug
	if m.Event == sdk.ITEM_LISTED {
		res := m.Payload.Payload.(*ItemListedRes)
		return res.Collection.Slug
	}
	if strings.HasPrefix(m.Topic, "collection:") {
		return strings.Split(m.Topic, ":")[1]
	}
	return ""
}

type HubIface interface {
	RUN()
}

type ClientIface interface {
	ServeWs(hub *Hub, w http.ResponseWriter, r *http.Request)
}
