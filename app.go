package main

import "fmt"

const nTeamMAX int = 20
const nMatchesMAX int = 760

type Member struct {
	id, name                           string
	totalKill, totalDeath, totalAssist int
}

type MemberStats struct {
	id                  string
	kill, death, assist int
}

type Team struct {
	id, name  string
	win, lose int
	members   [5]Member
}

type Match struct {
	date, month, year      int            // Date of the match occured
	tHomeId, tAwayId       string         // Team ID
	pHome, pAway           int            // Team Point
	mHomeStats, mAwayStats [5]MemberStats // Statistics of Members each teams
}

type tabTeams [nTeamMAX]Team
type tabMatches [nMatchesMAX]Match

var teams tabTeams
var matches tabMatches
var nTeams int
var nMatches int

func getIndexTeamFromName(name string) int {
	var i, idx int

	idx = -1
	for i = 0; i < nTeams && idx < 0; i++ {
		if teams[i].name == name {
			idx = i
		}
	}

	return idx
}

func getIndexTeamFromId(id string) int {
	var idx, left, mid, right int
	idx = -1
	left = 0
	right = nTeams - 1

	for left <= right && idx < 0 {
		mid = (left + right) / 2
		if id < teams[mid].id {
			right = mid - 1
		} else if id > teams[mid].id {
			left = mid + 1
		} else {
			idx = mid
		}
	}

	return idx
}

func showStandingTable() {
	var pass, idx, i int
	var t tabTeams
	var team Team

	t = teams
	pass = 1
	for pass <= nTeams-1 {
		idx = pass - 1
		i = pass
		for i < nTeams {
			if (t[idx].win - t[idx].lose) < (t[i].win - t[i].lose) {
				idx = i
			}
			i++
		}
		team = t[pass-1]
		t[pass-1] = t[idx]
		t[idx] = team
		pass++
	}

	fmt.Println("+-----+------------+--------------------------------+--------+--------+")
	fmt.Printf("| %-3s | %-10s | %-30s | %-6s | %-6s |\n", "Pos", "ID", "Nama Tim", "Menang", "Kalah")
	fmt.Println("+-----+------------+--------------------------------+--------+--------+")

	for i = 0; i < nTeams; i++ {
		fmt.Printf("| %-3d | %-10s | %-30s | %-6d | %-6d |\n", i+1, t[i].id, t[i].name, t[i].win, t[i].lose)
	}
	if nTeams != 0 {
		fmt.Println("+-----+------------+--------------------------------+--------+--------+")
	}
	fmt.Println()
}

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

// Opsi 1

func createTeamName(team *Team) {
	var name string
	var isDuplicate bool
	var i int

	Print("Buatlah nama untuk tim yang sedang dibuat: ", false)
	fmt.Scan(&name)

	if len(name) < 2 && len(name) > 30 {
		clearTerminal()
		Print("Nama tim tidak boleh kurang dari 2 dan lebih dari 30 karakter!", true)
		createTeamName(*&team)
		return
	}

	isDuplicate = false
	for i = 0; i < nTeams && !isDuplicate; i++ {
		isDuplicate = name == teams[i].name
	}

	if isDuplicate {
		clearTerminal()
		Print("Nama "+name+" sudah digunakan oleh tim lain!", true)
		createTeamName(*&team)
	} else {
		*&team.name = name
		*&team.id = fmt.Sprintf("TE%d", nTeams+1)
	}
}

func createTeamMembers(team *Team) {
	var i int
	var name string
	Print("Buatlah anggota tim "+team.name+"!", true)
	for i < 5 {
		Print(fmt.Sprintf("Anggota %d: ", i+1), false)
		fmt.Scan(&name)

		if len(name) < 2 || len(name) > 15 {
			Print("Nama anggota tidak boleh lebih dari 15 karakter atau kurang dari 2 karakter", true)
		} else {
			team.members[i].name = name
			team.members[i].id = fmt.Sprintf("%s-ME%d", team.id, i+1)
			i++
		}
	}
}

func createTeam() {
	var team Team
	createTeamName(&team)
	createTeamMembers(&team)

	teams[nTeams] = team
	nTeams++
}

// Opsi 2

func selectUpdateOptions(idx int) {
	var pick int
	Print(fmt.Sprintf("(%s) Apa yang hendak anda perbarui?", teams[idx].name), true)
	Print("1). Nama tim", true)
	Print("2). Pertandingan", true)
	Print("3). Kembali", true)
	Print("Pilih: ", false)
	fmt.Scan(&pick)

	clearTerminal()
	if pick < 1 || pick > 3 {
		Print(fmt.Sprintf("Tidak ada opsi %d dalam pilihan!", pick), true)
	} else {
		switch pick {
		case 1:
			updateTeamName(idx)
		case 2:
			if nTeams == 1 {
				Print("Anda tidak dapat memperbarui jadwal pertandingan jika hanya ada 1 tim dalam turnamen!", true)
			} else {
				selectUpdateMatchOptions(idx)
			}
		}
	}
	if pick != 3 {
		selectUpdateOptions(idx)
	}
}

func updateTeamName(idx int) {
	var name string
	var i int
	var isDuplicate bool

	Print(fmt.Sprintf("Buat nama pengganti '%s': ", teams[idx].name), false)
	fmt.Scan(&name)

	clearTerminal()
	if len(name) < 2 && len(name) > 30 {
		Print("Nama tim tidak boleh kurang dari 2 dan lebih dari 30 karakter!", true)
		updateTeamName(idx)
	} else {
		isDuplicate = false
		for i = 0; i < nTeams && !isDuplicate; i++ {
			if teams[idx].id != teams[i].id {
				isDuplicate = name == teams[i].name
			}
		}

		if isDuplicate {
			Print("Nama "+name+" sudah digunakan oleh tim lain!", true)
			updateTeamName(idx)
		} else {
			teams[idx].name = name
		}
	}
}

func addMatchDate(idx, indexOpponent int, month, date, year *int) {
	var i int
	var months [12]string
	var isAlready bool

	months = [12]string{"Januari", "Februari", "Maret", "April", "Mei", "Juni", "Juli", "Agustus", "September", "Oktober", "November", "Desember"}
	Print("Tentukan tanggal (format: bulan tanggal tahun): ", false)
	fmt.Scan(month, date, year)

	clearTerminal()
	if *month < 1 || *month > 12 {
		Print("Anda memasukkan bulan yang tidak valid!", true)
		addMatchDate(idx, indexOpponent, month, date, year)
		return
	}

	switch *month {
	case 1, 3, 5, 7, 8, 10, 12:
		if *date < 1 || *date > 31 {
			Print(fmt.Sprintf("Anda hanya dapat memasukkan tanggal antara 1 hingga 31 untuk bulan %s!", months[*month-1]), true)
			addMatchDate(idx, indexOpponent, month, date, year)
			return
		}
	case 4, 6, 9, 11:
		if *date < 1 || *date > 30 {
			Print(fmt.Sprintf("Anda hanya dapat memasukkan tanggal antara 1 hingga 30 untuk bulan %s!", months[*month-1]), true)
			addMatchDate(idx, indexOpponent, month, date, year)
			return
		}
	case 2:
		if *date < 1 || ((*date > 28 && *year%4 != 0) || (*date > 29 && *year%4 == 0)) {
			if *year%4 == 0 {
				Print(fmt.Sprintf("Anda hanya dapat memasukkan tanggal antara 1 hingga 29 untuk bulan %s!", months[*month-1]), true)
			} else {
				Print(fmt.Sprintf("Anda hanya dapat memasukkan tanggal antara 1 hingga 28 untuk bulan %s!", months[*month-1]), true)
			}
			addMatchDate(idx, indexOpponent, month, date, year)
			return
		}
	}

	isAlready = false
	for i = 0; i < nMatches && !isAlready; i++ {
		if ((matches[i].tHomeId == teams[idx].id && matches[i].tAwayId == teams[indexOpponent].id) || (matches[i].tHomeId == teams[indexOpponent].id && matches[i].tAwayId == teams[idx].id)) && matches[i].month == *month && matches[i].date == *date && matches[i].year == *year {
			isAlready = true
		}
	}

	if isAlready {
		Print(
			fmt.Sprintf(
				"Sudah terdapat sebuah pertandingan pada tanggal %d %s %d untuk tim %s dan %s!",
				*date,
				months[*month],
				*year,
				teams[idx].name,
				teams[indexOpponent].name,
			),
			true,
		)
		addMatchDate(idx, indexOpponent, month, date, year)
		return
	}
}

func addMatch(idx int) {
	var i, j, k, indexOpponent int
	var totalMatches, pick int
	var date, month, year int
	var t tabTeams

	Print("Daftar tim yang dapat menjadi lawan (maksimal pertemuan dengan tim yang sama sebanyak 4x):", true)
	for i = 0; i < nTeams; i++ {
		if teams[idx].id != teams[i].id {
			for j = 0; j < nMatches; j++ {
				if (matches[j].tHomeId == teams[idx].id && matches[j].tAwayId == teams[i].id) || (matches[j].tHomeId == teams[i].id && matches[j].tAwayId == teams[idx].id) {
					totalMatches++
				}
			}
			Print(fmt.Sprintf("%d) %s (%dx pertemuan)", k+1, teams[i].name, totalMatches), true)
			totalMatches = 0

			t[k] = teams[i]
			k++
		}
	}
	Print("Pilih tim nomor berapa: ", false)
	fmt.Scan(&pick)

	if pick < 1 || pick > k {
		clearTerminal()
		Print("Masukan angka yang valid!", true)
		addMatch(idx)
		return
	}

	indexOpponent = -1
	for i = 0; i < nTeams && indexOpponent < 0; i++ {
		if t[pick-1].id == teams[i].id {
			indexOpponent = i
		}
	}
	for i = 0; i < nMatches; i++ {
		if (matches[i].tHomeId == teams[idx].id && matches[i].tAwayId == teams[indexOpponent].id) || (matches[i].tHomeId == teams[indexOpponent].id && matches[i].tAwayId == teams[idx].id) {
			totalMatches++
		}
	}
	if totalMatches >= 4 {
		clearTerminal()
		Print(teams[idx].name+" sudah berhadapan dengan "+teams[indexOpponent].name+" sebanyak 4x, yang mana ini adalah jumlah maksimal!", true)
		addMatch(idx)
		return
	}

	matches[nMatches].tHomeId = teams[idx].id
	matches[nMatches].tAwayId = teams[indexOpponent].id
	for i = 0; i < 5; i++ {
		matches[nMatches].mHomeStats[i].id = teams[idx].members[i].id
		matches[nMatches].mAwayStats[i].id = teams[indexOpponent].members[i].id
	}
	clearTerminal()
	addMatchDate(idx, indexOpponent, &month, &date, &year)
	matches[nMatches].month = month
	matches[nMatches].date = date
	matches[nMatches].year = year
	nMatches++
	clearTerminal()
}

func pickMatch(idx int, pick *int, mc *tabMatches) {
	var i, index int
	var nMc int

	index = -1
	fmt.Println("Daftar Pertandingan", teams[idx].name)
	for i = 0; i < nMatches; i++ {
		if matches[i].tHomeId == teams[idx].id {
			index = getIndexTeamFromId(matches[i].tAwayId)
		}
		if matches[i].tAwayId == teams[idx].id {
			index = getIndexTeamFromId(matches[i].tHomeId)
		}

		if index >= 0 && ((matches[i].tHomeId == teams[idx].id && matches[i].tAwayId == teams[index].id) || (matches[i].tHomeId == teams[index].id && matches[i].tAwayId == teams[idx].id)) {
			fmt.Printf("1). %s ", teams[idx].name)
			if teams[idx].id == matches[i].tHomeId {
				fmt.Printf("(%d) vs (%d) ", matches[i].pHome, matches[i].pAway)
			} else if teams[idx].id == matches[i].tAwayId {
				fmt.Printf("(%d) vs (%d) ", matches[i].pAway, matches[i].pHome)
			}
			fmt.Printf("%s - %d/%d/%d", teams[index].name, matches[i].date, matches[i].month, matches[i].year)
			if matches[i].pHome == matches[i].pAway {
				fmt.Print(" (Belum)\n")
			} else {
				fmt.Print("\n")
			}

			mc[nMc] = matches[i]
			nMc++
		}
		index = -1
	}
	fmt.Println()

	Print("Pilih pertandingan yang hendak anda perbarui (ketik '0' jika hendak kembali): ", false)
	fmt.Scan(pick)

	if *pick < 0 || *pick > nMc {
		clearTerminal()
		Print("Anda memilih pertandingan yang tidak valid!", true)
		fmt.Println()
		pickMatch(idx, pick, mc)
	}
}

func updateMatchScore(idx int, match Match) {
	var i int
	var index, indexMatch int
	var score int

	if teams[idx].id == match.tHomeId {
		index = getIndexTeamFromId(match.tAwayId)
	}
	if teams[idx].id == match.tAwayId {
		index = getIndexTeamFromId(match.tHomeId)
	}

	indexMatch = -1
	for i = 0; i < nMatches && indexMatch < 0; i++ {
		if match.tHomeId == matches[i].tHomeId && match.tAwayId == matches[i].tAwayId && match.date == matches[i].date && match.month == matches[i].month && match.year == matches[i].year {
			indexMatch = i
		}
	}

	Print(fmt.Sprintf("Tentukan skor untuk %s (Range skor 0 - 2. Jika skor-nya 1 atau 0, %s akan secara otomatis mendapatkan 2 poin): ", teams[idx].name, teams[index].name), false)
	fmt.Scan(&score)

	if score < 0 || score > 2 {
		clearTerminal()
		Print("Skor yang kamu masukkan tidak valid!", true)
		updateMatchScore(idx, match)
		return
	}

	if matches[indexMatch].pHome != 0 || matches[indexMatch].pAway != 0 {
		if matches[indexMatch].pHome == 2 {
			if teams[idx].id == matches[indexMatch].tHomeId {
				teams[idx].win--
				teams[index].lose--
			} else {
				teams[idx].lose--
				teams[index].win--
			}
		} else {
			if teams[idx].id == matches[indexMatch].tHomeId {
				teams[idx].lose--
				teams[index].win--
			} else {
				teams[idx].win--
				teams[index].lose--
			}
		}
	}
	if score < 2 {
		clearTerminal()
		if teams[idx].id == matches[indexMatch].tHomeId {
			matches[indexMatch].pHome = score
			matches[indexMatch].pAway = 2
		} else if teams[idx].id == matches[indexMatch].tAwayId {
			matches[indexMatch].pAway = score
			matches[indexMatch].pHome = 2
		}
	} else {
		if teams[idx].id == matches[indexMatch].tHomeId {
			matches[indexMatch].pHome = score
		} else if teams[idx].id == matches[indexMatch].tAwayId {
			matches[indexMatch].pAway = score
		}
		Print(fmt.Sprintf("Tentukan skor untuk %s (Range skor 0 - 1): ", teams[index].name), false)
		fmt.Scan(&score)

		clearTerminal()
		if score < 0 || score > 1 {
			Print("Skor yang kamu masukkan tidak valid!", true)
			updateMatchScore(idx, match)
			return
		}

		if teams[index].id == matches[indexMatch].tHomeId {
			matches[indexMatch].pHome = score
		} else if teams[index].id == matches[indexMatch].tAwayId {
			matches[indexMatch].pAway = score
		}
	}

	if matches[indexMatch].pHome > matches[indexMatch].pAway {
		if teams[idx].id == matches[indexMatch].tHomeId {
			teams[idx].win++
			teams[index].lose++
		} else if teams[index].id == matches[indexMatch].tHomeId {
			teams[index].win++
			teams[idx].lose++
		}
	} else if matches[indexMatch].pHome < matches[indexMatch].pAway {
		if teams[idx].id == matches[indexMatch].tAwayId {
			teams[idx].win++
			teams[index].lose++
		} else if teams[index].id == matches[indexMatch].tAwayId {
			teams[index].win++
			teams[idx].lose++
		}
	}
}

func updateMatchDate(idx int, match Match) {
	var i, indexMatch, indexOpponent int
	var month, date, year int
	var months [12]string
	var isAlready bool

	months = [12]string{"Januari", "Februari", "Maret", "April", "Mei", "Juni", "Juli", "Agustus", "September", "Oktober", "November", "Desember"}
	for i = 0; i < nMatches; i++ {
		if match.tHomeId == matches[i].tHomeId && match.tAwayId == matches[i].tAwayId && match.date == matches[i].date && match.month == matches[i].month && match.year == matches[i].year {
			indexMatch = i
		}
	}
	switch teams[idx].id {
	case match.tHomeId:
		indexOpponent = getIndexTeamFromId(match.tAwayId)
	case match.tAwayId:
		indexOpponent = getIndexTeamFromId(match.tHomeId)
	}

	Print("Masukkan tanggal pertandingan (format: bulan tanggal tahun): ", false)
	fmt.Scan(&month, &date, &year)

	clearTerminal()
	if month < 1 || month > 12 {
		Print("Anda memasukkan bulan yang tidak valid!", true)
		updateMatchDate(idx, match)
		return
	}

	switch month {
	case 1, 3, 5, 7, 8, 10, 12:
		if date < 1 || date > 31 {
			Print(fmt.Sprintf("Anda hanya dapat memasukkan tanggal antara 1 hingga 31 untuk bulan %s!", months[month-1]), true)
			updateMatchDate(idx, match)
			return
		}
	case 4, 6, 9, 11:
		if date < 1 || date > 30 {
			Print(fmt.Sprintf("Anda hanya dapat memasukkan tanggal antara 1 hingga 30 untuk bulan %s!", months[month-1]), true)
			updateMatchDate(idx, match)
			return
		}
	case 2:
		if date < 1 || ((date > 28 && year%4 != 0) || (date > 29 && year%4 == 0)) {
			if year%4 == 0 {
				Print(fmt.Sprintf("Anda hanya dapat memasukkan tanggal antara 1 hingga 29 untuk bulan %s!", months[month-1]), true)
			} else {
				Print(fmt.Sprintf("Anda hanya dapat memasukkan tanggal antara 1 hingga 28 untuk bulan %s!", months[month-1]), true)
			}
			updateMatchDate(idx, match)
			return
		}
	}

	isAlready = false
	for i = 0; i < nMatches && !isAlready; i++ {
		if ((matches[i].tHomeId == teams[idx].id && matches[i].tAwayId == teams[indexOpponent].id) || (matches[i].tHomeId == teams[indexOpponent].id && matches[i].tAwayId == teams[idx].id)) && matches[i].month == month && matches[i].date == date && matches[i].year == year {
			isAlready = true
		}
	}

	if isAlready {
		Print(
			fmt.Sprintf(
				"Sudah terdapat sebuah pertandingan pada tanggal %d %s %d untuk tim %s dan %s!",
				date,
				months[month],
				year,
				teams[idx].name,
				teams[indexOpponent].name,
			),
			true,
		)
		updateMatchDate(idx, match)
		return
	}

	matches[indexMatch].date = date
	matches[indexMatch].month = month
	matches[indexMatch].year = year
}

func updateMatchKDA(idx int, match Match) {
	var i int
	var indexMatch, indexOpponent int
	var kill, death, assist int
	var pick, pickIndex int

	indexMatch = -1
	for i = 0; i < nMatches && indexMatch < 0; i++ {
		if match.tHomeId == matches[i].tHomeId && match.tAwayId == matches[i].tAwayId && match.date == matches[i].date && match.month == matches[i].month && match.year == matches[i].year {
			indexMatch = i
		}
	}

	Print("KDA (Kill, Death, Assist) pemain-pemain dari tim mana yang hendak ada perbarui?", true)
	if teams[idx].id == match.tHomeId {
		indexOpponent = getIndexTeamFromId(match.tAwayId)
		Print("1). "+teams[idx].name, true)
		Print("2). "+teams[indexOpponent].name, true)
	}
	if teams[idx].id == match.tAwayId {
		indexOpponent = getIndexTeamFromId(match.tHomeId)
		Print("1). "+teams[indexOpponent].name, true)
		Print("2). "+teams[idx].name, true)
	}
	Print("Pilih (ketik '0' jika membatalkan proses): ", false)
	fmt.Scan(&pick)

	clearTerminal()
	if pick < 0 || pick > 2 {
		Print("Anda memasukkan angka yang tidak valid!", true)
		updateMatchKDA(idx, match)
		return
	}

	switch pick {
	case 1:
		Print("Ubahlah KDA pemain ", false)
		if teams[idx].id == match.tHomeId {
			fmt.Print(teams[idx].name, " (format: kill death assist)!\n")
			pickIndex = idx
		} else {
			fmt.Print(teams[indexOpponent].name, " (format: kill death assist)!\n")
			pickIndex = indexOpponent
		}
		for i = 0; i < 5; i++ {
			teams[pickIndex].members[i].totalKill -= match.mHomeStats[i].kill
			teams[pickIndex].members[i].totalDeath -= match.mHomeStats[i].death
			teams[pickIndex].members[i].totalAssist -= match.mHomeStats[i].assist

			Print(fmt.Sprintf("%s: ", teams[pickIndex].members[i].name), false)
			fmt.Scan(&kill, &death, &assist)

			matches[indexMatch].mAwayStats[i].kill = kill
			matches[indexMatch].mAwayStats[i].death = death
			matches[indexMatch].mAwayStats[i].assist = assist

			teams[pickIndex].members[i].totalKill += kill
			teams[pickIndex].members[i].totalDeath += death
			teams[pickIndex].members[i].totalAssist += assist
		}
	case 2:
		Print("Ubahlah KDA pemain ", false)
		if teams[idx].id == match.tAwayId {
			fmt.Print(teams[idx].name, " (format: kill death assist)!\n")
			pickIndex = idx
		} else {
			fmt.Print(teams[indexOpponent].name, " (format: kill death assist)!\n")
			pickIndex = indexOpponent
		}
		for i = 0; i < 5; i++ {
			teams[pickIndex].members[i].totalKill -= match.mAwayStats[i].kill
			teams[pickIndex].members[i].totalDeath -= match.mAwayStats[i].death
			teams[pickIndex].members[i].totalAssist -= match.mAwayStats[i].assist

			Print(fmt.Sprintf("%s: ", teams[pickIndex].members[i].name), false)
			fmt.Scan(&kill, &death, &assist)

			matches[indexMatch].mAwayStats[i].kill = kill
			matches[indexMatch].mAwayStats[i].death = death
			matches[indexMatch].mAwayStats[i].assist = assist

			teams[pickIndex].members[i].totalKill += kill
			teams[pickIndex].members[i].totalDeath += death
			teams[pickIndex].members[i].totalAssist += assist
		}
	}

	clearTerminal()
	if pick != 0 {
		updateMatchKDA(idx, match)
	}
}

func pickMethod(idx int, pick *int) {
	Print(fmt.Sprintf("Apa yang hendak anda perbarui:"), true)
	Print("1). Skor", true)
	Print("2). Tanggal pertandingan", true)
	Print("3). KDA pemain", true)
	Print("4). Kembali", true)
	Print("Pilih: ", false)
	fmt.Scan(pick)

	if *pick < 1 || *pick > 4 {
		clearTerminal()
		Print("Pilihan anda tidak valid!", true)
		pickMethod(idx, pick)
	}
}

func updateMatch(idx int) {
	var indexMatch, pick int
	var mc tabMatches

	pickMatch(idx, &indexMatch, &mc)
	clearTerminal()
	if indexMatch == 0 {
		return
	}

	pickMethod(idx, &pick)
	clearTerminal()
	if pick == 4 {
		updateMatch(idx)
		return
	}

	switch pick {
	case 1:
		updateMatchScore(idx, mc[indexMatch-1])
	case 2:
		updateMatchDate(idx, mc[indexMatch-1])
	case 3:
		updateMatchKDA(idx, mc[indexMatch-1])
	}

	if pick != 4 {
		updateMatch(idx)
	}
}

func deleteMatch(idx int) {
	var action string
	var i, indexMatch, indexOpponent, pick int
	var nMc int
	var mc tabMatches

	for i = 0; i < nMatches; i++ {
		if teams[idx].id == matches[i].tHomeId || teams[idx].id == matches[i].tAwayId {
			Print(fmt.Sprintf("%d). %s ", nMc+1, teams[idx].name), false)
			if teams[idx].id == matches[i].tHomeId {
				indexOpponent = getIndexTeamFromId(matches[i].tAwayId)
				fmt.Printf("(%d) vs (%d) ", matches[i].pHome, matches[i].pAway)
			} else {
				indexOpponent = getIndexTeamFromId(matches[i].tHomeId)
				fmt.Printf("(%d) vs (%d) ", matches[i].pAway, matches[i].pHome)
			}
			fmt.Printf("%s - %d/%d/%d\n", teams[indexOpponent].name, matches[i].date, matches[i].month, matches[i].year)

			mc[nMc] = matches[i]
			nMc++
		}
	}

	Print("Pilih (ketik '0' untuk kembali): ", false)
	fmt.Scan(&pick)

	clearTerminal()
	if pick < 0 || pick > nMc {
		Print("Angka yang anda masukan tidak valid!", true)
		deleteMatch(idx)
		return
	}

	if pick == 0 {
		return
	}

	Print(fmt.Sprintf("%s ", teams[idx].name), false)
	if teams[idx].id == mc[pick-1].tHomeId {
		indexOpponent = getIndexTeamFromId(mc[pick-1].tAwayId)
		fmt.Printf("(%d) vs (%d) ", mc[pick-1].pHome, mc[pick-1].pAway)
	}
	if teams[idx].id == mc[pick-1].tAwayId {
		indexOpponent = getIndexTeamFromId(mc[pick-1].tHomeId)
		fmt.Printf("(%d) vs (%d) ", mc[pick-1].pAway, mc[pick-1].pHome)
	}
	fmt.Printf("%s - %d/%d/%d\n", teams[indexOpponent].name, mc[pick-1].date, mc[pick-1].month, mc[pick-1].year)
	Print("Apakah anda yakin hendak menghapus pertandingan ini? Ketik 'ya' untuk menghapus, dan lainnya untuk membatalkan: ", false)
	fmt.Scan(&action)

	clearTerminal()
	if action == "ya" || action == "yA" || action == "Ya" || action == "YA" {
		switch teams[idx].id {
		case mc[pick-1].tHomeId:
			if mc[pick-1].pHome > mc[pick-1].pAway {
				teams[idx].win--
				teams[indexOpponent].lose--
			} else if mc[pick-1].pHome < mc[pick-1].pAway {
				teams[idx].lose--
				teams[indexOpponent].win--
			}
		case mc[pick-1].tAwayId:
			if mc[pick-1].pHome > mc[pick-1].pAway {
				teams[idx].lose--
				teams[indexOpponent].win--
			} else if mc[pick-1].pHome < mc[pick-1].pAway {
				teams[idx].win--
				teams[indexOpponent].lose--
			}
		}
		for i = 0; i < 5; i++ {
			if teams[idx].id == mc[pick-1].tHomeId {
				teams[idx].members[i].totalKill -= mc[pick-1].mHomeStats[i].kill
				teams[idx].members[i].totalDeath -= mc[pick-1].mHomeStats[i].death
				teams[idx].members[i].totalAssist -= mc[pick-1].mHomeStats[i].assist

				teams[indexOpponent].members[i].totalKill -= mc[pick-1].mAwayStats[i].kill
				teams[indexOpponent].members[i].totalDeath -= mc[pick-1].mAwayStats[i].death
				teams[indexOpponent].members[i].totalAssist -= mc[pick-1].mAwayStats[i].assist
			} else {
				teams[idx].members[i].totalKill -= mc[pick-1].mAwayStats[i].kill
				teams[idx].members[i].totalDeath -= mc[pick-1].mAwayStats[i].death
				teams[idx].members[i].totalAssist -= mc[pick-1].mAwayStats[i].assist

				teams[indexOpponent].members[i].totalKill -= mc[pick-1].mHomeStats[i].kill
				teams[indexOpponent].members[i].totalDeath -= mc[pick-1].mHomeStats[i].death
				teams[indexOpponent].members[i].totalAssist -= mc[pick-1].mHomeStats[i].assist
			}
		}

		for i = 0; i < nMatches; i++ {
			if mc[pick-1].tHomeId == matches[i].tHomeId && mc[pick-1].tAwayId == matches[i].tAwayId && mc[pick-1].date == matches[i].date && mc[pick-1].month == matches[i].month && mc[pick-1].year == matches[i].year {
				indexMatch = i
			}
		}

		for i = indexMatch; i < nMatches; i++ {
			matches[i] = matches[i+1]
		}
		nMatches--
	}

	if nMatches != 0 {
		deleteMatch(idx)
	}
}

func selectUpdateMatchOptions(idx int) {
	var pick, i int
	var totalMatches int

	for i = 0; i < nMatches; i++ {
		if matches[i].tHomeId == teams[idx].id || matches[i].tAwayId == teams[idx].id {
			totalMatches++
		}
	}

	Print(fmt.Sprintf("(%s) Pilih salah-satu opsi di bawah ini:", teams[idx].name), true)
	Print("1). Tambah Pertandingan", true)
	Print("2). Perbarui Pertandingan", true)
	Print("3). Hapus Pertandingan", true)
	Print("4). Kembali", true)
	Print("Pilih: ", false)

	fmt.Scan(&pick)
	clearTerminal()

	if pick < 1 || pick > 4 {
		Print("Kamu memilih pilihan di luar opsi!", true)
	} else {
		switch pick {
		case 1:
			addMatch(idx)
		case 2:
			if totalMatches == 0 {
				Print(teams[idx].name+" tidak mempunyai jadwal pertandingan apapun untuk diperbarui!", true)
			} else {
				updateMatch(idx)
			}
		case 3:
			if totalMatches == 0 {
				Print(teams[idx].name+" tidak mempunyai jadwal pertandingan apapun untuk dihapus!", true)
			} else {
				deleteMatch(idx)
			}
		}
	}
	if pick != 4 {
		selectUpdateMatchOptions(idx)
	}
}

func selectTeam(idx *int) {
	var name string
	fmt.Println("Daftar Tim Dalam Turnamen")
	showStandingTable()

	Print("Tim mana yang hendak anda perbarui (ketik '0' untuk kembali): ", false)
	fmt.Scan(&name)

	if name != "0" {
		*idx = getIndexTeamFromName(name)
		clearTerminal()
		if *idx == -1 {
			Print("Tim yang hendak anda cari tidak ditemukan!", true)
			selectTeam(*&idx)
		}
	} else {
		*idx = -2
	}
}

func updateTeam() {
	var index int
	selectTeam(&index)
	clearTerminal()

	if index != -2 {
		selectUpdateOptions(index)
		updateTeam()
	}
}

// Opsi 3

func deleteTeamWait(idx int, action *string) {
	Print("Apakah anda yakin hendak menghapus "+teams[idx].name+"?", true)
	Print("Seluruh pertandingan yang melibatkan "+teams[idx].name+" akan ikut terhapus pula!", true)
	Print("Ketik 'ya' jika yakin hendak menghapus, ketik yang lain jika tidak jadi: ", false)
	fmt.Scan(action)
}

func deleteTeam() {
	var i, j, k, l int
	var index, indexOpponent int
	var action, newId string
	var nNewMatches int
	var newMatches tabMatches
	var team Team

	fmt.Println("Daftar Tim-Tim dalam Turnamen")
	for i = 0; i < nTeams; i++ {
		fmt.Printf("%d). %s\n", i+1, teams[i].name)
	}
	fmt.Println()

	Print("Tim mana yang hendak anda hapus (ketik '0' jika anda hendak meng-cancel proses): ", false)
	fmt.Scan(&index)

	clearTerminal()
	if index < 0 || index > nTeams {
		Print("Anda memilih tim yang tidak ada di dalam turnamen!", true)
		fmt.Println()
		deleteTeam()
		return
	}
	if index == 0 {
		return
	}

	index--
	team = teams[index]
	deleteTeamWait(index, &action)
	clearTerminal()

	if action == "ya" || action == "Ya" || action == "yA" || action == "YA" {
		for i = 0; i < nMatches; i++ {
			if matches[i].tHomeId != team.id && matches[i].tAwayId != team.id {
				newMatches[nNewMatches] = matches[i]
				nNewMatches++
			} else {
				if matches[i].pHome > matches[i].pAway {
					switch team.id {
					case matches[i].tHomeId:
						indexOpponent = getIndexTeamFromId(matches[i].tAwayId)
						teams[indexOpponent].lose--
					case matches[i].tAwayId:
						indexOpponent = getIndexTeamFromId(matches[i].tHomeId)
						teams[indexOpponent].win--
					}
				} else if matches[i].pHome < matches[i].pAway {
					switch team.id {
					case matches[i].tHomeId:
						indexOpponent = getIndexTeamFromId(matches[i].tAwayId)
						teams[indexOpponent].win--
					case matches[i].tAwayId:
						indexOpponent = getIndexTeamFromId(matches[i].tHomeId)
						teams[indexOpponent].lose--
					}
				}
				for j = 0; j < 5; j++ {
					if teams[indexOpponent].id == matches[i].tHomeId {
						teams[indexOpponent].members[j].totalKill -= matches[i].mHomeStats[j].kill
						teams[indexOpponent].members[j].totalDeath -= matches[i].mHomeStats[j].death
						teams[indexOpponent].members[j].totalAssist -= matches[i].mHomeStats[j].assist
					} else {
						teams[indexOpponent].members[j].totalKill -= matches[i].mAwayStats[j].kill
						teams[indexOpponent].members[j].totalDeath -= matches[i].mAwayStats[j].death
						teams[indexOpponent].members[j].totalAssist -= matches[i].mAwayStats[j].assist
					}
				}
			}
		}
		matches = newMatches
		nMatches = nNewMatches

		for i = index; i < nTeams; i++ {
			teams[i] = teams[i+1]
			newId = fmt.Sprintf("TE%d", i+1)

			for j = 0; j < nMatches; j++ {
				if teams[i].id == matches[j].tHomeId {
					matches[j].tHomeId = newId

					for k = 0; k < 5; k++ {
						for l = 0; l < 5; l++ {
							if teams[i].members[k].id == matches[j].mHomeStats[k].id {
								teams[i].members[k].id = fmt.Sprintf("%s-ME%d", newId, k+1)
								matches[j].mHomeStats[k].id = teams[i].members[k].id
							}
						}
					}
				}
				if teams[i].id == matches[j].tAwayId {
					matches[j].tAwayId = newId

					for k = 0; k < 5; k++ {
						for l = 0; l < 5; l++ {
							if teams[i].members[k].id == matches[j].mAwayStats[k].id {
								teams[i].members[k].id = fmt.Sprintf("%s-ME%d", newId, k+1)
								matches[j].mAwayStats[k].id = teams[i].members[k].id
							}
						}
					}
				}
			}
			teams[i].id = newId
		}
		nTeams--
	}

	if nTeams > 0 {
		deleteTeam()
	}
}

// Opsi 4

func viewStandingDetails() {
	var name string
	var pick string
	var index int

	Print("Masukan nama tim mana yang hendak anda lihat statistik-nya: ", false)
	fmt.Scan(&name)

	index = getIndexTeamFromName(name)
	if index < 0 {
		Print("Sistem tidak menemukan tim bernama: "+name, true)
		viewStandingDetails()
		return
	}
	clearTerminal()

	var members [5]Member
	var temp Member
	var pass, idx, i int

	members = teams[index].members
	pass = 1
	for pass <= 4 {
		idx = pass - 1
		i = pass
		for i < 5 {
			if (members[idx].totalKill*3 + members[idx].totalAssist - members[idx].totalDeath) < (members[i].totalKill*3 + members[i].totalAssist - members[i].totalDeath) {
				idx = i
			}
			i++
		}
		temp = members[pass-1]
		members[pass-1] = members[idx]
		members[idx] = temp
		pass++
	}
	fmt.Println("Statistik Anggota", name)
	fmt.Println("+------------+-----------------+--------+--------+--------+")
	fmt.Printf("| %-10s | %-15s | %-6s | %-6s | %-6s |\n", "ID", "Anggota", "Kill", "Death", "Assist")
	fmt.Println("+------------+-----------------+--------+--------+--------+")
	for i = 0; i < 5; i++ {
		fmt.Printf("| %-10s | %-15s | %-6d | %-6d | %-6d |\n", members[i].id, members[i].name, members[i].totalKill, members[i].totalDeath, members[i].totalAssist)
	}
	fmt.Println("+------------+-----------------+--------+--------+--------+")
	fmt.Println()

	Print("Ketik apapun untuk melanjutkan: ", false)
	fmt.Scan(&pick)
}

func viewStanding() {
	var teamName string
	var pick string
	var i, j, k int

	var players [5 * 20]Member
	var tempPlayer Member
	var nPlayers int

	fmt.Println("Klasemen Turnamen")
	showStandingTable()

	if nTeams != 0 {

		// Tabel pemain terbaik

		for i = 0; i < nTeams; i++ {
			for j = 0; j < 5; j++ {
				players[nPlayers] = teams[i].members[j]
				nPlayers++
			}
		}

		for i = 1; i < nPlayers; i++ {
			tempPlayer = players[i]
			for j = i; j > 0 && (tempPlayer.totalKill*3+tempPlayer.totalAssist-tempPlayer.totalDeath) > (players[j-1].totalKill*3+players[j-1].totalAssist-players[j-1].totalDeath); j-- {
				players[j] = players[j-1]
			}
			players[j] = tempPlayer
		}

		if nPlayers > 5 {
			nPlayers = 5
		}

		fmt.Println("=== Top 5 Pemain Dalam Turnamen ===")
		fmt.Println("🏅Pemain terbaik (MVP):", players[0].name)
		fmt.Println("+-----+------------------+-----------------+--------------------------------+--------+--------+--------+")
		fmt.Printf("| %-3s | %-16s | %-15s | %-30s | %-6s | %-6s | %-6s |\n", "Pos", "ID", "Pemain", "Tim", "Kill", "Death", "Assist")
		fmt.Println("+-----+------------------+-----------------+--------------------------------+--------+--------+--------+")
		for i = 0; i < nPlayers; i++ {
			teamName = "-"
			for j = 0; j < nTeams && teamName == "-"; j++ {
				for k = 0; k < 5; k++ {
					if players[i].id == teams[j].members[k].id {
						teamName = teams[j].name
					}
				}
			}

			fmt.Printf(
				"| %-3d | %-16s | %-15s | %-30s | %-6d | %-6d | %-6d |\n",
				i+1,
				players[i].id,
				players[i].name,
				teamName,
				players[i].totalKill,
				players[i].totalDeath,
				players[i].totalAssist,
			)
		}
		fmt.Println("+-----+------------------+-----------------+--------------------------------+--------+--------+--------+")
		fmt.Println()

		//

		Print("Apakah anda hendak melihat statistik tim lebih detail (ketik 'ya' untuk melihat, dan selain itu untuk kembali ke menu utama): ", false)
		fmt.Scan(&pick)

		if pick == "ya" || pick == "Ya" || pick == "yA" || pick == "YA" {
			viewStandingDetails()
			clearTerminal()
			viewStanding()
			return
		}
	} else {
		Print("Belum ada tim dalam klasemen, ketik apapun untuk melanjutkan: ", false)
		fmt.Scan(&pick)
	}
}

// Opsi 5

func isEarlier(tempMatch, incomingMatch Match) bool {
	return (tempMatch.year < incomingMatch.year) || (tempMatch.year == incomingMatch.year && tempMatch.month < incomingMatch.month) || (tempMatch.year == incomingMatch.year && tempMatch.month == incomingMatch.month && tempMatch.date < incomingMatch.date)
}

func viewStandingSchedule() {
	var i, j int
	var random string
	var indexHome, indexAway int
	var months [12]string

	months = [12]string{"Januari", "Febuari", "Maret", "April", "Mei", "Juni", "Juli", "Agustus", "September", "Oktober", "November", "Desember"}

	var incomingMatches [nMatchesMAX]Match
	var tempMatch Match
	var nMc int

	for i = 0; i < nMatches; i++ {
		if matches[i].pHome == 0 && matches[i].pAway == 0 {
			incomingMatches[nMc] = matches[i]
			nMc++
		}
	}

	for i = 1; i < nMc; i++ {
		tempMatch = incomingMatches[i]
		for j = i; j > 0 && isEarlier(tempMatch, incomingMatches[j-1]); j-- {
			incomingMatches[j] = incomingMatches[j-1]
		}
		incomingMatches[j] = tempMatch
	}

	if nMc > 5 {
		nMc = 5
	}

	fmt.Printf("Daftar %d jadwal pertandingan terdekat!\n", nMc)
	fmt.Println("+--------------------------------+--------------------------------+----------------------+")
	fmt.Printf("| %-30s | %-30s | %-20s |\n", "Tim Home", "Tim Away", "Tanggal")
	fmt.Println("+--------------------------------+--------------------------------+----------------------+")
	for i = 0; i < nMc; i++ {
		indexHome = getIndexTeamFromId(incomingMatches[i].tHomeId)
		indexAway = getIndexTeamFromId(incomingMatches[i].tAwayId)
		fmt.Printf(
			"| %-30s | %-30s | %-20s |\n",
			teams[indexHome].name,
			teams[indexAway].name,
			fmt.Sprintf("%d %s %d", incomingMatches[i].date, months[incomingMatches[i].month-1], incomingMatches[i].year),
		)
	}
	fmt.Println("+--------------------------------+--------------------------------+----------------------+")
	fmt.Println()

	Print("Ketik apapun untuk melanjutkan: ", false)
	fmt.Scan(&random)
	clearTerminal()
}

// Main Menu

func mainMenu() {
	var selection int
	var isBlocked bool

	Print(fmt.Sprintf("Saat ini anda memiliki %d dari %d slot tim dalam klasemen!", nTeams, nTeamMAX-nTeams), true)
	Print("Silakan pilih salah-satu opsi yang hendak anda lakukan:", true)
	Print("1). Tambah Tim", true)
	Print("2). Perbarui Tim", true)
	Print("3). Hapus Tim", true)
	Print("4). Lihat Klasemen", true)
	Print("5). Lihat Jadwal Pertandingan", true)
	Print("6). Keluar", true)
	Print("Pilihan anda: ", false)
	fmt.Scan(&selection)

	clearTerminal()
	isBlocked = false
	if selection < 1 || selection > 6 {
		Print("-----------------------------------------------", true)
		Print("Warning‼️", true)
		Print(fmt.Sprintf("Tidak ada opsi %d dalam pilihan!", selection), true)
		Print("-----------------------------------------------", true)
		isBlocked = true
	}

	if nTeams == 0 && (selection == 2 || selection == 3) {
		Print("-----------------------------------------------", true)
		Print("Warning‼️", true)
		Print("Opsi 2 dan 3 belum tersedia jika belum ada tim yang terdaftar!", true)
		Print("-----------------------------------------------", true)
		isBlocked = true
	}

	if nTeams == nTeamMAX && selection == 1 {
		Print("-----------------------------------------------", true)
		Print("Warning‼️", true)
		Print("Anda sudah tidak memiliki slot dalam klasemen untuk menambahkan tim!", true)
		Print("-----------------------------------------------", true)
		isBlocked = true
	}

	if !isBlocked {
		switch selection {
		case 1:
			createTeam()
		case 2:
			updateTeam()
		case 3:
			deleteTeam()
		case 4:
			viewStanding()
		case 5:
			viewStandingSchedule()
		}
	}

	clearTerminal()
	if selection != 5 {
		mainMenu()
	}
}

func main() {
	Print("Selamat datang di Turnamen anda!", true)
	mainMenu()
}
