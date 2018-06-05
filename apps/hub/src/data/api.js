import React from 'react';

export function CreateRecord(record) {
    console.log("Create called")
    console.dir(record)
    console.log(JSON.stringify({
        record
    }))
    try{

        
        //TODO It's passing it as {"record":{"title": "a Title"}} etc. so service doesn't recognise it.


        return fetch('http://localhost:8000/v1/record', {
            method:'POST',
            headers: {'Content-Type':'application/json'},
            body: JSON.stringify({
                record
            })
        }).then(response => {
            return response
        })
        
    } catch(err) {
        return null
    }
    
}