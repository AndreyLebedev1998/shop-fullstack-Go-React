package create

// Работа с обычной формой

/* func CreateProduct(w http.ResponseWriter, r *http.Request, sql *sql.DB) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Parse error", http.StatusBadRequest)
		return
	}

	priceStr := r.FormValue("price")

	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		http.Error(w, "Invalid price", http.StatusBadRequest)
		return
	}

	product := models.NewProduct{
		ProductName: r.FormValue("name"),
		Price:       price,
		CategoryId:  r.FormValue("category_id"),
		ImageUrl:    r.FormValue("img_url"),
	}

	fmt.Printf("%+v\n", product)

	w.Write([]byte("OK"))
} */
