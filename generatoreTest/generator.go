package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"time"
)

// Funzione per convertire un numero in binario con una lunghezza limitata (fino a 4 bit)
func toBinary(n int, maxBits int) string {
	return fmt.Sprintf(fmt.Sprintf("%%0%db", maxBits), n)
}

// Funzione per generare un nome binario lungo per il comando 'a'
func toLongBinary(n int) string {
	// Genera binario lungo (ad esempio fino a 12 bit)
	return fmt.Sprintf("%b", n)
}

func generateFile(filename string, numBlocks int, includeR bool, includeO bool, includeA bool, includeT bool) {
	file, err := os.Create(filename)
	if err != nil {
		fmt.Println("Errore nella creazione del file:", err)
		return
	}
	defer file.Close()

	// Inizializza il generatore di numeri casuali
	rand.Seed(time.Now().UnixNano())

	// Scrivi 'c' all'inizio
	file.WriteString("c\n")

	// Genera i blocchi
	for i := 0; i < numBlocks; i++ {
		// Scrivi un blocco di comandi (r, o, a, t, S)
		numCommands := rand.Intn(5) + 3 // Un blocco avrà tra 3 e 7 comandi

		for j := 0; j < numCommands; j++ {
			command := rand.Intn(4) // 0 = r, 1 = o, 2 = a, 3 = t
			x := rand.Intn(201) - 100
			y := rand.Intn(201) - 100
			z := rand.Intn(201) - 100
			w := rand.Intn(201) - 100

			// Comando r (raggio) con binario breve (max 4 bit)
			if command == 0 && includeR {
				file.WriteString(fmt.Sprintf("r %d %d %s\n", x, y, toBinary(rand.Intn(16), 4))) // Limitiamo a 15 (4 bit)
			} else if command == 1 && includeO {
				// Comando o (ostacolo)
				file.WriteString(fmt.Sprintf("o %d %d %d %d\n", x, y, z, w))
			} else if command == 2 && includeA {
				// Comando a (azione) con nome binario lungo
				// Il comando 'a' deve essere nel formato a x y w dove w è un nome binario lungo
				longBinaryName := toLongBinary(rand.Intn(1024)) // Limita a 1024 per un binario lungo
				file.WriteString(fmt.Sprintf("a %d %d %s\n", x, y, longBinaryName))
			} else if command == 3 && includeT {
				// Comando t con il terzo valore in binario breve (max 4 bit)
				file.WriteString(fmt.Sprintf("t %d %d %s\n", x, y, toBinary(rand.Intn(16), 4))) // Limitiamo a 15 (4 bit)
			}
		}

		// Scrivi 'S' alla fine di ogni blocco
		file.WriteString("S\n")
	}

	// Scrivi 'f' alla fine
	file.WriteString("f\n")
}

func main() {
	// Definisci i parametri da linea di comando
	includeR := flag.Bool("r", true, "Include comandi 'r' nel file")
	includeO := flag.Bool("o", true, "Include comandi 'o' nel file")
	includeA := flag.Bool("a", true, "Include comandi 'a' nel file")
	includeT := flag.Bool("t", true, "Include comandi 't' nel file")
	numBlocks := flag.Int("b", 1000, "Numero di blocchi da generare")
	filename := flag.String("f", "test_file.txt", "Nome del file di output")

	// Parsea i parametri
	flag.Parse()

	// Genera il file con i parametri forniti
	generateFile(*filename, *numBlocks, *includeR, *includeO, *includeA, *includeT)

	fmt.Printf("File '%s' generato con successo!\n", *filename)
}
