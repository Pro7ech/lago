package bigint_crt

//go test ./bigint_crt -v
//go test -bench=. ./bigint_crt -v

import (
	"testing"
	//"math"
	"lago/bigint"
	"math/bits"
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

var CRT_PARAMS = make([]bigint.Int, Q_FACTORS_LEN)
var N = bigint.NewInt(1)

//Precomputes the CRT_PARAMS used in CRT_INV
func init(){

	var qi bigint.Int
	var Nqi bigint.Int
	var Nqi_INV bigint.Int

	//CRT_PARAMS := make([]bigint.Int, Q_FACTORS_LEN)
	//N := bigint.NewInt(1)

	for i:= 0 ; i<int(Q_FACTORS_LEN) ; i++{
		CRT_PARAMS[i] = *bigint.NewInt(1)
	}

	//Computs the product N = q0*q1...qi
	for _, q := range Q_FACTORS{

		qi.SetInt(int64(q))		
		N.Mul(N,&qi)
	}
	//CRT_PARAMS = {(N/qi) * (1/(N/qi)) mod qi, ... }

	for i, q := range Q_FACTORS{

		qi.SetInt(int64(q))

		Nqi.Div(N,&qi)

		Nqi_INV.Inv(&Nqi,&qi)

		CRT_PARAMS[i].Mul(&CRT_PARAMS[i],Nqi.Mul(&Nqi,&Nqi_INV))

	}

	
}




type arg_create_crt_64 struct{
	v uint64
	want []uint64
}

var create_crt_vectors_64 = []arg_create_crt_64{
	{uint64(0), []uint64{0,0,0,0}},
	{uint64(1), []uint64{1,1,1,1}},
	{x_0, []uint64{x_0%q_0,x_0%q_1,x_0%q_2,x_0%q_3}},
	{x_1, []uint64{x_1%q_0,x_1%q_1,x_1%q_2,x_1%q_3}},
}


func TestCreate_crt_64(t *testing.T){

	var z int_64_crt
	var y int_64_crt

	for i, testPair := range create_crt_vectors_64 {

		y.bigint_64_crt = testPair.want

		if !z.EQUAL(NewInt_64_crt(testPair.v, &Q_FACTORS, &Q_FACTORS_LEN), &y){
			t.Errorf("Error creating crt vectors pair %v",i)
		}
	}
}

func BenchmarkCreate_crt_64(b *testing.B){
	v := uint64(112233445566778899)
	for i:=0 ; i< b.N; i++{
		NewInt_64_crt(v, &Q_FACTORS, &Q_FACTORS_LEN)
	}
}

func TestRecombine_crt_64(t *testing.T){

	for i, testPair := range create_crt_vectors_64 {

		outputTest := NewInt_64_crt(testPair.v, &Q_FACTORS, &Q_FACTORS_LEN).CRT_INV(N,&CRT_PARAMS)

		var expectedResult = bigint.NewInt(int64(testPair.v))
		
		if !outputTest.EqualTo(expectedResult){
			t.Errorf("Error crt recombine pair %v", i)
		}
	}

}

func BenchmarkRecombine_crt_64(b *testing.B){

	vectors := NewInt_64_crt(uint64(112233445566778899), &Q_FACTORS, &Q_FACTORS_LEN)

	for i:=0 ; i< b.N; i++{
		vectors.CRT_INV(N,&CRT_PARAMS)
		
	}
}



type arg_create_crt_bigInt struct{
	v *bigint.Int
	want []uint64
}

var create_crt_vectors_bigInt = []arg_create_crt_bigInt{
	{bigint.NewInt(0), []uint64{0,0,0,0}},
	{bigint.NewInt(1), []uint64{1,1,1,1}},
	{bigint.NewInt(6492755906530261339), []uint64{980274044, 1870478686, 188788141, 1203567290}},
	{bigint.NewIntFromString("106163969574508232672974664942102088094"), []uint64{1408544492, 445903048, 2008607093, 1081950402}},
}


func TestCreate_crt_bigInt(t *testing.T){

	var y int_64_crt

	for i, testPair := range create_crt_vectors_bigInt {

		y.bigint_64_crt = testPair.want

		z := NewInt_big_crt(testPair.v, &Q_FACTORS, &Q_FACTORS_LEN)

		if !z.EQUAL(z, &y){
			t.Errorf("Error creating crt vectors pair %v",i)
		}
	}
}



func BenchmarkCreate_crt_bigInt(b *testing.B){
	v := bigint.NewIntFromString("106163969574508232672974664942102088094")
	for i:=0 ; i< b.N; i++{
		NewInt_big_crt(v, &Q_FACTORS, &Q_FACTORS_LEN)
	}
}


func TestRecombine_crt_bigInt(t *testing.T){

	for i, testPair := range create_crt_vectors_bigInt {

		outputTest := NewInt_big_crt(testPair.v, &Q_FACTORS, &Q_FACTORS_LEN).CRT_INV(N,&CRT_PARAMS)

		var expectedResult = testPair.v.Mod(testPair.v,N)
		
		if !outputTest.EqualTo(expectedResult){
			t.Errorf("Error crt recombine pair %v",i)
		}
	}

}


func BenchmarkRecombine_crt_bigInt(b *testing.B){

	vectors := NewInt_big_crt(bigint.NewIntFromString("106163969574508232672974664942102088094"), &Q_FACTORS, &Q_FACTORS_LEN)

	for i:=0 ; i< b.N; i++{
		vectors.CRT_INV(N,&CRT_PARAMS)
		
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



func TestADD_64_crt(t *testing.T){

	var want int_64_crt

	for i, testPair := range add_crt_vectors {
		x := NewInt_64_crt(testPair.x, &Q_FACTORS, &Q_FACTORS_LEN)
		y := NewInt_64_crt(testPair.y, &Q_FACTORS, &Q_FACTORS_LEN)

		x.ADD_64(x,y)

		want.bigint_64_crt = testPair.want

		if !x.EQUAL(x,&want) {
			t.Errorf("Error ADD_64_crt test pair %v",i)

		}
	}
}


func BenchmarkADD_64_crt(b *testing.B){

	x := NewInt_64_crt(x_0, &Q_FACTORS, &Q_FACTORS_LEN)
	y := NewInt_64_crt(x_1, &Q_FACTORS, &Q_FACTORS_LEN)

	for i:=0 ; i< b.N; i++{
		x.ADD_64(x,y)
		
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



func Benchmark_mulmod_32(b* testing.B) {

	x := uint64(2106880038)
	y := uint64(1479843154)
	q := uint64(2748007003)


	for i:=0 ; i< b.N; i++{
		
		mulmod_32(x,y,q)
	}

}

func Benchmark_mulmod_bigint(b* testing.B) {

	x := bigint.NewInt(112233445566778899)
	y := bigint.NewInt(998877665544332211)
	q := bigint.NewInt(1152921504606489097)

	for i:=0 ; i< b.N; i++ {
		mulmod_bigint(x,y,q)
	}


}

func Test_mulmod1_64(t *testing.T){

	for i, testPair := range mulmod_vectors {
		z := mulmod1(&testPair.x,&testPair.y,&testPair.q)

		if z != testPair.want{
			t.Errorf("Error mulmod1_64 pair %v",i)
		}
	}

}
	

func Benchmark_mulmod1_64(b *testing.B){
	//x := uint64(112233445566778899)
	//y := uint64(998877665544332211)
	//q := uint64(1152921504606489097)

	x := uint64(112233445566778899)
	y := uint64(998877665544332211)
	q := uint64(1152921504606489097)

	for i:=0 ; i< b.N; i++{
		mulmod1(&x,&y,&q)

	}
}



func Test_mulmod2_64(t *testing.T){

	for i, testPair := range mulmod_vectors {
		z := mulmod2(&testPair.x,&testPair.y,&testPair.q)

		if z != testPair.want{
			t.Errorf("Error mulmod1_64 pair %v",i)
		}
	}

}

func Benchmark_mulmod2_64(b *testing.B){
	x := uint64(112233445566778899)
	y := uint64(998877665544332211)
	q := uint64(1152921504606489097)

	for i:=0 ; i< b.N; i++{
		mulmod2(&x,&y,&q)
	}
}


func Test_mulmod3_64(t *testing.T){

	for i, testPair := range mulmod_vectors {
		N := uint64(bits.LeadingZeros64(testPair.q))
    	mask := uint64((1<<N) - 1)

		z := mulmod3(&testPair.x,&testPair.y,&testPair.q,N,mask)

		if z != testPair.want{
			t.Errorf("Error mulmod1_64 pair %v",i)
		}
	}

}

func Benchmark_mulmod3_64(b *testing.B) {
	x := uint64(112233445566778899)
	y := uint64(998877665544332211)
	q := uint64(1152921504606489097)
	N := uint64(bits.LeadingZeros64(q))
    mask := uint64((1<<N) - 1)

    for i:=0 ; i< b.N; i++{
		mulmod3(&x,&y,&q,N,mask)
	}
	
}



type arg_mul_32_crt struct{
	x uint64
	y uint64
	want []uint64
}

var MUL_32_crt_vectors = []arg_mul_32_crt{
	{uint64(0), uint64(0), []uint64{0,0,0,0}},
	{uint64(0), uint64(1), []uint64{0,0,0,0}},
	{x_0, x_1, []uint64{537387374, 826233593, 692217772, 1742695417}}, //2041167892
}


func TestMUL_crt_32(t *testing.T){

	var want int_64_crt

	for i, testPair := range MUL_32_crt_vectors {
		x := NewInt_64_crt(testPair.x, &Q_FACTORS, &Q_FACTORS_LEN)
		y := NewInt_64_crt(testPair.y, &Q_FACTORS, &Q_FACTORS_LEN)

		x.MUL_32(x,y)

		want.bigint_64_crt = testPair.want

		if !x.EQUAL(x,&want) {
			t.Errorf("Error MUL_crt_32 test pair %v",i)

		}
	}
}


func BenchmarkMUL_crt_32(b *testing.B){

	x := NewInt_64_crt(x_0, &Q_FACTORS, &Q_FACTORS_LEN)
	y := NewInt_64_crt(x_1, &Q_FACTORS, &Q_FACTORS_LEN)

	for i:=0 ; i< b.N; i++{
		x.MUL_32(x,y)
		
	}
}