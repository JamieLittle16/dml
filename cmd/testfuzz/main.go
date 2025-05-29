package main

import (
	"fmt"
	"os"
	"sort"

	"dml/internal/colour"
)

func main() {
	fmt.Println("DML Dynamic Fuzz Level Test")
	fmt.Println("===========================")
	fmt.Println()

	// Test with common colors
	colorTests := []string{
		"white", "#FFFFFF",
		"black", "#000000",
		"red", "#FF0000",
		"green", "#00FF00",
		"blue", "#0000FF",
		"cyan", "#00FFFF",
		"magenta", "#FF00FF",
		"yellow", "#FFFF00",
		"grey", "#808080",
		"pink", "#FFC0CB",
		"orange", "#FFA500",
		"purple", "#800080",
	}

	for _, colorName := range colorTests {
		hex := colour.ToHex(colorName)
		fuzz := colour.GetFuzzLevel(hex)
		fmt.Printf("Color: %-10s  Hex: %-8s  Fuzz Level: %s\n", colorName, hex, fuzz)
	}

	fmt.Println("\nTest with varying brightness levels (grayscale):")
	fmt.Println("================================================")
	
	// Generate grayscale values and their fuzz levels
	grayscaleTests := make(map[string]string)
	for i := 0; i <= 255; i += 25 {
		hex := fmt.Sprintf("#%02X%02X%02X", i, i, i)
		fuzz := colour.GetFuzzLevel(hex)
		grayscaleTests[hex] = fuzz
	}

	// Sort by hex value and print
	var hexCodes []string
	for hex := range grayscaleTests {
		hexCodes = append(hexCodes, hex)
	}
	sort.Strings(hexCodes)

	for _, hex := range hexCodes {
		fuzz := grayscaleTests[hex]
		fmt.Printf("Grayscale: %-8s  Fuzz Level: %s\n", hex, fuzz)
	}

	// Test with varying saturation levels (red with different saturation)
	fmt.Println("\nTest with varying saturation levels (red):")
	fmt.Println("=========================================")
	
	saturationTests := make(map[string]string)
	for i := 0; i <= 255; i += 25 {
		hex := fmt.Sprintf("#FF%02X%02X", i, i)
		fuzz := colour.GetFuzzLevel(hex)
		saturationTests[hex] = fuzz
	}

	// Sort and print
	hexCodes = nil
	for hex := range saturationTests {
		hexCodes = append(hexCodes, hex)
	}
	sort.Strings(hexCodes)

	for _, hex := range hexCodes {
		fuzz := saturationTests[hex]
		fmt.Printf("Red saturation: %-8s  Fuzz Level: %s\n", hex, fuzz)
	}

	if len(os.Args) > 1 {
		// Allow command-line testing of specific colors
		for _, arg := range os.Args[1:] {
			hex := colour.ToHex(arg)
			if hex == "" {
				fmt.Printf("Unknown color: %s\n", arg)
				continue
			}
			fuzz := colour.GetFuzzLevel(hex)
			fmt.Printf("\nColor: %-10s  Hex: %-8s  Fuzz Level: %s\n", arg, hex, fuzz)
		}
	}
}