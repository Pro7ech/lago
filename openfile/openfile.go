package openfile

import(
	"os"
	"bufio")

func OpenFile(filename string)(vs []string){

	testfile, _ := os.Open(filename)

	defer testfile.Close()

	scanner := bufio.NewScanner(testfile)

	for scanner.Scan(){
		vs = append(vs,scanner.Text())
	}

	return vs

}