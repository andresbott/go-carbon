import React, {useContext} from "react";
import {DirListContext} from "../../data/DirListContext";
import {useAPI} from "../../data/DirList2Context";

function DirList() {

    const api = useAPI()

    const handleClick = () => {
        api.setPath("/tmp")

    };


    return (
        <ul>
            <li>Directory list</li>
            <button onClick={handleClick}>change Dir</button>
        </ul>
    )
}

export default DirList;
