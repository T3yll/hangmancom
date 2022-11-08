package hangcom

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"
)

func hangman() {
	//Déclaration des variables / Lecture des fichiers
	words, err := ioutil.ReadFile("words.txt")
	save, err := ioutil.ReadFile("save.txt")
	if err != nil {
		log.Fatal(err)
	}
	var tabSave []string = Initialisation(save)
	randomWord := ""
	wordUnder := ""
	letterInput := ""
	wordInput := ""
	randomWordRepair := ""
	var answPlayer []string
	var saveTab []string
	alreadyTyped := false
	again := false
	cpt := 0
	tempS := ""
	tempArgs := ""
	life := 9
	boolFail := true
	startFromSave := false
	//Initialisation grace aux différentes méthodes
	var tabS []string = Initialisation(words)
	randomWord = RandomPick(tabS)
	randomWordRepair = Repair(randomWord)
	wordUnder = Reveal(randomWord)
	tabBWordUnderscore := []byte(wordUnder)
	tabB := []byte(randomWord)
	//Détection si lancement à partir d'une sauvegarde si oui on print différement
	if len(os.Args) > 1 {
		if os.Args[1] == "--startWith" {
			startFromSave = true
			for i := 0; i < len(os.Args[1]); i++ {
				tempArgs = tempArgs + string(os.Args[1][i])
			}
		}
	}
	if startFromSave == false {
		fmt.Print("\n")
		fmt.Print("Good Luck, you have 10 attempts.\n")
		fmt.Print("\n")
		for i := 0; i < len(tabBWordUnderscore); i++ {
			fmt.Print(string(tabBWordUnderscore[i]))
			fmt.Print(" ")
		}
		fmt.Print("\n")
		fmt.Print("\n")
	}
	//Lancement de la partie
	for life < 10 {
		//Détection de la vie si plus de vie GAME OVER
		if life == -1 || life < -1 {
			fmt.Print("Game Over\n")
			fmt.Print("\n")
			break
		}
		//Détection si lancement à partir d'une sauvegarde si oui on initialise les variables avec celles de la sauvegarde précédente
		if len(os.Args) > 1 {
			if tempArgs == "--startWith" {
				tempS = string(save[0])
				life, err = strconv.Atoi(tempS)
				if err != nil {
					log.Fatal(err)
				}
				tempS = string(save[1])
				tempWord := ""
				for i := 0; i < len(save); i++ {
					tempWord = tempWord + string(save[i])
					if save[i] == 10 {
						saveTab = append(saveTab, tempWord)
						tempWord = ""
					}
				}
				tempS = string(saveTab[1])
				cpt, err = strconv.Atoi(tempS)
				tempArgs = ""
				randomWord = tabSave[2]
				wordUnder = tabSave[3]
				tabBWordUnderscore = []byte(wordUnder)
				randomWordRepair = Repair(randomWord)
				tabB = []byte(randomWord)
			}
		}
		//Même chose ici puis print du good luck
		if startFromSave == true {
			tempS = string(saveTab[1])
			cpt, err = strconv.Atoi(tempS)
			tempArgs = ""
			randomWord = tabSave[2]
			wordUnder = tabSave[3]
			tabBWordUnderscore = []byte(wordUnder)
			tabB = []byte(randomWord)
			fmt.Print("\n")
			fmt.Printf("Good Luck, you have %v attempts.\n", life+1)
			fmt.Print("\n")
			for i := 0; i < len(tabBWordUnderscore); i++ {
				fmt.Print(string(tabBWordUnderscore[i]))
				fmt.Print(" ")
			}
			fmt.Print("\n")
			fmt.Print("\n")
			fmt.Print("\n")
			startFromSave = false
		}
		//Ici c'est la partie détection de l'input du joueur
		fmt.Print("Choose: ")
		fmt.Scan(&letterInput)
		fmt.Print("\n")
		wordInput = letterInput
		//Si le joueur écrit STOP la partie s'arrête et les paramètres suivants sont sauvegardés
		if wordInput == "STOP" {
			fmt.Print("\n")
			fmt.Print("Game Saved in save.txt.")
			fmt.Print("\n")
			file, err := os.Create("save.txt")
			if err != nil {
				log.Fatal(err)
			}
			defer file.Close()
			//Life
			_, err = file.WriteString(fmt.Sprintf("%d\n", life))
			if err != nil {
				fmt.Printf("error writing string: %v", err)
			}
			//Cpt
			_, err = file.WriteString(fmt.Sprintf("%d\n", cpt))
			if err != nil {
				fmt.Printf("error writing string: %v", err)
			}
			//RandomWord
			file.WriteString(randomWord)
			//tabB
			file.Write(tabBWordUnderscore)
			file.WriteString("\n")
			file.WriteString(randomWordRepair)
			break
		}
		//On check si le joueur à mit une bonne réponse
		boolFail = true
		//answPlayer correspond à la liste des reponses du joueur
		answPlayer = append(answPlayer, wordInput)
		//Ici nous sommes dans le cas ou le joueur répond par une seule lettre
		if len(letterInput) < 2 {
			for i := 0; i < len(answPlayer)-1; i++ {
				if wordInput == answPlayer[i] {
					alreadyTyped = true
				}
			}
			for i := 0; i < len(randomWord)-1; i++ {
				if alreadyTyped == true {
					fmt.Print("Already tried, please retry !")
					fmt.Print("\n")
					alreadyTyped = false
					boolFail = false
				}
				//ici le joueur a bien répondu la lettre correspond
				if letterInput == string(tabB[i]) {
					tabBWordUnderscore[i] = tabB[i]
					boolFail = false
				}
			}
			if boolFail == true {
				if alreadyTyped == false {
					//ici nous sommes dans le cas ou le joueur a effectué une mauvaise réponse qu'il n'avait encore jamais faites
					fmt.Printf("Not present in the word, %v attempts remaining", life)
					fmt.Print("\n")
					HangmanPositions(life)
					life--
				} else {
					fmt.Print("Already tried, please retry !")
					fmt.Print("\n")
					alreadyTyped = false
				}
			}
		} else {
			//Ici nous sommes dans le cas ou le joueur répond un miot
			for i := 0; i < len(answPlayer)-1; i++ {
				if wordInput == answPlayer[i] {
					alreadyTyped = true
				}
			}
			if wordInput != randomWordRepair {
				if alreadyTyped == false {
					fmt.Print("\n")
					fmt.Printf("Not present in the word, %v attempts remaining", life-1)
					fmt.Print("\n")
					HangmanPositions(life - 1)
					life = life - 2
				} else {
					fmt.Print("Already tried, please retry !")
					fmt.Print("\n")
					alreadyTyped = false
				}
			} else {
				//Ici nous sommes dans le cas ou le joueur va répondre un mot juste donc la partie se termine
				if tabBWordUnderscore[len(tabBWordUnderscore)-1] == 10 {
					for i := 0; i < len(tabBWordUnderscore)-1; i++ {
						fmt.Print(string(randomWordRepair[i]))
						fmt.Print(" ")
					}
				} else {
					for i := 0; i < len(tabBWordUnderscore); i++ {
						fmt.Print(string(randomWordRepair[i]))
						fmt.Print(" ")
					}
				}
				fmt.Print("\n")
				fmt.Print("\n")
				fmt.Print("Congrats !\n")
				fmt.Print("\n")
				break
			}
		}
		//Ici nous sommes dans la partie gestion affichage et update avant le prochain tour
		//La on check si il reste des underscore dans le mot donc si le jeu continu
		for j := 0; j < len(randomWord)-1; j++ {
			if string(tabBWordUnderscore[j]) != "_" {
				cpt++
				again = true
			}
			if cpt == len(randomWord)-1 {
				again = false
				break
			}
		}
		cpt = 0
		//Si il continue alors on print le mot avec les underscore
		if again == true {
			fmt.Print("\n")
			for i := 0; i < len(tabBWordUnderscore); i++ {
				fmt.Print(string(tabBWordUnderscore[i]))
				fmt.Print(" ")
			}
			fmt.Print("\n")
			fmt.Print("\n")
		} else {
			//Sinon on print le mot complété et la partie se termine
			for i := 0; i < len(tabBWordUnderscore); i++ {
				fmt.Print(string(tabBWordUnderscore[i]))
				fmt.Print(" ")
			}
			fmt.Print("\n")
			fmt.Print("\n")
			fmt.Print("Congrats !\n")
			fmt.Print("\n")
			break
		}
	}
}

// Cette fonction permet l'initialisation des mots dans un tableau en enlevant les retours à la ligne
func Initialisation(pfTab []byte) []string {
	var tabS []string
	tempWord := ""
	for i := 0; i < len(string(pfTab)); i++ {
		tempWord = tempWord + string(pfTab[i])
		if string(pfTab[i]) == "\n" {
			tabS = append(tabS, tempWord)
			tempWord = ""
		}
	}
	return tabS
}

// Cette fonction permet d'enlever le byte 10 à la fin des mots
func Repair(randomWord string) string {
	tabB := []byte(randomWord)
	wordRepair := ""
	for i := 0; i < len(randomWord)-1; i++ {
		wordRepair = wordRepair + string(tabB[i])
	}
	return wordRepair
}

// Cette fonction sert à choisir un mot au hasard dans un tableau
func RandomPick(pfTab []string) string {
	rand.Seed(time.Now().UnixNano())
	randomInt := rand.Intn(len(pfTab))
	return pfTab[randomInt]
}

// Cette fonction permet de réveler certaines lettres au hasard et de remplacer les autres par des underscore
func Reveal(pfWord string) string {
	tabB := []byte(pfWord)
	var tabS []string
	var finalTabS []string
	finalString := ""
	n := len(pfWord)/2 - 1
	for i := 0; i < len(pfWord)-1; i++ {
		if string(tabB[i]) != " " {
			finalTabS = append(finalTabS, string(tabB[i]))
		}
	}
	tabInt := rand.Perm(len(pfWord) - 1)[:n]
	for j := 0; j < len(finalTabS); j++ {
		tabS = append(tabS, "_")
	}
	for i := 0; i < len(tabInt); i++ {
		tabS[tabInt[i]] = finalTabS[tabInt[i]]
	}
	for i := 0; i < len(tabS); i++ {
		if tabS[i] != " " {
			finalString = finalString + tabS[i]
		}
	}
	return finalString
}

// Ici nous avons toute la gestion de l'affichage des positions du hangman qui ne sont pas faites à partir d'une lecture d'un fichier pour optimiser le programme
func HangmanPositions(n int) {
	if n == 9 {
		fmt.Print("\n")
		fmt.Print("\n")
		fmt.Print("\n")
		fmt.Print("\n")
		fmt.Print("\n")
		fmt.Print("\n")
		fmt.Print("=========")
		fmt.Print("\n")
		fmt.Print("\n")
	}
	if n == 8 {
		fmt.Print("\n")
		fmt.Print("      |")
		fmt.Print("\n")
		fmt.Print("      |")
		fmt.Print("\n")
		fmt.Print("      |")
		fmt.Print("\n")
		fmt.Print("      |")
		fmt.Print("\n")
		fmt.Print("      |")
		fmt.Print("\n")
		fmt.Print("=========")
		fmt.Print("\n")
		fmt.Print("\n")
	}
	if n == 7 {
		fmt.Print("\n")
		fmt.Print("  +---+")
		fmt.Print("\n")
		fmt.Print("      |")
		fmt.Print("\n")
		fmt.Print("      |")
		fmt.Print("\n")
		fmt.Print("      |")
		fmt.Print("\n")
		fmt.Print("      |")
		fmt.Print("\n")
		fmt.Print("      |")
		fmt.Print("\n")
		fmt.Print("=========")
		fmt.Print("\n")
		fmt.Print("\n")
	}
	if n == 6 {
		fmt.Print("\n")
		fmt.Print("  +---+")
		fmt.Print("\n")
		fmt.Print("  |   |")
		fmt.Print("\n")
		fmt.Print("      |")
		fmt.Print("\n")
		fmt.Print("      |")
		fmt.Print("\n")
		fmt.Print("      |")
		fmt.Print("\n")
		fmt.Print("      |")
		fmt.Print("\n")
		fmt.Print("=========")
		fmt.Print("\n")
		fmt.Print("\n")
	}
	if n == 5 {
		fmt.Print("\n")
		fmt.Print("  +---+")
		fmt.Print("\n")
		fmt.Print("  |   |")
		fmt.Print("\n")
		fmt.Print("  O   |")
		fmt.Print("\n")
		fmt.Print("      |")
		fmt.Print("\n")
		fmt.Print("      |")
		fmt.Print("\n")
		fmt.Print("      |")
		fmt.Print("\n")
		fmt.Print("=========")
		fmt.Print("\n")
		fmt.Print("\n")
	}
	if n == 4 {
		fmt.Print("\n")
		fmt.Print("  +---+")
		fmt.Print("\n")
		fmt.Print("  |   |")
		fmt.Print("\n")
		fmt.Print("  O   |")
		fmt.Print("\n")
		fmt.Print("  |   |")
		fmt.Print("\n")
		fmt.Print("      |")
		fmt.Print("\n")
		fmt.Print("      |")
		fmt.Print("\n")
		fmt.Print("=========")
		fmt.Print("\n")
		fmt.Print("\n")
	}
	if n == 3 {
		fmt.Print("\n")
		fmt.Print("  +---+")
		fmt.Print("\n")
		fmt.Print("  |   |")
		fmt.Print("\n")
		fmt.Print("  O   |")
		fmt.Print("\n")
		fmt.Print(" /|   |")
		fmt.Print("\n")
		fmt.Print("      |")
		fmt.Print("\n")
		fmt.Print("      |")
		fmt.Print("\n")
		fmt.Print("=========")
		fmt.Print("\n")
		fmt.Print("\n")
	}
	if n == 2 {
		fmt.Print("\n")
		fmt.Print("  +---+")
		fmt.Print("\n")
		fmt.Print("  |   |")
		fmt.Print("\n")
		fmt.Print("  O   |")
		fmt.Print("\n")
		fmt.Print(` /|\  |`)
		fmt.Print("\n")
		fmt.Print("      |")
		fmt.Print("\n")
		fmt.Print("      |")
		fmt.Print("\n")
		fmt.Print("=========")
		fmt.Print("\n")
		fmt.Print("\n")
	}
	if n == 1 {
		fmt.Print("\n")
		fmt.Print("  +---+")
		fmt.Print("\n")
		fmt.Print("  |   |")
		fmt.Print("\n")
		fmt.Print("  O   |")
		fmt.Print("\n")
		fmt.Print(` /|\  |`)
		fmt.Print("\n")
		fmt.Print(" /    |")
		fmt.Print("\n")
		fmt.Print("      |")
		fmt.Print("\n")
		fmt.Print("=========")
		fmt.Print("\n")
		fmt.Print("\n")
	}
	if n == 0 {
		fmt.Print("\n")
		fmt.Print("  +---+")
		fmt.Print("\n")
		fmt.Print("  |   |")
		fmt.Print("\n")
		fmt.Print("  O   |")
		fmt.Print("\n")
		fmt.Print(` /|\  |`)
		fmt.Print("\n")
		fmt.Print(` / \  |`)
		fmt.Print("\n")
		fmt.Print("      |")
		fmt.Print("\n")
		fmt.Print("=========")
		fmt.Print("\n")
		fmt.Print("\n")
	}
}
