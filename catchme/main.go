package main

import (
	"fmt"
	"math"
)

func main() {
	var encoded string
	fmt.Print("Enter string: ")
	fmt.Scanln(&encoded)

	result := findMinSumSequence(encoded)
	fmt.Println(fmt.Sprintf("Result: %v", result))
}

func findMinSumSequence(encoded string) string {
	// กำหนดค่าเริ่มต้นของผลรวมต่ำสุดเป็นค่ามากที่สุด
	minSum := math.MaxInt32
	var minSequence string

	// เรียกฟังก์ชัน generateSequences เพื่อสร้างลำดับที่เป็นไปได้ทั้งหมด
	// และคำนวณผลรวมที่ต่ำที่สุด
	generateSequences(encoded, "", 0, &minSum, &minSequence)

	return minSequence
}

func generateSequences(encoded string, currentSeq string, index int, minSum *int, minSequence *string) {

	if index == len(encoded)+1 { //done
		sum := 0
		for _, digit := range currentSeq {
			sum += int(digit - '0') //rune to int
		}

		// ถ้าผลรวมของลำดับนี้น้อยกว่าผลรวมต่ำสุดที่มีอยู่
		if sum < *minSum {
			*minSum = sum
			*minSequence = currentSeq // อัพเดตผลรวมต่ำสุดและลำดับที่มีผลรวมต่ำสุด
		}
		return
	}

	// ถ้าเป็นตำแหน่งแรก ให้ลองตัวเลข 0-2
	if index == 0 {
		for i := 0; i <= 2; i++ {
			newSeq := currentSeq + fmt.Sprintf("%d", i)
			generateSequences(encoded, newSeq, index+1, minSum, minSequence)
		}
		return
	}

	// หาเลขตัวก่อนหน้านี้ในลำดับ
	prevDigit := int(currentSeq[index-1] - '0')

	// ตัวหนังสือก่อนหน้า
	relation := encoded[index-1]
	switch relation {
	case 'L': // ซ้ายมากกว่าขวา
		// ลองตัวเลขที่น้อยกว่าตัวเลขที่ตำแหน่งก่อนหน้า
		for i := 0; i < prevDigit; i++ {
			newSeq := currentSeq + fmt.Sprintf("%d", i)
			generateSequences(encoded, newSeq, index+1, minSum, minSequence)
		}
	case 'R': // ซ้ายกว่าขวามากกว่า (Left < Right)
		// ลองตัวเลขที่มากกว่าตัวเลขที่ตำแหน่งก่อนหน้า
		for i := prevDigit + 1; i <= 9; i++ {
			newSeq := currentSeq + fmt.Sprintf("%d", i)
			generateSequences(encoded, newSeq, index+1, minSum, minSequence)
		}
	case '=': // ซ้ายเท่ากับขวา (Left = Right)
		// ลองตัวเลขที่เหมือนกับตัวเลขที่ตำแหน่งก่อนหน้า
		newSeq := currentSeq + fmt.Sprintf("%d", prevDigit)
		generateSequences(encoded, newSeq, index+1, minSum, minSequence)
	}
}
