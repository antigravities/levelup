package util

import "math"

// inspired by https://www.evanmiller.org/how-not-to-sort-by-average-rating.html

const confidence = 1.96

// RateNormal rates a suggestion normally
func RateNormal(positiveRatings, totalRatings int) float32 {
	return (float32(positiveRatings) / float32(totalRatings))
}

// RateWilson rates a suggestion using Wilson score magic
func RateWilson(positiveRatings, totalRatings float64) float64 {
	phat := 1 * positiveRatings / totalRatings
	z := confidence
	return (phat + confidence*confidence/(2*totalRatings) - z*math.Sqrt((phat*(1-phat)+confidence*confidence/(4*totalRatings))/totalRatings)) / (1 + confidence*confidence/totalRatings)
}
