import { defineStore } from 'pinia'

import axios from "axios"

export const useUserStore = defineStore("user", {
    state: () => ({
        isLoggedIn: false,
    }),
    getters: {

        getUsers(state){
            return state.users
        }
    },
    actions: {
        async fetchUsers() {
            try {
                const data = await axios.get('https://jsonplaceholder.typicode.com/users')
                this.users = data.data
            }
            catch (error) {
                alert(error)
                console.log(error)
            }
        },
        async login() {
            const data = {
                email: this.email,
                password: this.password,
            }
            axios
                .post('localhost:3000/backend/api/auth/signin', data)
                .then((res) => {
                    const userData = res.data
                    userData.user.token = userData.token

                    // this.$store.commit('setUserDetails', userData.user)
                    // this.$toasted.show('You have logged in successfully', {
                    //     position: 'top-center',
                    //     duration: 500,
                    //     type: 'success',
                    // })
                    // this.$router.push('/home')
                })
                .catch((err) => {
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
        },
        async bla(user,pass){
            console.log("user: "+user)
            console.log("pass: "+pass)

            const data = {
                user: user,
                password:pass,
            }
            const endpoint = import.meta.env.VITE_SERVER_URL_V0+"/sys/user/login"
            axios
                .post(endpoint,data)
                .then((res) => {
                const userData = res.data
                userData.user.token = userData.token

                // this.$store.commit('setUserDetails', userData.user)
                // this.$toasted.show('You have logged in successfully', {
                //     position: 'top-center',
                //     duration: 500,
                //     type: 'success',
                // })
                // this.$router.push('/home')
            })
                .catch((err) => {
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