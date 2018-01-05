package genetic

import "sort"
import "math/rand"
import "math"
import "fmt"
import "encoding/json"

const PopulationSize = 1000

const ReplacementRateByGeneration = 0.25

type Population []*Individual

type computedIndividual struct {
	individual        *Individual
	fitnessScore      float64
	normalizedFitness float64
}

func NewRandomPopulation() *Population {
	population := Population(make([]*Individual, 0, PopulationSize))
	for i := 0; i < PopulationSize; i++ {
		population = append(population, NewRandomIndividual())
	}

	return &population
}

func (this *Population) GetIndividuals() []*Individual {
	return []*Individual(*this)
}

func (this *Population) preComputePopulationWithFitness(expectedResult int) []*computedIndividual {
	computedPopulation := make([]*computedIndividual, 0, PopulationSize)

	// Computing fitness scores
	var globalFitness float64 = 0
	for i := 0; i < PopulationSize; i++ {
		computedPopulation = append(computedPopulation, &computedIndividual{
			individual:   (*this)[i],
			fitnessScore: (*this)[i].GetFitnessScore(expectedResult),
		})

		globalFitness += computedPopulation[i].fitnessScore
	}

	// Sorting population by score
	sort.Slice(computedPopulation, func(i int, j int) bool {
		return computedPopulation[i].fitnessScore > computedPopulation[j].fitnessScore
	})

	// Calculating normalized fitness score
	var previousFitness float64 = 0
	for i := 0; i < PopulationSize; i++ {
		relativeFitness := 1 - computedPopulation[i].fitnessScore/globalFitness
		computedPopulation[i].normalizedFitness = previousFitness
		previousFitness += relativeFitness
	}

	return computedPopulation
}

func (this *Population) NextGeneration(expectedResult int) *Population {
	computedPopulation := this.preComputePopulationWithFitness(expectedResult)

	newPopulation := Population{}

	// Creating new individuals
	newIndividuals := int(math.Floor(ReplacementRateByGeneration * PopulationSize))
	for i := 0; i < newIndividuals; i++ {
		var selectedIndividualA, selectedIndividualB *Individual

		targetFitnessA := rand.Float64()
		for j := 0; j < PopulationSize && computedPopulation[j].normalizedFitness < targetFitnessA; j++ {
			selectedIndividualA = computedPopulation[j].individual
		}

		targetFitnessB := rand.Float64()
		for j := 0; j < PopulationSize && computedPopulation[j].normalizedFitness < targetFitnessB; j++ {
			selectedIndividualB = computedPopulation[j].individual
		}

		newPopulation = append(newPopulation, selectedIndividualA.Breed(selectedIndividualB))
	}

	// Adding best individuals from the previous generation
	for i := 0; i < PopulationSize-newIndividuals; i++ {
		newPopulation = append(newPopulation, computedPopulation[i].individual)
	}

	return &newPopulation
}

func (this *Population) DoGenerations(generationsCount, expectedResult int) *Population {
	population := this
	for i := 0; i < generationsCount; i++ {
		population = population.NextGeneration(expectedResult)
	}
	return population
}

func (this *Population) PrintResults() {
	uniqueResults := make(map[string]int)
	for _, individual := range this.GetIndividuals() {
		uniqueResults[individual.GetCompiledGenome()] = individual.GetResult()
	}

	formattedMap, _ := json.MarshalIndent(uniqueResults, "", "	")
	fmt.Println(string(formattedMap))
}
