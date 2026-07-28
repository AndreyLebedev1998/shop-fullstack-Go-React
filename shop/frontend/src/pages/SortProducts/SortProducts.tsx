import { useState, type FC } from "react"
import { indicators } from "../../constants/constants"
import { getProductsForCategories, sortProduct } from "../../api/products-api/products-api"
import type { Products } from "../../../types/types"
import "./sort-products.css"

interface TSortProducts {
    setShowSortBlock: (show: boolean) => void
    setProducts: (products: Products[]) => void
    id: string | undefined
}

const SortProducts: FC<TSortProducts> = ({setShowSortBlock, setProducts, id}) => {
    const localStorageIndicator = localStorage.getItem("products_sort_indicator")
    const [selected, setSelected] = useState<string | null>(null)
    const [loadingIndicator, setLoadingIndicator] = useState<string | null>(null)
    const handleSortProducts = (indicator: string) => {
        sortProduct(indicator).then((data) => {
            setSelected(indicator)
            setLoadingIndicator(null)
            localStorage.setItem("products_sort_indicator", indicator)
            setShowSortBlock(false)
            setProducts(data)
        }).finally(() => {
            setShowSortBlock(false)
        })
    }

    const handleWithoutIndictatorGetProducts = () => {
        if (id) {
            getProductsForCategories(id).then((data) => {
                localStorage.removeItem("products_sort_indicator")
                setProducts(data)
            }).finally(() => {
                setShowSortBlock(false)
            })
        }
    }

    return (
        <div className="sort-wrapper">
            <div className="sort-panel">
                <div className="sort-panel__header">
                    <span>Сортировка</span>
                    <i
                        className="bi bi-x-lg sort-panel__close"
                        onClick={() => setShowSortBlock(false)}
                    />
                </div>
                <div className="sort-panel__list">
                    {indicators.map((indicator) => {
                        const isActive = selected === (indicator || localStorageIndicator) || indicator === localStorageIndicator
                        const isLoading = loadingIndicator === indicator
                        return (
                            <button
                                key={indicator}
                                className={`sort-option ${isActive ? "sort-option--active" : ""}`}
                                onClick={() => handleSortProducts(indicator)}
                                disabled={isLoading}
                            >
                                <span className="sort-option__radio">
                                    {isActive && <span className="sort-option__radio-dot" />}
                                </span>
                                <span className="sort-option__label">
                                    {indicator}
                                </span>
                            </button>
                        )
                    })}
                </div>
                <button className="sort-reset" onClick={handleWithoutIndictatorGetProducts}>Сбросить</button>
            </div>
        </div>
    )
}

export default SortProducts
