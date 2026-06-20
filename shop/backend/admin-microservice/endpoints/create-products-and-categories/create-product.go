package create

import (
	"admin-microservice/models"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	product "github.com/AndreyLebedev1998/shop-gRPC-product"
)

func CreateProduct(w http.ResponseWriter, r *http.Request, client product.ProductsServiceClient) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	priceStr := r.FormValue("price")
	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil || price < 0 {
		http.Error(w, "Invalid price", http.StatusBadRequest)
		return
	}

	categoryIdStr := r.FormValue("category_id")
	categoryId, err := strconv.Atoi(categoryIdStr)
	if err != nil {
		http.Error(w, "Invalid category_id", http.StatusBadRequest)
		return
	}

	nameStr := r.FormValue("name")
	if nameStr == "" {
		http.Error(w, "name is not defined", http.StatusBadRequest)
		return
	}

	quantityStr := r.FormValue("availability_of_pieces")
	quantity, err := strconv.Atoi(quantityStr)
	if err != nil || quantity < 0 {
		http.Error(w, "Invalid availability_of_pieces", http.StatusBadRequest)
		return
	}

	subcategoryIdStr := r.FormValue("subcategory_id")
	subcategoryId, err := strconv.Atoi(subcategoryIdStr)
	if err != nil || quantity < 0 {
		http.Error(w, "Invalid subcategory id", http.StatusBadRequest)
		return
	}

	var newProduct models.NewProduct = models.NewProduct{
		ProductName:          nameStr,
		Price:                price,
		CategoryId:           categoryId,
		ImageUrl:             nil,
		AvailabilityOfPieces: quantity,
		SubcategoryId:        subcategoryId,
	}

	var savedFilePath string

	file, header, err := r.FormFile("image_url")
	if err != nil {
		if err != http.ErrMissingFile {
			http.Error(w, "Invalid file upload", http.StatusBadRequest)
			return
		}
	} else {
		defer file.Close()

		uploadDir := "./uploads-products-images"
		if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {

			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}

		ext := filepath.Ext(header.Filename)
		safeName := fmt.Sprintf("%d_%s", time.Now().UnixNano(), rand.Text())
		destPath := filepath.Join(uploadDir, safeName+ext)
		savedFilePath = destPath

		dst, err := os.Create(destPath)
		if err != nil {
			http.Error(w, "Failed to save file", http.StatusInternalServerError)
			return
		}

		defer dst.Close()

		if _, err := io.Copy(dst, file); err != nil {
			http.Error(w, "Failed to save file", http.StatusInternalServerError)
			return
		}

		imageURL := fmt.Sprintf("/uploads-products-images/%s", filepath.Base(destPath))
		newProduct.ImageUrl = &imageURL
	}

	var imageURL string
	if newProduct.ImageUrl != nil {
		imageURL = *newProduct.ImageUrl
	}

	var newProductGrpc = product.NewProduct{
		ProductName:          newProduct.ProductName,
		Price:                newProduct.Price,
		CategoryId:           int64(newProduct.CategoryId),
		ImageUrl:             imageURL,
		AvailabilityOfPieces: int64(newProduct.AvailabilityOfPieces),
		SubcategoryId:        int64(newProduct.SubcategoryId),
	}
	ctx := r.Context()

	resp, err := client.CreateProduct(ctx, &newProductGrpc)
	if err != nil {
		http.Error(w, "server error", http.StatusInternalServerError)
		if savedFilePath != "" {
			_ = os.Remove(savedFilePath)
		}
		return
	}

	w.Header().Set("Content-Type", "application-json")
	json.NewEncoder(w).Encode(resp)
}
