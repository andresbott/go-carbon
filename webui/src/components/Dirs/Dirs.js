import React, {memo, useEffect, useMemo} from 'react';
import {usePath} from "../../data/FileList/FileList";


const Dirs = () => {

    console.count("render")

    const data = usePath(state => state.data)
    const loaded = usePath(state => state.loaded)

    const loadPath = usePath(state => state.loadPath)
    useEffect(() => {
        const a = loadPath("/")
        console.log("ajax")
    }, []);


    if (!loaded) {
        console.log("loading")
        return <h2>loading...</h2>
    } else {
        console.log("loaded")
    }

    
    // const data = {};

    return (
        <MemoizedDirsContent files={data}></MemoizedDirsContent>
    );
};


export default Dirs;


const DirList = ({files}) => {
    console.count("render files")


    return (
        <>
            <h2>Directories</h2>
            <pre>{JSON.stringify(files, null, 2)}</pre>
        </>
    );
}

const MemoizedDirsContent = memo(
    DirList,
    (prevProps, newProps) => newProps.loading !== false //condition to determine when you want to update component
);
