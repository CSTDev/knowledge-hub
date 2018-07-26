export function GetRecords(bounds){
    let queryString = "minLat=" + bounds._southWest.lat + "&minLng=" + bounds._southWest.lng + "&maxLat=" + bounds._northEast.lat + "&maxLng=" + bounds._northEast.lng
    return fetch(window.APP_CONFIG.API_URL + '/record?' + queryString, {
        method: 'GET'
    }).then(response => {
        if(!response.ok)
            throw Error(response.status)
        return response
    }).catch(function(response){
        return response
    });
}

export function CreateRecord(record) {
        return fetch(window.APP_CONFIG.API_URL + '/record', {
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
    return fetch(window.APP_CONFIG.API_URL + '/record/' + record.id, {
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

export function DeleteRecord(id){
    return fetch(window.APP_CONFIG.API_URL + '/record/' + id, {
        method: 'DELETE',
        headers: {'Content-Type':'application/json'},
    }).then(response => {
        if(!response.ok){
            throw Error(response.statusText)
        }
        return response
    }).catch(function(){
        return null;
    })
}

export function LoadFields(){
    return fetch(window.APP_CONFIG.API_URL + '/field', {
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
    return fetch(window.APP_CONFIG.API_URL + '/field', {
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
    return fetch(window.APP_CONFIG.API_URL + '/field/' + fieldId, {
        method: 'DELETE'
    }).then(response => {
        if(!response.ok)
            throw Error(response.status)
        
        return response
    }).catch(function(response){
        return response
    });
}