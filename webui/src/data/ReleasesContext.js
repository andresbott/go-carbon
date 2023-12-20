// ADOBE CONFIDENTIAL
// ___________________
//
// Copyright 2021 Adobe
// All Rights Reserved.
//
// NOTICE: All information contained herein is, and remains
// the property of Adobe and its suppliers, if any. The intellectual
// and technical concepts contained herein are proprietary to Adobe
// and its suppliers and are protected by all applicable intellectual
// property laws, including trade secret and copyright laws.
// Dissemination of this information or reproduction of this material
// is strictly forbidden unless prior written permission is obtained
// from Adobe.

// react
import React, { useContext, useState } from "react";

// 21st century data structures
import Dict from "collections/dict";

import { processJson, useFetch } from "../util/fetch";
import Release from "./Release";

// custom context
export const ReleasesContext = React.createContext({});

// custom context hook to write and read selected release data all over the application
export function useReleasesContext() {
    return useContext(ReleasesContext);
}

// context component
export function ReleasesProvider({ children }) {
    // Releases List state variables
    const [releases, setReleases] = useState(new Dict());
    const [selectedReleaseId, setSelectedReleaseId] = useState(null);
    const [selected, setSelected] = useState(null);

    // context object
    const releaseContext = {
        releases: releases,
        setReleases: setReleases,
        selectedReleaseId: selectedReleaseId,
        setSelectedReleaseId: setSelectedReleaseId,
        selected: selected,
        setSelected: setSelected,
        updateSelectedStatus: setSelectedStatus,
    };

    // force propagate status change on the selected release
    function setSelectedStatus(status) {
        let clone = selected.clone();
        clone.statuses.setStatus(status.key, status.raw);
        setSelected(clone);
    }

    return (
        <ReleasesContext.Provider value={releaseContext}>
            {children}
        </ReleasesContext.Provider>
    );
}

// useReleasesFetch is a hook responsible for gathering the releases to be displayed,
// based on the passed arguments as filter
export function useReleasesFetch({
    imsToken = "",
    count = 10,
    start = 0,
    releases = null,
    showOnlyRCs = true,
} = {}) {
    let url = "/api/v0/releases?limit=" + count + "&start=" + start;
    if (releases !== null) {
        url += "&query=" + releases.replaceAll(/\s/g, "");
    }
    if (showOnlyRCs) {
        url += "&rcsonly=true";
    }
    let opts = {
        headers: {
            Authorization: `Bearer ${imsToken}`,
        },
    };

    return useFetch(url, opts, processReleases);
}

export async function processReleases(response) {
    let jsonData = await processJson(response);
    let releases = jsonData["_embedded"].releases.map((r) => new Release(r));
    return {
        totalNumberOfItems: jsonData["_totalNumberOfItems"],
        releases: releases,
    };
}

export async function processSingleRelease(response) {
    let jsonData = await processJson(response);
    return new Release(jsonData);
}
