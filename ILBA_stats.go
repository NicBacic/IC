package main

import (
"fmt"
"math/rand"
"sort"
"bufio"
"os"
)

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

type ByMakeSpan []Organization
type ByNum []Organization

func (a ByMakeSpan) Len() int          { return len(a) }
func (a ByMakeSpan) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByMakeSpan) Less(i, j int) bool { return a[i].totalMakeSpan < a[j].totalMakeSpan }

func (a ByNum) Len() int          { return len(a) }
func (a ByNum) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByNum) Less(i, j int) bool { return a[i].num < a[j].num }

func MakeJobs(n int, id int) []Jobs {
	jobs := make([]Jobs, n)

	for i := 0; i < n; i++ {
		jobTime := rand.Intn(1000) + 1
		jobs[i].timeToExecute = float32(jobTime)
		jobs[i].id = id
	}

	return jobs
}

func GetMakeSpan(O []Organization) []float32{
	makeSpan := make([]float32,len(O))
	for i:=0; i < len(O); i++ {
		makeSpan[i] = O[i].totalMakeSpan
	}
	return makeSpan
}

func CalculateMakeSpan(O Organization) float32 {
	for i:=0; i < len(O.p); i++ {
		if O.totalMakeSpan < O.p[i].localMakeSpan {
				O.totalMakeSpan = O.p[i].localMakeSpan
			}
	}
	return O.totalMakeSpan
}

func RecalculateMakeSpan(O []Organization) float32{
	var globalMakeSpan float32 = 0.0
	for N:=0; N < len(O); N++ {
		var currentMakeSpan float32 = 0.0
		for k:=0; k < len(O); k++{
			for p:=0; p < len(O[k].p); p++{
				for j:=0; j < len(O[k].p[p].jobs); j++ {
					if O[k].p[p].jobs[j].id == O[N].num {
						if currentMakeSpan < (O[k].p[p].jobs[j].timeToExecute + O[k].p[p].jobs[j].initTime){
							currentMakeSpan = O[k].p[p].jobs[j].timeToExecute + O[k].p[p].jobs[j].initTime
						}
					}
				}
			}
		}
		O[N].totalMakeSpan = currentMakeSpan
		if globalMakeSpan < O[N].totalMakeSpan {
			globalMakeSpan = O[N].totalMakeSpan
		}
	}
	return globalMakeSpan
}

func LocalSchedule(O []Organization, N int) float32 {
	var globalMakeSpan float32 = 0.0
	for k:=0; k < N; k++ {
		O[k] = Schedule(O[k])
		if globalMakeSpan < O[k].totalMakeSpan {
			globalMakeSpan = O[k].totalMakeSpan
		}
	}
	return globalMakeSpan
}

func Schedule(O Organization) Organization{
	j := 0

	//Schedule Jobs into J Processors
	for i := 0; i < len(O.jobs); i++ {
		if len(O.p[j].jobs) == 0 {
			O.p[j].localMakeSpan = 0
		}
		O.jobs[i].initTime = O.p[j].localMakeSpan
		O.p[j].localMakeSpan += O.jobs[i].timeToExecute
		O.p[j].jobs = append(O.p[j].jobs, O.jobs[i])
		j++
		if len(O.p) == j {
			j = 0
		}
	}

	//Calculate TotalMakeSpan
	O.totalMakeSpan = CalculateMakeSpan(O)
	
	return O
}

func CreateOrganizations(k int) []Organization{
	O := make([]Organization,k)
	

	for i:= 0; i < k; i++{
		O[i].num = i
		n:= rand.Intn(900) + 100//number of Jobs
		numProcessors := rand.Intn(10)

		O[i].p = make([]Processor,numProcessors + 1) // Plus 1 as rand can be 0.
	
		O[i].jobs = MakeJobs(n,i)

		O[i].totalMakeSpan = 0
	}

	return O
}

func EvaluateDifference(O []Organization) float32{
	var dif float32 = 0.0
	var best float32 = 0.0
	var min float32 = O[0].p[0].localMakeSpan
	
	for i:=0; i < len(O); i++{
		for j:=0; j < len(O[i].p); j++{
			if best < O[i].p[j].localMakeSpan {
				best = O[i].p[j].localMakeSpan
			}
			if min > O[i].p[j].localMakeSpan {
				min = O[i].p[j].localMakeSpan
			}
		}
	}
	dif = best - min
	return dif
}

func ILBA(){
	fmt.Println("Starting Instances...")

	//Create a random number of simulations
	simulacao := 100
	it := 0
	var N int = 5
	var aux int = 0
	stats := make([]float32,simulacao)
	globalMatrix := make([]float32,simulacao)
	minMatrix := make([]float32,simulacao)
	maxMatrix := make([]float32,simulacao)
	spanMatrix := make([]float32,simulacao)
	jobsMatrix := make([]int,simulacao)
	newMatrix := make([]float32,simulacao)
	oldMatrix := make([]float32,simulacao)
	dif := make([]float32,simulacao)

	f,err := os.Create("/tmp/ILBA")
	check(err)
	
	defer f.Close()
	f.Sync()
	outputWriter := bufio.NewWriter(f)

	f2,err2 := os.Create("/tmp/ILBAFive")
	check(err2)
	
	defer f2.Close()
	f2.Sync()
	outputFive := bufio.NewWriter(f2)
	
	fmt.Fprintf(outputWriter,"Iteracao Org Jobs Old New PerLocal PerGlobal MinJob MaxJob Dif MakeSpanDif\n")
	fmt.Fprintf(outputFive,"Iteracao Org Jobs Old New PerLocal PerGlobal MinJob MaxJob Dif MakeSpanDif\n")
	for it < simulacao { // Faz 10 simulacoes
		
		fmt.Printf("Starting Instance (%d)...\n\n", it)

		fmt.Printf("(%d) Organizations Created!!\n\n", N)

		//Create N Organizations
		O := CreateOrganizations(N) 
		
		//Local Scheduling

		oldMatrix[it] = LocalSchedule(O,N)
		oldMakeSpan := GetMakeSpan(O)
		//fmt.Printf("BEFORE ILBA\n")
		//PrintMakeSpan(O)

		//Sort Orgs by MakeSpan

		//PrintNumOrg(O) //Debug code: Before Sorting

		sort.Sort(ByMakeSpan(O))
		
		spanMatrix[it] = EvaluateDifference(O)
		//PrintNumOrg(O) //Debug code: After Sorting

		//PrintOrg(O,N)

		//Global Scheduling ILBA

		/*How ILBA works: For k = 2 until N (organizations), jobs from O(k) are rescheduled sequentially
		* and assigned to the less loaded organizations O(1) .... O(k).
		* Each job is rescheduled by ILBA either earlier or at the same time that the job was scheduled 
		* before the migration.
		* In other words, no job is delayed by ILBA, which guarantees that the local constraint is respected 
		* for MOSP(Cmax) and MOSP(Sum of Ci)*/	
		for k:=1; k < N; k++ {
			for i:=0; i < k; i++ {
				for pi:=0; pi < len(O[i].p); pi++{		
					for p := 0; p < len(O[k].p); p++ {
						for j := len(O[k].p[p].jobs)-1; j > -1; j-- { //Start from the job with the greatest initTime
							//If check jobs that can be rescheduled
							if O[k].p[p].jobs[j].initTime >= O[i].p[pi].localMakeSpan{ 
								O[k].p[p].jobs[j].initTime = O[i].p[pi].localMakeSpan //update job initTime
								 //assign the job to the processor Pi from Organization O(i)
								O[i].p[pi].jobs = append(O[i].p[pi].jobs, O[k].p[p].jobs[j])
								time := O[k].p[p].jobs[j].timeToExecute
								O[i].p[pi].localMakeSpan += time //update makeSpan of P(pi)

								//Remove job rescheduled from processor p[k] of O[k]
								copy(O[k].p[p].jobs[j:], O[k].p[p].jobs[j+1:])
								O[k].p[p].jobs[len(O[k].p[p].jobs)-1].timeToExecute = 0.0
								O[k].p[p].jobs[len(O[k].p[p].jobs)-1].initTime = 0.0 
								O[k].p[p].jobs = O[k].p[p].jobs[:len(O[k].p[p].jobs)-1]
								O[k].p[p].localMakeSpan = O[k].p[p].localMakeSpan - time
								//Finish removing
							} else {
								break; //No more Jobs from Processor p of O[k] can be rescheduled
							}
						}
					}	
				}
			}
			
		}
		
		newMatrix[it] = RecalculateMakeSpan(O)
		sort.Sort(ByNum(O))
		newMakeSpan := GetMakeSpan(O)
		stats[it],globalMatrix[it],jobsMatrix[it],minMatrix[it],maxMatrix[it] = Statistics(oldMakeSpan,newMakeSpan,O,oldMatrix[it],newMatrix[it])
		dif[it] = maxMatrix[it] - minMatrix[it]
		fmt.Fprintf(outputWriter,"%d %d %d %f %f %f %f %f %f %f %f\n",it,len(O),jobsMatrix[it],oldMatrix[it],newMatrix[it],stats[it],globalMatrix[it],minMatrix[it],maxMatrix[it],dif[it],spanMatrix[it])
		it++
		if it % 5 == 0 {
			var localMax float32 = 0.0
			var globalMax float32 = 0.0
			var jobsMax int = 0
			var minMax float32 = 0.0
			var maxMax float32 = 0.0
			var difMax float32 = 0.0
			var spanMax float32 = 0.0
			var oldMax float32 = 0.0
			var newMax float32 = 0.0
			for i:=aux; i < it; i++{
				localMax += stats[i]
				globalMax += globalMatrix[i]
				jobsMax += jobsMatrix[i]
				minMax += minMatrix[i]
				maxMax += maxMatrix[i]
				difMax += dif[i]
				spanMax += spanMatrix[i]
				oldMax += oldMatrix[i]
				newMax += newMatrix[i]
			}
			var meanLocal float32 = localMax / 5
			var meanGlobal float32 = globalMax / 5
			var meanJobs int = jobsMax / 5
			var meanMin float32 = minMax / 5
			var meanMax float32 = maxMax / 5
			var meanDif float32 = difMax / 5
			var meanSpan float32 = spanMax / 5
			var meanOld float32 = oldMax / 5
			var meanNew float32 = newMax / 5
			fmt.Fprintf(outputFive,"%d %d %d %f %f %f %f %f %f %f %f\n",it,N,meanJobs,meanOld,meanNew,meanLocal,meanGlobal,meanMin,meanMax,meanDif,meanSpan)
			aux = N
			N+=5
		}
	}
	outputFive.Flush()
	outputWriter.Flush()
}

func Statistics(oldMakeSpan []float32,newMakeSpan []float32, O []Organization,oldGlobalMakeSpan float32, newGlobalMakeSpan float32) (float32,float32,int,float32,float32){
	var best float32 = 0.0
	dataSub := make([]float32,len(newMakeSpan))
	var subTotal float32 = 0.0

	percentageImprovement := make([]float32,len(newMakeSpan))
	var totalPercentageImprovement float32 = 0.0
	jobsTotal := 0
	for i:=0; i < len(newMakeSpan); i++{
		dataSub[i] = oldMakeSpan[i] - newMakeSpan[i]
		subTotal += dataSub[i]
		if dataSub[i] > best {
			best = dataSub[i]
		}
		jobsTotal += len(O[i].jobs)
		percentageImprovement[i] = 100 - ((newMakeSpan[i] * 100) / (oldMakeSpan[i]))
		totalPercentageImprovement += percentageImprovement[i]
	}
	var averageJobs int = jobsTotal / len(O)
	var averageSub float32 = subTotal / float32(len(O))
	var averagePercentageImprovement float32 = totalPercentageImprovement / float32(len(O))

	var globalPercentageImprovement float32 = 100 - ((newGlobalMakeSpan * 100)/oldGlobalMakeSpan)

	//PRINT
	fmt.Printf("Total Jobs = %d\n",jobsTotal)
	fmt.Printf("Average Jobs = %d\n\n",averageJobs)

	fmt.Printf("Old Global MakeSpan = %f\n",oldGlobalMakeSpan)
	fmt.Printf("New Global MakeSpan = %f\n",newGlobalMakeSpan)
	fmt.Printf("Improvement in GlobalMake Span = %f\n", oldGlobalMakeSpan - newGlobalMakeSpan)
	fmt.Printf("Improvement in Percentage = %f\n",globalPercentageImprovement)
	
	fmt.Printf("Average Improvement in MakeSpan = %f\n",averageSub)
	fmt.Printf("Average Improvement in MakeSpan in Percentage = %f%%\n\n",averagePercentageImprovement)
	
	fmt.Printf("Total Improvement = %f\n\n",totalPercentageImprovement)
	
	fmt.Printf("Improvement per Organization\n")

	var minJob float32 = 10000
	var maxJob float32 = 0.0
	var totalExecTime float32 = 0.0
	meanExecTimePerOrg := make([]float32,len(O))
	

	for i:=0; i < len(O); i++ {
		//fmt.Printf("-------Organization (%d)\nJobs = (%d)\nOld MakeSpan = (%f)\nNew MakeSpan = (%f)\nImprovement = (%f)\nPercentage = (%f%%)\n\n",O[i].num,len(O[i].jobs),oldMakeSpan[i],newMakeSpan[i],dataSub[i],percentageImprovement[i])

		

		for j := 0; j < len(O[i].jobs); j++{
			
			totalExecTime += O[i].jobs[j].timeToExecute
			
		}
		meanExecTimePerOrg[i] = totalExecTime / float32(len(O[i].jobs))
		totalExecTime = 0
	}
	
	for i:= 0; i < len(meanExecTimePerOrg); i++ {
		if minJob > meanExecTimePerOrg[i]{
			minJob = meanExecTimePerOrg[i]
		}
		if maxJob < meanExecTimePerOrg[i]{
			maxJob = meanExecTimePerOrg[i]
		}
	}
	return averagePercentageImprovement,globalPercentageImprovement,jobsTotal,minJob,maxJob
} 

func PrintNumOrg(org []Organization){
	for i:=0; i < len(org); i++{
		fmt.Printf("(%d) ",org[i].num)
	}
	fmt.Printf("\n\n")
}

func PrintMakeSpan(O []Organization){
	for i:=0; i < len(O); i++{
		fmt.Printf("O(%d) makeSpan = (%f)\n\n",O[i].num,O[i].totalMakeSpan)
	}
}

func PrintAllInfoOrg(org []Organization, k int){
	for i:=0; i < k; i++{
		fmt.Printf("Organization (%d) has (%d) Processors, (%d) Jobs and MakeSpan (%f) \n",org[i].num,len(org[i].p),len(org[i].jobs),org[i].totalMakeSpan)
		for j:=0; j < len(org[i].p); j++ {
			fmt.Printf("Processor (%d) is number (%d) has (%d) Jobs and MakeSpan (%f)\n",j,org[i].p[j].num,len(org[i].p[j].jobs),org[i].p[j].localMakeSpan)
			for l:=0; l < len(org[i].p[j].jobs); l++ {
				fmt.Printf("Job (%d) has time to Execute(%f)\n",l,org[i].p[j].jobs[l].timeToExecute)
			}
			fmt.Printf("\n")
		}
		fmt.Printf("------------------------------\n")
	}
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main(){
	ILBA()
}
