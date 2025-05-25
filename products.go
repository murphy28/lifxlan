package lifxlan

import (
	_ "embed"
	"encoding/json"
	"fmt"
)

// Products json file is sourced from the official LIFX GitHub repository:
// https://github.com/LIFX/products/blob/master/products.json

type Vendor struct {
	VendorID int       `json:"vid"`
	Name     string    `json:"name"`
	Defaults Features  `json:"defaults"`
	Products []Product `json:"products"`
}

type Features struct {
	Hev               bool   `json:"hev"`
	Color             bool   `json:"color"`
	Chain             bool   `json:"chain"`
	Matrix            bool   `json:"matrix"`
	Relays            bool   `json:"relays"`
	Buttons           bool   `json:"buttons"`
	Infrared          bool   `json:"infrared"`
	Multizone         bool   `json:"multizone"`
	TemperatureRange  [2]int `json:"temperature_range"` // [min, max] in Kelvin
	ExtendedMultizone bool   `json:"extended_multizone"`
}

type Product struct {
	ProductID int      `json:"pid"`
	Name      string   `json:"name"`
	Features  Features `json:"features"`
	Upgrades  []struct {
		Major    int      `json:"major"`
		Minor    int      `json:"minor"`
		Features Features `json:"features"`
	}
}

var (
	//go:embed products.json
	productsJSON []byte

	ProductsByVendorIDAndProductID map[int]map[int]Product
)

func InitializeProducts() error {
	var vendors []Vendor
	if err := json.Unmarshal(productsJSON, &vendors); err != nil {
		return fmt.Errorf("failed to unmarshal products JSON: %w", err)
	}

	ProductsByVendorIDAndProductID = make(map[int]map[int]Product)

	for _, vendor := range vendors {
		vendorMap, exists := ProductsByVendorIDAndProductID[vendor.VendorID]
		if !exists {
			vendorMap = make(map[int]Product)
			ProductsByVendorIDAndProductID[vendor.VendorID] = vendorMap
		}

		for _, product := range vendor.Products {
			vendorMap[product.ProductID] = product
		}
	}

	return nil
}

// GetProduct returns the product information for a given vendor ID and product ID
func GetProduct(vendorID, productID int) (Product, error) {
	if vendorMap, exists := ProductsByVendorIDAndProductID[vendorID]; exists {
		if product, exists := vendorMap[productID]; exists {
			return product, nil
		}
	}
	return Product{}, fmt.Errorf("product with vendor ID %d and product ID %d not found", vendorID, productID)
}

func PrintProduct(product Product) {
	fmt.Printf("Product ID: %d\n", product.ProductID)
	fmt.Printf("Name: %s\n", product.Name)
	fmt.Printf("Features:\n")
	fmt.Printf("  HEV: %t\n", product.Features.Hev)
	fmt.Printf("  Color: %t\n", product.Features.Color)
	fmt.Printf("  Chain: %t\n", product.Features.Chain)
	fmt.Printf("  Matrix: %t\n", product.Features.Matrix)
	fmt.Printf("  Relays: %t\n", product.Features.Relays)
	fmt.Printf("  Buttons: %t\n", product.Features.Buttons)
	fmt.Printf("  Infrared: %t\n", product.Features.Infrared)
	fmt.Printf("  Multizone: %t\n", product.Features.Multizone)
	fmt.Printf("  Extended Multizone: %t\n", product.Features.ExtendedMultizone)
	fmt.Printf("  Temperature Range: %dK - %dK\n", product.Features.TemperatureRange[0], product.Features.TemperatureRange[1])
	if len(product.Upgrades) > 0 {
		fmt.Println("Upgrades:")
		for _, upgrade := range product.Upgrades {
			fmt.Printf("  Major: %d, Minor: %d, Features: %+v\n", upgrade.Major, upgrade.Minor, upgrade.Features)
		}
	}
}
