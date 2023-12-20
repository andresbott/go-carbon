// ComponentOne.js
import React, {useContext} from 'react';
import {MyContext} from "../data/MyContext";
import {APIContext, useAPI} from "../data/DirList2Context";


const ComponentOne = () => {
    const {setCount} = useContext(MyContext);
    const api = useAPI()
    console.log(api)
    // console.log(loading)
    // console.log(error)
    // console.log(error)

    const handleIncrement = () => {
        setCount((prevCount) => prevCount + 1);
    };

    const handleDecrement = () => {
        setCount((prevCount) => prevCount - 1);
    };

    return (
        <div>
            <h2>Component One</h2>
            <button onClick={handleIncrement}>Increment</button>
            <button onClick={handleDecrement}>Decrement</button>
        </div>
    );
};

export default ComponentOne;