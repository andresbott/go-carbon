export async function parseJsonResp(response) {

    if (response.ok){
        try {
            return response.json();
        }catch(err) {
            throw {
                name: "JsonParseErr",
                message: "Failed to parse JSON response",
                code: response.status
            }
        }
    }
    // in case of error response we still expect a json payload
    const errObj = async () => {
        const payload = await response.text();
        try{
            const jsonPayload = JSON.parse(payload)
            if (!jsonPayload.error){
                let msg = "";
                if (jsonPayload.hasOwnProperty("msg")){
                    msg =  jsonPayload.msg
                }
               return  {
                    name: "ServerErr",
                    message: `Status code (${response.status}) Message:${msg}`,
                    code: response.status
                }
            }else{
                return  {
                    name: "ServerErr",
                    message: `Status code (${response.status}) Message:${payload}`,
                    code: response.status
                }
            }
        }catch (err){
            return  {
                name: "ServerErr",
                message: `Status code (${response.status}) Message: Unknown error`,
                code: response.status
            }
        }
    }
    throw await errObj()
}