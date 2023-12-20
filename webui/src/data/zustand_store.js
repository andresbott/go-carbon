import create from 'zustand';

const useStore = create((set, get) => ({
        loadData: async (id) => {
            try {
                set({loadingData: true})
                const response = await fetch('https://dummyjson.com/products/' + id);
                set({data: response})
            } catch {
                console.log("error")
                // Todo show error
            } finally {
                set({loadingData: false})
            }
        },
        loadingData: false,
        data: {prod: null, cards: []},
        updateUser: (user) => {
            set({data: {...get().data, user: {...get().data?.user, ...user}}})
        },
        removeCard: (cardId) => {
            const cards = get().data.cards.filter(card => card.id !== cardId);
            set({data: {...get().data, cards}})
        }
    })
)

export default useStore;