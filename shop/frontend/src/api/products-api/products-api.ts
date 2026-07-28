import { config } from '../../config.ts'
import type { Categories, FavoriteProduct, ProductId, Products, RecommendationProduct, ResRemoveFavoriteProduct, SubcategoryProductType } from '../../../types/types.ts'

export async function getCategoriesWithSubcategories(): Promise<Categories[]> {
    const baseURL = config.PRODUCTS_BASE_URL
    const url = new URL('/categories', baseURL)
    try {
        const res = await fetch(url.toString())
        const data: Categories[]  = await res.json()
        return data 
    } catch (error) {
        console.error(error)
        return []
    }
}

export async function getProductsForCategories(id: string): Promise<Products[]> {
    const indicator = localStorage.getItem("products_sort_indicator") 
    const baseURL = config.PRODUCTS_BASE_URL
    const url = new URL('/products-for-categories', baseURL)
    url.searchParams.set("category_id", id)
    if (indicator) {
        url.searchParams.set("indicator", indicator)
    }
    try {
        const res = await fetch(url.toString())
        const data: Products[]  = await res.json()
        return data 
    } catch (error) {
        console.error(error)
        return []
    }
}

export async function getProductsForSubcategories(id: number): Promise<SubcategoryProductType[]> {
    const baseURL = config.PRODUCTS_BASE_URL
    const url = new URL('/products-for-subcategory', baseURL)
    const indicator = localStorage.getItem("products_sort_indicator") 
    url.searchParams.set("subcategory_id", String(id))
    if (indicator) {
        url.searchParams.set("indicator", indicator)
    }
    try {
        const res = await fetch(url.toString())
        const data: SubcategoryProductType[]  = await res.json()
        return data 
    } catch (error) {
        console.error(error)
        return []
    }
}

export async function addFavoriteProductForUser(productId: number, token: string): Promise<FavoriteProduct | null> {
    const baseURL = config.PRODUCTS_BASE_URL
    const url = new URL("/add-favorite-product-for-user", baseURL)
    const body: ProductId = {
        product_id: productId
    }
    try {
        const res = await fetch(url.toString(), {
            method: "POST",
            headers: { 'Authorization': `Bearer ${token}` },
            body: JSON.stringify(body)
        })
        const data = await res.json()
        return data
    } catch (error) {
        console.error(error)
        return null
    }
}

export async function getFavoriteProduucts(token: string): Promise<FavoriteProduct[]> {
    const baseUrl = config.PRODUCTS_BASE_URL
    const url = new URL("/get-favorite-products", baseUrl)
    try {
        const res = await fetch(url.toString(), {
            method: "GET",
            headers: { 'Authorization': `Bearer ${token}` }
        })
        const data = await res.json()
        return data
    } catch (error) {
        console.error(error)
        return []
    }
}

export async function removeFavoriteProduct(token: string, productId: number): Promise<ResRemoveFavoriteProduct | null> {
    const baseURL = config.PRODUCTS_BASE_URL
    const url = new URL("/remove-favorite-products", baseURL)
    url.searchParams.set("product_id", String(productId))

    try {
        const res = await fetch(url.toString(), {
            method: "DELETE",
            headers: { 'Authorization': `Bearer ${token}` }
        })

        const data = await res.json()
        return data
    } catch (error) {
        console.error(error)
        return null
    }
}

export async function getRecommendationsForUser(token: string): Promise<RecommendationProduct[]> {
    const baseURL = config.PRODUCTS_BASE_URL
    const url = new URL("/get-recommendations", baseURL)
    try {
        const res = await fetch(url.toString(), {
            method: "GET",
            headers: { 'Authorization': `Bearer ${token}` }
        })
        const data = await res.json()
        return data
    } catch (error) {
        console.error(error)
        return []
    }
}

export async function findProduct(symbols: string, categoryId?: string, indicator?: string): Promise<Products[]> {
    console.log(symbols)
    const baseURL = config.PRODUCTS_BASE_URL
    const url = new URL('/find-product-by-symbols', baseURL)
    url.searchParams.set("symbols", symbols)
    if (indicator) {
        url.searchParams.set("indicator", indicator)
    }
    if (categoryId) {
        url.searchParams.set("category_id", categoryId)
    }
    try {
        const res = await fetch(url.toString())
        const data: Products[]  = await res.json()
        return data 
    } catch (error) {
        console.error(error)
        return []
    }
}

export async function sortProduct(indicator: string, categoryId?: string, subcategoryId?: string): Promise<Products[]> {
    const baseURL = config.PRODUCTS_BASE_URL
    const url = new URL('/sort-products', baseURL)
    url.searchParams.set("indicator", indicator)
    if (categoryId) {
        url.searchParams.set("category_id", categoryId)
    }
    if (subcategoryId) {
        url.searchParams.set("subcategory_id", subcategoryId)
    }
    try {
        const res = await fetch(url.toString())
        const data: Products[]  = await res.json()
        return data 
    } catch (error) {
        console.error(error)
        return []
    }
}
