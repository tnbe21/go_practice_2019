package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"quiz/quiz"
	"time"
)

var cnt = 1
var finished chan bool
var correctCnt = 0

const MAX_QUIZ_COUNT = 3
const TIMEOUT_PER_QUIZ = 5

func main() {
	fmt.Println("クイズを開始します。")

	for cnt < MAX_QUIZ_COUNT {
		doQuiz()
	}

	fmt.Println("***")
	fmt.Fprintf(os.Stdout, "%d問中%d問正解\n", cnt, correctCnt)
	fmt.Println("***")
}

func doQuiz() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*TIMEOUT_PER_QUIZ)
	defer cancel()
	_quiz := quiz.GetQuiz()
	viewQuiz(ctx, _quiz)
	go scanAnswer(ctx, _quiz)
	time.Sleep(time.Second*(TIMEOUT_PER_QUIZ + 1))
}

func viewQuiz(ctx context.Context, _quiz quiz.Quiz) {
	var tickCnt = 0
	for {
		select {
		case <- ctx.Done():
			fmt.Println("不正解")
			cnt += 1
			return
		default:
			if tickCnt > 0 {
				fmt.Fprintf(os.Stdout, "\033[1A\r")
			}
			fmt.Fprintf(os.Stdout, "Q.%d %s %d\n> ", cnt, _quiz.GetBody(), TIMEOUT_PER_QUIZ-tickCnt)
			time.Sleep(time.Second)
			tickCnt += 1
			if tickCnt >= TIMEOUT_PER_QUIZ {
				return
			}
		}
	}
}

func scanAnswer(ctx context.Context, _quiz quiz.Quiz) {
	for {
		select {
		case <- ctx.Done():
			return
		default:
			scanner := bufio.NewScanner(os.Stdin)
			scanner.Scan()
			answer := scanner.Text()
			if _quiz.IsCorrect(answer) {
				correctCnt += 1
				fmt.Println("正解")
			} else {
				fmt.Println("不正解")
			}
			cnt += 1
		}
	}
}
