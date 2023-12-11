package status

type TableCell struct {
	Status  Status `json:"s"`
	Content string `json:"c"`
}

type Table struct {
	Header []string      `json:"h,omitempty"`
	Rows   [][]TableCell `json:"r"`
}

func (t *Table) DetailType() string {
	return "table"
}

func (t *Table) AppendRow(cells ...TableCell) {
	t.Rows = append(t.Rows, cells)
}
