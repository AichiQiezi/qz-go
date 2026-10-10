// render 把 interview/data/*.json 渲染成 interview/days/*.md，并重建 interview/README.md 索引。
//
// 用法（在仓库根目录执行）：
//
//	go run ./interview/render
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
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

// Day 与 Interview Drill 页面里 "days" 集合的文档结构一致
type Day struct {
	Date        string     `json:"date"`
	GeneratedAt string     `json:"generatedAt"`
	Questions   []Question `json:"questions"`
}

func main() {
	dir := flag.String("dir", "interview", "interview 目录")
	flag.Parse()

	files, err := filepath.Glob(filepath.Join(*dir, "data", "*.json"))
	if err != nil {
		log.Fatal(err)
	}
	if len(files) == 0 {
		log.Fatalf("no data files under %s", filepath.Join(*dir, "data"))
	}

	days := make([]Day, 0, len(files))
	for _, f := range files {
		day, err := loadDay(f)
		if err != nil {
			log.Fatalf("%s: %v", f, err)
		}
		days = append(days, day)
	}
	// 最新的一天排在最前面
	sort.Slice(days, func(i, j int) bool { return days[i].Date > days[j].Date })

	if err := os.MkdirAll(filepath.Join(*dir, "days"), 0o755); err != nil {
		log.Fatal(err)
	}
	for _, day := range days {
		path := filepath.Join(*dir, "days", day.Date+".md")
		if err := os.WriteFile(path, []byte(renderDay(day)), 0o644); err != nil {
			log.Fatal(err)
		}
	}
	if err := os.WriteFile(filepath.Join(*dir, "README.md"), []byte(renderIndex(days)), 0o644); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("rendered %d day(s) into %s\n", len(days), *dir)
}

func loadDay(path string) (Day, error) {
	var day Day
	b, err := os.ReadFile(path)
	if err != nil {
		return day, err
	}
	if err := json.Unmarshal(b, &day); err != nil {
		return day, err
	}
	if want := strings.TrimSuffix(filepath.Base(path), ".json"); day.Date != want {
		return day, fmt.Errorf("date %q does not match file name %q", day.Date, want)
	}
	if len(day.Questions) == 0 {
		return day, fmt.Errorf("no questions")
	}
	for i, q := range day.Questions {
		if q.Title == "" || q.Question == "" || q.Answer == "" {
			return day, fmt.Errorf("question %d is missing title, question or answer", i+1)
		}
	}
	return day, nil
}

func renderDay(day Day) string {
	var b strings.Builder
	fmt.Fprintf(&b, "# %s · 每日面试题 / Daily Interview Questions\n\n", day.Date)
	fmt.Fprintf(&b, "> %d questions · generated %s · [← 全部 / All days](../README.md)\n\n", len(day.Questions), day.GeneratedAt)
	writeTable(&b, day, "")

	for i, q := range day.Questions {
		b.WriteString("\n---\n\n")
		fmt.Fprintf(&b, "<a id=\"q%d\"></a>\n\n", i+1)
		fmt.Fprintf(&b, "## %d. %s\n\n", i+1, q.Title)
		fmt.Fprintf(&b, "**%s** · `%s` · `%s`\n\n", q.TitleZh, q.Topic, q.Difficulty)

		fmt.Fprintf(&b, "### Question\n\n%s\n\n", q.Question)
		fmt.Fprintf(&b, "### 题目\n\n%s\n\n", q.QuestionZh)

		// 答案默认折叠，先自己想再展开
		b.WriteString("<details>\n<summary><b>Answer · 参考答案</b></summary>\n\n")
		fmt.Fprintf(&b, "#### Answer\n\n%s\n\n", q.Answer)
		fmt.Fprintf(&b, "#### 参考答案\n\n%s\n\n", q.AnswerZh)
		b.WriteString("</details>\n\n")

		if len(q.Followups) > 0 {
			b.WriteString("### Follow-ups · 追问\n\n")
			for j, f := range q.Followups {
				fmt.Fprintf(&b, "- %s\n", f)
				if j < len(q.FollowupsZh) {
					fmt.Fprintf(&b, "  - %s\n", q.FollowupsZh[j])
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

func renderIndex(days []Day) string {
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

	for _, d := range days {
		fmt.Fprintf(&b, "\n## [%s](days/%s.md)\n\n", d.Date, d.Date)
		writeTable(&b, d, "days/"+d.Date+".md")
	}
	return b.String()
}

// writeTable 输出当天的题目列表；link 为空时链接到当前页面内的锚点
func writeTable(b *strings.Builder, day Day, link string) {
	b.WriteString("| # | Topic | Difficulty | Title | 标题 |\n| --- | --- | --- | --- | --- |\n")
	for i, q := range day.Questions {
		fmt.Fprintf(b, "| %d | %s | %s | [%s](%s#q%d) | %s |\n",
			i+1, cell(q.Topic), cell(q.Difficulty), cell(q.Title), link, i+1, cell(q.TitleZh))
	}
}

// cell 转义表格单元格中的竖线和换行
func cell(s string) string {
	s = strings.ReplaceAll(s, "|", `\|`)
	return strings.ReplaceAll(s, "\n", " ")
}
