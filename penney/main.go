// penney
package main

import (
	"fmt"
	"math/rand"
)

// Created: Fri Sep 29 00:27:06 2023

func main() {
	player1 := []string{"HHH", "HHT", "HTH", "HTT", "THH", "THT", "TTH", "TTT"}
	player2 := make([]string, 8)
	copy(player2, player1)
	for i, g := range player2 {
		c01 := g[:2]
		c1 := g[1:2]
		var r1 string
		if c1 == "H" {
			r1 = "T"
		} else {
			r1 = "H"
		}

		player2[i] = r1 + c01
	}
	fmt.Println(player1)
	fmt.Println(player2)

	count := []int{0, 0, 0, 0, 0, 0, 0, 0}
	ht := []string{"H", "T"}
	triplet := ""
	var player1Rslt, player2Rslt [8]struct {
		wins, flips, maxFlips int
	}
	for i := 0; i < 1000; i++ {
		triplet += ht[rand.Intn(2)]
		if len(triplet) > 3 {
			triplet = triplet[1:4]
		}
		for ci, cv := range count {
			cv++
			count[ci] = cv
			if cv >= 3 {
				if triplet == player1[ci] {
					player1Rslt[ci].wins++
					player1Rslt[ci].flips = cv
					if cv > player1Rslt[ci].maxFlips {
						player1Rslt[ci].maxFlips = cv
					}
					count[ci] = 0
				} else if triplet == player2[ci] {
					player2Rslt[ci].wins++
					player2Rslt[ci].flips = cv
					if cv > player2Rslt[ci].maxFlips {
						player2Rslt[ci].maxFlips = cv
					}
					count[ci] = 0
				}
			}
		}
	}
	for i, rslt1 := range player1Rslt {
		rslt2 := player2Rslt[i]
		tot := rslt1.wins + player2Rslt[i].wins
		fmt.Printf("%s: %4d %6.2f%%\n",
			player1[i], rslt1.wins, 100*float64(rslt1.wins)/float64(tot))
		fmt.Printf("%s: %4d %6.2f%%\n",
			player2[i], rslt2.wins, 100*float64(rslt2.wins)/float64(tot))
	}
}
