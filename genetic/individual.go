package genetic

import "math"
import "math/rand"
import "github.com/Knetic/govaluate"

const GeneSize = 4
const GenesByIndividual = 9

const BreedingMutationProbability = float32(0.1)

var PossibleGenes = map[byte]string{
	0x0: "0",
	0x1: "1",
	0x2: "2",
	0x3: "3",
	0x4: "4",
	0x5: "5",
	0x6: "6",
	0x7: "7",
	0x8: "8",
	0x9: "9",
	0xA: "+",
	0xB: "-",
	0xC: "*",
	0xD: "/",
}

var AlternateWithGenes = []byte{0xA, 0xB, 0xC, 0xD} // "+", "-", "*", "/"

type Individual uint64

func NewRandomIndividual() *Individual {
	individual := Individual(rand.Uint64())
	return &individual
}

func (this *Individual) GetGenome() uint64 {
	return uint64(*this)
}

func (this *Individual) GetCompiledGenome() string {
	compiledGenome := ""
	previousWasSymbol := true
	for i := GenesByIndividual - 1; i >= 0; i-- {
		gene := byte((*this >> uint(GeneSize*i)) & 0xF)
		if geneValue, exists := PossibleGenes[gene]; exists {

			// Checking if it's a symbol
			isSymbol := false
			for j := 0; j < len(AlternateWithGenes); j++ {
				if AlternateWithGenes[j] == gene {
					isSymbol = true
					break
				}
			}

			// Adding char if it alternates digits with operators
			if previousWasSymbol != isSymbol && !(isSymbol && i == 0) {
				compiledGenome += geneValue
				previousWasSymbol = isSymbol
			}
		}
	}
	return compiledGenome
}

func (this *Individual) GetResult() int {
	compiledGenome := this.GetCompiledGenome()

	if compiledGenome == "" {
		return math.MinInt32
	}

	expression, err := govaluate.NewEvaluableExpression(compiledGenome)
	if err != nil {
		return math.MinInt32
	}

	result, err := expression.Evaluate(nil)
	if err != nil {
		return math.MinInt32
	}

	return int(math.Floor(result.(float64)))
}

func (this *Individual) GetFitnessScore(expectedResult int) float64 {
	return 1 / math.Abs(float64(expectedResult-this.GetResult()))
}

func (individualA *Individual) Breed(individualB *Individual) *Individual {
	genomeSize := uint(GeneSize * GenesByIndividual)

	crossOverBit := uint(rand.Intn(int(genomeSize)))
	genomeASliceLength := crossOverBit
	genomeBSliceLength := genomeSize - crossOverBit
	bitsBeforeGenomeBSlice := (64 - genomeSize) + genomeASliceLength

	bitMaskA := uint64((0xFFFFFFFFFFFFFFFF >> genomeBSliceLength) << genomeBSliceLength)
	bitMaskB := uint64((0xFFFFFFFFFFFFFFFF << bitsBeforeGenomeBSlice) >> bitsBeforeGenomeBSlice)

	newGenome := (individualA.GetGenome() & bitMaskA) | (individualB.GetGenome() & bitMaskB)

	if rand.Float32() < BreedingMutationProbability {
		bitToMutate := uint(rand.Intn(int(genomeSize)))
		newGenome ^= uint64(1) << bitToMutate
	}

	newIndividual := Individual(newGenome)
	return &newIndividual
}
