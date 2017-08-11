package Test

import(
"sort"
)
//0800 7010321

func startParalelo(m int, n int64) bool{

	zipfDistribution(m,n,1.4267)

	sort.Sort(ByDMax(orgs))

	/*for k:=1; k < len(orgs); k++ {
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
	}*/

	for i:=0; i < len(orgs); i++{
		wg.Add(1)
		go YDSParalelo(orgs[i].Jobs,i)
   		
	}

	wg.Wait()

	return true
}
