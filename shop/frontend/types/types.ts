export type Products = {
    id: number
    product_name: string
    category_id: number
    price: number
    image_url: string
    availability_of_pieces: number
    subcategory_id: number
    category_name: string
}

export type Subcategories = {
    category_id: number
    category_name: string
    subcategory_id: number
}

export type Categories = {
    id: number
    category_name: string
    subcategories: Subcategories[]
}

export type Token = {
    token: string
}

export type AuthData = {
    user_id: number
    email: string
    lastname: string
    name: string
    phone: string
    authenticated: boolean
    token: string
}

export type tokenRecoveryPasswordType = {
    token: string
}

export type RegRequestData = {
    email: string
    name: string
    lastname: string
    phone: string
    password: string
}

export type UserId = {
    user_id: number
}

export type ErrorReg = {
    error: string
    user_id: null
}

export type LinksType = {
    link: string
    name: string
    active: boolean
}

export type ProductsInOrders = {
    product_id: number
    product_name: string
    category_id: number
    price: number
    image_url: string
    quantity: number 
    category_name: string
}

export type OrderType = {
    order_id: number,
    user_id: number,
    email: string,
    phone: string | null,
    status: string,
    total_price: number,
    created_at: string,
    updated_at: string,
    products: ProductsInOrders[]
}

export type ProductsSubcategories = {
    id: number,
    product_name: string,
    price: number,
    category_id: number,
    image_url: string,
    availability_of_pieces: number,
    subcategory_id: number,
    category_name: string
}

export type SubcategoryProductType = {
    id: number,
    product_name: string,
    price: number,
    category_id: number,
    image_url: string,
    availability_of_pieces: number,
    subcategory_id: number,
    category_name: string
}

export type ProductsInNewOrder = {
    product_id: number,
    quantity: number
}

export type NewOrderType = {
    user_id: number | null,
    email: string | null,
    phone: string | null,
    order_items: ProductsInOrders[]
}

export type ProductIdType = {
    product_id: number
}

export type MessageType = {
    message: string
}

export type CodeMatchingDataType = {
    email: string
    phone: string
    code: string
}

export type CodeMatchingTokenType = {
    token: string
}

export type RecoveryPasswordData = {
    email: string
    phone: string
    token: string
    new_password: string
    confirmation_password: string
}

export type MessageSuccess = {
    message: string
}

export type FiltersType = {
    category_id: string
    subcategory_id: string
    min_price: string
    max_price: string
}

export type InitialValuesForFilter = {
    min_price: number
    max_price: number
}

export type FavoriteProduct = {
    id: number
    product_name: string
    price: number
    category_id: number
    image_url: string
    availability_of_pieces: number
    subcategory_id: number
    user_id: number
}

export type ProductId = {
    product_id: number
}

export type ResRemoveFavoriteProduct = {
    user_id: number
    product_id: number
}

export type RecommendationProduct = {
  id: number;
  product_name: string;
  price: number;
  category_id: number;
  image_url: string;
  availability_of_pieces: number;
  subcategory_id: number;
  category_name: string;
  rating: string;
}
