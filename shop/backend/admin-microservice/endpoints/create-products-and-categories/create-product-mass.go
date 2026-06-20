package create

import (
	"admin-microservice/models"
	"crypto/rand"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	productGRPC "github.com/AndreyLebedev1998/shop-gRPC-product"
	"github.com/xuri/excelize/v2"
)

func ReadXLSXAndSaveImages(w http.ResponseWriter, r *http.Request, db *sql.DB, clientProduct productGRPC.ProductsServiceClient) {
	err := r.ParseMultipartForm(20 << 20)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	ctx := r.Context()

	xlsx, err := excelize.OpenReader(file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Println(xlsx.GetSheetList()[0])
	sheet := xlsx.GetSheetList()

	if len(sheet) == 0 {
		http.Error(w, "xlsx has no sheets", http.StatusBadRequest)
		return
	}

	rows, err := xlsx.GetRows(sheet[0])
	if err != nil {
		http.Error(w, "cannot read rows", http.StatusInternalServerError)
		return
	}

	var productReturns []*productGRPC.ReturnNewProduct

	for i, row := range rows {
		if i == 0 {
			continue
		}
		var numStartCell = i + 1
		var pictureColumn = "D"
		var coordinate string = fmt.Sprintf("%s%s", pictureColumn, strconv.Itoa(numStartCell))

		pictures, err := xlsx.GetPictures(sheet[0], coordinate)
		if err != nil {
			fmt.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var product models.NewProduct
		var product_name = row[0]
		var category_id = row[1]
		var price = row[2]
		var availability_of_pieces = row[4]
		var subcategory_id = row[5]

		categoryIdInt, err := strconv.Atoi(category_id)
		if err != nil {
			http.Error(w, "Invalid category id", http.StatusBadRequest)
			return
		}

		priceFloat, err := strconv.ParseFloat(price, 64)
		if err != nil {
			http.Error(w, "Invalid price", http.StatusBadRequest)
			return
		}

		quantity, err := strconv.Atoi(availability_of_pieces)
		if err != nil {
			http.Error(w, "Invalid availability_of_pieces", http.StatusBadRequest)
			return
		}

		subcategoryId, err := strconv.Atoi(subcategory_id)
		if err != nil {
			http.Error(w, "Invalid availability_of_pieces", http.StatusBadRequest)
			return
		}

		product.ProductName = product_name
		product.CategoryId = categoryIdInt
		product.Price = priceFloat
		product.ImageUrl = nil
		product.AvailabilityOfPieces = quantity
		product.SubcategoryId = subcategoryId

		uploadDir := "./uploads-products-images"
		safeName := fmt.Sprintf("%d_%s", time.Now().UnixNano(), rand.Text())

		var newProductGrpc = productGRPC.NewProduct{
			ProductName:          product.ProductName,
			Price:                product.Price,
			CategoryId:           int64(product.CategoryId),
			ImageUrl:             "",
			AvailabilityOfPieces: int64(product.AvailabilityOfPieces),
			SubcategoryId:        int64(product.SubcategoryId),
		}

		if len(pictures) == 0 {
			fmt.Println("no pictures found")
			resp, err := clientProduct.CreateProduct(ctx, &newProductGrpc)
			if err != nil {
				http.Error(w, "Error insert in to products", http.StatusInternalServerError)
				return
			}
			if resp != nil {
				var returnProduct = productGRPC.ReturnNewProduct{
					Id:                   resp.Id,
					ProductName:          resp.ProductName,
					Price:                resp.Price,
					CategoryId:           int64(resp.CategoryId),
					ImageUrl:             "",
					AvailabilityOfPieces: int64(resp.AvailabilityOfPieces),
					SubcategoryId:        int64(resp.SubcategoryId),
				}
				productReturns = append(productReturns, &returnProduct)
			}
			continue
		}

		err = os.MkdirAll(uploadDir, 0755)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		pic := pictures[0]
		var pathPicture = fmt.Sprintf("%s.%s", safeName, pic.Extension)
		filePath := filepath.Join(
			uploadDir,
			pathPicture,
		)

		resp, err := clientProduct.GetProductName(ctx, &productGRPC.ProductName{
			ProductName: product.ProductName,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if !resp.Result {
			err = os.WriteFile(filePath, pic.File, 0644)
			if err != nil {
				http.Error(w, "WRITE ERROR: "+err.Error(), http.StatusInternalServerError)
				return
			}
		}

		newProductGrpc.ImageUrl = pathPicture
		respCreate, err := clientProduct.CreateProduct(ctx, &newProductGrpc)
		if err != nil {
			if strings.Contains(err.Error(), "23505") || strings.Contains(err.Error(), "unique") {
				fmt.Printf("Product %s already exists, skipping\n", product.ProductName)
				continue
			}
			if !resp.Result {
				os.Remove(filePath)
			}
			http.Error(w, "Error insert in to products width image", http.StatusInternalServerError)
			return
		}
		if respCreate != nil {
			var returnProduct = productGRPC.ReturnNewProduct{
				Id:                   respCreate.Id,
				ProductName:          respCreate.ProductName,
				Price:                respCreate.Price,
				CategoryId:           int64(respCreate.CategoryId),
				ImageUrl:             respCreate.ImageUrl,
				AvailabilityOfPieces: int64(respCreate.AvailabilityOfPieces),
				SubcategoryId:        int64(respCreate.SubcategoryId),
			}
			productReturns = append(productReturns, &returnProduct)
		}
	}

	w.Header().Set("Content-Type", "application-json")
	json.NewEncoder(w).Encode(productReturns)
}
