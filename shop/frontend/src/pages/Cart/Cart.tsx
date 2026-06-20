import { useDispatch, useSelector } from "react-redux";
import type { RootState } from "../../store/store";
import { config } from "../../config";
import "./cart.css";
import { useTranslation } from "react-i18next";
import { addToCart, cleaningOrder, removeProductFromCart } from "../../store/cartSlice";
import type { NewOrderType, ProductIdType, ProductsInOrders } from "../../../types/types";
import { useEffect, useState } from "react";
import { Alert, Button, Form, InputGroup } from "react-bootstrap";
import { IMaskInput } from "react-imask";
import { createOrder } from "../../api/orders-api/orders-api";

const Cart = () => {
  const dispatch = useDispatch();
  const { t } = useTranslation();
  const cart = useSelector((state: RootState) => state.cart.cart);
  const user = useSelector((state: RootState) => state.auth.user);
  const [email, setEmail] = useState<string>("");
  const [phone, setPhone] = useState<string>("");
  const [error, setError] = useState(false);
  const [show, setShow] = useState(false);
  const totalPrice = cart
    ? cart?.order_items.reduce((acc, el) => (acc +=(el.price * el.quantity)), 0)
    : 0;

  useEffect(() => {
    localStorage.setItem("cart", JSON.stringify(cart));
  }, [cart]);

  const addProductToCart = (product: ProductsInOrders) => {
    dispatch(addToCart(product));
  };

  const deleteProductInCart = (productId: ProductIdType) => {
    dispatch(removeProductFromCart(productId));
  };

  const createNewOrder = () => {
    if (cart) {
      if (user) {
        createOrder(cart).then(() => dispatch(cleaningOrder()))
      } else {
         const hasEmail = email.length > 0 && email.includes("@")
    const hasPhone = phone.length === 18

    if (!hasEmail && !hasPhone) {
        setError(true)
        setShow(true)
        return
    }
        const initialOrder: NewOrderType = {
          email: email,
          phone: phone,
          user_id: null,
          order_items: cart.order_items
        }
        createOrder(initialOrder).then(() => dispatch(cleaningOrder()))
      }
    }
  }

  return (
    <div className="cart-page">
      {cart?.order_items.length ? <div className="cart-wrapper">
        {error && show && (
            <Alert onClose={() => setShow(false)} dismissible variant="danger">
              {t("orders.order_error")}
            </Alert>
          )}
        {cart?.order_items.map((product) => {
          return (
            <div key={product.product_id} className="cart-product">
              <div className="cart-product__image-wrapper">
                <img
                  className="cart-product__image"
                  src={`${config.PRODUCTS_IMAGES_BASE_URL}/${product.image_url}`}
                  alt={product.product_name}
                />
              </div>

              <div className="cart-product__info">
                <div className="cart-product__name">{product.product_name}</div>

                <div className="cart-product__actions">
                  <button
                    className="cart-btn"
                    onClick={() =>
                      deleteProductInCart({ product_id: product.product_id })
                    }
                  >
                    −
                  </button>
                  <span className="cart-quantity">{product.quantity}</span>
                  <button
                    className="cart-btn"
                    onClick={() => addProductToCart(product)}
                  >
                    +
                  </button>
                </div>
              </div>
            </div>
          );
        })}
        {!user && (
          <div className="cart-info-data">
            <InputGroup>
              <InputGroup.Text id="basic-addon1">@</InputGroup.Text>
              <Form.Control
                name="email"
                autoComplete="email"
                placeholder={t("auth.email")}
                onChange={(e) => setEmail(e.target.value)}
              />
            </InputGroup>
            <div className="cart-guest-hint">{t("orders.or")}</div>
            <InputGroup>
              <IMaskInput
                mask="+{7} (000) 000-00-00"
                placeholder={t("auth.phone")}
                id="phone"
                className="form-control"
                onChange={(e) => setPhone(e.target.value)}
                value={phone}
              />
            </InputGroup>
            <div className="cart__total-price">
              <strong>{t("orders.total_amount")}:</strong> {totalPrice} ₽
            </div>
          </div>
        )}
        <div className="cart-action">
        <Button onClick={() => createNewOrder()} variant="primary">{t("orders.place_order")}</Button>
        </div>
      </div> : <div className="cart-empty">{t("orders.cart_is_empty")}</div>}
    </div>
  );
};

export default Cart;
