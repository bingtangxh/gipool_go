package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"golang.org/x/term"
)

const ansiReset = "\033[0m"

// pool.Half>=10 表示编年史/特殊混合池，使用不同的显示布局。
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
		return "火"
	case Hydro:
		return "水"
	case Anemo:
		return "风"
	case Electro:
		return "雷"
	case Dendro:
		return "草"
	case Cryo:
		return "冰"
	case Geo:
		return "岩"
	case VisionOther:
		return "其他"
	default:
		return "未知"
	}
}

func roleTypeName(ch CharEntry) string {
	switch {
	case ch.isFourStar():
		return "四星"
	case ch.Attrib == RoleTypeLimitedFiveStar:
		return "限定五星"
	case isUPPermanentFiveStar(ch):
		return "常驻五星"
	case ch.Attrib == RoleTypeTravelerAether || ch.Attrib == RoleTypeTravelerLumine:
		return "旅行者"
	case ch.Attrib == RoleTypeCollab:
		return "联动"
	default:
		return "其他/未知"
	}
}

// isUPPermanentFiveStar 检测是不是常驻五星
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

// runeWidth 计算一个字符在终端的视觉长度（CJK/全角是2，其他的是1）。
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

// visualLen 计算一个字符串在终端的视觉长度
func visualLen(s string) int {
	n := 0
	for _, r := range s {
		n += runeWidth(r)
	}
	return n
}

// getch 直接从标准输入以原始终端模式读取一个字节
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

// cls 清屏
func cls() { fmt.Print("\033[H\033[2J") }

// readIntInRange 从标准输入读取一个数字
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
		fmt.Printf("Invalid input. Please enter a number between %d and %d: ", min, max)
	}
}

// typeMenu 绘制一个带编号项目的框线菜单，并返回所选项的0基索引。
// 项目显示为 [1]..[N-1]，最后一项显示为 [0]。
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
	fmt.Printf("║ %s%s%s ║\n", strings.Repeat(" ", (contentWidth-visualLen(title))/2+(contentWidth-visualLen(title))%2), title, strings.Repeat(" ", (contentWidth-visualLen(title))/2))
	fmt.Printf("╟%s╢\n", mid)
	for _, label := range labels {
		fmt.Printf("║ %s%s ║\n", label, strings.Repeat(" ", contentWidth-visualLen(label)))
	}
	fmt.Printf("╚%s╝\n", top)

	fmt.Printf("\nPlease select an option and press ENTER: ")
	input := readIntInRange(0, n-1)
	if input == 0 {
		return n - 1
	}
	return input - 1
}

// choiceMenu 绘制一个框线菜单，用户只需按一个键选择。
// 项1-9显示为[1]-[9]，项10-35显示为[A]-[Z]，最后一项为[0]。
// 返回所选项的0基索引。
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
	fmt.Printf("║ %s%s%s ║\n", strings.Repeat(" ", (contentWidth-visualLen(title))/2+(contentWidth-visualLen(title))%2), title, strings.Repeat(" ", (contentWidth-visualLen(title))/2))
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
			index := int(b - '1')
			if index < n-1 {
				return index
			}
		case b >= 'A' && b <= 'Z':
			index := int(b-'A') + 9
			if index < n-1 {
				return index
			}
		case b >= 'a' && b <= 'z':
			index := int(b-'a') + 9
			if index < n-1 {
				return index
			}
		}
	}
}

// charNameCN 根据 ID 返回角色中文名，若为 0 则返回空字符串。
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

// printAllPools 打印所有卡池，包括版本、日期和角色列表。
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
			// fmt.Print(" | ")
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
	fmt.Print("\nPress ENTER to continue...")
	getch()
}

// printDaysofAllLimited5StarCharacters 打印所有限定五星角色距上次UP的天数。
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

	for _, index := range arrangedInOrderOfDays {
		ch := charMap[index]
		days := daysPassedSinceLastUP[index]
		if days == minInt32 || !ch.isLimitedOrUPFiveStar() {
			continue
		}
		printColoredPadded(ch.NameCN, ch.Vision, maxCN+1)
		fmt.Print(" ")
		printColoredPadded(ch.Name, ch.Vision, maxEN+1)
		fmt.Printf(" %d\n", days)
	}
	fmt.Print("\nPress ENTER to continue...")
	getch()
}

// printPoolLinkList 打印角色的卡池历史，每行 3 个条目。
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
		fmt.Println("This character have't been UP yet.")
	}
}

// choiceCharFromList 显示菜单以从给定索引中选择一个角色。
// 返回所选 charMap 索引，若选择“返回”则返回-1。
func choiceCharFromList(indices []int) int {
	if len(indices) == 0 {
		fmt.Println("No characters found matching the criteria, Press ENTER to continue...")
		getch()
		return -1
	}

	items := make([]string, len(indices)+1)
	maxCN := 0
	maxEN := 0
	for _, index := range indices {
		ch := charMap[index]
		if w := visualLen(ch.NameCN); w > maxCN {
			maxCN = w
		}
		if w := visualLen(ch.Name); w > maxEN {
			maxEN = w
		}
	}
	for i, index := range indices {
		ch := charMap[index]
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

// choiceOneCharacter 显示子菜单以筛选并选择一个角色。
// 返回所选角色的 charMap 索引，返回-1表示返回。
func choiceOneCharacter() int {
	visionNames := []string{"风元素", "岩元素", "雷元素", "草元素", "水元素", "火元素", "冰元素", "返回"}
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
		case 0: // 按中文名字数筛选
			cls()
			fmt.Print("Please type how long the Chinese name is and press ENTER: ")
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
			index := choiceCharFromList(filtered)
			if index >= 0 {
				return index
			}

		case 1: // 按英文名字母数筛选
			cls()
			fmt.Print("Please type how long the English name is and press ENTER: ")
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
			index := choiceCharFromList(filtered)
			if index >= 0 {
				return index
			}

		case 2: // 按神之眼类型筛选
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
			index := choiceCharFromList(filtered)
			if index >= 0 {
				return index
			}

		case 3: // 直接输入角色编号
			cls()
			fmt.Printf("Please type a char index number (0-%d): ", len(charMap)-1)
			reader := bufio.NewReader(os.Stdin)
			line, _ := reader.ReadString('\n')
			var index int
			if _, err := fmt.Sscanf(strings.TrimSpace(line), "%d", &index); err == nil {
				if index >= 0 && index < len(charMap) {
					return index
				}
			}
			fmt.Println("Invalid character index, press ENTER to continue...")
			getch()

		default: // 返回
			return -1
		}
	}
}

// showMainMenu 运行主交互菜单循环。
func showMainMenu() {
	items := []string{
		"显示所有卡池",
		"显示所有限定五星角色距上次UP的天数",
		"查询单个角色历史卡池",
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
			index := choiceOneCharacter()
			if index >= 0 {
				buildPoolLinkList(index)
				cls()
				ch := charMap[index]
				labels := []string{"角色编号", "中文名", "英文名", "神之眼", "角色品质"}
				maxLabel := 0
				for _, lb := range labels {
					if w := visualLen(lb); w > maxLabel {
						maxLabel = w
					}
				}
				fmt.Printf("%s : %d\n", padVisual(labels[0], maxLabel), index)
				fmt.Printf("%s : ", padVisual(labels[1], maxLabel))
				fmt.Println(colorizeByVision(ch.Vision, ch.NameCN))
				fmt.Printf("%s : ", padVisual(labels[2], maxLabel))
				fmt.Println(colorizeByVision(ch.Vision, ch.Name))
				fmt.Printf("%s : %s\n", padVisual(labels[3], maxLabel), visionName(ch.Vision))
				fmt.Printf("%s : %s\n", padVisual(labels[4], maxLabel), roleTypeName(ch))
				fmt.Print("\n")
				fmt.Println("UP History:")
				printPoolLinkList(poolLinkLists[index])
				fmt.Print("\nPress ENTER to continue...")
				getch()
			}
		case 3:
			cynoJoke()
		default: // 退出
			return
		}
	}
}

// printHelp 打印帮助/关于信息。
func printHelp() {
	fmt.Println("Genshin Impact Wish Pool Information Tool\n\nCopyright (c) 2025-2026 BingtangXH.\nMay the Anemo God bless you.")
}
