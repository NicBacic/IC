package main

import (
"fmt"
"math/rand"
"sort"
"sync"
"time"
"log"
)

type Jobs struct{
id int
org int
r int
d int
w int
offset float32
}

type Processor struct{
s float32
Jobs []Jobs
id int
}

type Interval struct{
start int
end int
Jobs []Jobs
intensidade float32
}

type Schedule struct{
id int
org int
releaseDate float32
deadLine float32
speed float32
}

type Organization struct{
id int
Jobs []Jobs	
schedule []Schedule
dMax int
dLinha int
}

type ByReleaseDate []Schedule
type ByDeadline []Jobs
type ByDistance []Jobs
type ByDMax []Organization
type ByDLinha []Organization

func (a ByReleaseDate) Len() int         { return len(a) }
func (a ByReleaseDate) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByReleaseDate) Less(i, j int) bool { return a[i].releaseDate < a[j].releaseDate }

func (a ByDeadline) Len() int          { return len(a) }
func (a ByDeadline) Swap(i, j int) { a[i], a[j] = a[j], a[i]}
func (a ByDeadline) Less(i, j int) bool { return a[i].d < a[j].d }

func (a ByDistance) Len() int { return len(a) }
func (a ByDistance) Swap(i,j int) {a[i], a[j] = a[j], a[i]}
func (a ByDistance) Less(i, j int) bool {return (a[i].d - a[i].r) > (a[j].d - a[j].r)}

func (a ByDMax) Len() int          { return len(a) }
func (a ByDMax) Swap(i, j int) { a[i], a[j] = a[j], a[i]}
func (a ByDMax) Less(i, j int) bool { return a[i].dMax < a[j].dMax }

func (a ByDLinha) Len() int          { return len(a) }
func (a ByDLinha) Swap(i, j int) { a[i], a[j] = a[j], a[i]}
func (a ByDLinha) Less(i, j int) bool { return a[i].dLinha < a[j].dLinha }


var orgs []Organization
var identifier int = 0
var wg sync.WaitGroup

func YDS(jobs []Jobs, m int, t bool){
	
	var intervaloDensidadeMax Interval
	var schedule []Schedule

	var auxSchedule Schedule

	var tamanhoIntervalo float32 = 0.0
	var offset float32 = 0.0
	var tamanhoTarefa float32 = 0.0
	var somaProcessamento float32 = 0.0

	i:=0
	
	for len(jobs) != 0 {

		interval := calculaIntervalos(jobs)

		// parte 02 do psedoCódigo
		intervaloDensidadeMax = intervaloDeDensidadeMaxima(interval) 


		//parte 03 ''
		sort.Sort(ByDeadline(intervaloDensidadeMax.Jobs))

		tamanhoIntervalo = float32(intervaloDensidadeMax.end - intervaloDensidadeMax.start)

		offset = float32(intervaloDensidadeMax.start)
		
		somaProcessamento = 0.0
		for j:=0; j < len(intervaloDensidadeMax.Jobs); j++ {
			somaProcessamento += float32(intervaloDensidadeMax.Jobs[j].w)
		}
		

		//parte 04		
		for j := 0; j < len(intervaloDensidadeMax.Jobs); j++ {

			tamanhoTarefa = (float32(intervaloDensidadeMax.Jobs[j].w) / somaProcessamento) * tamanhoIntervalo

			auxSchedule = Schedule{intervaloDensidadeMax.Jobs[j].id, intervaloDensidadeMax.Jobs[j].org,
						intervaloDensidadeMax.Jobs[j].offset + offset,
						intervaloDensidadeMax.Jobs[j].offset + offset + tamanhoTarefa, 					   								intervaloDensidadeMax.intensidade}
	
			schedule = append(schedule, auxSchedule)
			offset += tamanhoTarefa

			for k:=0; k < len(jobs); k++ {
				if jobs[k].id == intervaloDensidadeMax.Jobs[j].id {
					if len(jobs) > 1 {
						jobs = append(jobs[:k], jobs[k+1:]...)
					} else {
						jobs = jobs[k:k]
					}
					break
				}
			}
		}	

		//parte 05 Ajusta Tarefas
		if len(jobs) != 0 {
			jobs = ajustaTarefas(jobs, intervaloDensidadeMax.start, intervaloDensidadeMax.end)	
		}

		//fmt.Printf("Jobs restantes %v\n\n",jobs)
		i++
	}
	if t {
		orgs[m].schedule = schedule

		defer wg.Done()
	}
}

func trySchedule(jobs []Jobs, k int) bool {

	for len(jobs) > 0 {

		intervalo := calculaIntervalos(jobs)

		intervaloDensidadeMax := intervaloDeDensidadeMaxima(intervalo)

		auxJob := make([]Jobs, len(intervaloDensidadeMax.Jobs))

		copy(auxJob,intervaloDensidadeMax.Jobs)

		sort.Sort(ByDistance(intervaloDensidadeMax.Jobs))

		if orgs[0].dMax <= intervaloDensidadeMax.end {

			for i:=0; i < len(intervaloDensidadeMax.Jobs); i++ {
		
				tempJob := intervaloDensidadeMax.Jobs[i]

				dLinhaMax := orgs[0].dMax
				
				jobsScheduled := make([]Jobs, len(orgs[0].Jobs))

				copy(jobsScheduled,orgs[0].Jobs)

				nowSpeed := intervaloDeDensidadeMaxima(calculaIntervalos(auxJob)).intensidade 

				tempJob.r = dLinhaMax

				jobsScheduled = append(jobsScheduled, tempJob)

				

				scheduledSpeed := calculaIntervaloOrg(jobsScheduled, dLinhaMax)
			
				if scheduledSpeed < nowSpeed {
					orgs[0].Jobs = append(orgs[0].Jobs, tempJob)
			
					for i := len(jobs)-1; i >= 0; i-- {
						if jobs[i].id == tempJob.id {
							if len(jobs) > 1 {
								jobs = append(jobs[:i], jobs[i+1:]...)
							} else {
								jobs = jobs[i:i]
							}
							break
						}
					}

					for i := len(auxJob)-1; i >= 0; i-- {
						if auxJob[i].id == tempJob.id {
							if len(auxJob) > 1 {
								auxJob = append(auxJob[:i], auxJob[i+1:]...)
							} else {
								auxJob = auxJob[i:i]
							}
							break
						}
					}

					for i := len(orgs[k].Jobs)-1; i >= 0; i-- {
						if orgs[k].Jobs[i].id == tempJob.id {
							if len(orgs[k].Jobs) > 1 {
								orgs[k].Jobs = append(orgs[k].Jobs[:i], orgs[k].Jobs[i+1:]...)
							} else {
								orgs[k].Jobs = orgs[k].Jobs[i:i]
							}
							break
						}
					}
					dLinhaMax = tempJob.d
					orgs[0].dLinha = dLinhaMax

					sort.Sort(ByDLinha(orgs[:k]))
					if len(auxJob) == 0 {
						break
					}
				}  else {
					break
				}
			}
		}
		for j := 0; j < len(intervaloDensidadeMax.Jobs); j++ {

			for k:=0; k < len(jobs); k++ {
				if jobs[k].id == intervaloDensidadeMax.Jobs[j].id {
					if len(jobs) > 1 {
						jobs = append(jobs[:k], jobs[k+1:]...)
					} else {
						jobs = jobs[k:k]
					}
					break
				}
			}
		}

		if len(jobs) != 0 {
			jobs = pequenoAjuste(jobs, intervaloDensidadeMax.end)	
		}
	}
	return true
}

func pequenoAjuste(jobs []Jobs, d int) []Jobs{
	for i := 0; i < len(jobs); i++{
		jobs[i].r = d
	}
	return jobs
}


func calculaIntervaloOrg(jobs []Jobs, releaseDate int) float32 {
	var listaDeIntervalos []Interval
	start := 0
	end := 0
	releases := make(map[int]int)
	deadlines := make(map[int]int)
	workload := 0
	var densidade float32 = 0.0

	for i:=0; i < len(jobs); i++ {
		if jobs[i].r == releaseDate {
			releases[jobs[i].r] = jobs[i].r
		}
	}

	for i:=0; i < len(jobs); i++ {
		if jobs[i].r == releaseDate {
			deadlines[jobs[i].d] = jobs[i].d
		}
	}

	for key := range(releases) {
		for key2 := range(deadlines){
			if releases[key] < deadlines[key2] {

				start = releases[key]
				end = deadlines[key2]
				var listaDeTarefas []Jobs

				for k := 0; k < len(jobs); k++ {
					if jobs[k].r >= start && jobs[k].d <= end {
						listaDeTarefas = append(listaDeTarefas,jobs[k])
						workload += jobs[k].w
					}
				}

				densidade = float32(workload) / float32(end-start)
				interval := Interval{start,end,listaDeTarefas,densidade}
				listaDeIntervalos = append(listaDeIntervalos,interval)
				workload = 0
			}
		}
	}
	intervaloMax := intervaloDeDensidadeMaxima(listaDeIntervalos)
	return intervaloMax.intensidade
	
}

func ajustaTarefas(jobs []Jobs, start int, end int) []Jobs{
	distancia := end - start
	for i := 0; i < len(jobs); i++{
		//Tarefa depois do deadline
		if jobs[i].r >= end {
			jobs[i].d -= distancia
			jobs[i].r -= distancia
			jobs[i].offset += float32(distancia)
		//Tarefa antes do release date
		} else if jobs[i].d <= start {
			continue
		//Tarefa atravessa o intervalo
		} else if jobs[i].r < start && jobs[i].d > end{
			jobs[i].d -= distancia
		//Tarefa parcialmente contida no intervalo (deadline maior que release date do intervalo, mas release date menor)
		} else if jobs[i].d > start && jobs[i].r < start{
			jobs[i].d = start
		//Tarefa parcialmente contida no intervalo (release date maior que start, mas deadline fora do intervalo)
		} else if jobs[i].r < end && jobs[i].d > end {
			jobs[i].r = start
			jobs[i].d -= distancia
			jobs[i].offset += float32(distancia)
		}
		
	}
	return jobs
}

func calculaIntervalos(jobs []Jobs) []Interval{

	var listaDeIntervalos []Interval
	start := 0
	end := 0
	releases := make(map[int]int)
	deadlines := make(map[int]int)
	workload := 0
	var densidade float32 = 0.0

	for i:=0; i < len(jobs); i++ {
		releases[jobs[i].r] = jobs[i].r
	}

	for i:=0; i < len(jobs); i++ {
		deadlines[jobs[i].d] = jobs[i].d
	}

	for key := range(releases) {
		for key2 := range(deadlines){
			if releases[key] < deadlines[key2] {

				start = releases[key]
				end = deadlines[key2]
				var listaDeTarefas []Jobs

				for k := 0; k < len(jobs); k++ {
					if jobs[k].r >= start && jobs[k].d <= end {
						listaDeTarefas = append(listaDeTarefas,jobs[k])
						workload += jobs[k].w
					}
				}

				densidade = float32(workload) / float32(end-start)
				interval := Interval{start,end,listaDeTarefas,densidade}
				listaDeIntervalos = append(listaDeIntervalos,interval)
				workload = 0
			}
		}
	}
	
	return listaDeIntervalos
}

func calcSpeed(jobs []Jobs) float32{
	var speed float32 = 0.0
	
	somaProcessamento := len(jobs) //como w = 1, soma do Processamento é o número de tarefas do intervalo
	
	speed = float32(somaProcessamento) / float32(jobs[0].d - jobs[0].r)

	return speed
}

func intervaloDeDensidadeMaxima(interval []Interval) Interval{
	var max float32 = 0.0
	index := 0
	for i:=0; i < len(interval); i++{
		if max < interval[i].intensidade {
			max = interval[i].intensidade
			index = i
		}
	}
	return interval[index]
}


func createJobs(id int) []Jobs{
	n := rand.Intn(2) + 1
	jobs := make([]Jobs,n)
	for i:=0; i < n; i++{
		jobs[i].id = identifier
		jobs[i].org = id
		jobs[i].r = 0
		jobs[i].d = rand.Intn(10) + 1
		jobs[i].w = 1
		jobs[i].offset = 0.0
		identifier++
	}
	return jobs
}

func buildRandom(n int) bool{
	orgs = make([]Organization,n)

	/*for i:=0; i < n; i++ {
		orgs[i].id = i
		orgs[i].localEnergyConsumption = 0.0
		orgs[i].Jobs = createJobs(i)
	}*/
	orgs[0].id = 0
	orgs[0].Jobs = make([]Jobs,2)
	vector := []int{2,4}

	for i:=0; i < 2; i++ {
		orgs[0].Jobs[i].r = 0
		orgs[0].Jobs[i].w = 1
		orgs[0].Jobs[i].org = 1
		orgs[0].Jobs[i].id = i+1
		orgs[0].Jobs[i].d = vector[i]
	}

	orgs[0].dMax = 4

	//fmt.Println("Organization 0 Jobs = ",orgs[0].Jobs)

	orgs[1].id = 1
	orgs[1].Jobs = make([]Jobs,7)

	vector2 := []int{1,2,3,5,7,8,9}
	
	orgs[1].dMax = 9

	for i:=0; i < 7; i++ {
		orgs[1].Jobs[i].r = 0
		orgs[1].Jobs[i].w = 1
		orgs[1].Jobs[i].org = 2
		orgs[1].Jobs[i].id = i+2+1
		orgs[1].Jobs[i].d = vector2[i]
	}

	orgs[2].id = 2
	orgs[2].Jobs = make([]Jobs,6)

	vector3 := []int{3,3,3,9,9,9}
	
	orgs[2].dMax = 9

	for i:=0; i < 6; i++ {
		orgs[2].Jobs[i].r = 0
		orgs[2].Jobs[i].w = 1
		orgs[2].Jobs[i].org = 3
		orgs[2].Jobs[i].id = i+2+7+1
		orgs[2].Jobs[i].d = vector3[i]
	}

	orgs[3].id = 3
	orgs[3].Jobs = make([]Jobs,11)

	vector4 := []int{4,4,4,4,5,6,13,14,15,17,17}
	
	orgs[3].dMax = 17

	for i:=0; i < 11; i++ {
		orgs[3].Jobs[i].r = 0
		orgs[3].Jobs[i].w = 1
		orgs[3].Jobs[i].org = 4
		orgs[3].Jobs[i].id = i+2+7+6+1
		orgs[3].Jobs[i].d = vector4[i]
	}

	//fmt.Println("Organization 1 Jobs = ",orgs[1].Jobs)
	return true
}

func zipfDistribution(m int, n int64, c float64){
	orgs = make([]Organization,m)

	j:=0
	idJob := 1

	var mFloat float64 = float64(m)

	var uIntn uint64 = uint64(n)

	r := rand.New(rand.NewSource(n))
 	zipf := rand.NewZipf(r, c, mFloat, uIntn)

	for j < m {
		k := zipf.Uint64()
		if k > 0 {
			orgs[j].id = j
			numberOfJobs := rand.Intn(int(n)) + 1
			orgs[j].Jobs = make([]Jobs,numberOfJobs)
			orgs[j].dMax = 1

			for i:=0; i < numberOfJobs; i++ {
				orgs[j].Jobs[i].r = 0
				orgs[j].Jobs[i].w = 1
				orgs[j].Jobs[i].org = j
				orgs[j].Jobs[i].id = idJob
				deadline := int(k)
				orgs[j].Jobs[i].d = deadline
				if deadline > orgs[j].dMax {
					orgs[j].dMax = deadline
				}
				idJob++
			}
			j++
	
		}
	}
	
	//print()

}

func print(){
	for i:=0; i < len(orgs); i++{
		fmt.Printf("Organization %d\n DMax = %d\n",orgs[i].id,orgs[i].dMax)
		fmt.Println(orgs[i].Jobs)

		fmt.Printf("-----------------------------------------\n\n")
	}
}

func start() bool{
	defer timeTrack(time.Now(), "start")

	sort.Sort(ByDMax(orgs))
	//print()

	for i:=0; i < len(orgs); i++{
		sort.Sort(ByDeadline(orgs[i].Jobs))
		jobs := make([]Jobs, len(orgs[i].Jobs))
		copy(jobs,orgs[i].Jobs)
		sort.Sort(ByDistance(jobs))
		wg.Add(1)
		go YDS(jobs,i,true)
	}

	wg.Wait()

	/*fmt.Println("After Sorting")
	print()*/

	for k:=1; k < len(orgs); k++ {
		jobs := make([]Jobs, len(orgs[k].Jobs))
		copy(jobs,orgs[k].Jobs)
		trySchedule(jobs,k)
		for c:=0; c < k; c++ {
			orgs[c].dMax = orgs[c].dLinha
		}
		if len(orgs[k].Jobs) > 0 {
			orgs[k].dMax = orgs[k].Jobs[len(orgs[k].Jobs) -1].d
		}
		sort.Sort(ByDMax(orgs[:k+1]))
	}

	for i:=0; i < len(orgs); i++{
		wg.Add(1)
		go YDS(orgs[i].Jobs,i,true)
	}

	wg.Wait()

	
	/*for i:=0; i < len(orgs); i++{
		fmt.Printf("Schedule Organization %d\n",orgs[i].id)
		fmt.Println(orgs[i].schedule)

		fmt.Printf("-----------------------------------------\n\n")
	}*/

	return true
}

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s took %dms", name, elapsed.Nanoseconds()/1000)
}

func main(){
	m:=20
	var n int64 =15

	zipfDistribution(m,n,1.4267)
	start()
}
