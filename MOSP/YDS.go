package main

import (
"fmt"
"math/rand"
"sort"
)
/*************************************************

This implementation can be optimized by using pointers and list instead of vectors.

type Jobs struct {
timeToExecute float32
initTime float32
id int // We use ID to know from which Organization this job belongs
}

type Processor struct{
jobs []Jobs
num int
localMakeSpan float32
}

type Organization struct{
num int
jobs []Jobs
p []Processor
totalMakeSpan float32
}

************************************************/


/*Sobre o Algoritmo YDS: O teorema 3.1.2 sugere uma estratégia gulosa para encontrar um escalonamento ótimo iterativamente. Basta, a cada iteração, identificar o intervalo de densidade máxima I* bem com as tarefas Ji* e processá-las com velocidade delta I*. Feito isso, removemos o intervalo I* de consideração, isto é, nenhuma outra tarefa será executada naquele intervalo. Caso alguma tarefa tenha seu intervalo de execução parcialmente dentro do intervalo de densidade máxima, seu tempo de chegada ou prazo são ajustados, isto é, modificados para coincidir com o limite do intervalo a qual intersectam.

Pseudo código: 

ALGORITMO YDS (CONJUNTO DE TAREFAS J (J1, J2, ... , Jn))

01. enquanto j!= 0

02. I*, deltaI*, Ji* <- intervalo de densidade máxima, adensidade e o conjunto de tarefas do intervalo

03. execute as tarefas de JI*, com velocidade delta I* em I*, segundo a politica EDF

04. J = J-Ji*

05. remova o intervalo I* de consideração.

Sabendo disso vamos construir um ambiente que imite as condições de um ambiente real.

Teremos a princípio um conjunto de organizações com id e apenas um processador com uma dada velocidade s e com um conjunto de tarefas dentro dele

As tarefas terão um identificador, realease date, deadline e volume de processamento.

O volume de processamento de um intervalo é dado por Somatório de wi / módulo(min(release date) - max(deadline))*/

type Jobs struct{
id int
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
releaseDate float32
deadLine float32
speed float32
}

type ByReleaseDate []Schedule
type ByDeadline []Jobs

func (a ByReleaseDate) Len() int         { return len(a) }
func (a ByReleaseDate) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByReleaseDate) Less(i, j int) bool { return a[i].releaseDate < a[j].releaseDate }

func (a ByDeadline) Len() int          { return len(a) }
func (a ByDeadline) Swap(i, j int) { a[i], a[j] = a[j], a[i]}
func (a ByDeadline) Less(i, j int) bool { return a[i].d < a[j].d }


func random() (int,int,int){
	r:= rand.Intn(10)
	d:= rand.Intn(10) + r + 1
	w:= rand.Intn(10) + 1
	return r,d,w
}

func build(n int) Processor{

	//Define a Processor
	p := Processor{}
	p.Jobs = make([]Jobs,n)

	//Release date, deadline and weight
	/*r:= 0
	d:= 1
	w:= 1
	offset := 0.0
	for i := 0; i < n; i++ {
		r,d,w = random()
		job := Jobs{r,d,w,i,offset}
		p.Jobs[i] = job
	}*/
	
	/* Caso teste
	p.Jobs[0] = Jobs{1, 3, 6, 5,0}
        p.Jobs[1] = Jobs{2, 2, 6, 3,0}
        p.Jobs[2] = Jobs{3, 0, 8, 2,0}
        p.Jobs[3] = Jobs{4, 6, 14, 6,0}
        p.Jobs[4] = Jobs{5, 10, 14, 6,0}
        p.Jobs[5] = Jobs{6, 11, 17, 2,0}
        p.Jobs[6] = Jobs{7, 12, 17, 2,0}

	 yds.Task("t1",  0, 17,  5),
    yds.Task("t2",  1, 11,  3),
    yds.Task("t3", 12, 20,  4),
    yds.Task("t4",  7, 11,  2),
    yds.Task("t5",  1, 20,  4),
    yds.Task("t6", 14, 20, 12),
    yds.Task("t7", 14, 17,  4),
    yds.Task("t8",  1,  7,  2)
				*/


	p.Jobs[0] = Jobs{1, 0, 17, 5,0}
        p.Jobs[1] = Jobs{2, 1, 11, 3,0}
        p.Jobs[2] = Jobs{3, 12, 20, 4,0}
        p.Jobs[3] = Jobs{4, 7, 11, 2,0}
        p.Jobs[4] = Jobs{5, 1, 20, 4,0}
        p.Jobs[5] = Jobs{6, 14, 20, 12,0}
        p.Jobs[6] = Jobs{7, 14, 17, 4,0}
	p.Jobs[7] = Jobs{8, 1, 7, 2,0}

	/*
	p.Jobs[0] = Jobs{0,3,3,0,0.0}
	p.Jobs[1] = Jobs{1,4,2,1,0.0}
	p.Jobs[2] = Jobs{2,3,1,2,0.0}*/
	return p
}


func YDS(jobs []Jobs) []Schedule{
	fmt.Printf("\nIniciando YDS com %d tarefas\n", len(jobs))
	
	var intervaloDensidadeMax Interval
	var schedule []Schedule

	var auxSchedule Schedule

	var tamanhoIntervalo float32 = 0.0
	var offset float32 = 0.0
	var tamanhoTarefa float32 = 0.0
	var somaProcessamento float32 = 0.0

	i:=0
	
	for len(jobs) != 0 {
		fmt.Printf("Iteração %d\n\n",i)

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

			auxSchedule = Schedule{intervaloDensidadeMax.Jobs[j].id,intervaloDensidadeMax.Jobs[j].offset + offset,
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
	return schedule
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

func main(){
	p := build(8)
	fmt.Println(p.Jobs)
	schedule := YDS(p.Jobs)
	sort.Sort(ByReleaseDate(schedule))
	fmt.Println(schedule)
}
