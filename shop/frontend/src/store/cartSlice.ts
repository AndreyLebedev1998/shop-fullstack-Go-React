import { createSlice, type PayloadAction } from "@reduxjs/toolkit";
import type { NewOrderType, ProductsInOrders, ProductIdType } from "../../types/types";

interface NewOrderState {
    cart: NewOrderType | null
}

const cartFromStorage = localStorage.getItem("cart");

const initialState: NewOrderState = {
    cart: cartFromStorage ? JSON.parse(cartFromStorage) : null,
}

const cartSlice = createSlice({
    name: "cart",
    initialState,
    reducers: {
        initialOrder(state, action: PayloadAction<NewOrderType>) {
            if (!state.cart) {
                state.cart = action.payload
            }
        },
        addToCart(state, action: PayloadAction<ProductsInOrders>) {
            if (state.cart) {
                const productInCart = state.cart?.order_items.find((product) => product.product_id === action.payload.product_id)
                if (productInCart) {
                    state.cart.order_items = state.cart.order_items.map((product) => {
                        if (product.product_id === productInCart.product_id) {
                            return {...product, quantity: product.quantity += 1}
                        }
                        return product
                    })
                } else {
                    state.cart.order_items = [...state.cart.order_items, action.payload]
                }
            }
        },
        removeProductFromCart(state, action: PayloadAction<ProductIdType>) {
            const removeProduct = state.cart?.order_items.find((product) => product.product_id === action.payload.product_id)
            if (removeProduct) {
                if (state.cart) {
                    if (removeProduct.quantity >= 2) {
                        state.cart.order_items = state.cart.order_items.map((product) => {
                            if (product.product_id === removeProduct.product_id) {
                                return {...product, quantity: product.quantity -= 1}
                            }
                            return product
                        })
                    } else {
                        state.cart.order_items = state.cart?.order_items.filter((product) => product.product_id !== removeProduct.product_id)
                    }
                }
            }
        },
        cleaningOrder(state) {
            state.cart = null
            localStorage.removeItem("cart")
        }
    }
})

export const {addToCart, removeProductFromCart, initialOrder, cleaningOrder} = cartSlice.actions
export default cartSlice.reducer
