package main

import (
	"time"
	"fmt"
	"unsafe"
)

type StressCpuMethodInfo struct {
	name  			string			/* human readable form of stressor */
	stress 			func(string) 	/* the cpu method function */
}

var cpu_methods = []StressCpuMethodInfo {
	{ "ackermann", 	stress_cpu_ackermann,	},
	{ "bitops",		stress_cpu_bitops,		},
	// { "collatz",	stress_cpu_collatz,		},
	// { "crc16",		stress_cpu_crc16,		},
	// { "factorial",	stress_cpu_factorial,	},
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

func stress_cpu_method(method int) {
	cpu_methods[method].stress("lizhaolong");
}

func stress_cpu(interval time.Duration, cpuPercent float64) {
	bias := 0.0
	startTime := time.Now().UnixNano()
	nanoInterval := int64(interval/time.Nanosecond)
	fmt.Printf("[%d]=========nanoInterval\n", nanoInterval)
	for {
		if time.Now().UnixNano() - startTime > nanoInterval {
			break
		}

		startTime1 := time.Now().UnixNano()
		// Loops and methods may be specified later.
		for i := 0; i < len(cpu_methods); i++ {
			stress_cpu_method(i)
		}
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
}

func main() {

	stress_cpu(time.Duration(1000)*time.Second, 40)

}