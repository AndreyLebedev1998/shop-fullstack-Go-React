import { config } from '../../config.ts'
import type { Categories, Products, SubcategoryProductType } from '../../../types/types.ts'

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
    const baseURL = config.PRODUCTS_BASE_URL
    const url = new URL('/products-for-categories', baseURL)
    url.searchParams.set("category_id", id)
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
    url.searchParams.set("subcategory_id", String(id))
    try {
        const res = await fetch(url.toString())
        const data: SubcategoryProductType[]  = await res.json()
        return data 
    } catch (error) {
        console.error(error)
        return []
    }
}
