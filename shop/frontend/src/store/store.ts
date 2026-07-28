import { configureStore } from "@reduxjs/toolkit";
import authReducer from "./authSlice.ts"
import categoryReducer from "./categoriesSlice.ts"
import cartReducer from "./cartSlice.ts"
import productReducer from "./productsSlice.ts"

export const store = configureStore({
    reducer: {
        auth: authReducer,
        categories: categoryReducer,
        cart: cartReducer,
        products: productReducer
    }
})

export type RootState = ReturnType<typeof store.getState>
export type AppDispatch = typeof store.dispatch
