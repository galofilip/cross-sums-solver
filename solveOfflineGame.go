package main

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/disintegration/gift"
	"image"
	"image/jpeg"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

const startX = 46
const startY = 581

func main() {
	filePaths, err := bigImg2smallImgs()
	if err != nil {
		fmt.Println(err)
		return
	}
	nums, err := findNums(filePaths)
	if err != nil {
		fmt.Println(err)
		return
	}
	solvedNums := solve(nums)
        fmt.Println("in 10nseconds your screen will be pressed, ro press on nums please switch to the game to make the program work")
        time.Sleep(10 * time.Second)
	
	err = press(solvedNums)
	if err != nil {
		fmt.Println(err)
		return
	}
/*	for _, row := range solvedNums {
		for _, val := range row {
			fmt.Printf("%-3d", val)
		}
		fmt.Println()
	}*/
}

func bigImg2smallImgs() (string, error) {
	fmt.Println("in 10 seconds a screen shot will be taken, make sure your phone will be on the cross-sum game")
	time.Sleep(10 * time.Second)

	out, err := os.Create("./sc.jepg")
	if err != nil {
		return "", err
	}
	defer out.Close()

	cmd := exec.Command("adb", "exec-out", "screencap", "-j")
	cmd.Stdout = out
	err = cmd.Run()
	if err != nil {
		return "", err
	}


	if err != nil {
		return "", err
	}
	fmt.Println("the screenshot has been taken successfully\nsolving...")

	dir := "./sc.jepg"
//	dir := "/home/philips/Projects/goWork/big_3.jpg"
	file, err := os.Open(dir)
	if err != nil {
		return "", errors.New(fmt.Sprintf("Failed to open image: %v", err))
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		return "", errors.New(fmt.Sprintf("Failed to decode image: %v", err))
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

			err := createImages(i, j, imageNum, &filePaths)
			if err != nil {
				return filePaths, err
			}
		}
	}
	return filePaths, nil
}

func findNums(filePaths string) ([9][9]int, error) {
	cmd := exec.Command("./venv/bin/python", "recognition.py", filePaths)

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return [9][9]int{}, errors.New(fmt.Sprintf("Python script failed: %v\nPython Error Log:\n%s", err, stderr.String()))
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
				return [9][9]int{}, errors.New(fmt.Sprintf("error in casting python answer to int: %v", err))
			}
		}
	}
	return nums, nil
}

func createImages(i, j int, imageNum *image.RGBA, filePaths *string) (error){
	n := fmt.Sprintf("./temp/output_%d_%d.jpeg", i, j)
	*filePaths += n + ","
	outFile, err := os.Create(n)
	if err != nil {
		return errors.New(fmt.Sprintf("Failed to create a file for the black-white jpeg: %v", err))
	}
	defer outFile.Close()

	options := &jpeg.Options{Quality: 95}

	err = jpeg.Encode(outFile, imageNum, options)
	if err != nil {
		return errors.New(fmt.Sprintf("Failed to encode the black-white image into jpeg file: %v", err))
	}
	return nil
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
	fmt.Println("solved!")
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

func press(solvedNums [9][9]int)(error){
	for i, row := range solvedNums[1:] {
		for j, val := range row[1:] {
			if val == -1 {
				fmt.Println("i: ", i, ", j: ", j , " ,x: ", (startX+(j+1)*111+56), ", y: ", (startY+(i+1)*111+56))
				cmd := exec.Command("adb", "shell", "input", "tap", strconv.Itoa(startX+(i+1)*111+56), strconv.Itoa(startY+(j+1)*111+56))
				err := cmd.Run()
				if err != nil {
					return err
				}
				fmt.Scanln()
			}
		}
	}
	cmd := exec.Command("adb", "shell", "input", "tap", strconv.Itoa(612), strconv.Itoa(1850))
	err := cmd.Run()
	if err != nil {
		return err
	}
	for i, row := range solvedNums[1:] {
		for j, val := range row[1:] {
			if val == -1 {
				fmt.Println("i: ", i, ", j: ", j , " ,x: ", (startX+(j+1)*111+56), ", y: ", (startY+(i+1)*111+56))
				cmd := exec.Command("adb", "shell", "input", "tap", strconv.Itoa(startX+(i+1)*111+56), strconv.Itoa(startY+(j+1)*111+56))
				err := cmd.Run()
				if err != nil {
					return err
				}
				fmt.Scanln()
			}
		}
	}
	return nil
}

//take image
//enter numbers to arrays

//solve arrays
//choose in app the right numbers and delete all wrong numbers
