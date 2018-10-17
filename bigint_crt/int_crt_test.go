package bigint_crt

//go test ./bigint_crt -v
//go test -bench=. ./bigint_crt -v

import (
	"testing"
	//"math"
	"./bigint"
	//"github.com/Pro7ech/lago/bigint"
	//"fmt"
)



//All qi are equivalent to 1 mod 2n (here n = 1024)
const q_0 = uint64(2147493889)
const q_1 = uint64(2147555329)
const q_2 = uint64(2147565569)
const q_3 = uint64(2147573761) 

const x_0  = uint64(112233445566778899)
const x_1  = uint64(998877665544332211)

var Q_FACTORS = []uint64{q_0,q_1,q_2,q_3}
var Q_FACTORS_LEN = uint16(4)




type arg_create_crt struct{
	v uint64
	want []uint64
}

var create_crt_vectors = []arg_create_crt{
	{uint64(0), []uint64{0,0,0,0}},
	{uint64(1), []uint64{1,1,1,1}},
	{x_0, []uint64{x_0%q_0,x_0%q_1,x_0%q_2,x_0%q_3}},
	{x_1, []uint64{x_1%q_0,x_1%q_1,x_1%q_2,x_1%q_3}},
}


func TestCreate_crt(t *testing.T){

	var z int_64_crt
	var y int_64_crt

	for i, testPair := range create_crt_vectors {

		y.bigint_64_crt = testPair.want

		if !z.EQUAL(NewInt_64_crt(testPair.v, &Q_FACTORS, &Q_FACTORS_LEN), &y){
			t.Errorf("Error creating crt vectors pair %v",i)
		}
	}
}

func BenchmarkCreate_crt(b *testing.B){
	v := uint64(112233445566778899)
	for i:=0 ; i< b.N; i++{
		NewInt_64_crt(v, &Q_FACTORS, &Q_FACTORS_LEN)
	}
}



func TestRecombine_crt(t *testing.T){

	for i, testPair := range create_crt_vectors {

		outputTest := NewInt_64_crt(testPair.v, &Q_FACTORS, &Q_FACTORS_LEN).CRT_INV()

		var expectedResult = bigint.NewInt(int64(testPair.v))
		
		if !outputTest.EqualTo(expectedResult){
			t.Errorf("Error crt recombine pair %v", i)
		}
	}

}

func BenchmarkRecombine_crt(b *testing.B){

	vectors := NewInt_64_crt(uint64(112233445566778899), &Q_FACTORS, &Q_FACTORS_LEN)

	for i:=0 ; i< b.N; i++{
		vectors.CRT_INV()
		
	}
}



type arg_add_crt struct{
	x uint64
	y uint64
	want []uint64
}

var add_crt_vectors = []arg_add_crt{
	{uint64(0), uint64(0), []uint64{0,0,0,0}},
	{uint64(0), uint64(1), []uint64{1,1,1,1}},
	{x_0, x_1, []uint64{2041167892, 1162922470, 1193382793, 2113157951}}, //2041167892
}



func TestADD_crt(t *testing.T){

	var want int_64_crt

	for i, testPair := range add_crt_vectors {
		x := NewInt_64_crt(testPair.x, &Q_FACTORS, &Q_FACTORS_LEN)
		y := NewInt_64_crt(testPair.y, &Q_FACTORS, &Q_FACTORS_LEN)

		x.ADD(x,y)

		want.bigint_64_crt = testPair.want

		if !x.EQUAL(x,&want) {
			t.Errorf("Error ADD_crt test pair %v",i)

		}
	}
}


func BenchmarkADD_crt(b *testing.B){

	x := NewInt_64_crt(x_0, &Q_FACTORS, &Q_FACTORS_LEN)
	y := NewInt_64_crt(x_1, &Q_FACTORS, &Q_FACTORS_LEN)

	for i:=0 ; i< b.N; i++{
		x.ADD(x,y)
		
	}
}


type arg_mulmod struct{
	x uint64
	y uint64
	q uint64
	want uint64
}

var mulmod_vectors = []arg_mulmod{
	{uint64(0), uint64(0),uint64(1152921504606489097),  uint64(0)},
	{uint64(0), uint64(1),uint64(1152921504606489097), uint64(0)},
	{uint64(112233445566778899),uint64(998877665544332211),uint64(1152921504606489097), uint64(1103875254192881828)}, //2041167892
}


func Test_mulmod(t *testing.T){

	for i, testPair := range mulmod_vectors {
		z := mulmod(&testPair.x,&testPair.y,&testPair.q)

		if z != testPair.want{
			t.Errorf("Error mulmod_peasant pair %v",i)
		}
	}

}

func Benchmark_mulmod(b *testing.B){
	//x := uint64(112233445566778899)
	//y := uint64(998877665544332211)
	//q := uint64(1152921504606489097)

	x := uint64(2106880038)
	y := uint64(1479843154)
	q := uint64(2748007003)
	z := uint64(0)

	for i:=0 ; i< b.N; i++{
		z = (x*y)%q //mulmod(&x,&y,&q)

	}

	z += 1
}

//func Test_kmod(t *testing.T){
//
//	for i, testPair := range mulmod_vectors {
//		z := kmod(testPair.x,testPair.y,testPair.q)
//
//		if z != testPair.want{
//			t.Errorf("Error mulmod_peasant pair %v",i)
//		}
//	}
//
//}

func Benchmark_kmod(b *testing.B){
	x := uint64(112233445566778899)
	y := uint64(998877665544332211)
	q := uint64(1152921504606489097)

	for i:=0 ; i< b.N; i++{
		kmod(x,y,q)
	}
}




