import { defineStore } from 'pinia'

import axios from 'axios'
import router from '@/router/index.js'

export const useUserStore = defineStore('user', {
    state: () => ({
        _isLoggedIn: false,
        _isFirstLogin: true,
        _user: '',
        _error: '',
        _loading: false,
        _wrongPw: false
    }),
    getters: {
        isLoggedIn(state) {
            return state._isLoggedIn
        },
        isFirstLogin(state) {
            return state._isFirstLogin
        },
        userName: (state) => state._user,
        loading: (state) => state._loading,
        wrongPw: (state) => state._wrongPw
    },
    actions: {
        async checkState() {
            const endpoint = import.meta.env.VITE_SERVER_URL_V0 + '/user/status'
            return axios
                .get(endpoint)
                .then((res) => {
                    if (res.status === 200) {
                        this._user = res.data.user
                        this._isLoggedIn = res.data['logged-in']
                    } else {
                        console.log('err')
                        console.log(res)
                        // error?
                    }
                })
                .catch((err) => {
                    console.log(err)
                })

            // console.log(import.meta.env.VITE_NOT_SECRET_CODE)
        },
        setFirtLoginFalse() {
            this._isFirstLogin = false
        },
        // USER LOGOUT
        async logout() {
            const endpoint = import.meta.env.VITE_SERVER_URL_V0 + '/user/logout'

            this._loading = true
            axios
                .post(endpoint, '')
                .then((res) => {
                    this._user = ''
                    this._isLoggedIn = false
                    router.push('/login')
                })
                .catch((err) => {
                    console.log(err)
                    // todo propagate login error
                })
                .finally(() => {
                    this._loading = false
                })
        },

        // USER LOGIN
        async login(user, pass) {
            const data = {
                user: user,
                password: pass
            }
            const endpoint = import.meta.env.VITE_SERVER_URL_V0 + '/user/login'

            const authAxios = axios.create()
            authAxios.interceptors.response.use(
                (response) => {
                    return response
                },
                (error) => {
                    if (error.response.status === 401) {
                        console.log('auth NOT OK')
                        this._isLoggedIn = false
                        this._wrongPw = true
                    }
                    return error
                }
            )

            this._loading = true
            authAxios
                .post(endpoint, data)
                .then((res) => {
                    console.log(res.data)
                    if (res.status === 200) {
                        this._user = res.data.user
                        this._isLoggedIn = true
                        this._wrongPw = false
                        router.push('/app')
                    } else {
                        console.log('err')
                        console.log(res)
                        // error?
                    }
                })
                .catch((err) => {
                    console.log(err)
                    // todo propagate login error

                    console.log('h')
                    // this.$toasted.show(
                    //     'Please enter the correct details and try again',
                    //     err,
                    //     {
                    //         position: 'top-left',
                    //         duration: 200,
                    //         type: danger,
                    //     }
                    // )
                })
                .finally(() => {
                    this._loading = false
                })

            // console.log(import.meta.env.VITE_NOT_SECRET_CODE)
        }
    }
})
