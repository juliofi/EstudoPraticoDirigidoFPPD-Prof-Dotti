// por Fernando Dotti - fldotti.github.io - PUCRS - Escola Politécnica
// servidor com criacao dinamica de thread de servico
// Problema:
//   considere um servidor que recebe pedidos por um canal (representando uma conexao)
//   ao receber o pedido, sabe-se através de qual canal (conexao) responder ao cliente.
//   Abaixo uma solucao sequencial para o servidor.
// Exercicio
//   deseja-se tratar os clientes concorrentemente, e nao sequencialmente.
//   como ficaria a solucao ?

package main

import (
	"fmt"
	"math/rand"
)

const (
	NCL = 10
)

type Request struct {
	v      int
	ch_ret chan int
}

// ------------------------------------
// cliente
func cliente(i int, req chan Request) {
	var v, r int
	my_ch := make(chan int)
	for {
		v = rand.Intn(1000)
		req <- Request{v, my_ch}
		r = <-my_ch
		fmt.Println("cli: ", i, " req: ", v, "  resp:", r)
	}
}

// ------------------------------------
// servidor sequencial
func servidorSeq(in chan Request) {
	for {
		req := <-in
		fmt.Println("                       trataReq ", req)
		req.ch_ret <- req.v * 2 // responde  ao cliente
	}
}

// servidor concorrente
func servidorConc(in chan Request) {
	for {
		req := <-in
		go func(r Request) {
			fmt.Println("                       trataReq ", r)
			r.ch_ret <- r.v * 2 // responde ao cliente
		}(req) // Goroutine para tratar cada pedido de forma concorrente
	}
}

// ------------------------------------
// main
func main() {
	fmt.Println("------ Servidores Concorrente -------")
	serv_chan := make(chan Request) // cria o canal de comunicação com o servidor

	// Cria os clientes
	for i := 0; i < NCL; i++ {
		go cliente(i, serv_chan)
	}

	// Inicia o servidor concorrente
	servidorConc(serv_chan)
}
