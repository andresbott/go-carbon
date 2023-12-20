import './FileExplorer.css';
import {DirectoryListContext, useDirectoryListFetch} from "../../data/DirListContext";
import {useContext, useEffect} from "react";
import {useFetch} from "../../util/fetch";
import DirList from "../../components/DirList/DirList";
import ComponentOne from "../../components/ComponentOne";
import ComponentTwo from "../../components/componentTwo";
import ZustandComp from "../../components/zustandComponent";


function App() {
    // const { data, setData } = useContext(DirectoryListContext);
    // const { data: apiData, loading, error } = useFetch('http://localhost:8080/api/v0/fe');
    //
    // // Set the context data when the API data is loaded
    // useEffect(() => {
    //     if (!loading && !error) {
    //         console.log(apiData)
    //         setData(apiData);
    //     }
    // }, [loading, error, apiData]);
    //
    // if (loading) {
    //     console.log("loading")
    //     return <p>Loading...</p>;
    // }
    //
    // if (error) {
    //     console.log("error")
    //     return <p>Error occurred: {error.message}</p>;
    // }


    // const {loading, error, data} = useDirectoryListFetch()
    // console.log({loading,error, data})

    // if (loading) {
    //     return <p>Loading...</p>;
    // }
    //
    // if (error) {
    //     return <p>Error occurred: {error.message}</p>;
    // }
    return (
        <div className="FileExplorer">
            <header>
                head
            </header>
            <div className="menu">
                left side
                {/*<DirList/>*/}
            </div>
            <div className="main">
                right side
                <div>
                    {/*<ComponentOne/>*/}
                    {/*<ComponentTwo/>*/}
                    <ZustandComp/>
                </div>


            </div>
            <footer>
                foot
            </footer>
        </div>
    );
}


export default App;
