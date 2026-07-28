import { createSlice, type PayloadAction } from "@reduxjs/toolkit";
import type { FavoriteProduct, ProductId } from "../../types/types";

interface FavoriteProductState {
    favoriteProduct: FavoriteProduct[]
}

const initialState: FavoriteProductState = {
    favoriteProduct: [],
}

const productsSlice = createSlice({
    name: 'products',
    initialState,
    reducers: {
        setFavoriteProducts(state, action: PayloadAction<FavoriteProduct[]>) {
            state.favoriteProduct = action.payload
        },
        addFavoriteProduct(state, action: PayloadAction<FavoriteProduct>) {
            state.favoriteProduct = [...state.favoriteProduct, action.payload]
        },
        deleteFavoriteProduct(state, action: PayloadAction<ProductId>) {
            state.favoriteProduct = state.favoriteProduct.filter((product) => product.id != action.payload.product_id)
        }
    }
})

export const {setFavoriteProducts, addFavoriteProduct, deleteFavoriteProduct} = productsSlice.actions
export default productsSlice.reducer
