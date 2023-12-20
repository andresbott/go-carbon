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

// custom components
import ReleaseStatuses from "./ReleaseStatuses";

// Release represents a release from release registry
export default class Release {
    constructor(json) {
        this._json = json;
        this._statuses = new ReleaseStatuses(this._json?.status);
    }

    // return a clone of the release
    clone() {
        // let that = this;
        let clone = Object.assign(
            Object.create(Object.getPrototypeOf(this)),
            this
        );
        clone._statuses = new ReleaseStatuses(this._json?.status);
        return clone;
    }

    // return the main Release Id
    get id() {
        return this._json?.ssg?.id;
    }

    // return an object {k:v} where keys represent the ssg environments and the values the respective ids
    get ssgIds() {
        let obj = {};
        if (this._json?.ssg?.id !== "undefined") {
            obj = { prod: this._json?.ssg?.id };
        }
        return obj;
    }

    // return the creation date of the release
    get createdAt() {
        let d = new Date(this._json?.createdAt);
        return d.toUTCString();
    }

    // return the quickstart commit id based on the string in base image
    get qsCommit() {
        let base = this._json?.source?.ssgRelease?.baseImage;

        let commit = base.split(":");
        if (commit.length === 2) {
            return commit[1];
        }
        return "";
    }

    // get details about the components stored in SSG for this release
    get ssgRelease() {
        let ssgRelease = {
            baseImage: "",
            buildContextUrl: "",
            buildImage: "",
            chartUrl: "",
            dispatcherImage: "",
            emptyAuthorImage: "",
            emptyPublishImage: "",
            publishFarmerImage: "",
        };

        // check if the keys in ssgRelease exist in _json.source.ssgRelease and replace
        Object.keys(ssgRelease).forEach((key) => {
            if (this._json?.source?.ssgRelease[key] !== "undefined") {
                ssgRelease[key] = this._json.source.ssgRelease[key];
            }
        });
        return ssgRelease;
    }

    // get an array of test as stored in SSG
    get ssgTests() {
        if (this._json?.source?.ssgRelease?.tests !== "undefined") {
            // we assume that the json payload has the structure: [{id:"..","type":".."},...]
            return this._json.source.ssgRelease.tests;
        } else {
            return [];
        }
    }

    // get an array of addons as stored in SSG
    get ssgAddons() {
        if (this._json?.source?.ssgRelease?.testedWithAddons !== "undefined") {
            // we assume that the addon structure is correct, if this changes we can make the needed transformations here
            return this._json.source.ssgRelease.testedWithAddons;
        } else {
            return [];
        }
    }

    // get array of tools defined for release
    get ssgTools() {
        if (this._json?.source?.ssgRelease?.tools !== "undefined") {
            return this._json.source.ssgRelease.tools;
        }
        return [];
    }

    // return a list of the IT test executed on this release
    get itTests() {
        if (this._json?.source?.tests?.it !== "undefined") {
            // we assume that the json payload has the structure: {"it":[{"artifact":"...",...},{...}] }
            return this._json?.source?.tests?.it;
        } else {
            return {};
        }
    }

    // return a list of the IT test executed on this release
    get uiTests() {
        if (this._json?.source?.tests?.ui !== "undefined") {
            // we assume that the json payload has the structure: {"it":[{"artifact":"...",...},{...}] }
            return this._json?.source?.tests?.ui;
        } else {
            return {};
        }
    }

    // return a list of the IT test executed on this release
    get securityTests() {
        if (this._json?.source?.tests?.sec !== "undefined") {
            // we assume that the json payload has the structure: {"it":[{"artifact":"...",...},{...}] }
            return this._json?.source?.tests?.sec;
        } else {
            return {};
        }
    }

    // get array of bundles defined for release
    get bundles() {
        if (this._json?.source?.bundles !== "undefined") {
            return this._json.source.bundles;
        }
        return [];
    }

    // return the raw json as formatted string
    get rawJson() {
        return JSON.stringify(this._json, null, 2);
    }

    // return array of statuses (ReleaseStatus) of the release
    get statuses() {
        return this._statuses;
    }

    // get patch information
    get patch() {
        if (this._json?.patch !== "undefined") {
            // we assume that the json payload has the structure: [{id:"..","type":".."},...]
            return this._json.patch;
        } else {
            return {};
        }
    }

    // get an array of patches for the current release
    get patchedBy() {
        if (this._json?.patchedBy !== "undefined") {
            return this._json.patchedBy;
        } else {
            return [];
        }
    }

    // get an array of patch branches for the current release
    get patchBranches() {
        if (this._json?.patchBranches !== "undefined") {
            return this._json.patchBranches;
        } else {
            return [];
        }
    }

    // get current visibility
    get visibility() {
        if (this._json?.status?.visibility) {
            return this._json.status.visibility;
        } else {
            return {};
        }
    }
}
