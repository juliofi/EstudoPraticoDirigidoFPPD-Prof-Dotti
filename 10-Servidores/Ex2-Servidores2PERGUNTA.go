// por Fernando Dotti - fldotti.github.io - PUCRS - Escola Politécnica
// servidor com criacao dinamica de thread de servico
// Problema:
//   considere um servidor que recebe pedidos por um canal (representando uma conexao)
//   ao receber o pedido, sabe-se através de qual canal (conexao) responder ao cliente.
//   Abaixo uma solucao sequencial para o servidor.
// Exercicio
//   deseja-se tratar os clientes concorrentemente, e nao sequencialmente.
//   como ficaria a solucao ?
// Veja abaixo a resposta ...
//   quantos clientes podem estar sendo tratados concorrentemente ?
//
// Exercicio:
//   agora suponha que o seu servidor pode estar tratando no maximo 10 clientes concorrentemente.
//   como voce faria ?
//

package main

import (
	"fmt"
	"math/rand"
)

const (
	NCL  = 100
	Pool = 10
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
// servidor
// thread de servico calcula a resposta e manda direto pelo canal de retorno informado pelo cliente
func trataReq(id int, req Request) {
	fmt.Println("                                 trataReq ", id)
	req.ch_ret <- req.v * 2
}

// servidor que dispara threads de servico
func servidorConc(in chan Request) {
	// servidor fica em loop eterno recebendo pedidos e criando um processo concorrente para tratar cada pedido
	var j int = 0
	for {
		j++
		req := <-in
		go trataReq(j, req)
	}
}

// servidor que limita o número de threads ativas
func servidorLimitado(in chan Request, limit int) {
	sem := make(chan struct{}, limit) // semáforo para limitar o número de goroutines

	for {
		req := <-in
		sem <- struct{}{} // bloqueia se o número de threads ativas for igual ao limite
		go func(r Request) {
			fmt.Println("                                 trataReq")
			r.ch_ret <- r.v * 2
			<-sem // libera uma thread
		}(req)
	}
}

// ------------------------------------
// main
func main() {
	fmt.Println("------ Servidor com Limite de Threads -------")
	serv_chan := make(chan Request)
	go servidorLimitado(serv_chan, Pool) // Limite de 10 threads simultâneas
	for i := 0; i < NCL; i++ {
		go cliente(i, serv_chan)
	}
	<-make(chan int) // Espera indefinidamente
}
