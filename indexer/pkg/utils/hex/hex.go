package hex

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/common/hexutil"
)

type Hex []byte

func (h Hex) HexStr() string {
	return h.String()
}

func (h Hex) String() string {
	return hexutil.Encode(h)
}

func (h Hex) Bytes() []byte {
	return []byte(h)
}

func (h Hex) NoPrefixHex() string {
	return h.String()[2:]
}

func (h Hex) EqualTo(o Hex) bool {
	if len(h) != len(o) {
		return false
	}
	for idx := 0; idx < len(o); idx++ {
		if h[idx] != o[idx] {
			return false
		}
	}
	return true
}

func (h *Hex) ToHex(value any) (err error) {
	switch v := value.(type) {
	case []byte:
		*h = v
		return
	case string:
		if len(v) > 2 {
			*h, err = hexutil.Decode(v)
		}
		return
	default:
		return fmt.Errorf("can not convert %T to HexStr", value)
	}
}

func (h Hex) MarshalJSON() ([]byte, error) {
	return json.Marshal(h.String())
}

func (h *Hex) UnmarshalJSON(value []byte) error {
	if len(value) == 0 {
		return nil
	}
	var s string
	if err := json.Unmarshal(value, &s); err != nil {
		return err
	}
	return h.ToHex(s)
}

func (h Hex) ToDB() ([]byte, error) {
	return []byte(h), nil
}

func (h *Hex) FromDB(value []byte) error {
	*h = value
	return nil
}

func HexstrToHex(str string) Hex {
	var h Hex
	_ = h.ToHex(str)
	return h
}

func IntToHex(data uint64) Hex {
	var h Hex
	_ = h.ToHex(fmt.Sprintf("0x%02x", data))
	return h
}

func stripBytes(input []byte) []byte {
	var result []byte
	for _, v := range input {
		if v != 0 {
			result = append(result, v)
		}
	}
	return result
}

func TrimHexStrAndDecodeToStr(hexStr string) string {
	var s string
	if len(hexStr) > 64*2+2 {
		s = hexStr[2+128:]
	} else {
		s = hexStr[2:]
	}
	v, _ := hexutil.Decode("0x" + s)
	return strings.TrimSpace(string(stripBytes(v)))
}
