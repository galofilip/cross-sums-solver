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

const startX = 44
const startY = 582

/*
1. double tap to switch to game 	:)
2. take screen shot(start of loop) 	:)
3. scroll up				:)
4. double tap to switch to termux	:)
5. solve				:)
6. double tap to switch to game		:)
7. make press				:)
8. scroll up				:)
9. tap on 3 lines			:)
10. scroll to delete game 
11. press on circel
12. press offlineGames in home
13. press game
14. start level(end loop)
*/

func main() {
	fmt.Println("in 5 seconds a screen shot will be taken, make sure your phone will be on the cross-sum game")
	time.Sleep(1 * time.Second)
	touch(252, 2290)
	time.Sleep(1 * time.Second)
	touch(252, 2290)
	time.Sleep(2 * time.Second)

	for range 2 {
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
		fmt.Println("in 4 seconds your screen will be pressed, to press on nums please switch to the game to make the program work")
		touch(252, 2290)
		time.Sleep(1 * time.Second)
		touch(252, 2290)
		time.Sleep(2 * time.Second)

		err = press(solvedNums)
		if err != nil {
			fmt.Println(err)
			return
		}
		
		scroll(550, 2334, 550, 2115)
		time.Sleep(1 * time.Second)
		touch(252, 2290)
		time.Sleep(1 * time.Second)
		scroll(1020, 1500, 1020, 500)

		time.Sleep(1 * time.Second)
		touch(540, 2290)
		time.Sleep(1 * time.Second)
		touch(950, 1100)
		time.Sleep(4 * time.Second)

		touch(815, 1950)
		time.Sleep(1 * time.Second)
		touch(575, 1900)
	}
}

func bigImg2smallImgs() (string, error) {
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
	fmt.Println("the screenshot has been taken successfully\nsolving...")
	
	scroll(550, 2334, 550, 2115)
	time.Sleep(1 * time.Second)
	touch(252, 2290)
	time.Sleep(1 * time.Second)
	touch(252, 2290)
	time.Sleep(2 * time.Second)

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
	r := image.Rect(startX, startY, startX+1008, startY+1008)
	g := gift.New(
		gift.Crop(r),
		gift.Invert(),
		gift.Threshold(65.0),
	)
	dstImg := image.NewRGBA(g.Bounds(img.Bounds()))
	g.Draw(dstImg, img)

	filePaths := ""

	for i := range 9 {
		for j := range 9 {
			if i == 0 && j == 0 {
				continue
			}
			r := image.Rect(j*112, i*112, (j+1)*112, (i+1)*112)
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
	for con(nums) {
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
				err := touch(startX+(j+1)*112+56, startY+(i+1)*112+56)
				if err != nil {
					return err
				}
			}
		}
	}
	err := touch(612, 1850)
	if err != nil {
		return err
	}
	for i, row := range solvedNums[1:] {
		for j, val := range row[1:] {
			if val == 0 {
				err := touch(startX+(j+1)*112+56, startY+(i+1)*112+56)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func con(solvedNums [9][9]int)(bool){
	for _, row := range solvedNums {
		for _, val := range row {
			if val != -1 && val != 0{
				return true
			}
		}
	}
	return false
}

func touch(x, y int)(error){
	cmd := exec.Command("adb", "shell", "input", "tap", strconv.Itoa(x), strconv.Itoa(y))
	err := cmd.Run()
	return err
}

func scroll(sX, sY, eX, eY int)(error){
	cmd := exec.Command("adb", "shell", "input", "swipe", strconv.Itoa(sX), strconv.Itoa(sY), strconv.Itoa(eX), strconv.Itoa(eY), "200")
	err := cmd.Run()
	return err
}
