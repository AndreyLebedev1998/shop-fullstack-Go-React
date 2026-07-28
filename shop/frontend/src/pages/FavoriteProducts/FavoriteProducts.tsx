import { useTranslation } from 'react-i18next';
import { useDispatch, useSelector } from 'react-redux';
import type { RootState } from '../../store/store';
import './favorite-products.css'
import ProductsList from '../../components/ProductsList/ProductsList';
import type { NewOrderType, ProductIdType, Products, ProductsInOrders } from '../../../types/types';
import { addToCart, initialOrder, removeProductFromCart } from '../../store/cartSlice';
import { Navigate } from 'react-router-dom';

const FavoriteProducts = () => {
    const dispatch = useDispatch()
    const favoriteProducts = useSelector(
    (state: RootState) => state.products.favoriteProduct,
  );
  const cart = useSelector((state: RootState) => state.cart.cart);
  const user = useSelector((state: RootState) => state.auth.user);
  const token = useSelector((state: RootState) => state.auth.token);
  const { t } = useTranslation();

   if (!token) {
    return <Navigate to={"/"} />;
  }

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
        <div className="favorite-products-page">
            <h3>{t("header.offcanvas.favorite_products")}</h3>
            <ProductsList products={favoriteProducts} deleteProductInCart={deleteProductInCart} addProductToCart={addProductToCart} isSwiper={false} />
        </div>
)
}

export default FavoriteProducts
