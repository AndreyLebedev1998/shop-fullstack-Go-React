import { BrowserRouter, Link, Route, Routes } from "react-router-dom";
import Header from "./components/Header/Header";
import Home from "./pages/Home/Home.tsx";
import Authorization from "./pages/Authorization/Authorization.tsx";
import "./App.css";
import { useEffect, useState } from "react";
import type { RootState } from "./store/store.ts";
import { useDispatch, useSelector } from "react-redux";
import { getMe } from "./api/auth-api/auth-api.ts";
import { logout, setAuth } from "./store/authSlice.ts";
import Orders from "./pages/Orders/Orders.tsx";
import CategoryProducts from "./pages/CategoryProducts/CategoryProducts.tsx";
import Cart from "./pages/Cart/Cart.tsx";
import RecoveryPassword from "./pages/RecoveryPassword/RecoveryPassword.tsx";
import NewPassword from "./pages/NewPassword/NewPassword.tsx";

function App() {
  const dispatch = useDispatch();
  const token = useSelector((state: RootState) => state.auth.token);
  const [choiceSubcategory, setChoceSubcategory] = useState<number>(0);
  const cart = useSelector((state: RootState) => state.cart.cart)
  const quantityInCart = cart && cart.order_items.length !== 0 && cart.order_items.reduce((acc, product) => (
    acc += product.quantity
  ), 0)

  useEffect(() => {
    const savedToken = token || localStorage.getItem("token");
    if (savedToken) {
      getMe(savedToken).then((data) => {
        if (data) {
          dispatch(setAuth(data));
        } else {
          dispatch(logout());
        }
      });
    }
  }, [dispatch, token]);

  return (
    <>
      <BrowserRouter>
        <div className="shop">
          <Header setChoceSubcategory={setChoceSubcategory} />
          <Routes>
            <Route path="/" element={<Home />} />
            <Route path="/authorization" element={<Authorization />} />
            <Route path="/orders" element={<Orders />} />
            <Route
              path="/products-categories/:id"
              element={
                <CategoryProducts
                  choiceSubcategory={choiceSubcategory}
                  setChoceSubcategory={setChoceSubcategory}
                />
              }
            />
            <Route path="/cart" element={<Cart />} />
            <Route path="/recovery-password" element={<RecoveryPassword />} />
            <Route path="/new-password" element={<NewPassword/>} />
          </Routes>
          <Link to="/cart">
            <i className="bi bi-cart">
              {quantityInCart && <span className="quantity-in-cart">{quantityInCart}</span>}
            </i>
          </Link>
        </div>
      </BrowserRouter>
    </>
  );
}

export default App;
