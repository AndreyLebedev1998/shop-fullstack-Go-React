import { createSlice, type PayloadAction } from "@reduxjs/toolkit";
import type { Categories } from "../../types/types";

interface CategoriesState {
    categories: Categories[] | null
}

const initialState: CategoriesState = {
    categories: null,
}

const categoriesSlice = createSlice({
    name: 'auth',
    initialState,
    reducers: {
        setSubcategories(state, action: PayloadAction<Categories[]>) {
            state.categories = action.payload
        },
    }
})

export const {setSubcategories} = categoriesSlice.actions
export default categoriesSlice.reducer
