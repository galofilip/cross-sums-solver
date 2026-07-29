package main

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/disintegration/gift"
)

func main() {
	filePaths := bigImg2smallImgs()
	if filePaths == "" {
		return
	}
	nums := findNums(filePaths)
	solvedNums := solve(nums)
	for _, row := range solvedNums {
		for _, val := range row {
			fmt.Printf("%-3d", val)
		}
		fmt.Println()
	}
}

func bigImg2smallImgs() string {
	fmt.Println("hello world")
	//dir := "/data/data/com.termux/files/home/storage/dcim/Screenshots/Screenshot_20260717_094107_Offline Games.jpg"
	dir := "/home/philips/Projects/goWork/big_3.jpg"
	file, err := os.Open(dir)
	if err != nil {
		fmt.Printf("Failed to open image: %v", err)
		return ""
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		fmt.Printf("Failed to decode image: %v", err)
		return ""
	}
	file.Close()
	r := image.Rect(46, 581, 1046, 1581)
	g := gift.New(
		gift.Crop(r),
		gift.Threshold(50.0),
	)
	dstImg := image.NewRGBA(g.Bounds(img.Bounds()))
	g.Draw(dstImg, img)

	filePaths := ""

	for i := range 9 {
		for j := range 9 {
			if i == 0 && j == 0 {
				continue
			}
			r := image.Rect(j*111, i*111, (j+1)*111, (i+1)*111)
			g := gift.New(gift.Crop(r))
			imageNum := image.NewRGBA(g.Bounds(dstImg.Bounds()))
			g.Draw(imageNum, dstImg)

			createImages(i, j, imageNum, &filePaths)
		}
	}
	return filePaths
}

func findNums(filePaths string) [9][9]int {
	cmd := exec.Command("./venv/bin/python", "recognition.py", filePaths)

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		fmt.Printf("Python script failed: %v\nPython Error Log:\n%s", err, stderr.String())
		return [9][9]int{}
	}

	resultStr := strings.TrimSpace(out.String())

	resultStr = resultStr[:len(resultStr)]
	strNums := strings.Split(resultStr, ",")

	var nums [9][9]int
	for i, _ := range nums {
		for j, _ := range nums {
			if i == 0 && j == 0 {
				nums[i][j] = -1
				continue
			}
			first := strNums[0]
			strNums = strNums[1:]
			nums[i][j], err = strconv.Atoi(first)
			if err != nil {
				fmt.Printf("error in casting python answer to int: %v", err)
				return [9][9]int{}
			}
		}
	}
	return nums
}

func createImages(i, j int, imageNum *image.RGBA, filePaths *string) {
	n := fmt.Sprintf("temp/output_%d_%d.jpeg", i, j)
	*filePaths += n + ","
	outFile, err := os.Create(n)
	if err != nil {
		fmt.Printf("Failed to create a file for the black-white jpeg: %v", err)
		return
	}
	defer outFile.Close()

	options := &jpeg.Options{Quality: 95}

	err = jpeg.Encode(outFile, imageNum, options)
	if err != nil {
		fmt.Printf("Failed to encode the black-white image into jpeg file: %v", err)
	}
}

func solve(nums [9][9]int) [9][9]int {
	for range 10 {
		for i, num := range nums[1:] {
			times := [9]int{0, 0, 0, 0, 0, 0, 0, 0, 0}
			recSolve(num, 1, &times)
			for j, x := range times[1:] {
				if x == 0 {
					nums[i+1][j+1] = -1
				} else if x == times[0] && nums[i+1][j+1] != 0 {
					nums[i+1][0] -= nums[i+1][j+1]
					nums[0][j+1] -= nums[i+1][j+1]
					nums[i+1][j+1] = 0
				}
			}
		}
		for i := range 8 {
			num := [9]int{}
			for j := range 9 {
				num[j] = nums[j][i+1]
			}
			times := [9]int{0, 0, 0, 0, 0, 0, 0, 0, 0}
			recSolve(num, 1, &times)
			for j, x := range times[1:] {
				if x == 0 {
					nums[j+1][i+1] = -1
				} else if x == times[0] && nums[j+1][i+1] != 0 {
					nums[0][i+1] -= nums[j+1][i+1]
					nums[j+1][0] -= nums[j+1][i+1]
					nums[j+1][i+1] = 0
				}
			}
		}
	}
	return nums
}

// need to know how to work with -1 and 0(ignore them)
func recSolve(line [9]int, i int, times *[9]int) {
	if line[0] < 0 || (i > 8 && line[0] > 0) {
		return
	}

	if line[0] == 0 {
		for j, n := range line {
			if n == 0 {
				times[j]++
			}
		}
		return
	}

	n := line[i]

	if line[i] != 0 {
		line[i] = -1
		recSolve(line, i+1, times)
		line[i] = n
	}

	if line[i] != -1 {
		line[0] = line[0] - line[i]
		line[i] = 0
		recSolve(line, i+1, times)
	}
}

//take image
//enter numbers to arrays

//solve arrays
//choose in app the right numbers and delete all wrong numbers
