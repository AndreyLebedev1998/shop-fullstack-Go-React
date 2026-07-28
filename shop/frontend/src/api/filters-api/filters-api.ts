import type { FiltersType, InitialValuesForFilter, Products } from '../../../types/types.ts'
import { config } from '../../config.ts'

export const getProductsFromFilters = async (data: FiltersType): Promise<Products[]> => {
    const baseUrl = config.PRODUCTS_BASE_URL
    const path = "/get-products-for-filters"
    const url = new URL(path, baseUrl);

    (Object.keys(data) as Array<keyof FiltersType>).forEach((el) => {
        if (data[el]) {
            url.searchParams.set(String(el), data[el])
        }
    })

    try {
        const res = await fetch(url.toString())
        const data = await res.json()
        return data
    } catch (error) {
        console.error(error)
        return []
    }
}

export const getInitialValuesForFilter = async (categoryId: string, subcategoryId: number): Promise<InitialValuesForFilter | null> => {
    const baseUrl = config.PRODUCTS_BASE_URL
    const path = "/get-initial-values-for-filter"
    const url = new URL(path, baseUrl)
    url.searchParams.set("category_id", String(categoryId))
    if (subcategoryId != 0) {
        url.searchParams.set("subcategory_id", String(subcategoryId))
    }

    try {
        const res = await fetch(url.toString())
        const data = await res.json()
        return data
    } catch (error) {
        console.error(error)
        return null
    }
}
