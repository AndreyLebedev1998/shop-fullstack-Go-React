import { useEffect, useState, type FC } from "react";
import "./filters.css";
import { useTranslation } from "react-i18next";
import { Button } from "react-bootstrap";
import { getProductsFromFilters } from "../../api/filters-api/filters-api";
import type { FiltersType, InitialValuesForFilter, Products } from "../../../types/types";

interface TFilters {
  categoryId: string
  subcategoryId: string
  initialValuesForFilter: InitialValuesForFilter | null
  setProducts: (products: Products[]) => void
}

const Filters: FC<TFilters> = ({categoryId, subcategoryId, initialValuesForFilter, setProducts}) => {
  const { t } = useTranslation();
  const [minPrice, setMinPrice] = useState("");
  const [maxPrice, setMaxPrice] = useState("")

  const submit = () => {
    const dataFilter: FiltersType = {
      min_price: String(minPrice),
      max_price: String(maxPrice),
      subcategory_id: subcategoryId,
      category_id: categoryId,
    }
    getProductsFromFilters(dataFilter).then((data) => {
      setProducts(data)
    })
  }

  useEffect(() => {
  if (!initialValuesForFilter) return;

  // eslint-disable-next-line react-hooks/set-state-in-effect
  setMinPrice(String(initialValuesForFilter.min_price));
  setMaxPrice(String(initialValuesForFilter.max_price));
}, [initialValuesForFilter]);

  return (
    <div className="filters filters-wrapper">
      <div className="filters-header"></div>
      <div className="filter-body">
        <div className="filter-section-title">{t("filters.price")}</div>
        <div className="price-filter">
          <div className="price-filter__row">
            <label className="price-filter__label" htmlFor="price-from">
              {t("filters.from")}
            </label>
            <input id="price-from" className="price-filter__input" value={minPrice} onChange={(e) => setMinPrice(e.target.value)}/>
          </div>
          <div className="price-filter__row">
            <label className="price-filter__label" htmlFor="price-to">
              {t("filters.to")}
            </label>
            <input id="price-to" className="price-filter__input" value={maxPrice} onChange={(e) => setMaxPrice(e.target.value)}/>
          </div>
        </div>
      </div>
      <div className="filter-actions">
        <Button onClick={submit}>Найти</Button>
      </div>
    </div>
  );
};

export default Filters;
