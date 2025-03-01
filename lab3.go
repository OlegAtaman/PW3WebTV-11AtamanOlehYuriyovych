package main

import (
	"fmt"
	"html/template"
	"log"
	"math"
	"net/http"
	"strconv"
)

// Шаблон HTML-сторінки
const tpl = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Калькулятор</title>
    <link rel="stylesheet" type="text/css" href="/static/style.css">
</head>
<body>
    <div class="container">
        <h2>Веб-калькулятор</h2>
        
        <div class="calc-container">
            <form method="post">
                <label>Pс, МВт:</label>
                <input type="text" name="pc" required><br>
                <label>sigma1, МВт:</label>
                <input type="text" name="sigma1" required><br>
                <label>Вартість, грн/кВт*год:</label>
                <input type="text" name="cost" required><br>
                <button type="submit">Обчислити</button>
            </form>
            {{if .Result}}
            <pre class="result">{{.Result}}</pre>
            {{end}}
        </div>
    </div>
</body>
</html>
`

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		r.ParseForm()
		data := make(map[string]interface{})
		data["Result"] = calculateProfit(r)

		tmpl, _ := template.New("calc").Parse(tpl)
		tmpl.Execute(w, data)
		return
	}
	tmpl, _ := template.New("calc").Parse(tpl)
	tmpl.Execute(w, nil)
}

func calculateProfit(r *http.Request) string {
	pc, _ := strconv.ParseFloat(r.FormValue("pc"), 64)
	sigma1, _ := strconv.ParseFloat(r.FormValue("sigma1"), 64)
	cost, _ := strconv.ParseFloat(r.FormValue("cost"), 64)

	p := 5.0
	p1 := (1 / (sigma1 * math.Sqrt(2*math.Pi))) * math.Exp(-math.Pow(p-pc, 2)/(2*math.Pow(sigma1, 2)))

	sigmaW1 := integrate(func(x float64) float64 {
		return (1 / (sigma1 * math.Sqrt(2*math.Pi))) * math.Exp(-math.Pow(x-pc, 2)/(2*math.Pow(sigma1, 2)))
	}, 4.75, 5.25, 1000)

	W1 := pc * 24 * sigmaW1
	P := W1 * cost
	Sh := pc * 24 * (1 - sigmaW1) * cost
	NetProfit := (P - Sh) * 1000

	return fmt.Sprintf("p1: %.6f\nsigmaW1: %.6f\nW1: %.6f\nП: %.6f\nШ: %.6f\nЧистий прибуток: %.6f", p1, sigmaW1, W1, P, Sh, NetProfit)
}

func integrate(f func(float64) float64, a, b float64, n int) float64 {
	dx := (b - a) / float64(n)
	sum := 0.0
	for i := 0; i < n; i++ {
		x := a + float64(i)*dx
		sum += f(x)
	}
	return sum * dx
}

func main() {
	http.HandleFunc("/", indexHandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	fmt.Println("Сервер запущено на http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
