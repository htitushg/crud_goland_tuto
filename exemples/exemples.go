package exemples

import "fmt"

func main() {
	moyenne := 15

	switch {
	case moyenne < 9:
		fmt.Printf("Votre moyenne est médiocre.\n")
	case moyenne < 10:
		fmt.Printf("Votre moyenne est passable.\n")
	case moyenne < 14:
		fmt.Printf("Votre moyenne est Bonne.\n")
	default:
		fmt.Printf("Votre moyenne est excellente.\n")
	}
}
