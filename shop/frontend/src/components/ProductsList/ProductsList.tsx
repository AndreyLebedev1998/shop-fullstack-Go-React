import { useDispatch, useSelector } from "react-redux";
import type { RootState } from "../../store/store";
import type { ProductIdType, Products } from "../../../types/types";
import { useRef, type FC } from "react";
import { Button, Card } from "react-bootstrap";
import { config } from "../../config";
import { useTranslation } from "react-i18next";
import type { Swiper as SwiperType } from "swiper";
import { Swiper, SwiperSlide } from "swiper/react";
import "./product-list.css";
import { Navigation } from "swiper/modules";
import {
  addFavoriteProductForUser,
  removeFavoriteProduct,
} from "../../api/products-api/products-api";
import {
  addFavoriteProduct,
  deleteFavoriteProduct,
} from "../../store/productsSlice";

interface TProducts {
  products: Products[];
  deleteProductInCart: (productId: ProductIdType) => void;
  addProductToCart: (products: Products) => void;
  isSwiper: boolean;
}

const ProductsList: FC<TProducts> = ({
  products,
  deleteProductInCart,
  addProductToCart,
  isSwiper,
}) => {
  const { t } = useTranslation();
  const dispatch = useDispatch();
  const cart = useSelector((state: RootState) => state.cart.cart);
  const swiperRef = useRef<SwiperType | null>(null);
  const token = useSelector((state: RootState) => state.auth.token);
  const favoriteProducts = useSelector(
    (state: RootState) => state.products.favoriteProduct,
  );

  const handleAddFavoriteProduct = (productId: number) => {
    if (token) {
      addFavoriteProductForUser(productId, token).then((data) => {
        if (data) {
          dispatch(addFavoriteProduct(data));
        }
      });
    }
  };

  const handleRemoveFavoriteProduct = (productId: number) => {
    if (token) {
      removeFavoriteProduct(token, productId).then((data) => {
        if (data) {
          dispatch(deleteFavoriteProduct({ product_id: data.product_id }));
        }
      });
    }
  };

  return (
    <div className="products-list-wrapper">
      {isSwiper ? (
        <Swiper
          spaceBetween={16}
          slidesPerView="auto"
          className="swiper-product"
          modules={[Navigation]}
          navigation={true}
          onSwiper={(swiper) => {
            swiperRef.current = swiper;
          }}
        >
          {products.map((product) => {
            const productInOrder = cart?.order_items.find(
              (el) => el.product_id === product.id,
            );
            const isFavoriteProduct = favoriteProducts.find(
              (el) => el.id === product.id,
            );
            const productInCart = cart?.order_items.find(
              (el) => el.product_id === product.id,
            );
            const isDisabled =
              product.availability_of_pieces === productInCart?.quantity;
            const isOutOfStock = product.availability_of_pieces === 0;
            return (
              <SwiperSlide key={product.id} className="swiper-product-slide">
                <Card className="product-categories">
                  {isOutOfStock && (
                    <div className="out-of-stock-badge">
                      {t("product_page.not_available")}
                    </div>
                  )}
                  {token ? (
                    <div className="product-categories__heart-wrapper">
                      <i
                        className={`bi ${
                          isFavoriteProduct
                            ? "bi-heart-fill favorite-product active-favorite-product"
                            : "bi-heart favorite-product"
                        }`}
                        onClick={() =>
                          isFavoriteProduct
                            ? handleRemoveFavoriteProduct(product.id)
                            : handleAddFavoriteProduct(product.id)
                        }
                      />
                    </div>
                  ) : null}
                  <Card.Img
                    className="product-categories__img-product"
                    variant="top"
                    src={`${config.PRODUCTS_IMAGES_BASE_URL}/${product.image_url}`}
                  />
                  <Card.Body>
                    <Card.Title>{product.product_name}</Card.Title>
                    <Card.Text>{product.price} ₽</Card.Text>
                    <div className="product-actions">
                      {isOutOfStock ? null : productInOrder &&
                        productInOrder.quantity ? (
                        <div className="add-to-cart-actions">
                          <button
                            className="cart-btn"
                            onClick={() =>
                              deleteProductInCart({ product_id: product.id })
                            }
                          >
                            <span className="minus">-</span>
                          </button>
                          <span className="cart-quantity">
                            {productInOrder.quantity}
                          </span>
                          <button
                            className="cart-btn"
                            onClick={() => addProductToCart(product)}
                            disabled={isDisabled}
                          >
                            <span className="plus">+</span>
                          </button>
                        </div>
                      ) : (
                        <Button
                          onClick={() => addProductToCart(product)}
                          variant="primary"
                          disabled={isDisabled}
                        >
                          {t("orders.add_to_cart")}
                        </Button>
                      )}
                    </div>
                  </Card.Body>
                </Card>
              </SwiperSlide>
            );
          })}
        </Swiper>
      ) : (
        products.length > 0 &&
        products.map((product) => {
          const productInOrder = cart?.order_items.find(
            (el) => el.product_id === product.id,
          );
          const isFavoriteProduct = favoriteProducts.find(
            (el) => el.id === product.id,
          );
          const productInCart = cart?.order_items.find(
            (el) => el.product_id === product.id,
          );
          const isDisabled =
            product.availability_of_pieces === productInCart?.quantity;
          const isOutOfStock = product.availability_of_pieces === 0;
          return (
            <Card key={product.id} className="product-categories">
              {isOutOfStock && (
                <div className="out-of-stock-badge">
                  {t("product_page.not_available")}
                </div>
              )}

              {token ? (
                <div className="product-categories__heart-wrapper">
                  <i
                    className={`bi ${
                      isFavoriteProduct
                        ? "bi-heart-fill favorite-product active-favorite-product"
                        : "bi-heart favorite-product"
                    }`}
                    onClick={() =>
                      isFavoriteProduct
                        ? handleRemoveFavoriteProduct(product.id)
                        : handleAddFavoriteProduct(product.id)
                    }
                  />
                </div>
              ) : null}
              <Card.Img
                className="product-categories__img-product"
                variant="top"
                src={`${config.PRODUCTS_IMAGES_BASE_URL}/${product.image_url}`}
              />
              <Card.Body>
                <Card.Title>{product.product_name}</Card.Title>
                <Card.Text>{product.price} ₽</Card.Text>
                <div className="product-actions">
                  {isOutOfStock ? null : productInOrder &&
                    productInOrder.quantity ? (
                    <div className="add-to-cart-actions">
                      <button
                        className="cart-btn"
                        onClick={() =>
                          deleteProductInCart({ product_id: product.id })
                        }
                      >
                        <span className="minus">-</span>
                      </button>
                      <span className="cart-quantity">
                        {productInOrder.quantity}
                      </span>
                      <button
                        className="cart-btn"
                        onClick={() => addProductToCart(product)}
                        disabled={isDisabled}
                      >
                        <span className="plus">+</span>
                      </button>
                    </div>
                  ) : (
                    <Button
                      onClick={() => addProductToCart(product)}
                      variant="primary"
                      disabled={isDisabled}
                    >
                      {t("orders.add_to_cart")}
                    </Button>
                  )}
                </div>
              </Card.Body>
            </Card>
          );
        })
      )}
    </div>
  );
};

export default ProductsList;
