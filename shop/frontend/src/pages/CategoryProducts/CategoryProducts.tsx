import { useEffect, useRef, useState, type FC } from "react";
import { useParams } from "react-router-dom";
import { Swiper, SwiperSlide } from "swiper/react";
import type { Swiper as SwiperType } from "swiper";
import { Navigation } from "swiper/modules";
import { getProductsForCategories } from "../../api/products-api/products-api";
import { Button, Card } from "react-bootstrap";
import { config } from "../../config";
import type { NewOrderType, ProductIdType, Products, ProductsInOrders } from "../../../types/types";
import { useTranslation } from "react-i18next";
import "./category-products.css";
import { useDispatch, useSelector } from "react-redux";
import type { RootState } from "../../store/store";
import ProductsSubcategories from "../../components/ProductsSubcategories/ProductsSubcategories";
import { addToCart, initialOrder, removeProductFromCart } from "../../store/cartSlice";

interface CategoryProductsType {
  choiceSubcategory: number,
  setChoceSubcategory: (num: number) => void
}

const CategoryProducts: FC<CategoryProductsType> = ({choiceSubcategory, setChoceSubcategory}) => {
  const { id } = useParams();
  const { t } = useTranslation();
  const [products, setProducts] = useState<Products[]>([]);
  const dispatch = useDispatch()
  const categories = useSelector(
    (state: RootState) => state.categories.categories,
  );
  const necessaryCategory = categories?.find(
    (category) => String(category?.id) === id,
  );
  const subcategories = necessaryCategory?.subcategories || [];
  const swiperRef = useRef<SwiperType | null>(null);
  const user = useSelector((state: RootState) => state.auth.user)
  const cart = useSelector((state: RootState) => state.cart.cart)

  useEffect(() => {
    if (id) {
      getProductsForCategories(id).then((data) => setProducts(data));
    }
  }, [id]);

 useEffect(() => {
  if (choiceSubcategory === 0) {
    swiperRef.current?.slideTo(0);
  }
}, [choiceSubcategory]);

  const addProductToCart = (product: Products) => {
    const newProductInCart: ProductsInOrders = {
      product_id: product.id,
      product_name: product.product_name,
      category_id: product.category_id,
      price: product.price,
      image_url: product.image_url,
      quantity: 1,
      category_name: ""
    }
    const initOrder: NewOrderType = {
      email: user?.email || null,
      phone: user?.phone || null,
      user_id: user?.user_id || null,
      order_items: []
    }
    dispatch(initialOrder(initOrder))
    dispatch(addToCart(newProductInCart))
  }

  const deleteProductInCart = (productId: ProductIdType) => {
     dispatch(removeProductFromCart(productId))
  }
  
  useEffect(() => {
    localStorage.setItem("cart", JSON.stringify(cart))
  }, [cart])

  return (
    <div className="products-categories-wrapper">
      <Swiper
        spaceBetween={50}
        slidesPerView={4}
        className="swiper-product-subcategories"
        onSlideChange={() => console.log("slide change")}
        modules={[Navigation]}
        navigation={true}
        onSwiper={(swiper) => {
  swiperRef.current = swiper;
}}
      >
        {subcategories &&
          subcategories.map((subcategory) => {
            return (
              <SwiperSlide
                key={subcategory.subcategory_id}
                className={subcategory.subcategory_id === choiceSubcategory ? "active-products-swiper" : ""}
              >
                <div onClick={() => setChoceSubcategory(subcategory.subcategory_id)} className="subcategory-name">
                  {t(`subcategories.${subcategory.category_name}`)}
                </div>
              </SwiperSlide>
            );
          })}
      </Swiper>
      {choiceSubcategory === 0 ? <div className="products-categories">
        {products.length > 0 &&
          products.map((product) => {
            const productInOrder = cart?.order_items.find((el) => el.product_id === product.id)
            return (
              <Card key={product.id} className="product-categories">
                <Card.Img
                  className="product-categories__img-product"
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
      </div> : <ProductsSubcategories choiceSubcategory={choiceSubcategory} addProductToCart={addProductToCart} deleteProductInCart={deleteProductInCart}/>}
    </div>
  );
};

export default CategoryProducts;
