package models

type PackageType string

const (
	Package PackageType = "пакет"
	Box     PackageType = "коробка"
	Film    PackageType = "пленка"

	PackageCost = 5
	BoxCost     = 20
	FilmCost    = 1

	PackageWeightLimit = 10
	BoxWeightLimit     = 30
)

type PackagingType struct {
	Type           PackageType
	AdditionalCost float64
}

func NewPackagingType(packagingType PackageType, additionalCost float64) *PackagingType {
	return &PackagingType{
		Type:           packagingType,
		AdditionalCost: additionalCost,
	}
}

func ToPackageType(s string) PackageType {
	switch s {
	case "пакет":
		return Package
	case "коробка":
		return Box
	case "пленка":
		return Film
	default:
		return ""
	}
}
