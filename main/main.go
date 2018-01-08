package main

import "Seb-C/genetic-algorithm-experiment/genetic"

func main() {
	genetic.
		NewRandomPopulation().
		DoGenerations(100, 321).
		PrintResults()
}
