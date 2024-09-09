package main

import (
	"fmt"
	"math/rand"
	"time"
)

const (
	NJ = 5 // número de jogadores
	M  = 4 // número de cartas por jogador (exceto o que inicia com M+1)
)

type carta string // carta é representada por uma string

var (
	ch        [NJ]chan carta // canais de comunicação entre os jogadores
	deck      []carta        // baralho de cartas
	baterChan chan int       // canal para sinalizar quando alguém bate
)

func init() {
	rand.Seed(time.Now().UnixNano()) // inicializa o gerador aleatório
	baterChan = make(chan int, NJ)   // canal para sinalizar batida
}

// função auxiliar para criar um baralho de cartas
func criarBaralho() []carta {
	var cartas []carta
	for i := 0; i < M; i++ {
		for j := 0; j < NJ; j++ {
			cartas = append(cartas, carta(fmt.Sprintf("%c", 'A'+j))) // cria M cartas de cada tipo
		}
	}
	cartas = append(cartas, "@")                                                              // adiciona um coringa
	rand.Shuffle(len(cartas), func(i, j int) { cartas[i], cartas[j] = cartas[j], cartas[i] }) // embaralha as cartas
	return cartas
}

// verifica se o jogador tem todas as cartas iguais
func temCartasIguais(mao []carta) bool {
	for i := 1; i < len(mao); i++ {
		if mao[i] != mao[0] {
			return false
		}
	}
	return true
}

func jogador(id int, in chan carta, out chan carta, cartasIniciais []carta) {
	mao := cartasIniciais // estado local: cartas na mão
	nroDeCartas := len(mao)

	for {
		if nroDeCartas == M+1 { // jogador com M+1 cartas joga
			fmt.Printf("Jogador %d joga com %v\n", id, mao)
			// escolhe uma carta aleatória para passar adiante
			cartaParaSair := mao[rand.Intn(nroDeCartas)]
			out <- cartaParaSair // envia a carta para o próximo jogador
			// remove a carta da mão
			for i, c := range mao {
				if c == cartaParaSair {
					mao = append(mao[:i], mao[i+1:]...)
					break
				}
			}
			nroDeCartas--

		} else { // jogador com M cartas recebe uma nova carta
			cartaRecebida := <-in
			fmt.Printf("Jogador %d recebeu a carta %s\n", id, cartaRecebida)
			mao = append(mao, cartaRecebida)
			nroDeCartas++

			// Verifica se pode bater
			if temCartasIguais(mao) {
				fmt.Printf("Jogador %d bateu!\n", id)
				baterChan <- id
			}
		}
	}
}

func main() {
	// Criação do baralho e embaralhamento
	deck = criarBaralho()

	// Inicialização dos canais
	for i := 0; i < NJ; i++ {
		ch[i] = make(chan carta)
	}

	// Distribuição das cartas para os jogadores
	for i := 0; i < NJ; i++ {
		var cartasIniciais []carta
		if i == 0 {
			cartasIniciais = deck[:M+1] // o primeiro jogador começa com M+1 cartas
			deck = deck[M+1:]
		} else {
			cartasIniciais = deck[:M] // os demais começam com M cartas
			deck = deck[M:]
		}
		go jogador(i, ch[i], ch[(i+1)%NJ], cartasIniciais) // cria processos conectados circularmente
	}

	// Simulação do jogo
	idDorminhoco := -1
	for i := 0; i < NJ; i++ {
		id := <-baterChan
		fmt.Printf("Jogador %d bateu na posição %d!\n", id, i+1)
		if i == NJ-1 {
			idDorminhoco = id
		}
	}

	fmt.Printf("Jogador %d é o dorminhoco!\n", idDorminhoco)
}
