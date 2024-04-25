package main

import (
	"fmt"
	"sync"
	"time"
)

const (
	maxCustomers    = 20
	waitingRoomSize = 16
	sofaSize        = 4
	numBarbers      = 3
)

// BarberShop struct representa a barbearia.
type BarberShop struct {
	mutex            sync.Mutex
	customers        int
	waitingRoom      chan int
	sofa             chan int
	barberChair      chan int
	barberPillow     chan int
	cash             chan int
	receipt          chan int
	waitingRoomCount int
	sofaCount        int
	doneCutting      chan int
}

// NovoBarberShop cria uma nova instância de uma barbearia.
func NovoBarberShop() *BarberShop {
	return &BarberShop{
		waitingRoom:      make(chan int, waitingRoomSize),
		sofa:             make(chan int, sofaSize),
		barberChair:      make(chan int, numBarbers),
		barberPillow:     make(chan int, numBarbers),
		cash:             make(chan int, 1),
		receipt:          make(chan int, 1),
		waitingRoomCount: 0,
		sofaCount:        0,
		doneCutting:      make(chan int, numBarbers),
	}
}

// customer representa um cliente na barbearia.
func (b *BarberShop) customer(id int, wg *sync.WaitGroup) {
	defer wg.Done()

	// Tenta entrar na sala de espera.
	b.mutex.Lock()
	if b.customers == maxCustomers {
		b.mutex.Unlock()
		fmt.Printf("Cliente %d encontrou a barbearia cheia e foi embora.\n", id)
		return
	}
	b.customers++
	b.waitingRoomCount++
	b.mutex.Unlock()

	// Entra na sala de espera.
	b.waitingRoom <- id
	fmt.Printf("Cliente %d entrou na sala de espera. Sala de espera: %d clientes.\n", id, b.waitingRoomCount)

	// Tenta sentar no sofá.
	<-b.waitingRoom
	b.mutex.Lock()
	b.waitingRoomCount--
	b.sofaCount++
	b.mutex.Unlock()
	b.sofa <- id
	fmt.Printf("Cliente %d sentou no sofá. Sofá: %d clientes.\n", id, b.sofaCount)

	// Tenta sentar na cadeira do barbeiro.
	<-b.sofa
	b.mutex.Lock()
	b.sofaCount--
	b.mutex.Unlock()
	b.barberChair <- id
	fmt.Printf("Cliente %d sentou na cadeira do barbeiro.\n", id)

	// Aguarda o barbeiro e obtém o corte de cabelo.
	<-b.barberChair
	b.barberPillow <- id
	<-b.doneCutting

	// Paga e aguarda o recibo.
	b.cash <- id
	<-b.receipt
	fmt.Printf("Cliente %d pagou e saiu da barbearia.\n", id)

	b.mutex.Lock()
	b.customers--
	b.mutex.Unlock()
}

// barber representa um barbeiro na barbearia.
func (b *BarberShop) barber(id int) {
	for {
		// Aguarda por um cliente.
		customerPos := <-b.barberPillow
		fmt.Printf("Barbeiro %d está cortando o cabelo do cliente %d.\n", id, customerPos)
		time.Sleep(time.Millisecond * 500) // Simula o corte de cabelo.
		b.doneCutting <- customerPos

		// Aguarda pelo pagamento.
		paymentPos := <-b.cash
		fmt.Printf("Barbeiro %d está recebendo o pagamento do cliente %d.\n", id, paymentPos)
		time.Sleep(time.Millisecond * 300) // Simula o recebimento do pagamento.
		b.receipt <- paymentPos
	}
}

func main() {
	barbershop := NovoBarberShop()
	var wg sync.WaitGroup

	// Inicia os barbeiros.
	for i := 0; i < numBarbers; i++ {
		go barbershop.barber(i + 1)
	}

	// Inicia os clientes.
	for i := 0; i < maxCustomers; i++ {
		wg.Add(1)
		go barbershop.customer(i+1, &wg)
		time.Sleep(time.Millisecond * 200) // Intervalo entre a chegada de novos clientes.
	}

	wg.Wait()
	fmt.Println("Todos os clientes foram atendidos.")
}
