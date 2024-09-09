package main

import (
	"fmt"
)

// Número de nós
const N = 10

// Tamanho do buffer de cada canal de entrada
const channelBufferSize = 5

// Topologia é uma matriz NxN onde 1 em [i][j] indica a presença da aresta do nó i para o j
type Topology [N][N]int

// Estrutura da mensagem
type Message struct {
	id       int    // Identificador da mensagem
	source   int    // Nó de origem
	receiver int    // Nó de destino
	data     string // Dados da mensagem
	route    []int  // Rota percorrida (pilha de nós)
}

// Um canal de entrada para cada nó
type inputChan [N]chan Message

// Estrutura que define cada nó
type nodeStruct struct {
	id               int
	topo             Topology
	inCh             inputChan
	received         map[int]bool // Mapeia se a mensagem já foi recebida
	receivedMessages []Message    // Mensagens recebidas pelo destino
}

// Difusão ou broadcast - Um nó envia a mensagem para TODOS seus vizinhos
func (n *nodeStruct) broadcast(m Message) {
	for j := 0; j < N; j++ { // Para todo vizinho j em N
		if n.topo[n.id][j] == 1 { // Se há conexão entre n.id e j
			n.inCh[j] <- m // Envia a mensagem para o nó j
		}
	}
}

// Função para retornar a mensagem pela rota inversa
func (n *nodeStruct) retornaMensagem(m Message) {
	// Desempilha o próximo nó na rota
	if len(m.route) > 0 {
		next := m.route[len(m.route)-1]
		m.route = m.route[:len(m.route)-1] // Remove o último da pilha
		fmt.Printf("   %d retornando para %d (msg %d)\n", n.id, next, m.id)
		n.inCh[next] <- m // Envia a mensagem de volta ao nó anterior
	}
}

// Cada nó recebe a matriz de conectividade e os canais de entrada de todos os processos
// Cada nó lê o seu canal de entrada, e repassa a mensagem para os vizinhos ou a retorna pela rota inversa
func (n *nodeStruct) nodo() {
	fmt.Println(n.id, " ativo!")
	for {
		m := <-n.inCh[n.id] // Recebe uma mensagem

		// Caso a mensagem seja para este nó
		if m.receiver == n.id {
			n.receivedMessages = append(n.receivedMessages, m)
			fmt.Printf("                                   %d recebeu de %d msg %d: %s\n", n.id, m.source, m.id, m.data)

			// Se a mensagem é de ida, envia uma resposta
			if m.id > 0 {
				fmt.Printf("                                   %d enviando resposta para %d\n", n.id, m.source)
				// Inverte a rota para a resposta
				go n.retornaMensagem(Message{-m.id, n.id, m.source, "resposta", m.route})
			}
		} else { // Não é para mim, repassa ou responde
			// Verifica se a mensagem já foi recebida
			if _, achou := n.received[m.id]; !achou {
				// Adiciona o nó atual na rota (empilha)
				m.route = append(m.route, n.id)
				fmt.Printf("%d repassa msg %d de %d para %d\n", n.id, m.id, m.source, m.receiver)
				n.received[m.id] = true // Marca como recebida
				go n.broadcast(m)       // Repassa a mensagem
			}
		}
	}
}

// Função que gera carga de mensagens
func carga(nodoInicial int, inCh chan Message) {
	for i := 1; i < 2; i++ {
		inCh <- Message{i, nodoInicial, 9, "requisição", []int{}} // Mensagem com rota vazia
	}
}

// ------------------------------------------------------------------------------------------------
// no main: montagem da topologia, criação de canais, inicialização de nós e geração de mensagens
// ------------------------------------------------------------------------------------------------
func main() {
	var topo Topology
	// Topologia: bidirecional
	topo = [N][N]int{
		{0, 1, 0, 0, 0, 0, 0, 0, 0, 0}, // 0 - 1
		{1, 0, 1, 0, 0, 0, 0, 0, 0, 0}, // 1 - 2
		{0, 1, 0, 1, 0, 0, 0, 0, 0, 0}, // 2 - 3
		{0, 0, 1, 0, 1, 0, 0, 0, 0, 0}, // 3 - 4
		{0, 0, 0, 1, 0, 1, 0, 0, 0, 1}, // 4 - 5 e 4 - 9
		{0, 0, 0, 0, 1, 0, 1, 0, 0, 0}, // 5 - 6
		{0, 0, 0, 0, 0, 1, 0, 1, 0, 0}, // 6 - 7
		{0, 0, 0, 0, 0, 0, 1, 0, 1, 0}, // 7 - 8
		{0, 0, 0, 0, 0, 0, 0, 1, 0, 1}, // 8 - 9
		{0, 0, 0, 0, 0, 1, 0, 0, 1, 0}, // 9
	}

	var inCh inputChan
	for i := 0; i < N; i++ {
		inCh[i] = make(chan Message, channelBufferSize)
	}

	// Lança todos os nós
	for id := 0; id < N; id++ {
		n := nodeStruct{id, topo, inCh, make(map[int]bool), []Message{}}
		go n.nodo() // Lança cada nó como uma goroutine
	}

	// Gera a carga inicial de mensagens
	go carga(0, inCh[0]) // Envia a partir do nó 0

	<-make(chan struct{}) // Bloqueia para manter o programa rodando
}
