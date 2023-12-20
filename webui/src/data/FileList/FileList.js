import {create} from 'zustand';
import {useEffect} from "react";

// const APIURL = "http://localhost:8080/api/v0/"
const _ApiUrl_ = "https://dummyjson.com/products/1"

const usePath = create((set, get) => ({
        loadPath: async (path) => {
            try {
                set({loaded: false})
                const response = await fetch(_ApiUrl_);
                const json = await response.json();
                set({data: json})
                // await new Promise(r => setTimeout(r, 1000));
            } catch (error) {
                set({error: error})
            } finally {
                set({loaded: true})
            }

        },

        // todo add a timestamp with the last loaded and only reload if some time has passed
        path: "/",
        loaded: false,
        error: null,
        data: {},
        // updateUser: (user) => {
        //     set({data: {...get().data, user: {...get().data?.user, ...user}}})
        // },
        // removeCard: (cardId) => {
        //     const cards = get().data.cards.filter(card => card.id !== cardId);
        //     set({data: {...get().data, cards}})
        // }
    })
)


export {usePath}