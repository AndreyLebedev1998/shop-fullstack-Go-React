import React, { useEffect, useState } from "react"
import { useTranslation } from "react-i18next"
import "./recommendations.css"
import type { NewOrderType, ProductIdType, Products, ProductsInOrders } from "../../../types/types"
import type { RootState } from "../../store/store"
import { useDispatch, useSelector } from "react-redux"
import { addToCart, initialOrder, removeProductFromCart } from "../../store/cartSlice"
import ProductsList from "../ProductsList/ProductsList"
import { getRecommendationsForUser } from "../../api/products-api/products-api"
import { Accordion } from "react-bootstrap"

const Recommendations = () => {
    const { t } = useTranslation()
    const dispatch = useDispatch()
    const [products, setProducts] = useState<Products[]>([])
    const cart = useSelector((state: RootState) => state.cart.cart);
    const user = useSelector((state: RootState) => state.auth.user);
    const token = useSelector((state: RootState) => state.auth.token);

    useEffect(() => {
        if (token) {
            getRecommendationsForUser(token).then((data) => setProducts(data))
        }
    }, [token])

    const addProductToCart = (product: Products) => {
          const productInCart = cart?.order_items.find((el) => el.product_id === product.id)
          const isDisabled = product.availability_of_pieces === productInCart?.quantity
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

    return (
        <div className="recommendations-block">
            <Accordion defaultActiveKey="0">
      <Accordion.Item eventKey="0">
        <Accordion.Header><h3>{t("home.recommendations")}</h3></Accordion.Header>
        <Accordion.Body>
            <ProductsList products={products} deleteProductInCart={deleteProductInCart} addProductToCart={addProductToCart} isSwiper={true} />
        </Accordion.Body>
      </Accordion.Item>
    </Accordion>
        </div>
    )
}

export default React.memo(Recommendations)
