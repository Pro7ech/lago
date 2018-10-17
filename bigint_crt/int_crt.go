package bigint_crt

import (
	"./bigint"
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

			if a.bigint_64_crt[i]>q{

				a.bigint_64_crt[i] -= q
			}
	}

	return a
}

func (* int_64_crt) MUL(a,b *int_64_crt) *int_64_crt{

	for i, q := range *b.q_factors{

			a.bigint_64_crt[i] = mulmod(&a.bigint_64_crt[i],&b.bigint_64_crt[i],&q)
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


func mulmod(A,B,Q *uint64) uint64{

	a,b,q := *A,*B,*Q

	if (a|b)>>32 == 0 {
		return (a*b)%q
	}


	result := uint64(0)

	

	for b>0{
		if b&1 == 1{
			result += a
			if result>q{
				result -=q
			}
		}

		a <<= 1
		for a>q{
			a -= q
		}
		b >>= 1
	}

	return result
}

func kmod(a,b,q uint64) uint64{

	if (a|b)>>32 == 0 {
		return (a*b)%q
	}

	a0 := a>>32
	a1 := a & 0xFFFFFFFF
	b0 := b>>32
	b1 := b & 0xFFFFFFFF

	x0 := (a0*b0)%q
	x1 := (a1*b0)%q
	x2 := (a0*b1)%q
	x3 := (a1*b1)%q

	return (x0 + x1 + x2 + x3)%q

}


// Inverse CRT mapping. Takes a crt_representation with 64 bits vectors and return the bigInt recomposition
func (this int_64_crt) CRT_INV() *bigint.Int{

	//First we need to convert all elements of the crt representation from 64bits to 
	//bit int

	var qi bigint.Int
	var C bigint.Int
	var tmp bigint.Int
	var tmp_inv bigint.Int

	var tmp_q bigint.Int


	f       := bigint.NewInt(1)
	result  := bigint.NewInt(0)


	for _, q := range *this.q_factors{

		tmp_q.SetInt(int64(q))		
		f.Mul(f,&tmp_q)

	}

	for i, q := range *this.q_factors{

		//qi_minus2 := NewInt(q-2) compute inverse with a^(q-2) % q

		qi.SetInt(int64(q))
		C.SetInt(int64(this.bigint_64_crt[i]))

		tmp.Div(f,&qi)

		tmp_inv.Inv(&tmp,&qi)

		C.Mul(&C,&tmp)
		C.Mul(&C,&tmp_inv)

		result.Add(result,&C)
		
	}

	return result.Mod(result,f)


}




