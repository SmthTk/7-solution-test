package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func maxPathSum(triangle [][]int) int {
	if len(triangle) == 0 {
		return 0
	}

	height := len(triangle)

	//เริ่มจาก แถวรองสุดท้ายเพื่อ หาค่าที่มากที่สุดของ แถวข้างล่าง ไล่ขึ้นไปจนบนสุด
	for i := height - 2; i >= 0; i-- {
		for j := 0; j < len(triangle[i]); j++ {
			triangle[i][j] += max(triangle[i+1][j], triangle[i+1][j+1]) //ค่าของตัวเอง + ค่าของสองแถวล่างของ สามเหลี่ยม แล้วไล่ขึ้นไป สุดท้ายที่ตำแหน่ง บนสุดจะได้ค่ามากที่สุด ไม่มีการวนไปทางอื่น
		}
	}

	return triangle[0][0]
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func main() {
	url := "https://raw.githubusercontent.com/7-solutions/backend-challenge/main/files/hard.json"
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error fetching JSON:", err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading JSON:", err)
		return
	}

	var triangle [][]int
	if err := json.Unmarshal(body, &triangle); err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}

	result := maxPathSum(triangle)
	fmt.Println("Maximum path sum:", result)
}
