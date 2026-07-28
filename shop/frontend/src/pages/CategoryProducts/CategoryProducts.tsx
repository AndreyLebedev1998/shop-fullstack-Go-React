import { useEffect, useRef, useState, type FC } from "react";
import { useParams } from "react-router-dom";
import { Swiper, SwiperSlide } from "swiper/react";
import type { Swiper as SwiperType } from "swiper";
import { Navigation } from "swiper/modules";
import {
  getProductsForCategories,
  findProduct,
} from "../../api/products-api/products-api";
import { Accordion, Form, InputGroup } from "react-bootstrap";
import type {
  InitialValuesForFilter,
  NewOrderType,
  ProductIdType,
  Products,
  ProductsInOrders,
} from "../../../types/types";
import { useTranslation } from "react-i18next";
import "./category-products.css";
import { useDispatch, useSelector } from "react-redux";
import type { RootState } from "../../store/store";
import ProductsSubcategories from "../../components/ProductsSubcategories/ProductsSubcategories";
import {
  addToCart,
  initialOrder,
  removeProductFromCart,
} from "../../store/cartSlice";
import Filters from "../../components/Filters/Filters";
import { getInitialValuesForFilter } from "../../api/filters-api/filters-api";
import ProductsList from "../../components/ProductsList/ProductsList";
import SortProducts from "../SortProducts/SortProducts";

interface CategoryProductsType {
  choiceSubcategory: number;
  setChoceSubcategory: (num: number) => void;
}

const CategoryProducts: FC<CategoryProductsType> = ({
  choiceSubcategory,
  setChoceSubcategory,
}) => {
  const { id } = useParams();
  const { t } = useTranslation();
  const [products, setProducts] = useState<Products[]>([]);
  const dispatch = useDispatch();
  const categories = useSelector(
    (state: RootState) => state.categories.categories,
  );
  const necessaryCategory = categories?.find(
    (category) => String(category?.id) === id,
  );
  const subcategories = necessaryCategory?.subcategories || [];
  const swiperRef = useRef<SwiperType | null>(null);
  const user = useSelector((state: RootState) => state.auth.user);
  const cart = useSelector((state: RootState) => state.cart.cart);
  const [showFilter, setShowFilter] = useState(false);
  const [showSortBlock, setShowSortBlock] = useState(false)
  const [initialValuesForFilter, setInitialValuesForFilter] =
    useState<InitialValuesForFilter | null>(null);
  const debounceRef = useRef<ReturnType<typeof setTimeout> | null>(null);

  useEffect(() => {
    if (id) {
      getProductsForCategories(id).then((data) => setProducts(data));
    }
  }, [id]);

  useEffect(() => {
    if (id) {
      getInitialValuesForFilter(id, choiceSubcategory).then((data) => {
        if (data) {
          setInitialValuesForFilter(data);
        }
      });
    }
  }, [id, choiceSubcategory]);

  useEffect(() => {
    if (choiceSubcategory === 0) {
      swiperRef.current?.slideTo(0);
    }
  }, [choiceSubcategory]);

  const addProductToCart = (product: Products) => {
    const productInCart = cart?.order_items.find(
      (el) => el.product_id === product.id,
    );
    const isDisabled =
      product.availability_of_pieces === productInCart?.quantity;
    if (isDisabled) return;
    const newProductInCart: ProductsInOrders = {
      product_id: product.id,
      product_name: product.product_name,
      category_id: product.category_id,
      price: product.price,
      image_url: product.image_url,
      quantity: 1,
      category_name: "",
    };
    const initOrder: NewOrderType = {
      email: user?.email || null,
      phone: user?.phone || null,
      user_id: user?.user_id || null,
      order_items: [],
    };
    dispatch(initialOrder(initOrder));
    dispatch(addToCart(newProductInCart));
  };

  const deleteProductInCart = (productId: ProductIdType) => {
    dispatch(removeProductFromCart(productId));
  };

  useEffect(() => {
    localStorage.setItem("cart", JSON.stringify(cart));
  }, [cart]);

  const handleFindProducts = (symbols: string) => {
    if (debounceRef.current) clearTimeout(debounceRef.current);

    debounceRef.current = setTimeout(() => {
      if (symbols.trim() !== "") {
        findProduct(symbols, id).then((data) => setProducts(data));
      } else {
        if (id) {
          getProductsForCategories(id).then((data) => setProducts(data));
        }
      }
    }, 300);
  };

  return (
    <div className="products-categories-wrapper">
      <Swiper
        spaceBetween={50}
        slidesPerView={4}
        slidesOffsetBefore={44}
        className="swiper-products-subcategories"
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
                className={
                  subcategory.subcategory_id === choiceSubcategory
                    ? "active-products-swiper"
                    : ""
                }
              >
                <div
                  onClick={() =>
                    setChoceSubcategory(subcategory.subcategory_id)
                  }
                  className="subcategory-name"
                >
                  {t(`subcategories.${subcategory.category_name}`)}
                </div>
              </SwiperSlide>
            );
          })}
      </Swiper>
      {choiceSubcategory === 0 ? (
        <Accordion defaultActiveKey="0">
          <Accordion.Item eventKey="0">
            <Accordion.Header>{t("helpers.search")}</Accordion.Header>
            <Accordion.Body>
              <InputGroup className="mb-3">
                <InputGroup.Text id="basic-addon1">
                  <i className="bi bi-search"></i>
                </InputGroup.Text>
                <Form.Control
                  type="text"
                  id="inputPassword5"
                  aria-describedby="passwordHelpBlock"
                  onChange={(event) => handleFindProducts(event.target.value)}
                />
              </InputGroup>
            </Accordion.Body>
          </Accordion.Item>
        </Accordion>
      ) : null}
      {choiceSubcategory === 0 ? (
        <div className="products-categories">
          {showFilter ? (
            <i
              className="bi bi-filter-circle"
              onClick={() => setShowFilter(false)}
            ></i>
          ) : (
            <i className="bi bi-filter" onClick={() => setShowFilter(true)}></i>
          )}
          {showFilter ? (
            <Filters
              categoryId={id || ""}
              subcategoryId={""}
              initialValuesForFilter={initialValuesForFilter}
              setProducts={setProducts}
            />
          ) : null}
          <ProductsList
            products={products}
            deleteProductInCart={deleteProductInCart}
            addProductToCart={addProductToCart}
            isSwiper={false}
          />
          {showSortBlock && <SortProducts setShowSortBlock={setShowSortBlock} setProducts={setProducts} id={id}/>}
          {showSortBlock ? <i className="bi bi-sort-up" onClick={() => setShowSortBlock(false)}></i> : <i className="bi bi-sort-down"  onClick={() => setShowSortBlock(true)}></i>}
        </div>
      ) : (
        <ProductsSubcategories
          categoryId={Number(id)}
          choiceSubcategory={choiceSubcategory}
          addProductToCart={addProductToCart}
          deleteProductInCart={deleteProductInCart}
          initialValuesForFilter={initialValuesForFilter}
        />
      )}
    </div>
  );
};

export default CategoryProducts;
