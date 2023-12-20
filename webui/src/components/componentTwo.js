// ComponentTwo.js
import React, { useContext } from 'react';
import {MyContext} from "../data/MyContext";

const ComponentTwo = () => {
    const { count, setCount } = useContext(MyContext);

    const handleReset = () => {
        setCount(0);
    };

    return (
        <div>
            <h2>Component Two</h2>
            <p>Count from Component One: {count}</p>
            <button onClick={handleReset}>Reset</button>
        </div>
    );
};

export default ComponentTwo;