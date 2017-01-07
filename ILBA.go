import (
"math/rand"
"sort"
)

//This implementation can be optimized by using pointers and list instead of vectors.

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

func (a ByMakeSpan) Len() int          { return len(a) }
func (a ByMakeSpan) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByMakeSpan) Less(i, j int) bool { return a[i].totalMakeSpan < a[j].totalMakeSpan }

func MakeJobs(n int, id int) []Jobs {
	jobs := make([]Jobs, n)

	for i := 0; i < n; i++ {
		jobTime := rand.Intn(1000) + 1
		jobs[i].timeToExecute = float32(jobTime)
		jobs[i].id = id
	}

	return jobs
}

func CalculateMakeSpan(O Organization) float32 {
	for i:=0; i < len(O.p); i++ {
		if O.totalMakeSpan < O.p[i].localMakeSpan {
				O.totalMakeSpan = O.p[i].localMakeSpan
			}
	}
	return O.totalMakeSpan
}

/*func RecalculateMakeSpan(O []Organization) float32{
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
}*/

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
		numProcessors := rand.Intn(10) + 1

		O[i].p = make([]Processor,numProcessors)
	
		O[i].jobs = MakeJobs(n,i)

		O[i].totalMakeSpan = 0
	}

	return O
}

func ILBA(){
	var N int = 5//inicial number of Organizations
	
	//Create N Organizations
	O := CreateOrganizations(N) 
		
	//Local Scheduling

	LocalSchedule(O,N)

	//Sort Orgs by MakeSpan

	sort.Sort(ByMakeSpan(O))

	//Global Scheduling ILBA

		/*How ILBA works: For k = 2 until N (organizations), jobs from O(k) are rescheduled sequentially
		* and assigned to the less loaded organizations O(1) .... O(k).
		* Each job is rescheduled by ILBA either earlier or at the same time that the job was scheduled 
		* before the migration.
		* In other words, no job is delayed by ILBA, which guarantees that the local constraint is respected 
		* for MOSP(Cmax) and MOSP(Sum of Ci)*/	

	for k:=1; k < N; k++ { //for each k organization
		for i:=0; i < k; i++ { 
			for pi:=0; pi < len(O[i].p); pi++{ //Search if any job can be scheduled into Processor P(i)	
				for p := 0; p < len(O[k].p); p++ { //Jobs from Processor P(k)
					for j := len(O[k].p[p].jobs)-1; j > -1; j-- { //Start from the job with the greatest initTime

						//If check jobs that can be rescheduled
						if O[k].p[p].jobs[j].initTime >= O[i].p[pi].localMakeSpan { 
							O[k].p[p].jobs[j].initTime = O[i].p[pi].localMakeSpan //update job initTime

							 //assign the job to the processor Pi from Organization O(i)
							O[i].p[pi].jobs = append(O[i].p[pi].jobs, O[k].p[p].jobs[j])
							time := O[k].p[p].jobs[j].timeToExecute
							O[i].p[pi].localMakeSpan += time //update makeSpan of Processor(pi)

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
	//Finish ILBA. You can recalculate each local makespan by calling RecalculateMakeSpan
}
