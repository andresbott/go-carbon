import React, {useEffect} from 'react';
import useStore from "../data/zustand_store";

const ZustandComp = () => {
    const {data, loadData, loadingData} = useStore();

    console.log(loadingData)
    // loadData(1)
    useEffect(() => {
        const d = loadData(1);


    }, []);

    console.log(data)

    return (
        <div>
            <h2>Zustand</h2>
        </div>
    );
};

export default ZustandComp;