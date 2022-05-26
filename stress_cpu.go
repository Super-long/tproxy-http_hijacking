package main

import (
	"time"
	"fmt"
	"unsafe"
	"math/rand"
	"math/big"
	"math"
	"github.com/mjibson/go-dsp/fft"
	"github.com/howeyc/crc16"
)

type StressCpuMethodInfo struct {
	name  			string			/* human readable form of stressor */
	stress 			func(string) 	/* the cpu method function */
}

var cpu_methods = []StressCpuMethodInfo {
	{ "ackermann", 	stress_cpu_ackermann,	},
	{ "bitops",		stress_cpu_bitops,		},
	{ "collatz",	stress_cpu_collatz,		},
	{ "crc16",		stress_cpu_crc16,		},
	{ "factorial",	stress_cpu_factorial,	},
	{ "fft", 		stress_cpu_fft,         },
	{ "pi", 		stress_cpu_pi,			}, 
	{ "fibonacci",	stress_cpu_fibonacci,	},
}

func stress_cpu_factorial(name string) {
	var f float64 = 1.0
	var precision float64 = 1.0e-6

	for n := 1; n < 150; n++ {
		np1 := float64(n + 1)
		fact := math.Round(math.Exp(math.Gamma(np1)))
		var dn float64

		f *= float64(n);

		/* Stirling */
		if (f - fact) / fact > precision {
			fmt.Println("%s: Stirling's approximation of factorial(%d) out of range\n",
				name, n);
		}

		/* Ramanujan */
		dn = float64(n);
		fact = math.SqrtPi * math.Pow((dn / float64(math.E)), dn)
		fact *= math.Pow((((((((8 * dn) + 4)) * dn) + 1) * dn) + 1.0/30.0), (1.0/6.0));
		if ((f - fact) / fact > precision) {
			fmt.Println("%s: Ramanujan's approximation of factorial(%d) out of range\n",
				name, n);
		}
	}
}

func stress_cpu_fft(name string) {
	var buffer [128]float64
	for i := 0; i < 128; i++ {
		buffer[i] = float64(i%64)
	}
	for i := 0; i < 8; i++ {
		fft.FFTReal(buffer[:])
	}
}

func stress_cpu_fibonacci(name string) {
	var fn_res uint64 = 0xa94fad42221f2702
	var f1 uint64 = 1
	var f2 uint64 = 1
	var fn uint64 = 1

	for !(fn & 0x8000000000000000 != 0) {
		fn = f1 + f2
		f1 = f2
		f2 = fn
	}

	if fn_res != fn {
		fmt.Printf("%s: fibonacci error detected, summation or assignment failure\n", name);
	}
}


func stress_cpu_collatz(name string) {
	var n uint64 = 989345275647
	var i int
	for i = 0; n != 1; i++ {
		if n&1 != 0 {
			n = (3 * n) + 1
		} else {
			n = n / 2
		}
	}

	if i != 1348 {
		fmt.Printf("%s: error detected, failed collatz progression\n", name)
	}
}

func stress_cpu_crc16(name string) {
	var randomBuffer [4096]byte
	rand.Read(randomBuffer[:])
	for i := 0; i < 8; i++ {
		crc16.ChecksumIBM(randomBuffer[:])
	}
}

func ackermann(m uint32, n uint32) uint32 {
	if m == 0 {
		return n + 1
	} else if n == 0 {
		return ackermann(m - 1, 1)
	} else {
		return ackermann(m - 1, ackermann(m, n - 1))
	}
}

func stress_cpu_ackermann(name string) {
	a := ackermann(3, 7);

	if a != 0x3fd {
		fmt.Printf("%s: ackermann error detected, ackermann(3,9) miscalculated\n", name);
	}
}

func stress_cpu_bitops(name string) {
	var i_sum uint32 = 0
	var sum uint32 = 0x8aac0aab

	for i := 0; i < 16384; i++ {
		{
			var r uint32 = uint32(i)
			var v uint32 = uint32(i)
			var s uint32 = uint32((unsafe.Sizeof(v) * 8) - 1)
			for v >>= 1; v != 0; v, s = v>>1, s-1 {
				r <<= 1
				r |= v & 1
			}
			r <<= s
			i_sum += r
		}
		{
			/* parity check */
			var v uint32 = uint32(i)

			v ^= v >> 16
			v ^= v >> 8
			v ^= v >> 4
			v &= 0xf
			i_sum += (0x6996 >> v) & 1
		}
		{
			/* Brian Kernighan count bits */
			var v uint32 = uint32(i)
			var j uint32 = uint32(i)

			for j = 0; v != 0; j++ {
				v &= v - 1
			}
			i_sum += j
		}
		{
			/* round up to nearest highest power of 2 */
			var v uint32 = uint32(i - 1)

			v |= v >> 1
			v |= v >> 2
			v |= v >> 4
			v |= v >> 8
			v |= v >> 16
			i_sum += v
		}
	}
	if i_sum != sum {
		fmt.Printf("%s: bitops error detected, failed bitops operations\n", name)
	}
}

func stress_cpu_pi(name string) {
	digits := big.NewInt(1000)
	unity := big.NewInt(0)
	unity.Exp(big.NewInt(10), digits, nil)
	pi := big.NewInt(0)
	four := big.NewInt(4)
	pi.Mul(four, pi.Sub(pi.Mul(four, arccot(5, unity)), arccot(239, unity)))
}

func arccot(x int64, unity *big.Int) *big.Int {
	bigx := big.NewInt(x)
	xsquared := big.NewInt(x*x)
	sum := big.NewInt(0)
	sum.Div(unity, bigx)
	xpower := big.NewInt(0)
	xpower.Set(sum)
	n := int64(3)
	zero := big.NewInt(0)
	sign := false
	
	term := big.NewInt(0)
	for {
		xpower.Div(xpower, xsquared)
		term.Div(xpower, big.NewInt(n))
		if term.Cmp(zero) == 0 {
			break
		}
		if sign {
			sum.Add(sum, term)
		} else {
			sum.Sub(sum, term)
		}
		sign = !sign
		n += 2
	}
	return sum
}


func stress_cpu_method(method int) {
	cpu_methods[method].stress("lizhaolong");
}

func stress_cpu(interval time.Duration, cpuPercent float64) {
	bias := 0.0
	startTime := time.Now().UnixNano()
	nanoInterval := int64(interval/time.Nanosecond)
	fmt.Printf("[%d]=========nanoInterval\n", nanoInterval)
	time_exec := make(map[string]time.Duration, 8)
	iterator := 0
	for {
		if time.Now().UnixNano() - startTime > nanoInterval {
			break
		}

		startTime1 := time.Now().UnixNano()
		iterator++
		// Loops and methods may be specified later.
		for i := 0; i < len(cpu_methods); i++ {
			t1 := time.Now()
			stress_cpu_method(i)
			t2 := time.Now()
			time_exec[cpu_methods[i].name] += t2.Sub(t1)
			//fmt.Println(cpu_methods[i].name, t2.Sub(t1))

		}
		//fmt.Println("=============================")
		endTime1 := time.Now().UnixNano()
		//fmt.Println(startTime1, endTime1, cpuPercent)
		delay := ((100 - cpuPercent) * float64(endTime1 - startTime1) / cpuPercent)
		//fmt.Printf("delay : [%f], bias : [%f]\n", delay, bias)
		delay -= bias
		if delay <= 0.0 {
			bias = 0.0;
		} else {
			startTime2 := time.Now().UnixNano()
			time.Sleep(time.Duration(delay) * time.Nanosecond)
			endTime2 := time.Now().UnixNano()
			bias = float64(endTime2 - startTime2) - delay
		}
	}
	for k, _ := range time_exec {
		time_exec[k] /= time.Duration(iterator)
		fmt.Println(k, time_exec[k])
	}
	fmt.Println(iterator)
}

func main() {
	stress_cpu(time.Duration(20)*time.Second, 40)
}