package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"golang.org/x/term"
)

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
	os.Stdin.Read(buf[:])
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
		line, _ := reader.ReadString('\n')
		var n int
		_, err := fmt.Sscanf(strings.TrimSpace(line), "%d", &n)
		if err == nil && n >= min && n <= max {
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

// printAllPools prints all wish pools with their version, dates, and character lists.
func printAllPools() {
	cls()
	for _, pool := range wishPool {
		fmt.Printf("%d.%d.%d\t%04d.%02d.%02d\t%04d.%02d.%02d\t",
			pool.Major, pool.Minor, pool.Half,
			pool.StartY, pool.StartM, pool.StartD,
			pool.EndY, pool.EndM, pool.EndD)

		first := true
		for _, id := range pool.Up5 {
			if id == 0 {
				break
			}
			if !first {
				fmt.Print(", ")
			}
			fmt.Print(charNameCN(id))
			first = false
		}
		fmt.Print("\t")

		first = true
		for _, id := range pool.Up4 {
			if id == 0 {
				break
			}
			if !first {
				fmt.Print(", ")
			}
			fmt.Print(charNameCN(id))
			first = false
		}
		fmt.Println()
	}
	fmt.Print("\n按回车继续...")
	getch()
}

// printDaysofAllLimited5StarCharacters prints days since last UP for limited 5-star characters.
func printDaysofAllLimited5StarCharacters() {
	cls()
	for _, idx := range arrangedInOrderOfDays {
		ch := charMap[idx]
		days := daysPassedSinceLastUP[idx]
		if days == minInt32 || !ch.isLimitedOrUPFiveStar() {
			continue
		}
		fmt.Printf("%-12s %-20s %d 天\n", ch.NameCN, ch.Name, days)
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
	for i, idx := range indices {
		ch := charMap[idx]
		items[i] = fmt.Sprintf("%-8s %s", ch.NameCN, ch.Name)
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
				fmt.Printf("角色: %s (%s)\n", ch.NameCN, ch.Name)
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
