import { useEffect, useState, type FC } from "react";
import { getProductsForSubcategories } from "../../api/products-api/products-api";
import "./subcategory-products.css";
import type {
  InitialValuesForFilter,
  ProductIdType,
  Products
} from "../../../types/types";
import Filters from "../Filters/Filters";
import ProductsList from "../ProductsList/ProductsList.tsx";

interface ProductsSubcategoriesType {
  categoryId: number,
  choiceSubcategory: number;
  addProductToCart: (product: Products) => void;
  deleteProductInCart: (productId: ProductIdType) => void;
  initialValuesForFilter: InitialValuesForFilter | null
}

const ProductsSubcategories: FC<ProductsSubcategoriesType> = ({
  categoryId,
  choiceSubcategory,
  addProductToCart,
  deleteProductInCart,
  initialValuesForFilter,
}) => {
  const [products, setProducts] = useState<Products[]>([]);
  const [showFilter, setShowFilter] = useState(false)
  useEffect(() => {
    getProductsForSubcategories(choiceSubcategory).then((data) =>
      setProducts(data),
    );
  }, [choiceSubcategory]);
  return (
    <div className="products-subcategories">
      {showFilter ? <i className="bi bi-filter-circle" onClick={() => setShowFilter(false)}></i> : <i className="bi bi-filter" onClick={() => setShowFilter(true)}></i>}
      {showFilter ? <Filters categoryId={String(categoryId) || ""} subcategoryId={String(choiceSubcategory) || ""} initialValuesForFilter={initialValuesForFilter} setProducts={setProducts}/> : null}
      <ProductsList products={products} deleteProductInCart={deleteProductInCart} addProductToCart={addProductToCart} isSwiper={false} />
    </div>
  );
};

export default ProductsSubcategories;
