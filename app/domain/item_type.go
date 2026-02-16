package domain

type ItemType string

const (
	Goods   ItemType = "goods"
	Service ItemType = "service"
)

type ItemStructure string

const (
	SingleItem       ItemStructure = "single"
	ContainsVariants ItemStructure = "variants"
)
