//Autore: 29407A Corrado Francesco Emanuele
//Progetto del corso "Algoritmi e strutture dati", sessione d'esame Gennaio/Febbraio 2025

// Piano: Insieme dei punti {(x, y) ∈ Z x Z }
// Passo Unitario Orizzontale: segmento orizzontale che collega due punti di coordinate (x, y) e (x+1, y)
// Passo Unitario Verticale: segmento verticale che collega due punti di coordinate (x, y) e (x, y+1)
// Percorso: Successione S = p1, p2, ..., pk di k passi unitati orizzontali e/o verticali tali che per ogni i ∈ {1, ..., k−1}
// pi e pi+1 hanno in comune solo un vertice e nessun punto è comune a più di due dei passi unitari. La lunghezza di S è k.
// Un percorso è detto libero se non incontra ostacoli.
// Distanza D(A, B): distanza tra due punti A e B è data dalla lunghezza minima dei percordi che li collegano.
// es. Se A = (xA, yA) e B = (xB, yB) allora D(A, B) = |xB − xA| + |yB − yA|

// Automa: Identificato univocamente da un nome η, che è una stringa finita sull'alfabeto {0, 1}.
// (η = b1 b2... bn per qualche intero positivo n e bi ∈ {0, 1} per ogni i ∈ {1, ..., n}).
// posizione P(η): posizione di un'automa in un dato momento specificata dando le coordinate (x0, y0) del punto del piano
// in cui si trova l'automa. Quindi, P (η) = (x, y) significa che η si trova nel punto (x, y); in un punto può esservi più di un automa.

// Ostacoli: Insieme di punti contenuti in un rettangolo con lati verticali e orizzontali. Un rettangolo si indica come
// R(x0, y0, x1, y1), dove x0, x1, y0, y1 sono interi e (x0, y0) e (x1, y1) sono rispettivamente le coordinate dei vertici 
// in basso a sinistra e in alto a destra del rettangolo. È possibile che gli ostacoli si sovrappongano; nessun automa potrà 
// posizionarsi in un punto appartenente a qualche ostacolo.

// Una sorgente può emettere segnali di richiamo rappresentati da stringhe finite α sull’alfabeto {0,1}. All'atto dell'emissione del
// segnale ogni automa determina se tale segnale lo riguarda o meno. l’automa di nome η deve rispondere al richiamo del 
// segnale α se e solo se α è un prefisso di η, vale a dire η = α b1 ... bn. Tra questi automi, si sposteranno verso la 
// sorgente tutti e soli gli automi che hanno distanza minima dalla sorgente e che possono muoversi lungo un percorso libero di distanza minima.
// Più precisamente, supponiamo che una sorgente posizionata in un punto del piano emetta il richiamo α e che l’insieme di automi 
// a cui è diretto il richiamo sia {η1, ..., ηk} (dunque α è prefisso di ciascun ηi). Per ogni i, indichiamo con di la 
// distanza dell’automa ηi dalla sorgente α e definiamo: d = min{di | 1 ≤ i ≤ k}. Sia ora A ⊆ {η1 , . . . , ηk } l’insieme degli automi 
// ηi tali che di = d e per cui esiste un percorso libero di lunghezza di. Allora, per ogni ηj ∈ A si pone P (ηj ) = (x, y);
// per ogni altro automa η, la posizione P (η) rimane immutata.
// Il piano può essere molto grande rispetto al numero di automi e ostacoli contenuti: non è efficiente utilizzare una matrice per rappresentare il piano.

package main

import(
	"fmt"
	"bufio"
	"os"
	"strings"
	"strconv"
	"sort"
	"container/heap"
)

// STRUTTURE DATI e relative funzioni

// Struttura dati per rappresentare un punto nel piano come coordinate cartesiane vX e vY
type Punto struct {
	x, y int
}

// Struttura dati per rappresentare un automa con la sua posizione nel piano tramite un punto
type Automa struct {
	p Punto
}

// Struttura dati per rappresentare un ostacolo come due punti che rappresentano i vertici opposti di un rettangolo
type Ostacolo struct {
	p1, p2 Punto
}

// Struttura dati per rappresentare un nodo di un Trienode implementato per la ricerca dei prefissi degli automi
type TrieNode struct {
	children map[rune]*TrieNode
	isEnd    bool
	automi   []string
}

// Crea un nuovo nodo del Trienode
func newTrieNode() *TrieNode {
	return &TrieNode{
		children: make(map[rune]*TrieNode),
		automi:   make([]string, 0),
	}
}

// Aggiunge un nome di automa al Trienode
func (node *TrieNode) insert(name string) {
	current := node
	for _, char := range name {
		if _, exists := current.children[char]; !exists {
			current.children[char] = newTrieNode()
		}
		current = current.children[char]
	}
	current.isEnd = true
	current.automi = append(current.automi, name)
}

// Restituisce tutti gli automi che hanno un dato prefisso
func (node *TrieNode) searchPrefix(prefix string) []string {
	current := node
	for _, char := range prefix {
		if _, exists := current.children[char]; !exists {
			return []string{}
		}
		current = current.children[char]
	}
	return current.collectAll()
}

// Restituisce tutti gli automi a partire da un nodo dato
func (node *TrieNode) collectAll() []string {
	var result []string
	if node.isEnd {
		result = append(result, node.automi...)
	}
	for _, child := range node.children {
		result = append(result, child.collectAll()...)
	}
	return result
}

// Struttura dati per rappresentare un percorso di un automa
type stateRoad struct {
	x, y, steps, dir, turns int
}

// Struttura dati per rappresentare un percorso di un automa con priorità
type orderedRoad struct {
	ox, oy, odir, priority int
}

// Struttura dati per rappresentare i nodi visitati durante la ricerca di un percorso
type visitedNodes struct {
	vX, vY, vDir int
}

// Struttura dati per rappresentare un percorso di un automa con distanza
type DistanceRoad struct {
	name string
	distance int
}

// Struttura dati per rappresentare una coda con priorità tra i percorsi
type PriorityQueue []DistanceRoad

// Implementazione delle funzioni dell'interfaccia heap per la PriorityQueue
func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].distance < pq[j].distance
}

func (pq PriorityQueue) Swap(i, j int){
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *PriorityQueue) Push(x interface{}) {
	*pq = append(*pq, x.(DistanceRoad))
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[:n-1]
	return item
}

// Struttura dati per rappresentare un piano con automi, ostacoli, posizioni occupate da automi ed ostacoli e nodo trienoide per la ricerca dei prefissi degli automi
type Piano struct {
	automi   map[string]Automa
	ostacoli []Ostacolo
	occupati map[Punto]bool
	posAutomi map[Punto][]string
	trie     *TrieNode
}

// Definizione di un tipo piano
type piano = *Piano

// ALGORITMI ED OPERAZIONI

// Crea un piano vuoto (eliminando l'eventuale piano già esistente)
func newPiano() piano {
	return &Piano{
		automi:   make(map[string]Automa),
		ostacoli: make([]Ostacolo, 0),
		occupati: make(map[Punto]bool),
		posAutomi: make(map[Punto][]string),
		trie:     newTrieNode(),
	}
}

// Stampa le posizioni di tutti gli automi η tali che α è un prefisso di η, secondo quanto indicato nelle Specifiche di Implementazione.
func (p piano) posizioni(a string) {
	fmt.Println("(")
	for _, nome := range p.trie.searchPrefix(a) {
		automa := p.automi[nome]
		fmt.Printf("%s: %d,%d\n", nome, automa.p.x, automa.p.y)
	}
	fmt.Println(")")
}

// Stampa l'elenco degli automi seguito dall'elenco degli ostacoli
func (p piano) stampa() {
	fmt.Println("(")
	for nome, automa := range p.automi {
		fmt.Printf("%s: %d,%d\n", nome, automa.p.x, automa.p.y)
	}
	fmt.Println(")")
	fmt.Println("[")
	for _, ost := range p.ostacoli {
		fmt.Printf("(%d,%d)(%d,%d)\n", ost.p1.x, ost.p1.y, ost.p2.x, ost.p2.y)
	}
	fmt.Println("]")
}

// Stampa un carattere che indica cosa si trova nella posizione (x, y) del piano: A se è un automa, O se è un ostacolo, E se la posizione è vuota.
func (p piano) stato(x, y int) {
	point := Punto{x, y}
	if automi, ok := p.posAutomi[point]; ok && len(automi) > 0 {
		fmt.Println("A")
		return
	}
	if p.isOstacolo(x, y) {
		fmt.Println("O")
		return
	}
	fmt.Println("E")
}

// Verifica se un punto è occupato da un ostacolo
func (p piano) isOstacolo(x, y int) bool {
	return p.occupati[Punto{x, y}]
}

// Se il punto (x, y) è contenuto in un ostacolo, non esegue nessuna operazione. Altrimenti, se non esiste alcun automa di nome η,
// crea un nuovo automa di nome η e lo posiziona nel punto (x, y). Se η esiste già, lo riposizione nel punto (x, y).
func (p piano) automa(x, y int, nome string) {
	if p.isOstacolo(x, y) {
		return
	}
	if automa, ok := p.automi[nome]; ok {
		point := Punto{automa.p.x, automa.p.y}
		if automi, ok := p.posAutomi[point]; ok {
			slice := removeElement(automi, nome)
			if len(slice) == 0 {
				delete(p.posAutomi, point)
			} else {
				p.posAutomi[point] = slice
			}
		}
	}
	p.automi[nome] = Automa{Punto{x, y}}
	p.posAutomi[Punto{x, y}] = append(p.posAutomi[Punto{x, y}], nome)
	p.trie.insert(nome)
}

// Rimuove un elemento da una slice
func removeElement(slice []string, element string) []string {
	for i, el := range slice {
		if el == element {
			return append(slice[:i], slice[i+1:]...)
		}
	}
	return slice
}

// Se i punti nel rettangolo R(x0, y0, x1, y1) non contengono alcun automa, inserisce nel piano l'ostacolo rappresentato da
// R(x0, y0, x1, y1), altrimenti non esegue nessuna operazione. Si può assumere che x0 < x1 e y0 < y1.
func (p piano) ostacolo(x0, y0, x1, y1 int) {
	for _, automa := range p.automi {
		if automa.p.x >= x0 && automa.p.x <= x1 && automa.p.y >= y0 && automa.p.y <= y1 {
			return
		}
	}
	ost := Ostacolo{Punto{x0, y0}, Punto{x1, y1}}
	p.ostacoli = append(p.ostacoli, ost)
	for x := x0; x <= x1; x++ {
		for y := y0; y <= y1; y++ {
			p.occupati[Punto{x, y}] = true
		}
	}
}

// Calcola la distanza di Manhattan tra due punti
func manhattanDist(p1, p2 Punto) int {
	abs := func(x int) int {
		if x < 0 {
			return -x
		}
		return x
	}
	return abs(p1.x - p2.x) + abs(p1.y - p2.y)
}

// Calcola i possibili percorsi di un automa da un punto a un altro
// La funzione:
// - Calcola i possibili percorsi di un automa da un punto a un altro
// - Restituisce i percorsi ordinati per priorità
func (p piano) possibleRoads(s stateRoad, x, y, dist int) []stateRoad {
    totalMoves := []struct {
        ox, oy, d int
    }{
		{0, 1, 2}, //up
        {1, 0, 1}, //right
		{0, -1, 0}, //down
        {-1, 0, 3}, //left
    }
    var moves []orderedRoad
    for _, move := range totalMoves {
        ox := s.x + move.ox
        oy := s.y + move.oy
        if p.isOstacolo(ox, oy) {
            continue
        }
        if manhattanDist(Punto{ox, oy}, Punto{x, y}) != dist-(s.steps + 1) {
            continue
        }
        priority := -manhattanDist(Punto{ox, oy}, Punto{x, y})
        moves = append(moves, orderedRoad{move.ox, move.oy, move.d, priority})
    }
    sort.Slice(moves, func(i, j int) bool {
        return moves[i].priority < moves[j].priority
    })
    var roads []stateRoad
    for _, move := range moves {
        nx := s.x + move.ox
        ny := s.y + move.oy
        turns := s.turns
        if s.dir != -1 && s.dir != move.odir {
            turns++
        }
        road := stateRoad{x: nx, y: ny, steps: s.steps + 1, dir: move.odir, turns: turns}
        roads = append(roads, road)
    }
    return roads
}

// Funzione BFS che trova il percorso minimo tra due punti e restituisce true se esiste un percorso, false altrimenti
func (p piano) findPercorso(ax, ay, bx, by int) (bool, int) {
    if p.isOstacolo(ax, ay) || p.isOstacolo(bx, by) {
        return false, -1
    }
    d := manhattanDist(Punto{ax, ay}, Punto{bx, by})
    if d == 0 {
        return true, 0
    }
    queue := []stateRoad{{ax, ay, 0, -1, 0}}
    visited := make(map[visitedNodes]int)
    minTurns := -1
    for len(queue) > 0 {
        current := queue[0]
        queue = queue[1:]
        if current.x == bx && current.y == by {
            if current.steps == d {
                if minTurns == -1 || current.turns < minTurns {
                    minTurns = current.turns
                }
            }
            continue
        }
        if current.steps >= d {
            continue
        }
        roads := p.possibleRoads(current, bx, by, d)
        for _, road := range roads {
            key := visitedNodes{vX: road.x, vY: road.y, vDir: road.dir}
            if v, ok := visited[key]; ok && v <= road.turns {
                continue
            }
            visited[key] = road.turns
            queue = append(queue, road)
        }
    }
    if minTurns == -1 {
        return false, -1
    }
    return true, minTurns
}

// Manda il richiamo e sposta gli automi con il prefisso dato verso il punto (x, y)
func (p piano) richiamo(x, y int, prefix string) {
	var pq PriorityQueue
	for _, nome := range p.trie.searchPrefix(prefix) {
		a := p.automi[nome]
		d := manhattanDist(Punto{a.p.x, a.p.y}, Punto{x, y})
		if ok, _ := p.findPercorso(a.p.x, a.p.y, x, y); ok {
			heap.Push(&pq, DistanceRoad{name: nome, distance: d})
		}
	}
	if pq.Len() == 0 {
		return
	}
	mDistance := pq[0].distance
	for pq.Len() > 0 && pq[0].distance == mDistance {
		item := heap.Pop(&pq).(DistanceRoad)
		name := item.name
		if automa, ok := p.automi[name]; ok {
			point := Punto{automa.p.x, automa.p.y}
			if automi, ok := p.posAutomi[point]; ok {
				slice := removeElement(automi, name)
				if len(slice) == 0 {
					delete(p.posAutomi, point)
				} else {
					p.posAutomi[point] = slice
				}
			}
		}
		p.automi[name] = Automa{Punto{x, y}}
		p.posAutomi[Punto{x, y}] = append(p.posAutomi[Punto{x, y}], name)
	}
}

// Stampa SI se esiste almeno un percorso libero da P(η) a (x, y) di lunghezza D(P(η), (x, y)), NO altrimenti.
// Se η non esiste o se il punto (x, y) è contenuto in un ostacolo, stampa NO.
func (p piano) esistePercorso(x, y int, a string) {
	automa, exists := p.automi[a]
	if !exists || p.isOstacolo(x, y) {
		fmt.Println("NO")
		return
	}
	exists, _ = p.findPercorso(automa.p.x, automa.p.y, x, y)
	if exists {
		fmt.Println("SI")
	} else {
		fmt.Println("NO")
	}
}

// Stampa la tortuosità ovvero il numero minimo di cambi di direzione che un automa di nome η deve compiere per raggiungere il punto (x, y).
func (p piano) tortuosita(x, y int, a string) {
	automa, present := p.automi[a]
	if !present || p.isOstacolo(x, y) {
		fmt.Println(-1)
		return
	}
	exists, turns := p.findPercorso(automa.p.x, automa.p.y, x, y)
	if exists {
		fmt.Println(turns)
	} else {
		fmt.Println(-1)
	}
}

// Applica al piano rappprsentato da p l'operazione associata alla stringa s, secondo quanto indicato nelle Specifiche di Implementazione.
func esegui(p piano, s string) piano {
	words := strings.Split(s, " ")
	key := words[0]
	switch key {
	case "c":
		*p = *newPiano()
	case "S":
		p.stampa()
	case "s":
		x, _ := strconv.Atoi(words[1])
		y, _ := strconv.Atoi(words[2])
		p.stato(x, y)
	case "a":
		x, _ := strconv.Atoi(words[1])
		y, _ := strconv.Atoi(words[2])
		p.automa(x, y, words[3])
	case "o":
		x0, _ := strconv.Atoi(words[1])
		y0, _ := strconv.Atoi(words[2])
		x1, _ := strconv.Atoi(words[3])
		y1, _ := strconv.Atoi(words[4])
		p.ostacolo(x0, y0, x1, y1)
	case "p":
		p.posizioni(words[1])
	case "r":
		x, _ := strconv.Atoi(words[1])
		y, _ := strconv.Atoi(words[2])
		p.richiamo(x, y, words[3])
	case "e":
		x, _ := strconv.Atoi(words[1])
		y, _ := strconv.Atoi(words[2])
		p.esistePercorso(x, y, words[3])
	case "t":
		x, _ := strconv.Atoi(words[1])
		y, _ := strconv.Atoi(words[2])
		p.tortuosita(x, y, words[3])
	case "f":
		os.Exit(0)
	}
	return p
}

// Legge da standard input una sequenza di comandi e li esegue sul piano p
func main() {
	p := &Piano{}
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		esegui(p, line)
	}
}