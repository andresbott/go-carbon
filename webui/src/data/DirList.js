
export default class DirectoryList {
    constructor(json) {
        this._data = json;

    }

    get files() {

        let files = [];

        for (let i = 0; i < this._data.items.length; i++) {
            // Access the value you want to add to the array
            const value = this._data[i].name;

            // Add the value to the array
            files.push(value);
        }

        return "some files"


    }
}
