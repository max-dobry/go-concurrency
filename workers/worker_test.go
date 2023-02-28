package workers

import (
	"bufio"
	"log"
	"os"
	"sync"
	"testing"
)

func BenchmarkWorkers(b *testing.B) {
	var domains []string

	file, err := os.Open("domains_test.txt")
	if err != nil {
		log.Println("can't open domains.txt file")
		b.Fatalf("error open file with domain names")
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		domains = append(domains, scanner.Text())
	}

	d := &disp{
		worker:          worker{},
		NumberOfWorkers: 3,
		wg:              &sync.WaitGroup{},
		urls:            domains,
	}

	b.Run("WorkerType1", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			d.wg.Add(d.NumberOfWorkers)
			d.jobChannel = make(chan Job, len(d.urls))
			strCh := make(chan string, len(d.urls))

			for i := 1; i <= d.NumberOfWorkers; i++ {
				go d.worker.Start(i, d.wg, d.jobChannel, strCh)
			}

			for i, v := range d.urls {
				d.jobChannel <- Job{
					ID:   i,
					Site: "http://" + v,
				}
			}
			close(d.jobChannel)
			d.wg.Wait()
			// close(d.ResultChannel)
		}
	})

	b.Run("WorkerType2", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			d.wg.Add(d.NumberOfWorkers)
			d.jobChannel = make(chan Job, len(d.urls))
			strCh := make(chan string, len(d.urls))

			for i := 1; i <= d.NumberOfWorkers; i++ {
				go d.worker.Start2(i, d.wg, d.jobChannel, strCh)
			}

			for i, v := range d.urls {
				d.jobChannel <- Job{
					ID:   i,
					Site: "https://" + v,
				}
			}
			close(d.jobChannel)
			d.wg.Wait()
			// close(d.ResultChannel)
		}
	})

	b.Run("WorkerType3", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			d.wg.Add(d.NumberOfWorkers)
			d.jobChannel = make(chan Job, len(d.urls))
			d.ResultChannel = make(chan JobResult, len(d.urls))

			for i := 1; i <= d.NumberOfWorkers; i++ {
				go d.worker.Start3(i, d.wg, d.jobChannel, d.ResultChannel)
			}

			for i, v := range d.urls {
				d.jobChannel <- Job{
					ID:   i,
					Site: "https://" + v,
				}
			}
			close(d.jobChannel)
			d.wg.Wait()
			// close(d.ResultChannel)
		}
	})
}