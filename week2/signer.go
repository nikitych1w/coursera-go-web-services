package week2

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

var start time.Time

// crc32(data)+"~"+crc32(md5(data))
func SingleHash(in, out chan interface{}){

	fmt.Println("SingleHash Start", time.Since(start))
	var wg = &sync.WaitGroup{}
	var m = &sync.Mutex{}
	//var crc32, md5, crc32md5 string

	for el := range in {
		//var workers = &sync.WaitGroup{}
		data := strconv.Itoa(el.(int))
		fmt.Println(data, " ======= SingleHash Iter", time.Since(start))

		wg.Add(1)
		go func(data string){
			defer wg.Done()

			ch1 := make(chan string)
			ch2 := make(chan string)

			//workers.Add(1)
			go func(){
				//defer workers.Done()
				crc32 := DataSignerCrc32(data)
				ch1 <- crc32
			}()

			//workers.Add(1)
			go func(){
				//defer workers.Done()
				m.Lock()
				md5 := DataSignerMd5(data)
				m.Unlock()
				crc32md5 := DataSignerCrc32(md5)
				ch2 <- crc32md5
			}()

			//workers.Wait()
			res := <-ch1 + "~" + <-ch2
			out <- res
		}(data)

	}

	wg.Wait()

	fmt.Println("SingleHash End", time.Since(start))
}

// crc32(th+data))
func MultiHash(in, out chan interface{}){
	fmt.Println("MultiHash Start",time.Since(start))
	var hashes [6]string
	var m = &sync.Mutex{}
	var wg = &sync.WaitGroup{}



	for el := range in {
		var workers = &sync.WaitGroup{}
		data := el.(string)
		fmt.Println(data, " ======= MultiHash Iter", time.Since(start))

		wg.Add(1)
		go func(data string) {
			defer wg.Done()
			for th := 0; th < 6; th++ {
				workers.Add(1)
				go func(th int) {
					defer workers.Done()
					res := DataSignerCrc32(strconv.Itoa(th) + data)
					m.Lock()
					hashes[th] = res
					m.Unlock()
				}(th)
			}

			workers.Wait()
			m.Lock()
			res := strings.Join(hashes[:], "")
			m.Unlock()
			out <- res
		}(data)

	}

	wg.Wait()
	fmt.Println("MultiHash End", time.Since(start))
}

func CombineResults(in, out chan interface{}){
	fmt.Println("CombineResults Start", time.Since(start))
	var result []string

	for el := range in {
		data := el.(string)
		fmt.Println(data, " ======= CombineResults Iter", time.Since(start))
		result = append(result, data)
	}

	sort.Strings(result)
	res := strings.Join(result, "_")
	out <- res
	fmt.Println("##################################################### res", res)
	fmt.Println("CombineResults End", time.Since(start))
}

func runJob(wg *sync.WaitGroup, j job, in, out chan interface{}) {
	defer wg.Done()
	defer close(out)
	j(in, out)
}

func ExecutePipeline(jobs... job){
	var wg = &sync.WaitGroup{}
	var in = make(chan interface{}, 100)
	var out = make(chan interface{}, 100)

	for _, j := range jobs {
		wg.Add(1)
		go runJob(wg, j, in, out)
		in, out = out, make(chan interface{}, 100)
	}

	wg.Wait()
}