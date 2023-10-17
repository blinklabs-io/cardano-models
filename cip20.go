package models

type Cip20Metadata struct {
	Num674 Num674 `cbor:"674,keyasint" json:"674"`
}

type Num674 struct {
	Msg []string `cbor:"msg" json:"msg"`
}
