export function CreateRecord(record) {
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

export function UpdateRecord(record){
    return fetch(process.env.REACT_APP_API_URL + '/record/' + record.id, {
        method:'PUT',
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

export function LoadFields(){
    return fetch(process.env.REACT_APP_API_URL + '/field', {
        method: 'GET'
    }).then(response => {
        if(!response.ok)
            throw Error(response.status)
        
        return response
    }).catch(function(response){
        return response
    });
}


export function UpdateFields(fields){
    return fetch(process.env.REACT_APP_API_URL + '/field', {
        method:'PUT',
        headers: {'Content-Type':'application/json'},
        body: JSON.stringify(
            fields
        )
    }).then(response => {
        if(!response.ok){
            throw Error(response.status)
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
}

export function DeleteField(fieldId){
    return fetch(process.env.REACT_APP_API_URL + '/field/' + fieldId, {
        method: 'DELETE'
    }).then(response => {
        if(!response.ok)
            throw Error(response.status)
        
        return response
    }).catch(function(response){
        return response
    });
}