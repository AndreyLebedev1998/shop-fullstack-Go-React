import { createSlice, type PayloadAction } from "@reduxjs/toolkit";
import type { AuthData, tokenRecoveryPasswordType } from "../../types/types";

interface AuthState {
    user: AuthData | null
    token: string | null
    tokenRecoveryPassword: tokenRecoveryPasswordType
    emailRecoveryPassword: string | null
}

const initialState: AuthState = {
    user: null,
    token: null,
    tokenRecoveryPassword: {token: ""},
    emailRecoveryPassword: null
}

const authSlice = createSlice({
    name: 'auth',
    initialState,
    reducers: {
        setAuth(state, action: PayloadAction<AuthData>) {
            if (!state.user) {
                state.user = {} as AuthData
            }
            state.user.authenticated = action.payload.authenticated
            state.user.name = action.payload.name
            state.user.lastname = action.payload.lastname
            state.user.email = action.payload.email
            state.user.phone = action.payload.phone
            state.user.user_id = action.payload.user_id
            state.token = action.payload.token
        },
        logout(state) {
            state.user = null
            state.token = null
            localStorage.removeItem("token")
        },
        setNewToken(state, action: PayloadAction<tokenRecoveryPasswordType> ) {
            state.tokenRecoveryPassword.token = action.payload.token
        },
        setEmailRecoveryPassword(state, action: PayloadAction<string>) {
            state.emailRecoveryPassword = action.payload
        }
    }
})

export const {setAuth, logout, setNewToken, setEmailRecoveryPassword} = authSlice.actions
export default authSlice.reducer
