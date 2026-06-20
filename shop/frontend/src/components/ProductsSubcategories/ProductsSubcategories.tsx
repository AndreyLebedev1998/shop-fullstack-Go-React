import { useEffect, useState, type FC } from "react";
import { getProductsForSubcategories } from "../../api/products-api/products-api";
import { Button, Card } from "react-bootstrap";
import { config } from "../../config";
import { useTranslation } from "react-i18next";
import "./subcategory-products.css";
import type {
  ProductIdType,
  Products,
  SubcategoryProductType,
} from "../../../types/types";
import type { RootState } from "../../store/store";
import { useSelector } from "react-redux";

interface ProductsSubcategoriesType {
  choiceSubcategory: number;
  addProductToCart: (product: Products) => void;
  deleteProductInCart: (productId: ProductIdType) => void;
}

const ProductsSubcategories: FC<ProductsSubcategoriesType> = ({
  choiceSubcategory,
  addProductToCart,
  deleteProductInCart,
}) => {
  const { t } = useTranslation();
  const [products, setProducts] = useState<SubcategoryProductType[]>([]);
  const cart = useSelector((state: RootState) => state.cart.cart);
  useEffect(() => {
    getProductsForSubcategories(choiceSubcategory).then((data) =>
      setProducts(data),
    );
  }, [choiceSubcategory]);
  return (
    <div className="products-subcategories">
      {products &&
        products.length > 0 &&
        products.map((product) => {
          const productInOrder = cart?.order_items.find(
            (el) => el.product_id === product.id,
          );
          return (
            <Card key={product.id} className="product-subcategories">
              <Card.Img
                className="products-subcategories__img-product"
                variant="top"
                src={`${config.PRODUCTS_IMAGES_BASE_URL}/${product.image_url}`}
              />
              <Card.Body>
                <Card.Title>{product.product_name}</Card.Title>
                <Card.Text>{product.price} ₽</Card.Text>
                <div className="product-actions">
                  {productInOrder && productInOrder.quantity ? (
                    <div className="add-to-cart-actions">
                      <button
                        className="cart-btn"
                        onClick={() =>
                          deleteProductInCart({ product_id: product.id })
                        }
                      >
                        −
                      </button>
                      <span className="cart-quantity">
                        {productInOrder.quantity}
                      </span>
                      <button
                        className="cart-btn"
                        onClick={() => addProductToCart(product)}
                      >
                        +
                      </button>
                    </div>
                  ) : (
                    <Button
                      onClick={() => addProductToCart(product)}
                      variant="primary"
                    >
                      {t("orders.add_to_cart")}
                    </Button>
                  )}
                </div>
              </Card.Body>
            </Card>
          );
        })}
    </div>
  );
};

export default ProductsSubcategories;
