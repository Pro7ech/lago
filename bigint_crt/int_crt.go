package bigint_crt

import (
	"./bigint"
	"math/bits"
)



type int_64_crt struct{

	bigint_64_crt []uint64
	q_factors  *[]uint64
	q_factors_len *uint16

}




// NewInt creates a new Int with a given int64 value.
func NewInt_64_crt(v uint64, Q_FACTORS *[]uint64, Q_FACTORS_LEN *uint16) *int_64_crt {

	a := make([]uint64,*Q_FACTORS_LEN)

	tmp := &int_64_crt{ a , Q_FACTORS , Q_FACTORS_LEN}

	for i, qi := range *Q_FACTORS{

			tmp.bigint_64_crt[i] = (v%qi)
	}

	return tmp
}

// Creates a new crt representation of a bigint integer
func NewInt_big_crt (v *bigint.Int, Q_FACTORS *[]uint64, Q_FACTORS_LEN *uint16) *int_64_crt {

	a := make([]uint64,*Q_FACTORS_LEN)

	tmp := &int_64_crt{ a , Q_FACTORS , Q_FACTORS_LEN}

	for i, qi := range *Q_FACTORS{

			var tmp_qi *bigint.Int

			tmp_qi.Value.SetInt64(int64(qi))

			v.Mod(v,tmp_qi)

			tmp.bigint_64_crt[i] = uint64(v.Int64())
	}

	return tmp
}


func (* int_64_crt) ADD(a,b *int_64_crt) *int_64_crt{

	for i, q := range *b.q_factors{

			a.bigint_64_crt[i] += b.bigint_64_crt[i]

			if a.bigint_64_crt[i]>q{

				a.bigint_64_crt[i] -= q
			}
	}

	return a
}




func (* int_64_crt) SUB(a,b *int_64_crt) *int_64_crt{

	for i, q := range *b.q_factors{

			a.bigint_64_crt[i] += q
			a.bigint_64_crt[i] -= b.bigint_64_crt[i]

			for a.bigint_64_crt[i]>q{a.bigint_64_crt[i] -= q}
	}

	return a
}

func (* int_64_crt) MUL_32(a,b *int_64_crt) *int_64_crt{

	for i, q := range *b.q_factors{

			a.bigint_64_crt[i] = (a.bigint_64_crt[i]*b.bigint_64_crt[i])%q
	}

	return a
}

func (* int_64_crt) MUL_64(a,b *int_64_crt) *int_64_crt{

	for i, q := range *b.q_factors{

			a.bigint_64_crt[i] = mulmod2(&a.bigint_64_crt[i],&b.bigint_64_crt[i],&q)
	}

	return a
}




func (* int_64_crt) EQUAL(a, b *int_64_crt) bool {

	x := a.bigint_64_crt
	y := b.bigint_64_crt

    if len(x) != len(y) {
        return false
    }
    for i, v := range x {
        if v != y[i] {
            return false
        }
    }
    return true
}


func mulmod_32(a,b,q uint64) uint64 {
	return (a*b)%q
	
}



func mulmod1(A,B,Q *uint64) uint64{

	a,b,q := *A,*B,*Q

	if (a>=q) { a %= q}
	if (b>=q) { b %= q}
	if (bits.LeadingZeros64(a)+bits.LeadingZeros64(b)) > 64 {return (a*b)%q}
	if (a<b) { a,b = b,a}


	result := uint64(0)

	for b>0{
		if b&1 == 1{
			result += a
			if result>q {result -=q}
		}

		a <<= 1
		for a>q{ a -= q}
		b >>= 1
	}

	return result
}




func mulmod2(A,B,Q *uint64) uint64{
	a,b,q := *A,*B,*Q

	if (a>=q) { a %= q}
	if (b>=q) { b %= q}
	if (bits.LeadingZeros64(a)+bits.LeadingZeros64(b)) > 64 {return (a*b)%q}

	a0 := a>>32
	a1 := a & 0xFFFFFFFF
	b0 := b>>32
	b1 := b & 0xFFFFFFFF

	x0 := (a0*b0)
	x1 := (a1*b0) + (a0*b1)
	x2 := (a1*b1)

	for (x0>=q) {x0 -= q}
	for (x1>=q) {x1 -= q}
	for (x2>=q) {x2 -= q}

	for i:=0 ; i<32 ; i++{
		x0 <<= 2
		x1 <<= 1

		for (x0>=q) {x0 -= q}
		if (x1>=q) {x1 -= q}


	}


	return (x0 + x1 + x2)%q

}

func mulmod3(A,B,Q *uint64, N,mask uint64) uint64{

	a,b,q := *A,*B,*Q

    if (a >= q) {a %= q}
    if (b >= q) {b %= q}
    if (bits.LeadingZeros64(a)+bits.LeadingZeros64(b)) > 64 {return (a*b)%q}
    if (a<b) {a,b = b,a}

    result := uint64(0)
    
    for (a>0 && b>0){
        result = (result + (b&mask) * a) %q
        b>>=N
        a = (a<<N)%q
        
    }
    
    return result    
    
}



// Inverse CRT mapping. Takes a crt_representation with 64 bits vectors and return the bigInt recomposition
func (this int_64_crt) CRT_INV(N *bigint.Int, CRT_PARAMS *[]bigint.Int) *bigint.Int{

	//First we need to convert all elements of the crt representation from 64bits to 
	//bit int
	var C bigint.Int

	result  := bigint.NewInt(0)

	PARAMS := *CRT_PARAMS


	for i, _ := range *this.q_factors{

		C.SetInt(int64(this.bigint_64_crt[i]))

		C.Mul(&C,&PARAMS[i])

		result.Add(result,&C)
		
	}

	return result.Mod(result,N)


}




