package structs

const sockAddrDefault = "/var/run/klystron/klystron.sock"
const sockBufferSize64 = 212992
const sockBufferSize32 = 163840

type fontConfig struct {
	Name    string `json:"name"`
	Style   string `json:"style"`
	Size    int    `json:"size"`
	FontDir string `json:"fontDir"`
}

type pdfFile struct {
	Orientation string     `json:"orientation"`
	Unit        string     `json:"unit"`
	Format      string     `json:"format"`
	Font        fontConfig `json:"font"`
}

type Job struct {
	Name      string       `json:"name"`
	Directory string       `json:"directory"`
	Outputs   []pdfFile    `json:"outputs"`
	Sources   []DataSource `json:"sources"`
}

type DataSource struct {
	Path string `json:"path"`
	Data string `json:"data"`
}
