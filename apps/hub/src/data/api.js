export function CreateRecord(record) {
    console.log("Create called")

        return fetch(process.env.REACT_APP_API_URL + '/record', {
            method:'POST',
            headers: {'Content-Type':'application/json'},
            body: JSON.stringify(
                record
            )
        }).then(response => {
            if(!response.ok){
                throw Error(response.statusText)
            }
            return response
        }).catch(function(){
            return null
        })
}

export function UpdateField(fieldId, value){
    const fieldToUpdate = {
        id: fieldId,
        value: value
    }
    console.log(JSON.stringify(fieldToUpdate));
}

export function LoadFields(){
    return fetch(process.env.REACT_APP_API_URL + '/field', {
        method: 'GET'
    }).then(response => {
        if(!response.ok)
            throw Error(response.statusText)
        
        return response
    }).catch(function(){
        return null
    });
}