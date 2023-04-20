package hex

import (
	"encoding/json"
	"fmt"
	"math/big"
	"strings"
)

type BigInt big.Int

// String 转成10进制的数字
func (b BigInt) String() string {
	v := big.Int(b)
	return v.String()
}

func (b *BigInt) Uint64() uint64 {
	v := big.Int(*b)
	return v.Uint64()
}

func (b *BigInt) Scan(value interface{}) error {
	var t big.Int
	switch v := value.(type) {
	case []byte:
		t.SetBytes(v)
		*b = BigInt(t)
		return nil
	case string:
		t.SetString(v, 10)
		*b = BigInt(t)
		return nil
	}
	return fmt.Errorf("can't convert %T to hex.BigInt", value)
}

func (b *BigInt) EqualTo(a *BigInt) bool {
	a2 := big.Int(*a)
	b2 := big.Int(*b)

	if b2.Cmp(&a2) == 0 {
		return true
	}

	return false
}

func (b *BigInt) GreaterZero() bool {
	zero := big.NewInt(0)
	b2 := big.Int(*b)
	return b2.Cmp(zero) > 0
}

func (b *BigInt) FromDB(value []byte) error {
	return b.Scan(value)
}

func (b BigInt) ToDB() ([]byte, error) {
	var bInt = big.Int(b)
	return bInt.Bytes(), nil
}

func (b BigInt) Bytes() []byte {
	var bInt = big.Int(b)
	return bInt.Bytes()
}

func (b BigInt) MarshalJSON() ([]byte, error) {
	return json.Marshal(b.String())
}

func (b *BigInt) UnmarshalJSON(value []byte) error {
	var s string
	if err := json.Unmarshal(value, &s); err != nil {
		return err
	}
	return b.Scan(s)
}

func remove0x(str string) string {
	if strings.HasPrefix(str, "0x") || strings.HasPrefix(str, "0X") {
		return str[2:]
	}
	return str
}

func HexstrToBigInt(str string) BigInt {
	bigInt := new(big.Int)
	bigInt.SetString(remove0x(str), 16)
	return BigInt(*bigInt)
}

func IntstrToBigInt(str string) BigInt {
	bigInt := new(big.Int)
	bigInt.SetString(str, 10)
	return BigInt(*bigInt)
}
