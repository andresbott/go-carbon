import React, {useContext, useState, useEffect, createContext} from "react";

const APIContext = createContext({});

export function APIContextProvider({children}) {


    const [path, setPath] = useState("/")
    const url = "http://localhost:8080/api/v0/fe?path=" + path;
    // const {data: apiData, loading, error} = useFetch(url);

    const [data, setData] = useState(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);

    const DirectoryListContext = {
        loading: loading,
        error: error,
        data: data,
        path: path,
        setPath: setPath,
    }


    useEffect(() => {
        async function fetchData() {
            try {
                const response = await fetch(url);
                if (!response.ok) {
                    throw new Error('Failed to fetch data');
                }
                const json = await response.json();
                setData(json);
                setLoading(false);
            } catch (error) {
                setError(error.message);
                setLoading(false);
            }
        }

        fetchData();
    }, [path]);
    return (
        <APIContext.Provider value={{DirectoryListContext}}>
            {children}
        </APIContext.Provider>
    );
}

export function useAPI() {
    const context = useContext(APIContext);
    if (context === undefined) {
        throw new Error("Context must be used within a Provider");
    }
    return context.DirectoryListContext;
}
