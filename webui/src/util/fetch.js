import { useState, useEffect } from 'react';

// https://developer.mozilla.org/en-US/docs/Web/API/fetch
export function useFetch (resource,options = {}) {
    const [data, setData] = useState(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);

    useEffect(() => {
        if (!resource) {
            return;
        }

        const fetchData = async () => {
            try {
                console.log("fetch")
                const response = await fetch(resource);
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
        };

        fetchData();
    }, [resource]);

    return { data, loading, error };
}

