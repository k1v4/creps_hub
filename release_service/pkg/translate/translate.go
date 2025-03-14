package translate

var m = map[string]string{
	"январь":   "January",
	"февраль":  "February",
	"март":     "March",
	"апрель":   "April",
	"май":      "May",
	"июнь":     "June",
	"июль":     "July",
	"август":   "August",
	"сентябрь": "September",
	"октябрь":  "October",
	"ноябрь":   "November",
	"декабрь":  "December",
}

func Translate(rusMonth string) string {
	return m[rusMonth]
}
