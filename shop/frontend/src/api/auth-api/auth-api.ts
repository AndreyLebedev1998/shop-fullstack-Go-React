import { config } from '../../config.ts'
import type { AuthData, CodeMatchingDataType, CodeMatchingTokenType, ErrorReg, MessageSuccess, MessageType, RecoveryPasswordData, RegRequestData, Token, UserId } from '../../../types/types.ts'

export async function authorization(email: string, password: string, setErrorServer: (err : boolean) => void, setShow: (show: boolean) => void): Promise<Token | null> {
    const baseURL = config.AUTH_BASE_URL
    const url = new URL('/authorization', baseURL)
    try {
        const res = await fetch(url.toString(), {
            method: "POST",
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ email, password })
        })

        if (!res.ok) {
            setErrorServer(true)
            setShow(true)
            throw new Error(`HTTP Error: ${res.status}`)
        }

        const data: Token = await res.json()
        return data
    } catch (error) {
        console.error(error)
        return null
    }
}

export async function getMe(token: string): Promise<AuthData | null> {
    const baseURL = config.AUTH_BASE_URL
    const url = new URL('/get-me', baseURL)
    try {
        const res = await fetch(url.toString(), {
            method: "GET",
            headers: { 'Authorization': `Bearer ${token}` },
        })
        const data: AuthData = await res.json()
        return data
    } catch (error) {
        console.error(error)
        return null
    }
}

export async function registration(regData: RegRequestData): Promise<UserId | ErrorReg | null> {
    const baseURL = config.AUTH_BASE_URL
    const url = new URL('/registration', baseURL)
    try {
        const res = await fetch(url.toString(), {
            method: "POST",
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ email: regData.email, password: regData.password, lastname: regData.lastname, phone: regData.phone,  name: regData.name})
        })

        if (res.status === 409) {
            return { error: 'conflict', user_id: null } as ErrorReg
        }
        const data: AuthData = await res.json()
        return data
    } catch (error) {
        console.error(error)
        return null
    }
}

export async function sendCodeTg(email: string, setErrorMessage: (err: boolean) => void, setShowError: (show: boolean) => void, setShowErrorMessageChatId: (err: boolean) => void, setErrorMessageChatId: (err: boolean) => void, setSuccessMeaage: (data: boolean) => void): Promise<MessageType | null> {
    const baseURL = config.AUTH_BASE_URL
    const url = new URL('/send-code-from-tg', baseURL)
    try {
        const res = await fetch(url.toString(), {
            method: "POST",
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({email,})
        })
        if (res.status === 400) {
            setShowErrorMessageChatId(true)
            setErrorMessageChatId(true)
            return null
        }
        const data: MessageType = await res.json()
        return data
    } catch (error) {
        setErrorMessage(true)
        setShowError(true)
        setSuccessMeaage(false)
        console.error(error)
        return null
    }
}

export async function codeMatching(data: CodeMatchingDataType, setErrorCodeMathing: (err: boolean) => void, setShowErrorCodeMatching: (err: boolean) => void): Promise<CodeMatchingTokenType | null> {
    const baseUrl = config.AUTH_BASE_URL
    const path = '/code-mathing'
    const url = new URL(path, baseUrl)
    try {
        const res = await fetch(url.toString(), {
            method: "POST",
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(data)
        })

        const resData: CodeMatchingTokenType = await res.json()
        return resData
    } catch (error) {
        setErrorCodeMathing(true)
        setShowErrorCodeMatching(true)
        console.error(error)
        return null
    }
}

export async function recoveryPassword(data: RecoveryPasswordData): Promise<MessageSuccess | null> {
    const baseUrl = config.AUTH_BASE_URL
    const path = "/recovery-password"
    const url = new URL(path, baseUrl)
    try {
        const res = await fetch(url.toString(), {
            method: "POST",
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(data)
        })
        const resData = await res.json()
        return resData
    } catch (error) {
        console.error(error)
        return null
    }
}
