LAGO : Lattice Cryptography Library in GOlang
===========

This package provides a toolbox of lattice-based cryptographic primitives and Fan-Vercauteren (FV) homomorphic encryption for GO.
It's currently only used for research purposes.

Examples
---------

[main.go](https://github.com/dedis/student_18_lattices/blob/master/main.go) gives an example about how to use this package.
In each sub-package, you can also find testing go files, which also show the usage.

Notes
-------

[kyber](https://github.com/dedis/student_18_lattices/tree/master/kyber) is cloned from [Yawning/kyber](https://github.com/Yawning/kyber),
which is used to cross verify NTT transformation in [ntt_test.go](https://github.com/dedis/student_18_lattices/blob/master/polynomial/ntt_test.go)