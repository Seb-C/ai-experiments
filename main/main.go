package main

import "Seb-C/ai-experiments/genetic"

func main() {
	genetic.
		NewRandomPopulation().
		DoGenerations(100, 321).
		PrintResults()
}
