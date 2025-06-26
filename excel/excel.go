package main

import (
	"errors"
	"fmt"
	"github.com/xuri/excelize/v2"
	"log"
)

func main() {
	f, err := excelize.OpenFile("./excel/销售计划导入模版.xlsx")
	if err != nil {
		log.Fatal("打开文件失败:", err)
	}

	sheetName := f.GetSheetName(0)
	rowIter, err := f.Rows(sheetName)
	if err != nil {
		log.Fatal("读取 sheet 失败:", err)
	}

	var (
		headerRowIndex int = -1
		currentRow     int = 0
		dataRows       [][]string
		headerFound    bool = false
	)

	for rowIter.Next() {
		currentRow++
		row, err := rowIter.Columns()
		if err != nil {
			log.Fatalf("读取第 %d 行失败: %v", currentRow, err)
		}

		// 跳过说明行（通常合并 A1:C1，只 A1 有内容）
		if isMergedInstructionRow(row) {
			continue
		}

		// 找表头，只记录一次
		if !headerFound && isHeaderRow(row) {
			headerRowIndex = currentRow
			headerFound = true
			continue
		}

		// 收集数据行
		if headerFound && currentRow > headerRowIndex {
			if len(row) >= 3 && row[0] != "" && row[1] != "" && row[2] != "" {
				dataRows = append(dataRows, row)
			}
		}
	}

	if headerRowIndex == -1 {
		log.Fatal(errors.New("未找到表头"))
	}

	fmt.Println("✅ 表头在第", headerRowIndex, "行")
	for _, row := range dataRows {
		fmt.Printf("商店名称: %s, 日期: %s, 金额: %s\n", row[0], row[1], row[2])
	}
}

// 判断是否是说明行（合并 A1:C1 情况，A 列有内容，其余为空）
func isMergedInstructionRow(row []string) bool {
	fmt.Println(len(row))
	return len(row) >= 3 && row[0] != "" && row[1] == "" && row[2] == ""
}

// 判断是否是表头行
func isHeaderRow(row []string) bool {
	return len(row) >= 3 &&
		row[0] == "商店名称" &&
		row[1] == "日期" &&
		row[2] == "金额"
}
