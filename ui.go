package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"golang.org/x/term"
)

const ansiReset = "\033[0m"

// pool.Half >= 10 means chronicled/special mixed pools and uses different display layout.
const poolHalfSpecialThreshold = 10
const poolSplitLineWidth = 103
const roleAttribUPMin = 3
const roleAttribUPDivisor = 4
const roleAttribUPRemainder = 3

var poolSplitLine = strings.Repeat("-", poolSplitLineWidth)

var visionColor = map[uint8]string{
	VisionOther:   "\033[38;5;7m",
	Pyro:          "\033[38;5;9m",
	Hydro:         "\033[38;5;33m",
	Anemo:         "\033[38;5;43m",
	Electro:       "\033[38;5;99m",
	Dendro:        "\033[38;5;46m",
	Cryo:          "\033[38;5;159m",
	Geo:           "\033[38;5;220m",
	VisionUnknown: "\033[38;5;7m",
}

func colorizeByVision(vision uint8, text string) string {
	if text == "" {
		return text
	}
	color, ok := visionColor[vision]
	if !ok {
		color = visionColor[VisionOther]
	}
	return color + text + ansiReset
}

func visionName(vision uint8) string {
	switch vision {
	case Pyro:
		return "Pyro"
	case Hydro:
		return "Hydro"
	case Anemo:
		return "Anemo"
	case Electro:
		return "Electro"
	case Dendro:
		return "Dendro"
	case Cryo:
		return "Cryo"
	case Geo:
		return "Geo"
	case VisionOther:
		return "Other"
	default:
		return "Unknown"
	}
}

func roleTypeName(ch CharEntry) string {
	switch {
	case ch.isFourStar():
		return "4-star"
	case ch.Attrib == RoleTypeLimitedFiveStar:
		return "Limited 5-star"
	case isUPPermanentFiveStar(ch):
		return "UP 5-star (permanent)"
	case ch.Attrib == RoleTypeTravelerAether || ch.Attrib == RoleTypeTravelerLumine:
		return "Traveler"
	case ch.Attrib == RoleTypeCollab:
		return "Collab"
	default:
		return "Other"
	}
}

// isUPPermanentFiveStar identifies encoded permanent 5-star entries that had their own UP pool.
func isUPPermanentFiveStar(ch CharEntry) bool {
	return ch.Attrib > roleAttribUPMin && ch.Attrib%roleAttribUPDivisor == roleAttribUPRemainder
}

func padVisual(text string, width int) string {
	padding := width - visualLen(text)
	if padding <= 0 {
		return text
	}
	return text + strings.Repeat(" ", padding)
}

func printColoredPadded(text string, vision uint8, width int) {
	fmt.Print(colorizeByVision(vision, text))
	padding := width - visualLen(text)
	if padding > 0 {
		fmt.Print(strings.Repeat(" ", padding))
	}
}

// runeWidth returns the visual terminal width of a rune (2 for CJK/fullwidth, 1 otherwise).
func runeWidth(r rune) int {
	if r >= 0x1100 && r <= 0x115F ||
		r >= 0x2E80 && r <= 0x303F ||
		r >= 0x3040 && r <= 0x33FF ||
		r >= 0x3400 && r <= 0x4DBF ||
		r >= 0x4E00 && r <= 0x9FFF ||
		r >= 0xA000 && r <= 0xA4CF ||
		r >= 0xAC00 && r <= 0xD7AF ||
		r >= 0xF900 && r <= 0xFAFF ||
		r >= 0xFE30 && r <= 0xFE6F ||
		r >= 0xFF00 && r <= 0xFF60 ||
		r >= 0xFFE0 && r <= 0xFFE6 ||
		r >= 0x20000 && r <= 0x2A6DF {
		return 2
	}
	return 1
}

// visualLen returns the visual terminal width of a string.
func visualLen(s string) int {
	n := 0
	for _, r := range s {
		n += runeWidth(r)
	}
	return n
}

// getch reads a single byte from stdin using raw terminal mode.
func getch() byte {
	fd := int(os.Stdin.Fd())
	oldState, err := term.MakeRaw(fd)
	if err != nil {
		var b byte
		fmt.Scanf("%c", &b)
		return b
	}
	defer term.Restore(fd, oldState)
	var buf [1]byte
	if _, err := os.Stdin.Read(buf[:]); err != nil {
		return 0
	}
	return buf[0]
}

// cls clears the terminal screen.
func cls() {
	fmt.Print("\033[H\033[2J")
}

// readIntInRange reads an integer in [min, max] from stdin.
func readIntInRange(min, max int) int {
	reader := bufio.NewReader(os.Stdin)
	for {
		line, err := reader.ReadString('\n')
		if err != nil && line == "" {
			return min
		}
		var n int
		if _, scanErr := fmt.Sscanf(strings.TrimSpace(line), "%d", &n); scanErr == nil && n >= min && n <= max {
			return n
		}
		fmt.Printf("请输入 %d 到 %d 之间的数字: ", min, max)
	}
}

// typeMenu draws a box-drawing menu with numbered items and returns the 0-based index of the selection.
// Items are shown as [1]..[N-1] and [0] for the last item.
func typeMenu(items []string, title string) int {
	cls()
	n := len(items)

	labels := make([]string, n)
	for i, item := range items {
		if i == n-1 {
			labels[i] = "[0] " + item
		} else {
			labels[i] = fmt.Sprintf("[%d] %s", i+1, item)
		}
	}

	contentWidth := visualLen(title)
	for _, label := range labels {
		if w := visualLen(label); w > contentWidth {
			contentWidth = w
		}
	}

	top := strings.Repeat("═", contentWidth+2)
	mid := strings.Repeat("─", contentWidth+2)

	fmt.Printf("╔%s╗\n", top)
	fmt.Printf("║ %s%s ║\n", title, strings.Repeat(" ", contentWidth-visualLen(title)))
	fmt.Printf("╟%s╢\n", mid)
	for _, label := range labels {
		fmt.Printf("║ %s%s ║\n", label, strings.Repeat(" ", contentWidth-visualLen(label)))
	}
	fmt.Printf("╚%s╝\n", top)

	fmt.Printf("请输入选择: ")
	input := readIntInRange(0, n-1)
	if input == 0 {
		return n - 1
	}
	return input - 1
}

// choiceMenu draws a box-drawing menu requiring a single keystroke.
// Items 1-9 shown as [1]-[9], items 10-35 as [A]-[Z], last item as [0].
// Returns the 0-based index of the selected item.
func choiceMenu(items []string, title string) int {
	cls()
	n := len(items)

	labels := make([]string, n)
	for i, item := range items {
		if i == n-1 {
			labels[i] = "[0] " + item
		} else if i < 9 {
			labels[i] = fmt.Sprintf("[%d] %s", i+1, item)
		} else {
			labels[i] = fmt.Sprintf("[%c] %s", 'A'+rune(i-9), item)
		}
	}

	contentWidth := visualLen(title)
	for _, label := range labels {
		if w := visualLen(label); w > contentWidth {
			contentWidth = w
		}
	}

	top := strings.Repeat("═", contentWidth+2)
	mid := strings.Repeat("─", contentWidth+2)

	fmt.Printf("╔%s╗\n", top)
	fmt.Printf("║ %s%s ║\n", title, strings.Repeat(" ", contentWidth-visualLen(title)))
	fmt.Printf("╟%s╢\n", mid)
	for _, label := range labels {
		fmt.Printf("║ %s%s ║\n", label, strings.Repeat(" ", contentWidth-visualLen(label)))
	}
	fmt.Printf("╚%s╝\n", top)

	for {
		b := getch()
		switch {
		case b == '0':
			return n - 1
		case b >= '1' && b <= '9':
			idx := int(b - '1')
			if idx < n-1 {
				return idx
			}
		case b >= 'A' && b <= 'Z':
			idx := int(b-'A') + 9
			if idx < n-1 {
				return idx
			}
		case b >= 'a' && b <= 'z':
			idx := int(b-'a') + 9
			if idx < n-1 {
				return idx
			}
		}
	}
}

// charNameCN returns the Chinese name of a character by ID, or empty string for sentinel 0.
func charNameCN(id uint) string {
	if id == 0 || int(id) >= len(charMap) {
		return ""
	}
	return charMap[id].NameCN
}

func charByID(id uint) (CharEntry, bool) {
	index := int(id)
	if index >= len(charMap) {
		return CharEntry{}, false
	}
	return charMap[index], true
}

// printAllPools prints all wish pools with their version, dates, and character lists.
func printAllPools() {
	cls()
	maxCNVisualLen := 0
	for _, ch := range charMap {
		if w := visualLen(ch.NameCN); w > maxCNVisualLen {
			maxCNVisualLen = w
		}
	}
	if maxCNVisualLen == 0 {
		maxCNVisualLen = 1
	}

	for _, pool := range wishPool {
		if pool.Half >= poolHalfSpecialThreshold {
			fmt.Println(poolSplitLine)
		}
		fmt.Printf("%d.%d.%d\t%04d.%02d.%02d\t%04d.%02d.%02d\t",
			pool.Major, pool.Minor, pool.Half,
			pool.StartY, pool.StartM, pool.StartD,
			pool.EndY, pool.EndM, pool.EndD)

		fmt.Print(" | ")
		up5 := make([]uint, 0, 2)
		for _, id := range pool.Up5 {
			if id == 0 {
				break
			}
			up5 = append(up5, id)
		}
		if pool.Half < poolHalfSpecialThreshold {
			for i := 0; i < 2; i++ {
				if i < len(up5) {
					ch, ok := charByID(up5[i])
					if !ok {
						fmt.Print(strings.Repeat(" ", maxCNVisualLen+1))
						fmt.Print(" | ")
						continue
					}
					printColoredPadded(ch.NameCN, ch.Vision, maxCNVisualLen+1)
				} else {
					fmt.Print(strings.Repeat(" ", maxCNVisualLen+1))
				}
				fmt.Print(" | ")
			}
		} else {
			for i, id := range up5 {
				ch, ok := charByID(id)
				if !ok {
					continue
				}
				if i > 0 {
					fmt.Print(" ")
				}
				fmt.Print(colorizeByVision(ch.Vision, ch.NameCN))
			}
			fmt.Print(" | ")
		}

		for _, id := range pool.Up4 {
			if id == 0 {
				break
			}
			ch, ok := charByID(id)
			if !ok {
				continue
			}
			printColoredPadded(ch.NameCN, ch.Vision, maxCNVisualLen+1)
		}
		fmt.Println()
		if pool.Half >= poolHalfSpecialThreshold {
			fmt.Println(poolSplitLine)
		}
	}
	fmt.Print("\n按回车继续...")
	getch()
}

// printDaysofAllLimited5StarCharacters prints days since last UP for limited 5-star characters.
func printDaysofAllLimited5StarCharacters() {
	cls()
	maxCN := 0
	maxEN := 0
	for _, ch := range charMap {
		if !ch.isLimitedOrUPFiveStar() {
			continue
		}
		if w := visualLen(ch.NameCN); w > maxCN {
			maxCN = w
		}
		if w := visualLen(ch.Name); w > maxEN {
			maxEN = w
		}
	}

	for _, idx := range arrangedInOrderOfDays {
		ch := charMap[idx]
		days := daysPassedSinceLastUP[idx]
		if days == minInt32 || !ch.isLimitedOrUPFiveStar() {
			continue
		}
		printColoredPadded(ch.NameCN, ch.Vision, maxCN+1)
		fmt.Print(" ")
		printColoredPadded(ch.Name, ch.Vision, maxEN+1)
		fmt.Printf(" %d 天\n", days)
	}
	fmt.Print("\n按回车继续...")
	getch()
}

// printPoolLinkList prints a character's pool history, 3 entries per line.
func printPoolLinkList(head *PoolNode) {
	count := 0
	for node := head; node != nil; node = node.Next {
		if count > 0 {
			if count%3 == 0 {
				fmt.Println()
			} else {
				fmt.Print("\t")
			}
		}
		fmt.Printf("%d.%d.%d", node.Major, node.Minor, node.Half)
		count++
	}
	if count > 0 {
		fmt.Println()
	}
	if count == 0 {
		fmt.Println("(无记录)")
	}
}

// choiceCharFromList shows a menu to pick one character from the given indices.
// Returns the chosen charMap index, or -1 if the user chose "返回".
func choiceCharFromList(indices []int) int {
	if len(indices) == 0 {
		fmt.Println("没有找到符合条件的角色")
		fmt.Print("按回车继续...")
		getch()
		return -1
	}

	items := make([]string, len(indices)+1)
	maxCN := 0
	maxEN := 0
	for _, idx := range indices {
		ch := charMap[idx]
		if w := visualLen(ch.NameCN); w > maxCN {
			maxCN = w
		}
		if w := visualLen(ch.Name); w > maxEN {
			maxEN = w
		}
	}
	for i, idx := range indices {
		ch := charMap[idx]
		items[i] = padVisual(ch.NameCN, maxCN+1) + " " + padVisual(ch.Name, maxEN+1)
	}
	items[len(indices)] = "返回"

	var choice int
	if len(items) <= 35 {
		choice = choiceMenu(items, "选择角色")
	} else {
		choice = typeMenu(items, "选择角色")
	}
	if choice == len(indices) {
		return -1
	}
	return indices[choice]
}

// choiceOneCharacter shows a submenu to filter and select a character.
// Returns the charMap index of the selected character, or -1 to go back.
func choiceOneCharacter() int {
	visionNames := []string{"Anemo", "Geo", "Electro", "Dendro", "Hydro", "Pyro", "Cryo", "返回"}
	visionValues := []uint8{Anemo, Geo, Electro, Dendro, Hydro, Pyro, Cryo}

	subMenuItems := []string{
		"按照角色中文名有几个字筛选",
		"按照角色英文名有几个字母筛选",
		"按照角色神之眼类型筛选",
		"直接输入角色编号（高级）",
		"返回",
	}

	for {
		choice := typeMenu(subMenuItems, "选择角色筛选方式")

		switch choice {
		case 0: // Filter by Chinese name rune count
			cls()
			fmt.Print("输入中文名字数: ")
			reader := bufio.NewReader(os.Stdin)
			line, _ := reader.ReadString('\n')
			var length int
			if _, err := fmt.Sscanf(strings.TrimSpace(line), "%d", &length); err != nil || length <= 0 {
				continue
			}
			var filtered []int
			for i, ch := range charMap {
				if len([]rune(ch.NameCN)) == length {
					filtered = append(filtered, i)
				}
			}
			idx := choiceCharFromList(filtered)
			if idx >= 0 {
				return idx
			}

		case 1: // Filter by English name letter count
			cls()
			fmt.Print("输入英文名字母数: ")
			reader := bufio.NewReader(os.Stdin)
			line, _ := reader.ReadString('\n')
			var length int
			if _, err := fmt.Sscanf(strings.TrimSpace(line), "%d", &length); err != nil || length <= 0 {
				continue
			}
			var filtered []int
			for i, ch := range charMap {
				letterCount := 0
				for _, r := range ch.Name {
					if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') {
						letterCount++
					}
				}
				if letterCount == length {
					filtered = append(filtered, i)
				}
			}
			idx := choiceCharFromList(filtered)
			if idx >= 0 {
				return idx
			}

		case 2: // Filter by vision type
			vChoice := typeMenu(visionNames, "选择神之眼类型")
			if vChoice == len(visionNames)-1 {
				continue
			}
			vision := visionValues[vChoice]
			var filtered []int
			for i, ch := range charMap {
				if ch.Vision == vision {
					filtered = append(filtered, i)
				}
			}
			idx := choiceCharFromList(filtered)
			if idx >= 0 {
				return idx
			}

		case 3: // Direct index input
			cls()
			fmt.Printf("输入角色编号 (0-%d): ", len(charMap)-1)
			reader := bufio.NewReader(os.Stdin)
			line, _ := reader.ReadString('\n')
			var idx int
			if _, err := fmt.Sscanf(strings.TrimSpace(line), "%d", &idx); err == nil {
				if idx >= 0 && idx < len(charMap) {
					return idx
				}
			}
			fmt.Println("无效的角色编号")
			fmt.Print("按回车继续...")
			getch()

		default: // 返回
			return -1
		}
	}
}

// showMainMenu runs the main interactive menu loop.
func showMainMenu() {
	items := []string{
		"显示所有卡池",
		"显示所有限定五星角色距上次UP的天数",
		"查询特定角色的卡池历史",
		"显示一个赛诺冷笑话",
		"退出",
	}

	for {
		choice := typeMenu(items, "原神祈愿卡池信息工具")
		switch choice {
		case 0:
			printAllPools()
		case 1:
			printDaysofAllLimited5StarCharacters()
		case 2:
			idx := choiceOneCharacter()
			if idx >= 0 {
				buildPoolLinkList(idx)
				cls()
				ch := charMap[idx]
				labels := []string{"角色编号", "中文名", "英文名", "神之眼", "角色类型"}
				maxLabel := 0
				for _, lb := range labels {
					if w := visualLen(lb); w > maxLabel {
						maxLabel = w
					}
				}
				fmt.Printf("%s : %d\n", padVisual(labels[0], maxLabel), idx)
				fmt.Printf("%s : ", padVisual(labels[1], maxLabel))
				fmt.Println(colorizeByVision(ch.Vision, ch.NameCN))
				fmt.Printf("%s : ", padVisual(labels[2], maxLabel))
				fmt.Println(colorizeByVision(ch.Vision, ch.Name))
				fmt.Printf("%s : %s\n", padVisual(labels[3], maxLabel), visionName(ch.Vision))
				fmt.Printf("%s : %s\n", padVisual(labels[4], maxLabel), roleTypeName(ch))
				fmt.Printf("卡池历史:\n")
				printPoolLinkList(poolLinkLists[idx])
				fmt.Print("\n按回车继续...")
				getch()
			}
		case 3:
			cynoJoke()
		default: // 退出
			return
		}
	}
}

// printHelp prints the help/about information.
func printHelp() {
	fmt.Println("Genshin Impact Wish Pool Information Tool\n\nCopyright (c) 2025-2026 BingtangXH.\nMay the Anemo God bless you.")
}
