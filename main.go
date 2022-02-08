package main

import (
	"fmt"
	"sync"
)

var stop bool = false

func produkt(products chan<- int, max int, wg *sync.WaitGroup) {
	for i := 1; i <= max; i++ {
		products <- i
	}
	close(products)
	defer wg.Done()
}

func kobieta(id int, magazyn <-chan int, kosz chan<- int, typ string, wg *sync.WaitGroup) {
	for {
		if len(magazyn) < 1 {
			break
		}

		x := <-magazyn
		fmt.Printf("Kobieta %v wkłada do kosza %v %v \n", id, typ, x)

		kosz <- x
	}

	stop = true
	defer wg.Done()
}

func sluzba(id int, kosz_zniczy <-chan int, kosz_wiazanek <-chan int, wg *sync.WaitGroup) {
	for {
		if len(kosz_zniczy) < 2 || len(kosz_wiazanek) < 1 {
			if stop {
				break
			} else {
				continue
			}
		}

		znicz1 := <-kosz_zniczy
		znicz2 := <-kosz_zniczy
		wiazanka := <-kosz_wiazanek

		fmt.Printf("Donosiciel %v bierze znicz %v i %v oraz wiazanke %v\n", id, znicz1, znicz2, wiazanka)
	}

	defer wg.Done()
}

func main() {
	liczba_zniczy := 100
	liczba_wiazanek := 50

	liczba_babek := 2
	liczba_poslancow := 5

	magazyn_zniczy := make(chan int, liczba_zniczy)
	magazyn_wiazanek := make(chan int, liczba_wiazanek)

	wg1 := new(sync.WaitGroup)
	wg1.Add(2)

	go produkt(magazyn_zniczy, liczba_zniczy, wg1)
	go produkt(magazyn_wiazanek, liczba_wiazanek, wg1)

	wg1.Wait()

	kosz_zniczy := make(chan int, 10)
	kosz_wiazanek := make(chan int, 10)

	wg2 := new(sync.WaitGroup)
	wg2.Add(liczba_poslancow)
	wg2.Add(liczba_babek * 2)

	for i := 1; i <= liczba_babek; i++ {
		go kobieta(i, magazyn_zniczy, kosz_zniczy, "znicz", wg2)
		go kobieta(i, magazyn_wiazanek, kosz_wiazanek, "wiazanke", wg2)
	}

	for i := 1; i <= liczba_poslancow; i++ {
		go sluzba(i, kosz_zniczy, kosz_wiazanek, wg2)
	}

	wg2.Wait()
}
