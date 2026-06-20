import type { NewOrderType, OrderType } from "../../../types/types"
import { config } from "../../config"

export async function getOrdersAllTime(token: string, setError: (err: boolean) => void, setShow: (show: boolean) => void): Promise<OrderType[] | null> {
    const baseURL = config.ORDERS_BASE_URL
    const url = new URL('/get-orders-by-user', baseURL)
    try {
        const res = await fetch(url.toString(), {
            method: "POST",
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({token: token })
        })
        const data: OrderType[] = await res.json()
        return data
    } catch (error) {
        console.error(error)
        setError(true)
        setShow(true)
        return null
    }
}

export async function createOrder(newOrder: NewOrderType): Promise<OrderType | null> {
    const baseURL = config.ORDERS_BASE_URL
    const url = new URL('create-order', baseURL) 

    try {
        const res = await fetch(url.toString(), {
            method: "POST",
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(newOrder)
        })
        const data: OrderType = await res.json()
        return data
    } catch (error) {
        console.error(error)
        return null
    }
}
