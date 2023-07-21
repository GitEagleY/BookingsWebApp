package models

type TemplateData struct {
	StringMap map[string]string      //testing string map
	IntMap    map[string]int         //testing int map
	FloatMap  map[string]float32     //testing float map
	DataMap   map[string]interface{} //testing data map
	CSRFtoken string                 //cross site request forgery token
	Flash     string
	Warnaing  string
	Error     error
}
