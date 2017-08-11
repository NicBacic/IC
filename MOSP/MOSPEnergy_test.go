package Test

import(
"testing"
)

func benchmarkZipf(m int, n int64, b *testing.B){
	for i := 0; i < b.N; i++ {
		zipfDistribution(m,n,1.4267)
	}
}

func benchmarkMOSPEnergy(m int, n int64, b *testing.B){
	for i := 0; i < b.N; i++ {
		start(m,n)
	}
}

func benchmarkMOSPParalelo(m int, n int64, b *testing.B){
	for i := 0; i < b.N; i++ {
		startParalelo(m,n)
	}
}

func benchmarkMOSPBuffered(m int, n int64, b *testing.B){
	for i := 0; i < b.N; i++ {
		startBuffered(m,n)
	}
}

func benchmarkYDS(m int, n int64, b *testing.B){
	for i := 0; i < b.N; i++ {
		startYDS(m,n)
	}
}

func benchmarkIntervalo(m int, n int64, b *testing.B){
	for i := 0; i < b.N; i++ {
		startCalcula(m,n)
	}
}

func BenchmarkZipf0530(b *testing.B) { benchmarkZipf(5,30,b) }
func BenchmarkMOSPEnergy0530(b *testing.B)  { benchmarkMOSPEnergy(5,30, b) }
func BenchmarkMOSPParalelo0530(b *testing.B)  { benchmarkMOSPParalelo(5,30, b) }

func BenchmarkZipf3015(b *testing.B) { benchmarkZipf(30,15,b) }
func BenchmarkMOSPEnergy3015(b *testing.B)  { benchmarkMOSPEnergy(30,15, b) }
func BenchmarkMOSPParalelo3015(b *testing.B)  { benchmarkMOSPParalelo(30,15, b) }

func BenchmarkZipf2015(b *testing.B) { benchmarkZipf(20,15,b) }
func BenchmarkMOSPEnergy2015(b *testing.B)  { benchmarkMOSPEnergy(20,15, b) }
func BenchmarkMOSPParalelo2015(b *testing.B)  { benchmarkMOSPParalelo(20,15, b) }

func BenchmarkZipf2005(b *testing.B) { benchmarkZipf(20,5,b) }
func BenchmarkMOSPEnergy2005(b *testing.B)  { benchmarkMOSPEnergy(20,5, b) }
func BenchmarkMOSPParalelo2005(b *testing.B)  { benchmarkMOSPParalelo(20,5, b) }

func BenchmarkZipf1015(b *testing.B) { benchmarkZipf(10,15,b) }
func BenchmarkMOSPEnergy1015(b *testing.B)  { benchmarkMOSPEnergy(10,15, b) }
func BenchmarkMOSPParalelo1015(b *testing.B)  { benchmarkMOSPParalelo(10,15, b) }

func BenchmarkZipf1005(b *testing.B) { benchmarkZipf(10,5,b) }
func BenchmarkMOSPEnergy1005(b *testing.B)  { benchmarkMOSPEnergy(10,5, b) }
func BenchmarkMOSPParalelo1005(b *testing.B)  { benchmarkMOSPParalelo(10,5, b) }

func BenchmarkZipf300015(b *testing.B) { benchmarkZipf(3000,15,b) }
func BenchmarkMOSPEnergy300015(b *testing.B)  { benchmarkMOSPEnergy(3000,15, b) }
func BenchmarkMOSPParalelo300015(b *testing.B)  { benchmarkMOSPParalelo(3000,15, b) }
func BenchmarkMOSPBuffered300015(b *testing.B)  { benchmarkMOSPBuffered(3000,15, b) }

func BenchmarkZipf100015(b *testing.B) { benchmarkZipf(1000,15,b) }
func BenchmarkMOSPEnergy100015(b *testing.B)  { benchmarkMOSPEnergy(1000,15, b) }
func BenchmarkMOSPParalelo100015(b *testing.B)  { benchmarkMOSPParalelo(1000,15, b) }
func BenchmarkMOSPBuffered100015(b *testing.B)  { benchmarkMOSPBuffered(1000,15, b) }

func BenchmarkZipf30015(b *testing.B) { benchmarkZipf(300,15,b) }
func BenchmarkMOSPEnergy30015(b *testing.B)  { benchmarkMOSPEnergy(300,15, b) }
func BenchmarkMOSPParalelo30015(b *testing.B)  { benchmarkMOSPParalelo(300,15, b) }
func BenchmarkMOSPBuffered30015(b *testing.B)  { benchmarkMOSPBuffered(300,15, b) }

func BenchmarkYDS0130(b *testing.B)  { benchmarkYDS(1,30, b) }
func BenchmarkIntervalo0130(b *testing.B)  { benchmarkIntervalo(1,30, b) }
func BenchmarkZipf0130(b *testing.B) { benchmarkZipf(1,30,b) }

func BenchmarkYDS0115(b *testing.B)  { benchmarkYDS(1,15, b) }
func BenchmarkIntervalo0115(b *testing.B)  { benchmarkIntervalo(1,15, b) }
func BenchmarkZipf0115(b *testing.B) { benchmarkZipf(1,15,b) }

func BenchmarkYDS0105(b *testing.B)  { benchmarkYDS(1,5, b) }
func BenchmarkIntervalo0105(b *testing.B)  { benchmarkIntervalo(1,5, b) }
func BenchmarkZipf0105(b *testing.B) { benchmarkZipf(1,5,b) }

