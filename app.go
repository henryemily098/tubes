package main

import "fmt"

const NMAX int = 100

type Member struct {
	id, name            string
	kill, death, assist int
}

type Team struct {
	id, name           string
	win, lose          int
	members            [5]Member
	jadwalPertandingan [NMAX]int
}

var teams [NMAX]Team
var nTeams int

func Print(message string, newLine bool) {
	if newLine {
		fmt.Println("[SYSTEM]", message)
	} else {
		fmt.Print("[SYSTEM] ", message)
	}
}

func clearTerminal() {
	fmt.Print("\033[2J")
	fmt.Print("\033[H")
}

// Main Menu

func mainMenu(selection *int) {
	Print(fmt.Sprintf("Saat ini anda memiliki %d dari %d slot tim dalam klasemen!", nTeams, NMAX-nTeams), true)
	Print("Silakan pilih salah-satu opsi yang hendak anda lakukan:", true)
	Print("1). Tambah Tim", true)
	Print("2). Perbarui Tim", true)
	Print("3). Hapus Tim", true)
	Print("4). Lihat Klasemen", true)
	Print("5). Lihat Statistik Tim dan Pemainnya", true)
	Print("6). Keluar", true)
	Print("Pilihan anda: ", false)
	fmt.Scan(*&selection)

	clearTerminal()
	if *selection < 1 || *selection > 5 {
		Print("-----------------------------------------------", true)
		Print("Warning‼️", true)
		Print(fmt.Sprintf("Tidak ada opsi %d dalam pilihan!", *selection), true)
		Print("-----------------------------------------------", true)
		mainMenu(*&selection)
	}

	if nTeams == NMAX && *selection == 1 {
		Print("-----------------------------------------------", true)
		Print("Warning‼️", true)
		Print("Kamu sudah tidak memiliki slot dalam klasemen untuk menambahkan tim!", true)
		Print("-----------------------------------------------", true)
	}
}

// Opsi 1

func createTeam() {

}

func main() {
	var s int
	Print("Selamat datang di Turnamen anda!", true)
	mainMenu(&s)

	switch s {
	case 1:
		createTeam()
	}
}
