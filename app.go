package main

import (
	"fmt"
	"strings"
)

const NMAX int = 50

type Member struct {
	id   string
	name string
}

type Team struct {
	id, name string
	members  [5]Member
}

type StatisticTeam struct {
	id           string
	winner, lose int
}

type MemberStatistics struct {
	id                  string
	kill, assist, death int
}

type MatchResults struct {
	team_a_id, team_b_id                     string
	point_a, point_b                         int
	members_a_statistic, members_b_statistic [5]MemberStatistics
}

type StandingTournament struct {
	teams              [16]Team
	nTeams             int
	statisticTeams     [16]StatisticTeam
	hasilPertandingan  [120]MatchResults
	jumlahPertandingan int
}

type PlayOff struct {
	team_a_id, team_b_id, round              string
	point_a, point_b                         int
	members_a_statistic, members_b_statistic [5]MemberStatistics
	// 12 Teams = Round 1 -> Round 2 (1 team auto final, 2 team semifinal) -> SemiFinal -> Final
	// 8 Teams = QuarterFinal -> SemiFinal -> Final
	// 6 Teams = Round 1 (1 team auto final, 2 team semifinal) -> SemiFinal -> Final
	// 4 Teams = SemiFinal -> Final
}

type Tournament struct {
	id, name, champion_id        string
	system, nPlayOff, nStandings int
	standings                    [4]StandingTournament
	playOff                      [12]PlayOff
	finishGroup, finishPlayOff   bool
}

var nTournaments int
var tournaments [NMAX]Tournament

func factorial(n int) int {
	if n == 0 {
		return 1
	} else {
		return n * factorial(n-1)
	}
}

func combination(n, r int) int {
	return factorial(n) / (factorial(n-r) * factorial(r))
}

func generateID(letter string, number int) string {
	return fmt.Sprint(strings.ToUpper(letter), number)
}

func searchTournamentIndex(id string) int {
	var left, right, mid int
	var idx int

	left = 0
	right = nTournaments - 1
	idx = -1

	for left <= right && idx < 0 {
		mid = (left + right) / 2
		if id < tournaments[mid].id {
			right = mid - 1
		} else if id > tournaments[mid].id {
			left = mid + 1
		} else {
			idx = mid
		}
	}
	return idx
}

func findTournamentTeamsIndex(tournament Tournament, indexGroup int, id string) int {
	var left, right, mid int
	var idx int

	left = 0
	right = tournament.standings[indexGroup].nTeams - 1
	idx = -1

	for left <= right && idx < 0 {
		mid = (left + right) / 2
		if id < tournament.standings[indexGroup].teams[mid].id {
			right = mid - 1
		} else if id > tournament.standings[indexGroup].teams[mid].id {
			left = mid + 1
		} else {
			idx = mid
		}
	}
	return idx
}

func findTournamentMatchResultsIndex(tournament Tournament, indexGroup int, idA, idB string) int {
	var i, idx int
	var team_a_id, team_b_id string

	i = 0
	idx = -1
	for i < tournament.standings[indexGroup].jumlahPertandingan && idx < 0 {
		team_a_id = tournament.standings[indexGroup].hasilPertandingan[i].team_a_id
		team_b_id = tournament.standings[indexGroup].hasilPertandingan[i].team_b_id
		if (team_a_id == idA && team_b_id == idB) || (team_a_id == idB && team_b_id == idA) {
			idx = i
		}
		i++
	}
	return idx
}

func sortTournamentTeamsBasedStatistics(array [16]StatisticTeam, nTeams int) [16]StatisticTeam {
	var i, idx, pass int
	var temp StatisticTeam
	pass = 1
	for pass <= nTeams-1 {
		idx = pass - 1
		i = pass
		for i < nTeams {
			if array[idx].winner < array[i].winner || (array[idx].winner == array[i].winner && array[idx].lose > array[i].lose) {
				idx = i
			}
			i++
		}
		temp = array[pass-1]
		array[pass-1] = array[idx]
		array[idx] = temp
		pass++
	}
	return array
}

func clearTerminal() {
	fmt.Print("\033[2J")
	fmt.Print("\033[H")
}

func showTournamentTableList() {
	var i int
	var text string

	fmt.Println()
	fmt.Println("+------------+--------------------------------+----------------------+")
	fmt.Printf("│ %-10s │ %-30s │ %-20s │\n", "ID", "Nama Turnamen", "Sistem")
	fmt.Println("+------------+--------------------------------+----------------------+")
	for i = 0; i < nTournaments; i++ {
		text = fmt.Sprintf("│ %-10s │ %-30s │", tournaments[i].id, tournaments[i].name)
		switch tournaments[i].system {
		case 1:
			text += fmt.Sprintf(" %-20s │", "Klasemen")
		case 2:
			text += fmt.Sprintf(" %-20s │", "Group")
		}
		fmt.Println(text)
	}
	fmt.Println("+------------+--------------------------------+----------------------+")
	fmt.Println()
}

func showTournamentStandings(tournament Tournament, i int) {
	var teams [16]StatisticTeam
	var j, idx, nQualified int

	fmt.Println("+----------------------+--------------------------------+----------------------+----------------------+")
	fmt.Printf("│ %-20s │ %-30s │ %-20s │ %-20s │\n", "ID", "Nama Tim", "Winner - Lose", "Status Qualified")
	fmt.Println("+----------------------+--------------------------------+----------------------+----------------------+")
	teams = sortTournamentTeamsBasedStatistics(tournament.standings[i].statisticTeams, tournament.standings[i].nTeams)

	switch tournament.system {
	case 1:
		nQualified = tournament.nPlayOff
	case 2:
		nQualified = tournament.nPlayOff / tournament.nStandings
	}

	for j = 0; j < tournament.standings[i].nTeams; j++ {
		idx = findTournamentTeamsIndex(tournament, i, teams[j].id)
		if j+1 > nQualified {
			fmt.Printf(
				"│ %-20s │ %-30s │ %-20s │ %-20s │\n",
				teams[j].id,
				tournament.standings[i].teams[idx].name,
				fmt.Sprintf("%d - %d", teams[j].winner, teams[j].lose),
				"Eliminate",
			)
		} else {
			fmt.Printf(
				"│ %-20s │ %-30s │ %-20s │ %-20s │\n",
				teams[j].id,
				tournament.standings[i].teams[idx].name,
				fmt.Sprintf("%d - %d", teams[j].winner, teams[j].lose),
				"Qualify",
			)
		}
	}
	fmt.Println("+----------------------+--------------------------------+----------------------+----------------------+")
}

func Print(message string, newLine bool) {
	if newLine {
		fmt.Println("[SYSTEM] " + message)
	} else {
		fmt.Print("[SYSTEM] " + message)
	}
}

// Main Menu

func chooseAction(n *int) {
	Print("Apa yang hendak anda lakukan saat ini?", true)
	Print("1). Perbarui Turnamen", true)
	Print("2). View Turnamen", true)
	Print("3). Hapus Turnamen", true)
	Print("4). Panduan Pembuatan Turnamen", true)
	if nTournaments < NMAX {
		Print("5). Tambah Turnamen", true)
		Print("6). Exit", true)
	} else {
		Print("5). Exit", true)
	}
	Print("Masukkan pilihan anda: ", false)
	fmt.Scan(*&n)

	if *n < 1 || (nTournaments == NMAX && *n > 5) || (nTournaments < NMAX && *n > 6) {
		Print("Masukkan yang anda masukan tidak valid!", true)
		Print("-----------------------------------------------", true)
		chooseAction(*&n)
		return
	}

	if *n != 5 && *n != 6 && nTournaments == 0 {
		Print("Anda tidak memiliki data turnamen apapun untuk dilihat, diperbarui, ataupun dihapus!", true)
		Print("-----------------------------------------------", true)
		chooseAction(*&n)
		return
	}
	Print("-----------------------------------------------", true)
	clearTerminal()
}

func startingMessage() {
	Print("Selamat datang di administratif Turnamen Valorant Champions Tour E-Sports!", true)
	Print(fmt.Sprint("Saat ini terdapat ", nTournaments, " dari ", NMAX, " turnamen dalam database anda!"), true)
	if nTournaments == NMAX {
		Print("Anda sudah kehabisan slot turnamen! Silakan hapus 1 atau beberapa turnamen dari database anda jika hendak memasukkan turnamen baru!", true)
	} else {
		Print(fmt.Sprint("Anda memiliki sisa slot sebanyak ", NMAX-nTournaments, " untuk menyimpan data dari Turnamen anda!"), true)
	}
}

// Opsi 1

func updateGroupMatchPointsTeam(name string, rEnd int, point *int) {
	Print(fmt.Sprintf("Poin tim %s (range: 0 hingga %d): ", name, rEnd), false)
	fmt.Scan(*&point)

	if *point < 0 || *point > 2 {
		Print("Anda memasukkan angka di luar range!", true)
		updateGroupMatchPointsTeam(name, rEnd, *&point)
	}
}

func updateGroupMatchPoints(idx, indexGroup int, teamAId, teamBId string) {
	var index, teamAIndex, teamBIndex int
	var pointA, pointB int
	var i int

	teamAIndex = findTournamentTeamsIndex(
		tournaments[idx],
		indexGroup,
		teamAId,
	)
	teamBIndex = findTournamentTeamsIndex(
		tournaments[idx],
		indexGroup,
		teamBId,
	)
	index = findTournamentMatchResultsIndex(
		tournaments[idx],
		indexGroup,
		teamAId,
		teamBId,
	)

	Print(
		fmt.Sprintf(
			"Note: Jika anda memasukkan angka selain 2 untuk %s, tim %s akan otomatis mendapatkan 2 poin.",
			tournaments[idx].standings[indexGroup].teams[teamAIndex].name,
			tournaments[idx].standings[indexGroup].teams[teamBIndex].name,
		),
		true,
	)
	updateGroupMatchPointsTeam(
		tournaments[idx].standings[indexGroup].teams[teamAIndex].name,
		2,
		&pointA,
	)

	if pointA < 2 {
		pointB = 2
	} else {
		updateGroupMatchPointsTeam(
			tournaments[idx].standings[indexGroup].teams[teamBIndex].name,
			1,
			&pointB,
		)
	}

	if index < 0 {
		index = tournaments[idx].standings[indexGroup].jumlahPertandingan
		tournaments[idx].standings[indexGroup].jumlahPertandingan++
	} else {
		tournaments[idx].standings[indexGroup].statisticTeams[teamAIndex].winner--
		tournaments[idx].standings[indexGroup].statisticTeams[teamBIndex].lose--
	}

	tournaments[idx].standings[indexGroup].hasilPertandingan[index].team_a_id = teamAId
	tournaments[idx].standings[indexGroup].hasilPertandingan[index].team_b_id = teamBId
	tournaments[idx].standings[indexGroup].hasilPertandingan[index].point_a = pointA
	tournaments[idx].standings[indexGroup].hasilPertandingan[index].point_b = pointB

	if pointA > pointB {
		tournaments[idx].standings[indexGroup].statisticTeams[teamAIndex].winner++
		tournaments[idx].standings[indexGroup].statisticTeams[teamBIndex].lose++
	} else if pointA < pointB {
		tournaments[idx].standings[indexGroup].statisticTeams[teamAIndex].lose++
		tournaments[idx].standings[indexGroup].statisticTeams[teamBIndex].winner++
	}

	// A
	Print("Statistik Tim "+tournaments[idx].standings[indexGroup].teams[teamAIndex].name, true)
	for i = 0; i < 5; i++ {
		Print("Buatlah statisik KDA untuk "+tournaments[idx].standings[indexGroup].teams[teamAIndex].members[i].name, true)
		Print("Kill: ", false)
		fmt.Scan(&tournaments[idx].standings[indexGroup].hasilPertandingan[index].members_a_statistic[i].kill)
		Print("Death: ", false)
		fmt.Scan(&tournaments[idx].standings[indexGroup].hasilPertandingan[index].members_a_statistic[i].death)
		Print("Assist: ", false)
		fmt.Scan(&tournaments[idx].standings[indexGroup].hasilPertandingan[index].members_a_statistic[i].assist)
		tournaments[idx].standings[indexGroup].hasilPertandingan[index].members_a_statistic[i].id = tournaments[idx].standings[indexGroup].teams[teamAIndex].members[i].id
	}
	// B
	Print("---", true)
	Print("Statistik Tim "+tournaments[idx].standings[indexGroup].teams[teamBIndex].name, true)
	for i = 0; i < 5; i++ {
		Print("Buatlah statisik KDA untuk "+tournaments[idx].standings[indexGroup].teams[teamBIndex].members[i].name, true)
		Print("Kill: ", false)
		fmt.Scan(&tournaments[idx].standings[indexGroup].hasilPertandingan[index].members_b_statistic[i].kill)
		Print("Death: ", false)
		fmt.Scan(&tournaments[idx].standings[indexGroup].hasilPertandingan[index].members_b_statistic[i].death)
		Print("Assist: ", false)
		fmt.Scan(&tournaments[idx].standings[indexGroup].hasilPertandingan[index].members_b_statistic[i].assist)
		tournaments[idx].standings[indexGroup].hasilPertandingan[index].members_b_statistic[i].id = tournaments[idx].standings[indexGroup].teams[teamBIndex].members[i].id
	}
}

func updateGroupMatch(idx, indexGroup int) {
	var teamAId, teamBId string
	var teamAIndex, teamBIndex int

	Print("ID tim A yang akan bertarung: ", false)
	fmt.Scan(&teamAId)
	Print("ID tim B yang akan bertarung: ", false)
	fmt.Scan(&teamBId)

	if teamAId == teamBId {
		Print("Tim yang sama tidak dapat melawan tim-nya sendiri!", true)
		updateGroupMatch(idx, indexGroup)
	} else {
		teamAIndex = findTournamentTeamsIndex(
			tournaments[idx],
			indexGroup,
			teamAId,
		)
		teamBIndex = findTournamentTeamsIndex(
			tournaments[idx],
			indexGroup,
			teamBId,
		)

		if teamAIndex < 0 || teamBIndex < 0 {
			Print("Kami tidak menemukan tim dengan ID: ", false)
			if teamAIndex < 0 {
				fmt.Print(teamAId, " ")
			}
			if teamBIndex < 0 {
				fmt.Print(teamBId, " ")
			}
			fmt.Print("\n")
			updateGroupMatch(idx, indexGroup)
		} else {
			updateGroupMatchPoints(idx, indexGroup, teamAId, teamBId)
			clearTerminal()
			updateGroup(idx)
		}
	}
}

func updateGroup(idx int) {
	var i, pickGroup int
	var teamAIndex, teamBIndex int
	var choose string

	switch tournaments[idx].system {
	case 1:
		Print(fmt.Sprintf("Tersisa %d dari %d pertandingan yang belum terlaksana.",
			combination(tournaments[idx].standings[0].nTeams, 2)-tournaments[idx].standings[0].jumlahPertandingan,
			combination(tournaments[idx].standings[0].nTeams, 2)),
			true,
		)
		pickGroup = 0
	case 2:
		Print("Daftar grup yang belum ataupun sudah menyelesaikan seluruh pertandingan!", true)
		for i = 0; i < tournaments[idx].nStandings; i++ {
			Print(fmt.Sprintf("Group %d: %d/%d Pertandingan",
				i+1,
				tournaments[idx].standings[i].jumlahPertandingan,
				combination(tournaments[idx].standings[i].nTeams, 2)),
				true,
			)
		}
		Print("Group mana yang hendak anda perbarui: ", false)
		fmt.Scan(&pickGroup)

		if pickGroup < 1 || pickGroup > tournaments[idx].nStandings {
			Print("Anda memilih kehampaan!", true)
			updateGroup(idx)
			return
		}
		pickGroup--
	}

	clearTerminal()
	Print(fmt.Sprintf("Terdapat %d tim di dalam Group %d", tournaments[idx].standings[pickGroup].nTeams, pickGroup+1), true)
	fmt.Println()
	showTournamentStandings(tournaments[idx], pickGroup)
	fmt.Println()

	if tournaments[idx].standings[pickGroup].jumlahPertandingan > 0 {
		fmt.Printf("Pertandingan yang telah terjadi di grup %d:\n", pickGroup+1)
		fmt.Println("+----------------------------------------------------+-----------------+")
	}
	for i = 0; i < tournaments[idx].standings[pickGroup].jumlahPertandingan; i++ {
		teamAIndex = findTournamentTeamsIndex(
			tournaments[idx],
			pickGroup,
			tournaments[idx].standings[pickGroup].hasilPertandingan[i].team_a_id,
		)
		teamBIndex = findTournamentTeamsIndex(
			tournaments[idx],
			pickGroup,
			tournaments[idx].standings[pickGroup].hasilPertandingan[i].team_b_id,
		)
		fmt.Printf("│ %-50s │ %-15s │\n",
			tournaments[idx].standings[pickGroup].teams[teamAIndex].name+" vs "+tournaments[idx].standings[pickGroup].teams[teamBIndex].name,
			fmt.Sprintf("%d - %d", tournaments[idx].standings[pickGroup].hasilPertandingan[i].point_a, tournaments[idx].standings[pickGroup].hasilPertandingan[i].point_b),
		)

		if i == tournaments[idx].standings[pickGroup].jumlahPertandingan-1 {
			fmt.Println("+----------------------------------------------------+-----------------+")
			fmt.Println()
		}
	}

	Print("Ketik 'kembali' jika anda hendak kembali ke main menu, jika tetap melanjutkan ketik apapun: ", false)
	fmt.Scan(&choose)
	if strings.ToLower(choose) != "kembali" {
		updateGroupMatch(idx, pickGroup)
	} else {
		clearTerminal()
	}
}

func updatePlayOff(idx int) {
	var teams [16]StatisticTeam
	var teamIdx int
	var i, j int

	for i = 0; i < tournaments[idx].nStandings; i++ {
		teams = sortTournamentTeamsBasedStatistics(
			tournaments[idx].standings[i].statisticTeams,
			tournaments[idx].standings[i].nTeams,
		)
		fmt.Println("Daftar tim yang lolos playoff di Grup", i+1)
		fmt.Println("+-------+--------------------------------+--------------+")
		fmt.Printf(
			"│ %-5d │ %-30s │ %-12s │\n",
			"No.",
			"Nam Tim",
			"Menang - Kalah",
		)
		fmt.Println("+-------+--------------------------------+--------------+")
		for j = 0; j < (tournaments[idx].nPlayOff / 2); j++ {
			teamIdx = findTournamentTeamsIndex(tournaments[idx], i, teams[j].id)
			fmt.Printf(
				"│ %-5d │ %-30s │ %-12s │\n",
				j+1,
				tournaments[idx].standings[i].teams[teamIdx].name,
				fmt.Sprintf("%d - %d", teams[j].winner, teams[j].lose),
			)
		}
		fmt.Println("+-------+--------------------------------+--------------+")
	}
}

func updateSpecificTournament(idx int) {
	var pick int
	var text string

	text = "Saat ini anda berada pada tahap "
	if tournaments[idx].finishGroup {
		text += "PlayOff. "
	} else {
		switch tournaments[idx].system {
		case 1:
			text += "Klasemen. "
		case 2:
			text += "Group. "
		}
	}

	text += "Pilihlah salah-satu dari 3 pilihan ini"
	Print(text, true)

	if tournaments[idx].finishGroup {
		Print("1). Melanjutkan pembaruan pada babak playoff", true)
	} else {
		Print("1). Melanjutkan pembaruan pada babak grup", true)
	}
	Print("2). Kembali", true)
	Print("Pilih antara 2 opsi tersebut: ", false)
	fmt.Scan(&pick)

	switch pick {
	case 1:
		clearTerminal()
		switch tournaments[idx].finishGroup {
		case true:
			updatePlayOff(idx)
		case false:
			updateGroup(idx)
		}
	case 2:
		Print("-----------------------------------------------", true)
		clearTerminal()
	default:
		Print("Anda tidak memilih pilihan yang valid!", true)
		updateSpecificTournament(idx)
	}
}

func updateTournament() {
	var id string
	var idx int
	Print("Data turnamen yang ada di dalam database anda:", true)
	showTournamentTableList()

	Print("Masukkan ID untuk memilih turnamen yang hendak diperbarui (ketik 'cancel' jika hendak kembali ke main menu): ", false)
	fmt.Scan(&id)
	if strings.ToLower(id) != "cancel" {
		idx = searchTournamentIndex(id)
		if idx < 0 {
			Print("Anda memasukkan ID yang tidak ada di dalam database!", true)
			updateTournament()
		} else {
			if tournaments[idx].finishPlayOff {
				clearTerminal()
				Print("Turnamen "+tournaments[idx].name+" telah selesai diselenggarakan! Anda tidak bisa merubahnya.", true)
				updateTournament()
			} else {
				clearTerminal()
				updateSpecificTournament(idx)
				Print("-----------------------------------------------", true)

				clearTerminal()
				updateTournament()
			}
		}
	} else {
		clearTerminal()
		main()
	}
}

// Opsi 2

func viewTournamentWait() {
	var action string
	Print("Ketik 'next' ketika sudah selesai: ", false)
	fmt.Scan(&action)

	if strings.ToLower(action) != "next" {
		viewTournamentWait()
	}
}

func viewSpecificTournament(tournament Tournament) {
	var i, j, k int
	var idx, indexGroup int
	var totalKill, totalAssist, totalDeath int
	var id string

	fmt.Println()
	switch tournament.system {
	case 1:
		fmt.Println("Klasemen Turnamen", tournament.name)
	case 2:
		fmt.Println("Grup-Grup Turnamen", tournament.name)
		fmt.Println()
	}
	for i = 0; i < tournament.nStandings; i++ {
		if tournament.system == 2 {
			fmt.Println("Group", i+1)
		}
		showTournamentStandings(tournament, i)
	}

	Print("Masukkan ID tim jika kamu hendak melihat lebih detail anggota-anggota di dalamnya (ketik 'cancel' jika hendak kembali ke main menu): ", false)
	fmt.Scan(&id)

	if strings.ToLower(id) != "cancel" {
		idx = -1
		for i = 0; i < tournament.nStandings && idx < 0; i++ {
			for j = 0; j < tournament.standings[i].nTeams && idx < 0; j++ {
				idx = findTournamentTeamsIndex(tournament, i, id)
				indexGroup = i
			}
		}

		if idx >= 0 {
			clearTerminal()
			fmt.Println()
			fmt.Println(tournament.standings[indexGroup].teams[idx].name, "Members")
			fmt.Println("+----------------------+----------------------+------------+------------+------------+")
			fmt.Printf("│ %-20s │ %-20s │ %-10s │ %-10s │ %-10s │\n", "ID", "Username", "Kills", "Death", "Assist")
			fmt.Println("+----------------------+----------------------+------------+------------+------------+")
			for i = 0; i < 5; i++ {
				for j = 0; j < tournament.standings[indexGroup].jumlahPertandingan; j++ {
					for k = 0; k < 5; k++ {
						switch tournament.standings[indexGroup].teams[idx].members[i].id {
						case tournament.standings[indexGroup].hasilPertandingan[j].members_a_statistic[k].id:
							totalKill += tournament.standings[indexGroup].hasilPertandingan[j].members_a_statistic[k].kill
							totalDeath += tournament.standings[indexGroup].hasilPertandingan[j].members_a_statistic[k].death
							totalAssist += tournament.standings[indexGroup].hasilPertandingan[j].members_a_statistic[k].assist
						case tournament.standings[indexGroup].hasilPertandingan[j].members_b_statistic[k].id:
							totalKill += tournament.standings[indexGroup].hasilPertandingan[j].members_b_statistic[k].kill
							totalDeath += tournament.standings[indexGroup].hasilPertandingan[j].members_b_statistic[k].death
							totalAssist += tournament.standings[indexGroup].hasilPertandingan[j].members_b_statistic[k].assist
						}
					}
				}
				fmt.Printf("│ %-20s │ %-20s │ %-10d │ %-10d │ %-10d │\n",
					tournament.standings[indexGroup].teams[idx].members[i].id,
					tournament.standings[indexGroup].teams[idx].members[i].name,
					totalKill,
					totalDeath,
					totalAssist,
				)
				totalKill = 0
				totalDeath = 0
				totalAssist = 0
			}
			fmt.Println("+----------------------+----------------------+------------+------------+------------+")
			fmt.Println()
			viewTournamentWait()
		}

		clearTerminal()
		viewSpecificTournament(tournament)
	} else {
		clearTerminal()
	}
}

func viewTournament() {
	var idx int
	var id string

	Print("Data turnamen yang ada di dalam database anda:", true)

	fmt.Println()
	showTournamentTableList()
	fmt.Println()

	Print("Masukkan ID untuk memilih turnamen mana yang hendak anda lihat (ketik 'cancel' untuk kembali ke main menu): ", false)
	fmt.Scan(&id)

	if strings.ToLower(id) != "cancel" {
		clearTerminal()
		idx = searchTournamentIndex(id)
		if idx < 0 {
			Print("Turnamen dengan ID '"+id+"' tidak ditemukan!", true)
			viewTournament()
		} else {
			viewSpecificTournament(tournaments[idx])
			Print("-----------------------------------------------", true)

			clearTerminal()
			viewTournament()
		}
	} else {
		Print("-----------------------------------------------", true)
		clearTerminal()
		main()
	}
}

// Opsi 3

func deleteTournamentWait(idx int) bool {
	var answer string
	Print("Apakah anda yakin hendak menghapus "+tournaments[idx].name+" ("+tournaments[idx].id+") dari database anda (YA/TIDAK): ", false)
	fmt.Scan(&answer)

	if strings.ToLower(answer) == "ya" {
		return true
	} else if strings.ToLower(answer) == "tidak" {
		return false
	} else {
		return deleteTournamentWait(idx)
	}
}

func deleteTournament() {
	var id string
	var idx, i, j, k int
	var isDelete bool
	Print("Data turnamen yang ada di dalam database anda:", true)
	showTournamentTableList()

	Print("Masukkan ID turnamen yang hendak anda hapus (ketik 'cancel' jika anda membatalkan proses): ", false)
	fmt.Scan(&id)
	if strings.ToLower(id) == "cancel" {
		Print("-----------------------------------------------", true)
		main()
		return
	}

	idx = searchTournamentIndex(id)
	if idx < 0 {
		Print("ID yang anda masukkan tidak valid!", true)
		deleteTournament()
		return
	}

	isDelete = deleteTournamentWait(idx)
	if !isDelete {
		deleteTournament()
		return
	}

	for i = idx; i < nTournaments-1; i++ {
		tournaments[i] = tournaments[i+1]
	}
	nTournaments--

	for i = 0; i < nTournaments; i++ {
		tournaments[i].id = generateID("TO", i+1)
		for j = 0; j < tournaments[i].nStandings; j++ {
			for k = 0; k < tournaments[i].standings[j].nTeams; k++ {
				switch tournaments[i].system {
				case 1:
					tournaments[i].standings[j].teams[k].id = tournaments[i].id + "-" + generateID("TE", k+1)
				case 2:
					tournaments[i].standings[j].teams[k].id = tournaments[i].id + "-" + generateID("GR", j+1) + "-" + generateID("TE", k+1)
				}
				idx = findTournamentTeamsIndex(tournaments[i], j, tournaments[i].standings[j].teams[k].id)
				tournaments[i].standings[j].statisticTeams[idx].id = tournaments[i].standings[j].teams[k].id
			}
		}
	}

	Print("-----------------------------------------------", true)
	clearTerminal()
	main()
}

// Opsi 4

func guideTournament() {

}

// Opsi 5

func createTournamentSystem() {
	Print("Berikut adalah sistem untuk turnamen VALORANT yang hendak anda buat:", true)
	Print("1). Klasemen", true)
	Print("2). Group", true)
	Print("Sistem mana yang hendak anda buat: ", false)
	fmt.Scan(&tournaments[nTournaments].system)

	if tournaments[nTournaments].system < 1 || tournaments[nTournaments].system > 2 {
		Print("Sistem yang anda masukan tidak ada!", true)
		createTournamentSystem()
	}
}

func createTournamentGroups() {
	var pick int
	Print("Pilihan jumlah group yang diperbolehkan saat ini:", true)
	Print("1). 2 Group", true)
	Print("2). 4 Group", true)
	Print("Pilihan anda: ", false)
	fmt.Scan(&pick)

	switch pick {
	case 1:
		tournaments[nTournaments].nStandings = 2
	case 2:
		tournaments[nTournaments].nStandings = 4
	default:
		Print("Angka yang anda masukan tidak valid!", true)
		createTournamentGroups()
	}
}

func createTournamentGroupsTeams() {
	var totalTeams, i int
	Print("Jumlah tim per-grup akan menentukan jumlah tim yang akan lolos playoff, ini adalah sistemnya:", true)
	switch tournaments[nTournaments].nStandings {
	case 2:
		Print("▫ 8 tim per-grup (total 16 tim), terdapat 6 tim per-grup (total 12 tim) yang lolos ke playoff, akan ada 1 yang lolos ke final tanpa melewati semifinal", true)
		Print("▫ 6 tim per-grup (total 12 tim), terdapat 4 tim per-grup (total 8 tim) yang lolos ke playoff", true)
		Print("▫ 4 tim per-grup (total 8 tim), terdapat 2 tim per-grup (total 4 tim) yang lolos ke playoff", true)
	case 4:
		Print("▫ 4 tim per-grup (total 16 tim), terdapat 2 tim per-grup (total 8 tim) yang lolos ke playoff", true)
		Print("▫ 3 tim per-grup (total 12 tim), terdapat 2 tim per-grup (total 6 tim) yang lolos ke playoff, akan ada 1 yang lolos ke final tanpa melewati semifinal", true)
		Print("▫ 2 tim per-grup (total 8 tim), terdapat 1 tim per-grup (total 4 tim) yang lolos ke playoff", true)
	}
	Print("Masukan jumlah tim yang akan bertanding (minimal 8, maksimal 16, dan harus kelipatan 4): ", false)
	fmt.Scan(&totalTeams)

	if totalTeams < 8 || totalTeams > 16 || totalTeams%4 != 0 {
		Print("Angka yang anda masukan tidak valid!", true)
		createTournamentGroupsTeams()
		return
	}

	for i = 0; i < tournaments[nTournaments].nStandings; i++ {
		tournaments[nTournaments].standings[i].nTeams = totalTeams / tournaments[nTournaments].nStandings
	}

	switch tournaments[nTournaments].nStandings {
	case 2:
		tournaments[nTournaments].nPlayOff = totalTeams - 4
	case 4:
		tournaments[nTournaments].nPlayOff = totalTeams / 2
	}
}

func createTournamentTeams() {
	var i, j, k int
	switch tournaments[nTournaments].system {
	case 1:
		Print("Masukan jumlah tim yang akan bertanding (minimal 8, maksimal 16, dan harus kelipatan 4): ", false)
		fmt.Scan(&tournaments[nTournaments].standings[0].nTeams)

		if tournaments[nTournaments].standings[0].nTeams < 8 || tournaments[nTournaments].standings[0].nTeams > 16 || tournaments[nTournaments].standings[0].nTeams%4 != 0 {
			Print("Angka yang anda masukan tidak valid!", true)
			createTournamentTeams()
		}

		tournaments[nTournaments].nStandings = 1
		tournaments[nTournaments].nPlayOff = tournaments[nTournaments].standings[0].nTeams / 2
	case 2:
		createTournamentGroups()
		createTournamentGroupsTeams()
	}

	switch tournaments[nTournaments].system {
	case 1:
		Print("Buatlah tim-tim untuk klasemen!", true)
	case 2:
		Print("Buatlah tim-tim untuk grup!", true)
	}

	for i = 0; i < tournaments[nTournaments].nStandings; i++ {
		if tournaments[nTournaments].system == 2 {
			Print(fmt.Sprintf("Group %d", i+1), true)
		}
		for j = 0; j < tournaments[nTournaments].standings[i].nTeams; j++ {
			Print(fmt.Sprintf("Nama tim ke-%d: ", j+1), false)
			fmt.Scan(&tournaments[nTournaments].standings[i].teams[j].name)
			switch tournaments[nTournaments].system {
			case 1:
				tournaments[nTournaments].standings[i].teams[j].id = tournaments[nTournaments].id + "-" + generateID("TE", j+1)
			case 2:
				tournaments[nTournaments].standings[i].teams[j].id = tournaments[nTournaments].id + "-" + generateID("GR", i+1) + "-" + generateID("TE", j+1)
			}
			tournaments[nTournaments].standings[i].statisticTeams[j].id = tournaments[nTournaments].standings[i].teams[j].id

			for k = 0; k < 5; k++ {
				Print(fmt.Sprintf("Nama anggota ke-%d: ", k+1), false)
				fmt.Scan(&tournaments[nTournaments].standings[i].teams[j].members[k].name)
				tournaments[nTournaments].standings[i].teams[j].members[k].id = tournaments[nTournaments].standings[i].teams[j].id + "-" + generateID("ME", k+1)
			}
		}
	}
}

func createTournamentDefaultTeam(indexTeam, indexGroup int, name string, members [5]string) {
	var i int
	tournaments[nTournaments].standings[indexGroup].teams[indexTeam].name = name
	tournaments[nTournaments].standings[indexGroup].teams[indexTeam].id = tournaments[nTournaments].id + "-" + generateID("GR", indexGroup+1) + "-" + generateID("TE", indexTeam+1)
	tournaments[nTournaments].standings[indexGroup].statisticTeams[indexTeam].id = tournaments[nTournaments].standings[indexGroup].teams[indexTeam].id
	for i = 0; i < 5; i++ {
		tournaments[nTournaments].standings[indexGroup].teams[indexTeam].members[i].name = members[i]
		tournaments[nTournaments].standings[indexGroup].teams[indexTeam].members[i].id = tournaments[nTournaments].standings[indexGroup].teams[indexTeam].id + "-" + generateID("ME", i+1)
	}
}

func createTournamentDefault() {
	tournaments[nTournaments].name = "Pacific_Stage_1_2025"
	tournaments[nTournaments].id = generateID("TO", nTournaments+1)
	tournaments[nTournaments].system = 2
	tournaments[nTournaments].nStandings = 2
	tournaments[nTournaments].nPlayOff = 8

	tournaments[nTournaments].standings[0].nTeams = 6
	tournaments[nTournaments].standings[1].nTeams = 6

	// Group 1
	createTournamentDefaultTeam(
		0,
		0,
		"BOOM_Esports",
		[5]string{"dos9", "Famouz", "Shiro", "NcSlasher", "BerserX"},
	)
	createTournamentDefaultTeam(
		1,
		0,
		"DRX",
		[5]string{"MaKo", "free1ng", "HYUNMIN", "Estrella", "BeYN"},
	)
	createTournamentDefaultTeam(
		2,
		0,
		"GenG",
		[5]string{"Munchkin", "Foxy9", "Ash", "t3xture", "Suggest"},
	)
	createTournamentDefaultTeam(
		3,
		0,
		"Paper_Rex",
		[5]string{"mindfreak", "Jinggg", "f0rsakeN", "d4v41", "something"},
	)
	createTournamentDefaultTeam(
		4,
		0,
		"Global_Esports",
		[5]string{"patrickWHO", "UdoTan", "Kr1stal", "kellyS", "Papi"},
	)
	createTournamentDefaultTeam(
		5,
		0,
		"DetonatioN_FocusMe",
		[5]string{"Meiy", "Art", "Jinboong", "Akame", "gyen"},
	)

	// Group 2
	createTournamentDefaultTeam(
		0,
		1,
		"Rex_Regum_Qeon",
		[5]string{"crazyguy", "Monyet", "xffero", "Jemkin", "Kushy"},
	)
	createTournamentDefaultTeam(
		1,
		1,
		"TALON",
		[5]string{"Crws", "Killua", "thyy", "JitBoyS", "primmie"},
	)
	createTournamentDefaultTeam(
		2,
		1,
		"T1",
		[5]string{"stax", "Meteor", "Sylvan", "BuZz", "iZu"},
	)
	createTournamentDefaultTeam(
		3,
		1,
		"Nongshim_RedForce",
		[5]string{"Persia", "Francis", "margaret", "Dambi", "Ivy"},
	)
	createTournamentDefaultTeam(
		4,
		1,
		"ZETA_DIVISION",
		[5]string{"TenTen", "SugarZ3ro", "CLZ", "SyouTa", "Xdll"},
	)
	createTournamentDefaultTeam(
		5,
		1,
		"Team_Secret",
		[5]string{"JessieVash", "invy", "Wild0reoo", "2GE", "Jremy"},
	)
}

func createTournament() {
	var name string
	Print("Masukkan nama turnamen yang hendak anda buat (ketik 'cancel' untuk berhenti, setelah memasukkan nama, anda tidak dapat mundur): ", false)
	fmt.Scan(&name)

	if strings.ToLower(name) == "cancel" {
		Print("-----------------------------------------------", true)
		main()
		return
	}

	if strings.ToLower(name) == "default" {
		clearTerminal()
		Print("Kamu membuat turnamen template yang disediakan oleh sistem.", true)
		createTournamentDefault()
	} else {
		tournaments[nTournaments].name = name
		tournaments[nTournaments].id = generateID("TO", nTournaments+1)
		createTournamentSystem()
		createTournamentTeams()
		clearTerminal()
	}

	nTournaments++
	Print("-----------------------------------------------", true)
	main()
}

// Main Program

func main() {
	var n int
	startingMessage()
	chooseAction(&n)

	if (nTournaments == NMAX && n == 5) || (nTournaments < NMAX && n == 6) {
		return
	}

	switch n {
	case 1:
		updateTournament()
	case 2:
		viewTournament()
	case 3:
		deleteTournament()
	case 4:
		guideTournament()
	case 5:
		createTournament()
	}
}
