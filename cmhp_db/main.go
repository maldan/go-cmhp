package cmhp_db

type Record struct {
	Name   string `json:"name"`
	Offset int    `json:"offset"`
	Size   int    `json:"size"`
}

// 12345678912345678912345678932123

func ReadIndex(path string) {

}

func FindFile(path string, name string) {

}

func Add(path string, name string, content any) {

	//cmhp_file.Write("sas.index", "")
	//cmhp_file.Write("sas.db", "")
}
