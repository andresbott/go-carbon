import { defineStore } from 'pinia'

import axios from "axios"

export const useUserStore = defineStore("user", {
    state: () => ({
        _isLoggedIn: false,
    }),
    getters: {
        isLoggedIn(state){
            return state._isLoggedIn
        }
    },
    actions: {
        // async fetchUsers() {
        //     try {
        //         const data = await axios.get('https://jsonplaceholder.typicode.com/users')
        //         this.users = data.data
        //     }
        //     catch (error) {
        //         alert(error)
        //         console.log(error)
        //     }
        // },
        // async login() {
        //     const data = {
        //         email: this.email,
        //         password: this.password,
        //     }
        //     axios
        //         .post('localhost:8085/api/v0/user/login', data)
        //         .then((res) => {
        //             const userData = res.data
        //             userData.user.token = userData.token
        //
        //             // this.$store.commit('setUserDetails', userData.user)
        //             // this.$toasted.show('You have logged in successfully', {
        //             //     position: 'top-center',
        //             //     duration: 500,
        //             //     type: 'success',
        //             // })
        //             // this.$router.push('/home')
        //         })
        //         .catch((err) => {
        //             console.log(err)
        //             // this.$toasted.show(
        //             //     'Please enter the correct details and try again',
        //             //     err,
        //             //     {
        //             //         position: 'top-left',
        //             //         duration: 200,
        //             //         type: danger,
        //             //     }
        //             // )
        //         })
        // },
        async checkState(){
            https://medium.com/@bugintheconsole/axios-vue-js-3-pinia-a-comfy-configuration-you-can-consider-for-an-api-rest-a6005c356dcd

            const endpoint = import.meta.env.VITE_SERVER_URL_V0+"/user/status"

            axios
                .get(endpoint)
                .then((res) => {
                    const userData = res.data
                    // console.log(userData)
                    this._isLoggedIn = true
                    console.log(this.isLoggedIn)
                    // userData.user.token = userData.token

                    // this.$store.commit('setUserDetails', userData.user)
                    // this.$toasted.show('You have logged in successfully', {
                    //     position: 'top-center',
                    //     duration: 500,
                    //     type: 'success',
                    // })
                    // this.$router.push('/home')
                })
                .catch( (err) => {

                    console.log(err)
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

            // console.log(import.meta.env.VITE_NOT_SECRET_CODE)
        },
        async login(user, pass){
            const data = {
                user: user,
                password:pass,
            }
            const endpoint = import.meta.env.VITE_SERVER_URL_V0+"/user/login"

            const authAxios = axios.create();
            authAxios.interceptors.response.use(response => {
                return response;
            }, error => {
                if (error.response.status === 401) {
                    //place your reentry code
                    console.log("401")
                }
                return error;
            });

            authAxios
                .post(endpoint,data)
                .then((res) => {
                const userData = res.data
                    // console.log(userData)
                    this._isLoggedIn = true
                    console.log(this.isLoggedIn)
                // userData.user.token = userData.token

                // this.$store.commit('setUserDetails', userData.user)
                // this.$toasted.show('You have logged in successfully', {
                //     position: 'top-center',
                //     duration: 500,
                //     type: 'success',
                // })
                // this.$router.push('/home')
            })
                .catch( (err) => {

                    console.log(err)
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

            // console.log(import.meta.env.VITE_NOT_SECRET_CODE)
        },
    },

})