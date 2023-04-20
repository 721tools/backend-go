package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type RefillSpec struct {
	BlockStart      uint64 `json:"block_start,omitempty"`
	BlockEnd        uint64 `json:"block_end,omitempty"`
	ContractAddress string `json:"contract_address,omitempty"`
}

type RefillStatus struct {
	RunHeight uint64 `json:"run_height"`
}

// Refill 重新跑合约
type Refill struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   RefillSpec   `json:"spec,omitempty"`
	Status RefillStatus `json:"status,omitempty"`
}
