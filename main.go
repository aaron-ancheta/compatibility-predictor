package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"os"
	"reflect"
)

// A struct to store the input json
type InputJSON struct {
	Team       []Person `json:"team"`
	Applicants []Person `json:"applicants"`
}

// A struct to store each person's name and attributes
// either a team member or an applicant
type Person struct {
	Name       string     `json:"name"`
	Attributes Attributes `json:"attributes"`
}

// A struct to store the attribute name and attribute scores
type Attributes struct {
	Intelligence       float64 `json:"intelligence"`
	Strength           float64 `json:"strength"`
	Endurance          float64 `json:"endurance"`
	SpicyFoodTolerance float64 `json:"spicyFoodTolerance"`
}

// A struct to store the name and compatibility scores of each applicant
type ScoredApplicant struct {
	Name  string  `json:"name"`
	Score float64 `json:"score"`
}

// A struct that stores the output JSON
type OutputJSON struct {
	ScoredApplicants []ScoredApplicant `json:"scoredApplicants"`
}

func main() {
	// Initialize the struct to store the input json
	var input InputJSON
	data, err := os.ReadFile("input.json")
	if err != nil {
		log.Fatal("Error reading the input file:", err)
	}
	err = json.Unmarshal(data, &input)
	if err != nil {
		log.Fatal("Error unmarshaling the input JSON:", err)
	}

	// Normalize the values of attributes for each team members
	for i := range input.Team {
		NormalizeInput(&input.Team[i].Attributes, 10)
	}
	// Normalize the values of attributes for each applicants
	for i := range input.Applicants {
		NormalizeInput(&input.Applicants[i].Attributes, 10)
	}

	// Create an Attribute that contains the centroid for each attributes of the existing team members.
	centroid := TeamCentroid(input)
	// Print the centroid
	fmt.Printf("Centroid: %+.3v\n", centroid)

	fmt.Println()

	// Create a map to store euclidean distance for easier retrieval
	// keys: applicant's name (type -> strings)
	// values: Euclidean distance of applicant's attributes using the centroid values calculated above (type -> float64)
	euclideanDistance := make(map[string]float64)

	// Calculate Euclidean distance of each applicant's attributes and store them in the map euclideanDistance
	for _, applicant := range input.Applicants {
		euclideanDistance[applicant.Name] = EuclideanDistance(centroid, applicant.Attributes)
	}

	fmt.Println("Euclidean Distance of each applicants:")
	// Print out the euclideanDistance map
	for applicant, euclideanDistance := range euclideanDistance {
		fmt.Printf("Name: %s | Euclidean Distance: %.3v\n", applicant, euclideanDistance)
	}

	fmt.Println()
	fmt.Println("Average Fit Score (using Euclidean Distance): ")

	for applicant, euclideanDistance := range euclideanDistance {
		fmt.Printf("Applicant Name: %s   | Average fit score: %.3v\n", applicant, AverageFit(euclideanDistance))
	}

	fmt.Println()

	// Using the Team's centroid calculated above, calculate a gap vector
	// gap vector = [(gap for intelligence), (gap for strenght), (gap for endurance), (gap for spicy food tolerance)]
	gapVector := GapVector(centroid)

	fmt.Printf("Gap vector: %.3v\n", gapVector)

	// Calculate the total gap (add all the gaps)
	totalGap := 0.0
	for _, gap := range gapVector {
		totalGap += gap
	}

	fmt.Printf("Total Gap: %.3f\n", totalGap)

	fmt.Println()
	// Create a map that will store the contribution of applicants, which is calculated using the gap vector above
	// key = applicant's name (type -> string)
	// value = applicant's contribution (type -> float64)
	contribution := make(map[string]float64)
	// Calculate the each applicant's contribution based on the gap vector from the existing team
	for _, applicant := range input.Applicants {
		contribution[applicant.Name] = ApplicantContribution(gapVector, applicant.Attributes)
	}
	fmt.Println("Applicant's contribution:")
	// Print out the content of the contribution map
	for applicant, contribution := range contribution {
		fmt.Printf("Name: %s, Contribution: %.3v\n", applicant, contribution)
	}
	fmt.Println()
	fmt.Println("Gap Scores of Applicants:")

	// Print out the gap score of each applicants
	// gap scores = contribution / total gap
	for name, contribution := range contribution {
		fmt.Printf("Name of Applicant: %s  |  Gap score: %.3f\n", name, contribution/totalGap)
	}

	fmt.Println("Applicant Scores (using Average Fit and Gap Scores):")

	// Create a map to store applicant's overall score to
	applicantScore := make(map[string]float64)

	for _, applicant := range input.Applicants {
		// Get the euclidean distance of the applicant from the euclidean distance map above
		distance := euclideanDistance[applicant.Name]
		// Calculate the average fit score of the current applicant using the euclidean distance of that applicant
		averageFit := AverageFit(distance)

		// Calculate the gap filler score of the applicant using the the applicant's contribution from the contribution map above
		gapFillerScore := contribution[applicant.Name] / totalGap

		// Calculate the overall score of the current applicant using the average fit and gap filler score of the current applicant
		applicantScore[applicant.Name] = ApplicantScore(averageFit, gapFillerScore)
	}
	fmt.Println()
	// Print the final applicant scores
	for name, applicantScore := range applicantScore {
		fmt.Printf("Name: %s  |  Applicant Score: %.3f\n", name, applicantScore)
	}

	// Convert the applicantScore map to a slice of ScoredApplicant structs
	var scoredApplicants []ScoredApplicant
	for name, score := range applicantScore {
		scoredApplicants = append(scoredApplicants, ScoredApplicant{
			Name:  name,
			Score: score,
		})
	}

	// Create the output structure
	output := OutputJSON{
		ScoredApplicants: scoredApplicants,
	}

	// Marshal the output as JSON
	outputData, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		log.Fatal("Error marshaling output JSON:", err)
	}

	// Write to a file
	err = os.WriteFile("scored_applicants.json", outputData, 0644)
	if err != nil {
		log.Fatal("Error writing JSON file:", err)
	}

}

// TeamCentroid returns the overall mean of all the attributes of the existing team members
// The team's centroid returned by the function is in the form of Attribute struct
func TeamCentroid(data InputJSON) Attributes {
	var centroid Attributes // Stores the average of each attributes of the whole team
	count := len(data.Team) // Number of team members

	// Get the sum of each team member's attributes per category
	for _, member := range data.Team {
		centroid.Intelligence += member.Attributes.Intelligence
		centroid.Strength += member.Attributes.Strength
		centroid.Endurance += member.Attributes.Endurance
		centroid.SpicyFoodTolerance += member.Attributes.SpicyFoodTolerance
	}
	// return the average of the team's individual attributes
	return Attributes{
		Intelligence:       centroid.Intelligence / float64(count),
		Strength:           centroid.Strength / float64(count),
		Endurance:          centroid.Endurance / float64(count),
		SpicyFoodTolerance: centroid.SpicyFoodTolerance / float64(count),
	}
}

// EuclideanDistance returns the euclidean distance of the applicant's attributes and the team's centroid for each attributes
// the euclidean distance is in the form of float64
func EuclideanDistance(centroid, applicant Attributes) float64 {
	// return the square root of the square of team's centroid - applicant's attribute
	return math.Sqrt(
		math.Pow(float64(centroid.Intelligence-applicant.Intelligence), 2) +
			math.Pow(float64(centroid.Strength-applicant.Strength), 2) +
			math.Pow(float64(centroid.Endurance-applicant.Endurance), 2) +
			math.Pow(float64(centroid.SpicyFoodTolerance-applicant.SpicyFoodTolerance), 2))
}

// AverageFit returns the overall "vibe" score or how fit the applicant is to the team based on their attributes' euclidean distance
func AverageFit(distance float64) float64 {
	return 1 - (distance / math.Sqrt(float64(reflect.TypeOf(Attributes{}).NumField())))
}

// NormalizeInput normalizes the input attribute values for easier computations
// it also makes returning applicant scores within [0,1] range easier
func NormalizeInput(attribute *Attributes, divisor float64) {
	attribute.Intelligence = float64(attribute.Intelligence) / divisor
	attribute.Strength = float64(attribute.Strength) / divisor
	attribute.Endurance = float64(attribute.Endurance) / divisor
	attribute.SpicyFoodTolerance = float64(attribute.SpicyFoodTolerance) / divisor
}

// GapVector calculates the gap of the team's centroid, 1 being perfect
// Necessary to calculate the Gap Filler Score of the applicants
func GapVector(centroid Attributes) []float64 {
	return []float64{math.Max(0, 1-centroid.Intelligence), math.Max(0, 1-centroid.Strength), math.Max(0, 1-centroid.Endurance), math.Max(0, 1-centroid.SpicyFoodTolerance)}
}

// ApplicantContribution calculates the overall contribution of the applicant relative to the gap vector of the existing team
// returns the "how much can this applicant fill in the team's gap"
func ApplicantContribution(gapVector []float64, applicant Attributes) float64 {
	applicantAttributes := []float64{
		applicant.Intelligence,
		applicant.Strength,
		applicant.Endurance,
		applicant.SpicyFoodTolerance,
	}

	var sum float64
	for i := range gapVector {
		sum += (gapVector[i] * applicantAttributes[i])
	}
	return sum
}

// ApplicantScore returns the overall compatibility score of the applicant based on their Average fit score (vibe)
// and their Gap filler Score (what can they bring on the table)
// the score is weighted, with 80% on the vibe score and 20% on the contribution score (which can be adjusted based on hiring preference)
func ApplicantScore(averageFit, gapFillerScore float64) float64 {
	return (0.80 * averageFit) + (0.20 * gapFillerScore)
}
