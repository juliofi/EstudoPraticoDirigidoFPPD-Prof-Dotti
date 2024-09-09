// por Fernando Dotti - fldotti.github.io - PUCRS - Escola Politécnica
// Este é um exemplo muito simples de como modelar uma topologia com N nodos.
// A topologia é modelada por uma matriz de incidência.
// Basta adicionar 1s na matriz, na posição [i][j] para criar arestas direcionadas do nodo i para o nodo j.
// Nesta modelagem, cada nodo tem um canal de entrada.
// Com a função broadcast, um nodo manda para todos os vizinhos.
// j é dito vizinho de i se existe topology[i][j]=1
// ATENÇÃO:
//    1) A topologia criada no exemplo abaixo é um **grafo acíclico dirigido.**
//       Cada aresta tem somente uma direção (ex.: 0 manda para 1, mas 1 não manda para 0).
//       Assim, um nodo pode receber mais de uma vez uma mensagem,
//       mas nesta topologia, a mensagem não entra em ciclo.
// EXERCÍCIO:
//    1) Rode o exemplo. Note que cada mensagem é repassada mais de uma vez em alguns nodos.
//    2) Implemente a eliminação de duplicatas.
//       Cada nodo deve repassar a mensagem apenas a primeira vez que a recebe.

package main

import (
	"fmt"
)

// Número de nodos
const N = 10

// Topologia é uma matriz NxN onde 1 em [i][j] indica presença da aresta do nodo i para o j
type Topology [N][N]int

// O que é enviado entre nodos, pode adicionar campos nesta estrutura ...
type Message struct {
	id int // Identificador da mensagem - um número sequencial ...
}

// Um canal de entrada para cada nodo i
type inputChan [N]chan Message

// Estrutura que define um nodo
type nodeStruct struct {
	id   int
	topo Topology
	inCh inputChan
}

// Tamanho do buffer de cada canal de entrada
const channelBufferSize = 1

// Difusão ou broadcast - um nodo manda para TODOS os seus vizinhos no grafo
// O nodo `n`, conforme a topologia, usando canais do vetor `inCh`, manda mensagem para todos eles
func (n *nodeStruct) broadcast(m Message) {
	for j := 0; j < N; j++ { // Para todo vizinho j em N
		if n.topo[n.id][j] == 1 { // A matriz em [id][j] diz se o nodo `n` está conectado a `j`
			n.inCh[j] <- m // Envia a mensagem `m` para o canal de entrada do nodo `j`
		}
	}
}

// Função do nodo: lê seu canal de entrada e envia a mensagem para os nodos conectados
func (n *nodeStruct) nodo() {
	fmt.Println(n.id, "ativo!")
	for {
		m := <-n.inCh[n.id]              // Espera uma mensagem de entrada
		fmt.Println(n.id, "tratando", m) // Exibe a mensagem recebida
		n.broadcast(m)                   // Repassa a mensagem para todos os nodos conectados
	}
}

// ------------------------------------------------------------------------------------------------
// No main: montagem da topologia, criação de canais, inicialização de nodos e geração de mensagens
// ------------------------------------------------------------------------------------------------

func main() {
	var topo Topology
	// Se [i,j]==1, então o nodo i pode enviar para o nodo j pelo canal j.
	// Para alterar a topologia, basta adicionar 1s. Cada 1 é uma aresta direcional.
	// Para modelar comunicação em ambas direções entre i e j, então [i,j] e [j,i] devem ser 1.
	topo = [N][N]int{
		//  0  1  2  3  4  5  6  7  8  9       Aresta de    Para
		{0, 1, 0, 0, 0, 0, 0, 0, 0, 0}, // 0           0 -> 1
		{0, 0, 1, 0, 0, 0, 0, 0, 0, 0}, // 1           1 -> 2
		{0, 0, 0, 1, 0, 0, 0, 0, 0, 0}, // 2           2 -> 3
		{0, 0, 0, 0, 1, 0, 0, 0, 1, 0}, // 3           3 -> 4 e 3 -> 7
		{0, 0, 0, 0, 0, 1, 0, 0, 0, 1}, // 4           4 -> 5 e 4 -> 9
		{0, 0, 0, 0, 0, 0, 1, 0, 0, 0}, // 5           5 -> 6
		{0, 0, 0, 0, 0, 0, 0, 1, 0, 0}, // 6           6 -> 7
		{0, 0, 0, 0, 0, 0, 0, 0, 1, 0}, // 7           7 -> 8
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 1}, // 8           8 -> 9
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, // 9
	}

	var inCh inputChan // Cada nodo i tem um canal de entrada, chamado inCh[i]
	for i := 0; i < N; i++ {
		inCh[i] = make(chan Message, channelBufferSize) // Cria os canais de entrada
	}

	// Lança todos os nodos
	for id := 0; id < N; id++ {
		n := nodeStruct{id, topo, inCh}
		go n.nodo() // Cada nodo é executado como uma goroutine
	}

	// Gera mensagens de teste que serão "inundadas" na rede
	for i := 1; i < 2; i++ { // Gera uma mensagem de teste
		inCh[0] <- Message{i} // Envia uma mensagem a partir do nodo 0
		// time.Sleep(time.Second) // Pode adicionar delay, se necessário
	}

	<-make(chan struct{}) // Bloqueia a execução para manter o programa ativo
}
