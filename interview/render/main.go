// render 把面试题 JSON 渲染成 Markdown，并重建 README 索引。支持两种题集：
//
//   - days：每日题，interview/data/<日期>.json -> interview/days/<日期>.md，按日期倒序
//   - modules：Agent 速成课，interview/agent/data/<模块 id>.json -> interview/agent/modules/<模块 id>.md，按 order 排列
//
// 用法（在仓库根目录执行）：
//
//	go run ./interview/render
//	go run ./interview/render -dir interview/agent -kind modules
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

type Vocab struct {
	En string `json:"en"`
	Zh string `json:"zh"`
}

type Question struct {
	ID          string   `json:"id"`
	Topic       string   `json:"topic"`
	Difficulty  string   `json:"difficulty"`
	Title       string   `json:"title"`
	TitleZh     string   `json:"title_zh"`
	Question    string   `json:"question"`
	QuestionZh  string   `json:"question_zh"`
	Answer      string   `json:"answer"`
	AnswerZh    string   `json:"answer_zh"`
	Followups   []string `json:"followups"`
	FollowupsZh []string `json:"followups_zh"`
	Vocab       []Vocab  `json:"vocab"`
}

// Set 是一组题：每日题（与 Interview Drill 页面的 "days" 集合一致），
// 或 Agent 速成课的一个模块（与速成课页面的 "modules" 集合一致）
type Set struct {
	Date        string     `json:"date"`  // 每日题
	ID          string     `json:"id"`    // 模块
	Order       int        `json:"order"` // 模块
	Kind        string     `json:"kind"`  // 模块："core" 或 "daily"
	Title       string     `json:"title"` // 模块
	TitleZh     string     `json:"title_zh"`
	Summary     string     `json:"summary"`
	SummaryZh   string     `json:"summary_zh"`
	GeneratedAt string     `json:"generatedAt"`
	Questions   []Question `json:"questions"`
}

// key 是题集的文件名（不含 .json）
func (s Set) key(kind string) string {
	if kind == "modules" {
		return s.ID
	}
	return s.Date
}

func main() {
	dir := flag.String("dir", "interview", "题集目录，JSON 在其下的 data/ 中")
	kind := flag.String("kind", "days", `题集类型："days" 或 "modules"`)
	flag.Parse()
	if *kind != "days" && *kind != "modules" {
		log.Fatalf("unknown -kind %q", *kind)
	}

	files, err := filepath.Glob(filepath.Join(*dir, "data", "*.json"))
	if err != nil {
		log.Fatal(err)
	}
	if len(files) == 0 {
		log.Fatalf("no data files under %s", filepath.Join(*dir, "data"))
	}

	sets := make([]Set, 0, len(files))
	for _, f := range files {
		set, err := loadSet(f, *kind)
		if err != nil {
			log.Fatalf("%s: %v", f, err)
		}
		sets = append(sets, set)
	}
	// 每日题最新的一天排在最前面；模块按 order 排列
	sort.Slice(sets, func(i, j int) bool {
		if *kind == "modules" {
			return sets[i].Order < sets[j].Order
		}
		return sets[i].Date > sets[j].Date
	})

	sub := *kind
	if err := os.MkdirAll(filepath.Join(*dir, sub), 0o755); err != nil {
		log.Fatal(err)
	}
	for _, set := range sets {
		path := filepath.Join(*dir, sub, set.key(*kind)+".md")
		if err := os.WriteFile(path, []byte(renderSet(set, *kind)), 0o644); err != nil {
			log.Fatal(err)
		}
	}
	if err := os.WriteFile(filepath.Join(*dir, "README.md"), []byte(renderIndex(sets, *kind)), 0o644); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("rendered %d %s into %s\n", len(sets), *kind, *dir)
}

func loadSet(path, kind string) (Set, error) {
	var set Set
	b, err := os.ReadFile(path)
	if err != nil {
		return set, err
	}
	if err := json.Unmarshal(b, &set); err != nil {
		return set, err
	}
	if want := strings.TrimSuffix(filepath.Base(path), ".json"); set.key(kind) != want {
		return set, fmt.Errorf("%s key %q does not match file name %q", kind, set.key(kind), want)
	}
	if len(set.Questions) == 0 {
		return set, fmt.Errorf("no questions")
	}
	for i, q := range set.Questions {
		if q.Title == "" || q.Question == "" || q.Answer == "" {
			return set, fmt.Errorf("question %d is missing title, question or answer", i+1)
		}
	}
	return set, nil
}

// moduleTitle 是模块标题，核心模块带上序号
func moduleTitle(s Set) string {
	if s.Kind == "core" {
		return fmt.Sprintf("%02d. %s · %s", s.Order, s.Title, s.TitleZh)
	}
	return s.Title + " · " + s.TitleZh
}

func renderSet(set Set, kind string) string {
	var b strings.Builder
	if kind == "modules" {
		fmt.Fprintf(&b, "# %s\n\n", moduleTitle(set))
		fmt.Fprintf(&b, "%s\n\n> %s\n\n", set.Summary, set.SummaryZh)
		fmt.Fprintf(&b, "%d questions · generated %s · [← 目录 / Contents](../README.md)\n\n", len(set.Questions), set.GeneratedAt)
	} else {
		fmt.Fprintf(&b, "# %s · 每日面试题 / Daily Interview Questions\n\n", set.Date)
		fmt.Fprintf(&b, "> %d questions · generated %s · [← 全部 / All days](../README.md)\n\n", len(set.Questions), set.GeneratedAt)
	}
	writeTable(&b, set, "")

	for i, q := range set.Questions {
		b.WriteString("\n---\n\n")
		fmt.Fprintf(&b, "<a id=\"q%d\"></a>\n\n", i+1)
		fmt.Fprintf(&b, "## %d. %s\n\n", i+1, q.Title)
		fmt.Fprintf(&b, "**%s** · `%s` · `%s`\n\n", q.TitleZh, q.Topic, q.Difficulty)

		b.WriteString("### Question · 题目\n\n")
		b.WriteString(bilingual(q.ID+" question", q.Question, q.QuestionZh))

		// 答案默认折叠，先自己想再展开
		b.WriteString("<details>\n<summary><b>Answer · 参考答案</b></summary>\n\n")
		b.WriteString(bilingual(q.ID+" answer", q.Answer, q.AnswerZh))
		b.WriteString("</details>\n\n")

		if len(q.Followups) > 0 {
			b.WriteString("### Follow-ups · 追问\n\n")
			for j, f := range q.Followups {
				fmt.Fprintf(&b, "- %s\n", f)
				if j < len(q.FollowupsZh) {
					fmt.Fprintf(&b, "  > %s\n", q.FollowupsZh[j])
				}
			}
			b.WriteString("\n")
		}

		if len(q.Vocab) > 0 {
			b.WriteString("### Vocabulary · 词汇\n\n| English | 中文 |\n| --- | --- |\n")
			for _, v := range q.Vocab {
				fmt.Fprintf(&b, "| %s | %s |\n", cell(v.En), cell(v.Zh))
			}
		}
	}
	return b.String()
}

func renderIndex(sets []Set, kind string) string {
	if kind == "modules" {
		return renderModuleIndex(sets)
	}
	days := sets
	total := 0
	for _, d := range days {
		total += len(d.Questions)
	}

	var b strings.Builder
	b.WriteString("# 每日面试题 · Daily Interview Questions\n\n")
	b.WriteString("后端面试题（Go 为主、Java 为辅），中英双语，每天 5 道，覆盖 Go、System Design、Database、Middleware、Java & CS。\n")
	b.WriteString("每题包含题目、参考答案（默认折叠）、追问和词汇表。\n\n")
	fmt.Fprintf(&b, "共 %d 天、%d 道题。\n\n", len(days), total)
	b.WriteString("- 数据：`interview/data/YYYY-MM-DD.json`（与 Interview Drill 页面的 `days` 集合结构一致）\n")
	b.WriteString("- 重新生成 Markdown：在仓库根目录执行 `go run ./interview/render`\n")
	b.WriteString("- Agent 面试速成：[agent/README.md](agent/README.md)\n")

	for _, d := range days {
		fmt.Fprintf(&b, "\n## [%s](days/%s.md)\n\n", d.Date, d.Date)
		writeTable(&b, d, "days/"+d.Date+".md")
	}
	return b.String()
}

func renderModuleIndex(mods []Set) string {
	total := 0
	for _, m := range mods {
		total += len(m.Questions)
	}

	var b strings.Builder
	b.WriteString("# Agent 面试速成 · Agent Interview Crash Course\n\n")
	b.WriteString("面向后端工程师的 AI Agent 面试题，中英双语，按模块从 Agent 循环一路讲到生产环境；之后每天追加一组新题。\n")
	b.WriteString("每题包含题目、参考答案（默认折叠）、追问和词汇表，中文逐段放在英文下面。\n\n")
	fmt.Fprintf(&b, "共 %d 个模块、%d 道题。\n\n", len(mods), total)
	b.WriteString("- 数据：`interview/agent/data/<模块 id>.json`（与 Agent 面试速成页面的 `modules` 集合结构一致）\n")
	b.WriteString("- 重新生成 Markdown：在仓库根目录执行 `go run ./interview/render -dir interview/agent -kind modules`\n")

	for _, m := range mods {
		fmt.Fprintf(&b, "\n## [%s](modules/%s.md)\n\n", moduleTitle(m), m.ID)
		fmt.Fprintf(&b, "%s\n\n> %s\n\n", m.Summary, m.SummaryZh)
		writeTable(&b, m, "modules/"+m.ID+".md")
	}
	return b.String()
}

// writeTable 输出一组题的列表；link 为空时链接到当前页面内的锚点
func writeTable(b *strings.Builder, set Set, link string) {
	b.WriteString("| # | Topic | Difficulty | Title | 标题 |\n| --- | --- | --- | --- | --- |\n")
	for i, q := range set.Questions {
		fmt.Fprintf(b, "| %d | %s | %s | [%s](%s#q%d) | %s |\n",
			i+1, cell(q.Topic), cell(q.Difficulty), cell(q.Title), link, i+1, cell(q.TitleZh))
	}
}

// block 是题目 markdown 里的一个段落、列表或代码块
type block struct {
	kind  string // "p"、"ul" 或 "pre"
	text  string // 段落文本，或包含围栏的整个代码块
	items []string
}

var itemRE = regexp.MustCompile(`^\s*[-*] `)

func parseBlocks(src string) []block {
	var out []block
	lines := strings.Split(src, "\n")
	for i := 0; i < len(lines); {
		l := lines[i]
		switch {
		case strings.HasPrefix(l, "```"):
			start := i
			for i++; i < len(lines) && !strings.HasPrefix(lines[i], "```"); i++ {
			}
			i++ // 跳过结束围栏
			out = append(out, block{kind: "pre", text: strings.Join(lines[start:min(i, len(lines))], "\n")})
		case itemRE.MatchString(l):
			var items []string
			for ; i < len(lines) && itemRE.MatchString(lines[i]); i++ {
				items = append(items, itemRE.ReplaceAllString(lines[i], ""))
			}
			out = append(out, block{kind: "ul", items: items})
		case strings.TrimSpace(l) == "":
			i++
		default:
			var p []string
			for ; i < len(lines) && strings.TrimSpace(lines[i]) != "" && !strings.HasPrefix(lines[i], "```") && !itemRE.MatchString(lines[i]); i++ {
				p = append(p, lines[i])
			}
			out = append(out, block{kind: "p", text: strings.Join(p, " ")})
		}
	}
	return out
}

// bilingual 把中文逐段放在对应的英文段落或列表项下面（引用块）；
// 两边结构对不上时打印警告，并退回为先英文、后中文
func bilingual(name, en, zh string) string {
	a, z := parseBlocks(en), parseBlocks(zh)
	aligned := len(a) == len(z)
	for i := 0; aligned && i < len(a); i++ {
		aligned = a[i].kind == z[i].kind && len(a[i].items) == len(z[i].items)
	}
	if !aligned {
		log.Printf("warning: %s: English and Chinese blocks do not line up", name)
		return en + "\n\n" + zh + "\n\n"
	}

	var b strings.Builder
	for i, x := range a {
		y := z[i]
		switch x.kind {
		case "pre":
			b.WriteString(x.text + "\n\n")
			if y.text != x.text {
				b.WriteString(y.text + "\n\n")
			}
		case "ul":
			for j, it := range x.items {
				fmt.Fprintf(&b, "- %s\n  > %s\n", it, y.items[j])
			}
			b.WriteString("\n")
		default:
			fmt.Fprintf(&b, "%s\n\n> %s\n\n", x.text, y.text)
		}
	}
	return b.String()
}

// cell 转义表格单元格中的竖线和换行
func cell(s string) string {
	s = strings.ReplaceAll(s, "|", `\|`)
	return strings.ReplaceAll(s, "\n", " ")
}
