import React, {useContext, useEffect, useState} from "react";

import {useFetch} from "../util/fetch";
import { parseJsonResp } from "../util/json";
import DirectoryList from "./DirList";

export const ReleasesContext = React.createContext({});


// context component
export function ReleasesProvider({ children }) {
    // Releases List state variables
    const [releases, setReleases] = useState({});

    // context object
    const releaseContext = {
        releases: releases,
        setReleases: setReleases,
    };

    // force propagate status change on the selected release


    return (
        <ReleasesContext.Provider value={releaseContext}>
            {children}
        </ReleasesContext.Provider>
    );
}

// custom context
export const DirListContext = React.createContext({});


export function DirectoryListContextProvider({ children }) {
    const [files, setFiles] = useState([]);
    const [path, setPath] = useState("/");

    const DirectoryListContext ={
        files: files,
        setFiles: setFiles,
        path: path,
        setPath: setPath,
    }
    // function setter (status){
    //     // set some status
    // }
    //
    // now in value we can pass the context item  value={{ dataContext }}

    return (
        <DirListContext.Provider value={ DirectoryListContext }>
            {children}
        </DirListContext.Provider>
    );
}


// custom context hook to write and read selected release data all over the application
// export function useReleasesContext() {
//     return useContext(DirectoryListContext);
// }


// useDirectoryListFetch is a hook that fetches directories from the api
export function useDirectoryListFetch({path = "",} = {}){

    const  dirCtx = useContext(DirListContext);

    const url = "http://localhost:8080/api/v0/fe?path=" + path;
    const { data: apiData, loading, error } = useFetch(url);

    // Set the context data when the API data is loaded
    useEffect(() => {
        if (!loading && !error) {
            dirCtx.setFiles(apiData);
        }
        console.log("effect")
    }, [loading, error, apiData,dirCtx.path]);
    return { loading, error, dirCtx };

}

async function processDir(response) {
    let jsonData = await parseJsonResp(response);
    let releases = jsonData["_embedded"].releases.map((r) => new DirectoryList(r));
    return {
        totalNumberOfItems: jsonData["_totalNumberOfItems"],
        releases: releases,
    };
}
